package router

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/structs"
	"Gateway311/engine/telemetry"
)

const (
	rpcChanSize = 10
	rpcTimeout  = time.Second * 3 // 3 seconds
)

var (
	showRunTimes               = true
	showResponse               = true
	showResponseDetail         = false
	rpcSID             sidType = 1
)

func getRPCID() int64 {
	return rpcSID.get()
}

var serviceMethods = map[structs.NRequestType]string{
	structs.NRTUnknown:      "",
	structs.NRTServicesAll:  "Services.All",
	structs.NRTServicesArea: "Services.Area",
	structs.NRTCreate:       "Report.Create",
	structs.NRTSearchLL:     "Report.SearchLL",
	structs.NRTSearchDID:    "Report.SearchDID",
	structs.NRTSearchRID:    "Report.SearchRID",
}

var newResponse map[structs.NRequestType]func() interface{}

func init() {
	newResponse = make(map[structs.NRequestType]func() interface{})
	newResponse[structs.NRTServicesAll] = func() interface{} { return new(structs.NServicesResponse) }
	newResponse[structs.NRTServicesArea] = func() interface{} { return new(structs.NServicesResponse) }
	newResponse[structs.NRTCreate] = func() interface{} { return new(structs.NCreateResponse) }
	newResponse[structs.NRTSearchLL] = func() interface{} { return new(structs.NSearchResponse) }
	newResponse[structs.NRTSearchDID] = func() interface{} { return new(structs.NSearchResponse) }
	newResponse[structs.NRTSearchRID] = func() interface{} { return new(structs.NSearchResponse) }
}

// =======================================================================================
//                                      MANAGER
// =======================================================================================

// RPCCallMgr manages the RPC call(s) to service a request.
type RPCCallMgr struct {
	reqmgr        requester
	serviceMethod string
	results       chan structs.NRoute
	pending       int64 // The count of pending RPC calls.

	calls map[structs.NRoute]*rpcCall
	errs  []error
}

type requester interface {
	Routes() structs.NRoutes
	RType() structs.NRequestType
	Data() interface{}
	Processer() func(ndata interface{}) error
}

// NewRPCCallMgr returns a RPCCallMgr instance, with routes as per the requester interface.
func NewRPCCallMgr(reqmgr requester) (*RPCCallMgr, error) {
	if len(reqmgr.Routes()) == 0 {
		return nil, fmt.Errorf("no routes")
	}
	r := &RPCCallMgr{
		reqmgr:        reqmgr,
		serviceMethod: serviceMethods[reqmgr.RType()],
		results:       make(chan structs.NRoute, 5),
		calls:         make(map[structs.NRoute]*rpcCall),
	}

	for _, route := range reqmgr.Routes() {
		rpccall, err := newrpcCall(r, route)
		if err != nil {
			log.Errorf("%s", err.Error())
			return nil, err
		}
		r.calls[route] = rpccall
	}
	// log.Debug("New RPC Call Mgr:\n%s", r)

	return r, nil
}

func (r *RPCCallMgr) send() {
	for _, call := range r.calls {
		err := call.run()
		if err != nil {
			r.errs = append(r.errs, err)
			continue
		}
		r.incPending()
	}
}

func (r *RPCCallMgr) receive() {
	// Responses are serialized via the results channel
	// Collect responses via the "r.results" channel.
	if r.pending > 0 {
		var timedout bool
		timeout := common.TimeoutChan(rpcTimeout)
		for !timedout {
			select {
			case respKey := <-r.results:
				answer := r.calls[respKey]
				r.decPending()
				// log.Debug("%s", answer.String())
				telemetry.SendRPC(answer.response.(structs.NResponser).GetIDS(), "done", "", time.Now())
				if answer.err != nil {
					r.errs = append(r.errs, answer.err)
					log.Errorf("RPC call to: %q failed - %s", respKey, answer.err)
					break
				}
				// log.Debug("Answer: %s", answer.response)
				err := r.reqmgr.Processer()(answer.response)
				if err != nil {
					r.errs = append(r.errs, answer.err)
					log.Errorf("RPC call to: %q failed - %s", respKey, answer.err)
					break
				}

			case timedout = <-timeout:
			}

			if r.pending == 0 {
				break
			}
		}
		if timedout {
			for route, call := range r.calls {
				if !call.replied {
					log.Errorf("Adapter: %q timed out", route.SString())
				}
			}
		}
	}
	if showResponse {
		log.Debug("RPC Manager:%v", r)
	}
}

// Run executes all RPC calls.
func (r *RPCCallMgr) Run() error {
	startTime := time.Now()

	r.send()

	r.receive()

	if showRunTimes {
		log.Info("RPC Call: %q took: %s", r.serviceMethod, time.Since(startTime))
	}

	if len(r.errs) > 0 {
		var errs string
		for i, e := range r.errs {
			log.Errorf(e.Error())
			if i > 0 {
				errs = errs + "; "
			}
			errs = errs + e.Error()
		}
		return errors.New(errs)
	}
	return nil
}

