package router

import (
	"errors"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/structs"

	"github.com/davecgh/go-spew/spew"
)

var (
	servicesData ServicesData
)

const (
	rpcChanSize = 10
	rpcTimeout  = time.Second * 3 // 3 seconds
)

// =======================================================================================
//                                      RPC
// =======================================================================================

// ResponseProcesser is the interface for
type ResponseProcesser interface {
	Process(interface{}) error
}

type rpcCalls struct {
	service   string
	request   interface{}
	result    chan *rpcAdapterStatus
	processes int
	listIF    map[string]*rpcAdapterStatus
	process   ResponseProcesser
	errs      []error
}

// newRPCCall creates a new rpcCalls.  rpcCalls holds all information about an RPC call,
// including which Adapters are called, request, response status, etc.
func newRPCCall(service, areaID string, request interface{}, process ResponseProcesser) (*rpcCalls, error) {
	r := rpcCalls{
		service: service,
		request: request,
		result:  make(chan *rpc.Call, rpcChanSize),
		listIF:  make(map[string]*rpcAdapterStatus),
		process: process,
		errs:    make([]error, 0),
	}
	r.statusList(areaID)
	return &r, nil
}

//
func (r *rpcCalls) statusList(areaID string) error {
	if len(adapters.cityAdapters) == 0 {
		return fmt.Errorf("Area %q is not supported on this Gateway", areaID)
	}
	for _, adp := range adapters.cityAdapters[areaID] {
		rs, err := newAdapterStatus(adp, r.service)
		if err != nil {

		}
		r.listIF[adp.Name] = rs
	}
	return nil
}

// Run executes the RPC calls:
//
// 1. All RPC calls are started in go routines.
// 2. The go routines send the results back on the "result" channel.
// 3. The "result" channel is selected until the timeout period has expired.
// 4. Any results not received, and erroneous results, are error logged.
// 5. The "processResponse" interface is called with the results.
func (r *rpcCalls) run() error {
	// Launch all RPC calls.
	for i, a := range r.listIF {
		log.Debug("")
	}
	// 2. The go routines send the results back on the "result" channel.
	// 3. The "result" channel is selected until the timeout period has expired.
	// 4. Any results not received, and erroneous results, are error logged.
	// 5. The "processResponse" interface is called with the results.

	return nil
}

func (r *rpcCalls) callAdapter(name string) error {
	adpStat := r.listIF[name]
	rStat, err := newRPCStatus(r.service)
	if err != nil {
		msg := fmt.Sprintf("Unable to call adapter - %s", err)
		log.Error(msg)
		return errors.New(msg)
	}

	go func() {
		err := adpStat.adapter.Client.Call(r.service, r.request, rStat.resp)
		rStat.err = err
		r.result <- rStat
	}()

	return nil
}

func (r rpcCalls) String() string {
	ls := new(common.LogString)
	ls.AddS("RPC Call\n")
	ls.AddF("Service: %s\n", r.service)
	ls.AddF("Request: %p  done: %p\n", r.request, r.result)
	ls.AddF("Request: %s\n", spew.Sdump(r.request))
	ls.AddF("  Name Sent  Repl Type  Address \n")
	for _, v := range r.listIF {
		ls.AddF("%s\n", v)
	}
	ls.AddF("Processes: %d\n%s\n", r.processes, r.listIF)
	ls.AddF("Process interface: %p\n", r.process)
	if len(r.errs) == 0 {
		ls.AddS("NO ERRORS!\n")
	} else {
		ls.AddS("Errors\n")
		for _, e := range r.errs {
			ls.AddF("\t%s\n", e.Error())
		}
	}
	return ls.Box(80)
}

// --------------------------------- rpcAdapterStatus -----------------------------------
type rpcStatus struct {
	err  error
	resp interface{}
}

func newRPCStatus(service string) (*rpcStatus, error) {
	resp, err := makeResponseStruct(service)
	if err != nil {
		return nil, err
	}
	return &rpcStatus{
		resp: resp,
	}, nil
}

func makeResponseStruct(service string) (interface{}, error) {
	switch service {
	case "Service.All", "Service.Area":
		return new(structs.NServicesResponse), nil
	case "Create.Report":
		return new(structs.NCreateResponse), nil
	case "Search.DeviceID", "Search.Location":
		return new(structs.SearchResp), nil
	default:
		return nil, fmt.Errorf("Invalid request type: %q", service)
	}
}

