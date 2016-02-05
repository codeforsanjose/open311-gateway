package router

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/color"

	"Gateway311/engine/common"
	"Gateway311/engine/structs"
)

const (
	rpcChanSize = 10
	rpcTimeout  = time.Second * 3 // 3 seconds
)

var (
	showRunTimes = true
	showResponse = true
)

// =======================================================================================
//                                      RPC
// =======================================================================================

// NewRPCCall creates a new RPCCall.  RPCCall holds all information about an RPC call,
// including which Adapters are called, request, response status, etc.
//
// This call will set up an RPC call for a specific Adapter, or for any Adapter servicing
// the specified Area. Only one of the following should be used - the other parameter
// should be set to an empty string.
// adpID: adapter ID
// areaID: Area ID
func NewRPCCall(service string, request interface{}, process func(interface{}) error) (*RPCCall, error) {
	rpcCall := RPCCall{
		service: service,
		request: request,
		results: make(chan *rpcAdapterStatus, rpcChanSize),
		adpList: make(adapterRouteList),
		process: process,
		errs:    make([]error, 0),
	}

	// log.Debug("%+v", r)
	router, ok := request.(structs.NRouter)
	// log.Debug("%+v  ok: %t", router, ok)
	if !ok {
		return nil, fmt.Errorf("The request (type: %s) does not implement structs.NRouter", reflect.TypeOf(request))
	}
	// log.Debug("Building adpList...")
	routeMethods, ok := serviceMap[service]
	if !ok {
		return nil, fmt.Errorf("serviceMap does not exist for service: %s", service)
	}
	adpList, err := routeMethods.buildAdapterList(router)
	if err != nil {
		return nil, fmt.Errorf("Unable to prep RPCCall for request: %s - %s", reflect.TypeOf(request), err)
	}

	rpcCall.adpList = adpList

	// log.Debug("RPCCall: %s", r)
	return &rpcCall, nil
}

// ResponseProcesser is the interface for
type ResponseProcesser interface {
	Process(interface{}) error
}

// RPCCall represents an RPC call.  This may be calls to multiple Adapters.
type RPCCall struct {
	service   string
	request   interface{}
	results   chan *rpcAdapterStatus
	processes int
	adpList   adapterRouteList // Key: AdapterID
	process   func(interface{}) error
	errs      []error
}

// Run executes the RPC call(s).  It is synchronous - it will wait for all requestrs
// to return or timeout.
func (r *RPCCall) Run() error {
	// Send all the RPC calls in go routines.
	var sendTime time.Time
	if showRunTimes {
		sendTime = time.Now()
	}
	err := r.send()
	if err != nil {
		msg := fmt.Sprintf("Error starting RPC calls.")
		log.Error(msg)
		return errors.New(msg)
	}

	// Collect responses via the "r.results" channel.
	if r.processes > 0 {
		var timedout bool
		timeout := common.TimeoutChan(rpcTimeout)
		for !timedout {
			select {
			case answer := <-r.results:
				r.adpList[answer.route] = answer
				r.processes--
				if answer.err != nil {
					r.errs = append(r.errs, answer.err)
					log.Errorf("RPC call to: %q failed - %s", answer.adapter.ID, answer.err)
					break
				}
				// log.Debug("Answer: %s", answer.response)
				r.process(answer.response)

			case timedout = <-timeout:
			}

			if r.processes == 0 {
				break
			}
		}
		if timedout {
			for k, v := range r.adpList {
				if v == nil {
					log.Errorf("Adapter: %q timed out", k)
				}
			}
		}
	}
	if showRunTimes {
		log.Info("RPC Call: %q took: %s", r.service, time.Since(sendTime))
	}
	if showResponse {
		log.Debug("Response:%s", r)
	}
	return nil
}

func (r *RPCCall) send() error {
	for k, v := range r.adpList {
		if v.adapter.Connected {
			// Give the pointer to the AdapterStatus to the go routine.
			var pAdpStat *rpcAdapterStatus
			pAdpStat, r.adpList[k] = v, nil
			r.processes++
			go func(pAdpStat *rpcAdapterStatus, rqst interface{}) {
				// log.Debug("Calling adapter:\n%s\n", pAdpStat)
				// log.Debug("Request type: %T", rqst)
				var rqstCopy interface{}
				switch v := rqst.(type) {
				case *structs.NServicels
					rCopy := *v
					structs.NRequester(&rCopy).SetRoute(pAdpStat.route)
					rqstCopy = &rCopy
					log.Debug("Sending: %s", rCopy.String())
				case *structs.NCreatels
					rCopy := *v
					structs.NRequester(&rCopy).SetRoute(pAdpStat.route)
					rqstCopy = &rCopy
					log.Debug("Sending: %s", rCopy.String())
				case *structs.NSearchRequestLL:
					rCopy := *v
					structs.NRequester(&rCopy).SetRoute(pAdpStat.route)
					rqstCopy = &rCopy
					log.Debug("Sending: %s", rCopy.String())
				default:
					log.Errorf("Invalid type in send RPC: %T", rqst)
					return
				}

				pAdpStat.err = pAdpStat.adapter.Client.Call(r.service, rqstCopy, pAdpStat.response)
				r.results <- pAdpStat
			}(pAdpStat, r.request)
		} else {
			log.Warning("Skipping: %s - not connected!", v.adapter.ID)
		}
	}

	// log.Debug("After Run():\n%s\n", r)
	return nil
}