// -------------------------------- rpcmanager Interface ---------------------------------

func (r *RPCCallMgr) rType() structs.NRequestType {
	return r.reqmgr.RType()
}

func (r *RPCCallMgr) data() interface{} {
	return r.reqmgr.Data()
}

func (r *RPCCallMgr) service() string {
	return r.serviceMethod
}

func (r *RPCCallMgr) resultChan() chan structs.NRoute {
	return r.results
}

func (r *RPCCallMgr) incPending() int64 {
	return atomic.AddInt64((&r.pending), 1)
}

func (r *RPCCallMgr) decPending() int64 {
	return atomic.AddInt64((&r.pending), -1)
}

func (r *RPCCallMgr) processer() func(ndata interface{}) error {
	return r.reqmgr.Processer()
}

// -------------------------------- String ---------------------------------
func (r RPCCallMgr) String() string {
	ls := new(common.LogString)
	ls.AddS("RPCCallMgr\n")
	ls.AddF("ReqMgr: %p  type: %v\n", r.reqmgr, r.reqmgr.RType().String())
	ls.AddF("Pending: %v\n", r.pending)
	for _, v := range r.calls {
		ls.AddS(v.String())
	}
	for i, e := range r.errs {
		if i == 0 {
			ls.AddS("--Errors--\n")
		}
		ls.AddS(e.Error())
	}
	return ls.Box(90)
}

// =======================================================================================
//                                      RPC CALL
// =======================================================================================

// rpcCall manages the data for each RPC call.
type rpcCall struct {
	rpc rpcmanager
	adp AdpRPCer

	sync.Mutex
	id       int64
	route    structs.NRoute
	response interface{}
	sent     bool
	replied  bool
	err      error
}

// String returns a representation of an rpcCall instance.
func (r rpcCall) String() string {
	ls := new(common.LogString)
	ls.AddF("rpcCall - %d\n", r.id)
	ls.AddF("Service: %v\n", r.rpc.service())
	ls.AddF("Adapter: %v  Route: %v\n", r.adp.AdpID(), r.route.String())
	ls.AddF("Sent: %t  replied: %t\n", r.sent, r.replied)
	if r.err != nil {
		ls.AddF("Error: %v\n", r.err.Error())
	}
	return ls.Box(70)
}

type rpcmanager interface {
	rType() structs.NRequestType
	data() interface{}
	service() string
	resultChan() chan structs.NRoute
	incPending() int64
	decPending() int64
	processer() func(ndata interface{}) error
}

func newrpcCall(rpcmgr rpcmanager, route structs.NRoute) (*rpcCall, error) {
	r := &rpcCall{
		id:    getRPCID(),
		rpc:   rpcmgr,
		route: route,
	}
	// log.Debug("Route: %v", route.String())
	adp, err := GetRouteAdapter(r.route)
	if err != nil {
		return nil, err
	}
	r.adp = adp
	// log.Debug(r.String())
	return r, nil
}

func (r *rpcCall) setSent() {
	r.Lock()
	defer r.Unlock()
	r.sent = true
}

func (r *rpcCall) prepRPC() (rqstCopy interface{}, err error) {
	prep := func(d structs.NRequester) {
		d.SetID(0, r.id)
		d.SetRoute(r.route)
	}
	switch v := r.rpc.data().(type) {
	case *structs.NServiceRequest:
		rCopy := *v
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debug("Sending: %s", rCopy.String())
	case *structs.NCreateRequest:
		rCopy := *v
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debug("Sending: %s", rCopy.String())
	case *structs.NSearchRequestLL:
		rCopy := *v
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debug("Sending: %s", rCopy.String())
	case *structs.NSearchRequestDID:
		rCopy := *v
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debug("Sending: %s", rCopy.String())
	case *structs.NSearchRequestRID:
		rCopy := *v
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debug("Sending: %s", rCopy.String())
	default:
		msg := fmt.Sprintf("Invalid type in send RPC: %T", r.rpc.data())
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	return rqstCopy, nil
}

func (r *rpcCall) run() error {
	if r.adp.Connected() {
		// log.Debug("{%s::%v} Calling prepRPC...", r.route.String(), r.id)
		payload, err := r.prepRPC()
		if err != nil {
			return err
		}
		r.setSent()
		telemetry.SendRPC(payload.(structs.NRequester).GetIDS(), "open", r.route.SString(), time.Now())
		go func() {
			response := newResponse[r.rpc.rType()]()
			r.err = r.adp.Call(r.rpc.service(), payload, response)
			r.Lock()
			defer r.Unlock()
			r.response = response
			r.replied = true
			r.rpc.resultChan() <- r.route
		}()
	}
	return nil
}