// --------------------------------- rpcAdapterStatus -----------------------------------
type rpcAdapterStatus struct {
	adapter  *Adapter
	response interface{}
	sent     bool
	replied  bool
}

func (r *rpcAdapterStatus) callAdapter(result chan *rpc.Call) error {

	return nil
}

func newAdapterStatus(adp *Adapter, service string) (*rpcAdapterStatus, error) {
	aStat := &rpcAdapterStatus{
		adapter: adp,
		sent:    false,
		replied: false,
	}

	rs, err := makeResponseStruct(service)
	if err != nil {

	}

	return rs, nil
}

func (r rpcAdapterStatus) String() string {
	s := fmt.Sprintf("    %-5s %5t %5t %-6s  %s\n", r.adapter.Name, r.sent, r.replied, r.adapter.Type, r.adapter.Address)
	s += spew.Sdump(r.response) + "\n"
	return s
}

// ==============================================================================================================================
//                                      SERVICE CACHE
// ==============================================================================================================================

// ServicesData is the cache for services data
type ServicesData struct {
	list map[string]structs.NServices
	sync.RWMutex
}

func (sd *ServicesData) clearMap() {
	servicesData.list = make(map[string]structs.NServices)
}

// Displays the contents of the Spec_Type custom type.
func (sd ServicesData) String() string {
	ls := new(common.LogString)
	ls.AddS("ServicesData\n")
	for k, v := range sd.list {
		ls.AddF("<<<<<City: %s >>>>>\ns", k, v)
	}
	return ls.Box(90)
}

func (sd *ServicesData) merge(data *structs.NServices) error {
	for _, ns := range *data {
		if _, ok := sd.list[ns.AreaID]; !ok {
			log.Debug("Created City: %q", ns.AreaID)
			sd.list[ns.AreaID] = make(structs.NServices, 0)
		}
		sd.list[ns.AreaID] = append(sd.list[ns.AreaID], ns)
		// log.Debug("   Appending: %s - %s", ns.MID(), ns.Name)
	}
	return nil
}

// Refresh updates the Services cache by requesting fresh data from all connected Adapters.

func (sd *ServicesData) Refresh() error {
	serviceCall := "Service.All"
	log.Info("Refreshing Services List...")
	creq := structs.NServiceRequest{
		City: "",
	}
	// log.Debug("%+v\n", creq)

	var responses []interface{}
	for range adapters.Adapters {
		responses = append(responses, &structs.NServicesResponse{})
	}

	done, listIF, processes, _ := callRPCs(serviceCall, &creq, responses)
	log.Debug("  Processes: %d", processes)

	if processes > 0 {
		sd.Lock()
		defer sd.Unlock()
		timeout := common.TimeoutChan(rpcTimeout)

		timedout := false
		for !timedout {
			select {
			case answer := <-done:
				processes--
				if answer.Error != nil {
					log.Errorf("RPC call %q failed: %s\n", serviceCall, answer.Error)
					break
				}
				r := answer.Reply.(*structs.NServicesResponse)
				listIF[r.IFID].replied = true
				log.Debug("\tMerging Services list for %q", r.IFID)
				sd.merge(&r.Services)
			case timedout = <-timeout:
			}

			if processes <= 0 {
				break
			}
		}
		failed := false
		for k, v := range listIF {
			switch {
			case v.sent && !v.replied:
				log.Errorf("RPC call %q to: %s timed out.", serviceCall, k)
				failed = true
			case !v.sent:
				log.Errorf("Engine and Adapter configs are out of synch - received reply from: %s", k)
				failed = true
			}
		}
		if failed {
			log.Error("Service List refresh FAILED!")
		} else {
			log.Info("Service List refresh complete...")
		}
	}
	// log.Debug(spew.Sdump(sd))
	return nil
}

