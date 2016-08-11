package router

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/open311-gateway/engine/common"
	"github.com/open311-gateway/engine/structs"
	"github.com/open311-gateway/engine/telemetry"

	log "github.com/jeffizhungry/logrus"
)

const (
	rpcChanSize = 5
	rpcTimeout  = time.Second * 3 // 3 seconds
)

var (
	showRunTimes = true
	showResponse = true
)

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
		results:       make(chan structs.NRoute, rpcChanSize),
		calls:         make(map[structs.NRoute]*rpcCall),
	}

	for _, route := range reqmgr.Routes() {
		rpccall, err := newrpcCall(r, route)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		r.calls[route] = rpccall
	}

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
				telemetry.SendRPC(answer.response.(structs.NResponser).GetIDS(), "done", "", time.Now())
				if answer.err != nil {
					r.errs = append(r.errs, answer.err)
					log.WithFields(log.Fields{
						"method": r.serviceMethod,
						"route":  respKey.String(),
						"error":  answer.err,
					}).Error("RPC call failed.")
					break
				}
				err := r.reqmgr.Processer()(answer.response)
				if err != nil {
					r.errs = append(r.errs, err)
					log.WithFields(log.Fields{
						"method": r.serviceMethod,
						"route":  respKey.String(),
						"error":  err,
					}).Error("RPC Processor() failed.")
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
					log.WithFields(log.Fields{
						"method": r.serviceMethod,
						"route":  route.String(),
					}).Error("Adapter call timedout.")
				}
			}
		}
	}
	if showResponse {
		log.Debugf("RPC Manager:%v", r)
	}
}

// Run executes all RPC calls.
func (r *RPCCallMgr) Run() error {
	startTime := time.Now()

	// Initiate all RPC calls
	r.send()

	// Process responses
	r.receive()

	if showRunTimes {
		log.WithFields(log.Fields{
			"method":   r.serviceMethod,
			"duration": time.Since(startTime),
		}).Info("RPC call duration.")
	}

	if len(r.errs) > 0 {
		var errs string
		for i, e := range r.errs {
			log.WithFields(log.Fields{
				"method": r.serviceMethod,
				"error":  e.Error(),
			}).Error("RPC Run() failed.")
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
func (r *rpcCall) String() string {
	r.Lock()
	defer r.Unlock()
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
		id:    common.RPCID(),
		rpc:   rpcmgr,
		route: route,
	}
	adp, err := GetRouteAdapter(r.route)
	if err != nil {
		return nil, err
	}
	r.adp = adp
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
	switch data := r.rpc.data().(type) {
	case *structs.NServiceRequest:
		rCopy := *data
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debugf("Sending: %s", rCopy.String())
	case *structs.NCreateRequest:
		rCopy := *data
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debugf("Sending: %s", rCopy.String())
	case *structs.NSearchRequestLL:
		rCopy := *data
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debugf("Sending: %s", rCopy.String())
	case *structs.NSearchRequestDID:
		rCopy := *data
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debugf("Sending: %s", rCopy.String())
	case *structs.NSearchRequestRID:
		rCopy := *data
		prep(&rCopy)
		rqstCopy = &rCopy
		log.Debugf("Sending: %s", rCopy.String())
	default:
		msg := fmt.Sprintf("Invalid type in send RPC: %T", r.rpc.data())
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	return rqstCopy, nil
}

func (r *rpcCall) run() error {
	if r.adp.Connected() {
		payload, err := r.prepRPC()
		if err != nil {
			return err
		}
		r.setSent()
		telemetry.SendRPC(payload.(structs.NRequester).GetIDS(), "open", r.route.String(), time.Now())
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
