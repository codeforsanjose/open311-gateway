package router

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/structs"
)

const (
	rpcChanSize = 10
	rpcTimeout  = time.Second * 3 // 3 seconds
)

var (
	showRunTimes = true
	// routeMap     map[string]routeMapMethods
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
	r := RPCCall{
		service: service,
		request: request,
		results: make(chan *rpcAdapterStatus, rpcChanSize),
		listIF:  make(map[string]*rpcAdapterStatus),
		process: process,
		errs:    make([]error, 0),
	}

	log.Debug("%+v", r)
	router, ok := request.(structs.NRouter)
	log.Debug("%+v  ok: %t", router, ok)
	if !ok {
		return nil, fmt.Errorf("The request (type: %s) does not implement structs.Router", reflect.TypeOf(request))
	}
	log.Debug("Building listIF...")
	routeMethods, ok := routeMap[service]
	if !ok {
		return nil, fmt.Errorf("RouteMap does not exist for service: %s", service)
	}
	listIF, err := routeMethods.buildAdapterList(router)
	if err != nil {
		return nil, fmt.Errorf("Unable to prep RPCCall for request: %s - %s", reflect.TypeOf(request), err)
	}

	r.listIF = listIF

	log.Debug("RPCCall: %s", r)
	return &r, nil
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
	listIF    map[string]*rpcAdapterStatus // Key: AdapterID
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
				r.listIF[answer.adapter.ID] = answer
				r.processes--
				if answer.err != nil {
					r.errs = append(r.errs, answer.err)
					log.Errorf("RPC call to: %q failed - %s", answer.adapter.ID, answer.err)
					break
				}
				log.Debug("Answer: %s", answer.response)
				r.process(answer.response)

			case timedout = <-timeout:
			}

			if r.processes == 0 {
				break
			}
		}
		if timedout {
			for k, v := range r.listIF {
				if v == nil {
					log.Errorf("Adapter: %q timed out", k)
				}
			}
		}
	}
	if showRunTimes {
		log.Info("RPC Call: %q took: %s", r.service, time.Since(sendTime))
	}
	return nil
}

// setAdapters builds the list of Adapters that will be called with the Request.
func (r *RPCCall) setAdapters() error {

	return nil
}

func (r *RPCCall) send() error {
	for k, v := range r.listIF {
		if v.adapter.Connected {
			// Give the pointer to the AdapterStatus to the go routine.
			var pAdpStat *rpcAdapterStatus
			pAdpStat, r.listIF[k] = v, nil
			r.processes++
			go func(pas *rpcAdapterStatus) {
				// fmt.Printf("Inside go routine:\n%s\n%s\n", pas, r)
				pas.err = pas.adapter.Client.Call(r.service, r.request, pas.response)
				r.results <- pas
			}(pAdpStat)
		} else {
			log.Warning("Skipping: %s - not connected!", v.adapter.ID)
		}
	}

	fmt.Printf("After Run():\n%s\n", r)
	return nil
}

/*
func (r *RPCCall) adapter(adpID string) error {
	adp, err := GetAdapter(adpID)
	log.Debug("adp: %s", adp)
	rs, err := newAdapterStatus(adp, r.service)
	log.Debug("rs: %s", rs)
	if err != nil {
		return fmt.Errorf("Error creating Adapter list - %s", err)
	}
	r.listIF[adp.ID] = rs
	log.Debug("RPCCall: %s", r)
	return err
}

// statusList populates r.listIF with pointers to Adapters that service the specified
// Area.
func (r *RPCCall) statusList(areaID string) error {
	var al []*Adapter
	if strings.ToLower(areaID) == "all" {
		log.Debug("Using ALL adapters")
		for _, v := range adapters.Adapters {
			al = append(al, v)
		}
	} else {
		log.Debug("Using only adapters for areaID: %s", areaID)
		var ok bool
		al, ok = adapters.areaAdapters[areaID]
		if !ok {
			return fmt.Errorf("Area %q is not supported on this Gateway", areaID)
		}
	}
	for _, adp := range al {
		rs, err := newAdapterStatus(adp, r.service)
		if err != nil {
			return fmt.Errorf("Error creating Adapter list - %s", err)
		}
		r.listIF[adp.ID] = rs
	}
	return nil
}
*/

// --------------------------------- rpcAdapterStatus -----------------------------------
type rpcAdapterStatus struct {
	adapter  *Adapter
	response interface{}
	sent     bool
	replied  bool
	err      error
}

func newAdapterStatus(adp *Adapter, service string) (*rpcAdapterStatus, error) {
	aStat := &rpcAdapterStatus{
		adapter: adp,
		sent:    false,
		replied: false,
	}

	rs := routeMap[service].newResponse()
	aStat.response = rs
	log.Debug("aStat: %s", aStat)
	return aStat, nil
}

// func makeResponse(service string) (interface{}, error) {
// 	switch service {
// 	case "Service.All", "Service.Area":
// 		return new(structs.NServicesResponse), nil
// 	case "Create.Run":
// 		return new(structs.NCreateResponse), nil
// 	case "Search.DeviceID", "Search.Location":
// 		return new(structs.SearchResp), nil
// 	default:
// 		return nil, fmt.Errorf("Invalid request type: %q", service)
// 	}
// }

// ==============================================================================================================================
//                                      STRINGS
// ==============================================================================================================================

func (r RPCCall) String() string {
	ls := new(common.LogString)
	ls.AddS("RPC Call\n")
	ls.AddF("Service: %s\n", r.service)
	ls.AddF("Request: (%p)  results chan: (%p)\n", &r.request, r.results)
	ls.AddF("%s\n", r.request)
	ls2 := new(common.LogString)
	ls2.AddS("Adapters\n")
	ls2.AddF("         Name/Type/Address                 Sent  Repl     Response\n")
	for k, v := range r.listIF {
		ls2.AddF("  %4s: %s\n", k, v)
	}
	ls.AddF("%s", ls2.Box(80))
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
	return ls.Box(90)
}

func (r rpcAdapterStatus) String() string {
	s := fmt.Sprintf("%-30s     %5t %5t   (%s)%p", fmt.Sprintf("%s (%s @%s)", r.adapter.ID, r.adapter.Type, r.adapter.Address), r.sent, r.replied, reflect.TypeOf(r.response), r.response)
	// s += spew.Sdump(r.response) + "\n"
	return s
}