// callRPCs sends an API request to all connected Adapters. It returns a response
// channel of type *rpc.Call, then number of RPC calls to expect to return, and
// any immediate errors encountered with the call to rpc.Client.Go().
func (rc *rpcCalls) call(serviceMethod string, request interface{}, reply []interface{}) (chan *rpc.Call, map[string]*rpcCall, int, error) {
	var (
		done      chan *rpc.Call
		errs      string
		processes int
	)
	done = make(chan *rpc.Call, rpcChanSize)
	listIF := make(map[string]*rpcCall)

	for i, a := range adapters.Adapters {
		if a.Connected {
			listIF[a.Name] = &rpcCall{true, false}
			replyCall := a.Client.Go(serviceMethod, request, reply[i], done)
			processes++
			if replyCall.Error != nil {
				msg := fmt.Sprintf("Error calling adapter: %s: %s\n", a.Name, replyCall.Error)
				log.Error(msg)
				errs = errs + msg
			}
		}
	}
	if len(errs) > 0 {
		return done, listIF, processes, errors.New(errs)
	}
	return done, listIF, processes, nil
}

/*
type rpcCall struct {
	sent    bool
	replied bool
}

// Refresh updates the Services cache by requesting fresh data from all connected Adapters.
func (sd *ServicesData) Refresh() error {
	serviceCall := "Service.All"
	log.Info("Refreshing Services List...")
	creq := structs.NServiceRequest{
		City: "",
	}
	// log.Debug("%+v\n", creq)

	var responses []interface{}
	for range adapters.Adapters {
		responses = append(responses, &structs.NServicesResponse{})
	}

	done, listIF, processes, _ := callRPCs(serviceCall, &creq, responses)
	log.Debug("  Processes: %d", processes)

	if processes > 0 {
		sd.Lock()
		defer sd.Unlock()
		timeout := common.TimeoutChan(rpcTimeout)

		timedout := false
		for !timedout {
			select {
			case answer := <-done:
				processes--
				if answer.Error != nil {
					log.Errorf("RPC call %q failed: %s\n", serviceCall, answer.Error)
					break
				}
				r := answer.Reply.(*structs.NServicesResponse)
				listIF[r.IFID].replied = true
				log.Debug("\tMerging Services list for %q", r.IFID)
				sd.merge(&r.Services)
			case timedout = <-timeout:
			}

			if processes <= 0 {
				break
			}
		}
		failed := false
		for k, v := range listIF {
			switch {
			case v.sent && !v.replied:
				log.Errorf("RPC call %q to: %s timed out.", serviceCall, k)
				failed = true
			case !v.sent:
				log.Errorf("Engine and Adapter configs are out of synch - received reply from: %s", k)
				failed = true
			}
		}
		if failed {
			log.Error("Service List refresh FAILED!")
		} else {
			log.Info("Service List refresh complete...")
		}
	}
	// log.Debug(spew.Sdump(sd))
	return nil
}

// callRPCs sends an API request to all connected Adapters. It returns a response
// channel of type *rpc.Call, then number of RPC calls to expect to return, and
// any immediate errors encountered with the call to rpc.Client.Go().
func callRPCs(serviceMethod string, request interface{}, reply []interface{}) (chan *rpc.Call, map[string]*rpcCall, int, error) {
	var (
		done      chan *rpc.Call
		errs      string
		processes int
	)
	done = make(chan *rpc.Call, rpcChanSize)
	listIF := make(map[string]*rpcCall)

	for i, a := range adapters.Adapters {
		if a.Connected {
			listIF[a.Name] = &rpcCall{true, false}
			replyCall := a.Client.Go(serviceMethod, request, reply[i], done)
			processes++
			if replyCall.Error != nil {
				msg := fmt.Sprintf("Error calling adapter: %s: %s\n", a.Name, replyCall.Error)
				log.Error(msg)
				errs = errs + msg
			}
		}
	}
	if len(errs) > 0 {
		return done, listIF, processes, errors.New(errs)
	}
	return done, listIF, processes, nil
}
*/

// ==============================================================================================================================
//                                      MISC
// ==============================================================================================================================
// SplitMID breaks down an MID, and returns the IFID and AreaID.
func SplitMID(mid string) (string, string, error) {
	parts := strings.Split(mid, "-")
	log.Debug("MID: %+v\n", parts)
	if len(parts) != 4 {
		return "", "", fmt.Errorf("Invalid MID: %s", mid)
	}
	return parts[0], parts[1], nil
}

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

func init() {
	servicesData.clearMap()
}