// =======================================================================================
//                                      ADAPTER STATUS
// =======================================================================================
type rpcAdapterStatus struct {
	adapter  *Adapter
	route    structs.NRoute
	response interface{}
	sent     bool
	replied  bool
	err      error
}

func newAdapterStatus(adp *Adapter, service string, route structs.NRoute) (*rpcAdapterStatus, error) {
	aStat := &rpcAdapterStatus{
		adapter: adp,
		route:   route,
		sent:    false,
		replied: false,
	}

	rs := serviceMap[service].newResponse()
	aStat.response = rs
	// log.Debug("aStat: %s", aStat)
	return aStat, nil
}

// =======================================================================================
//                                      ADAPTER ROUTE MAP
// =======================================================================================

// adapterRouteMap is used in the RPC system to keep track of what is being sent to each Adapter
// Route, and the reply status and content.  Each RPCCall has an adpaterRouteMap instance.
type adapterRouteList map[structs.NRoute]*rpcAdapterStatus

func newAdapterRouteList() adapterRouteList {
	return make(adapterRouteList)
}

// ==============================================================================================================================
//                                      STRINGS
// ==============================================================================================================================

func (r RPCCall) String() string {
	ls := new(common.LogString)
	ls.AddS("RPC Call\n")
	ls.AddF("Service: %s\n", r.service)
	ls.AddF("%s", r.request)
	ls.AddS(r.adpList.String())
	ls.AddF("Processes: %d\n", r.processes)
	ls.AddF("Process interface: (%p)\n", r.process)
	if len(r.errs) == 0 {
		ls.AddS("Error: NONE!\n")
	} else {
		ls.AddS("Errors\n")
		for _, e := range r.errs {
			ls.AddF("\t%s\n", e.Error())
		}
	}
	return ls.Box(120)
}

func (r rpcAdapterStatus) String() string {
	ls := new(common.LogString)
	ls.AddS("rpcAdapterStatus\n")
	ls.AddS(" Name/Type/Address                 Sent  Repl    ResponseType                     Route\n")
	ls.AddF("%-30s     %5t %5t   %-30s  %s", fmt.Sprintf("%s (%s @%s)", r.adapter.ID, r.adapter.Type, r.adapter.Address), r.sent, r.replied, fmt.Sprintf("(%T)", r.response), r.route.String())
	return ls.Box(100)
}

func (r rpcAdapterStatus) StringResp() string {
	ls := new(common.LogString)
	ls.AddS("rpcAdapterStatus\n")
	ls.AddS(" Name/Type/Address                 Sent  Repl    ResponseType                     Route\n")
	ls.AddF("%-30s     %5t %5t   %-30s  %s", fmt.Sprintf("%s (%s @%s)", r.adapter.ID, r.adapter.Type, r.adapter.Address), r.sent, r.replied, fmt.Sprintf("(%T)", r.response), r.route.String())
	ls.AddF("%s", r.response)
	return ls.Box(100)
}

func (r rpcAdapterStatus) StringNH() string {
	s := fmt.Sprintf("%-30s     %5t %5t   %-30s  %s", fmt.Sprintf("%s (%s @%s)", r.adapter.ID, r.adapter.Type, r.adapter.Address), r.sent, r.replied, fmt.Sprintf("(%T)", r.response), r.route.String())
	return s
}

func (r rpcAdapterStatus) StringNHResp() string {
	s := fmt.Sprintf("%-30s     %5t %5t   %-30s  %s", fmt.Sprintf("%s (%s @%s)\n", r.adapter.ID, r.adapter.Type, r.adapter.Address), r.sent, r.replied, fmt.Sprintf("(%T)", r.response), r.route.String())
	s += stringResponse(r.response)
	return s
}

func stringResponse(r interface{}) string {
	switch v := r.(type) {
	case *structs.NServicesResponse:
		return v.String()
	case *structs.NCreateResponse:
		return v.String()
	case *structs.NSearchResponse:
		return v.String()
	default:
		return fmt.Sprintf("Cannot show type: %T\n", r)
	}
}

func (r adapterRouteList) String() string {
	ls := new(common.LogString)
	ls.AddS("adapterRouteList\n")
	ls.AddF("    Route/Type/Address              Sent  Repl    ResponseType                     Route        (match)\n")
	for route, adpStat := range r {
		routeMatch := color.GreenString("OK")
		if route != adpStat.route {
			routeMatch = color.RedString("Mismatch!")
		}
		ls.AddF("%s %s\n", adpStat.StringNH(), routeMatch)
		if showResponse {
			ls.AddS(stringResponse(adpStat.response))
		}
	}
	return ls.Box(110)
}
