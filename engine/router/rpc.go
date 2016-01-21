package router

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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
)

// =======================================================================================
//                                      RPC
// =======================================================================================

// newRPCCall creates a new rpcCalls.  rpcCalls holds all information about an RPC call,
// including which Adapters are called, request, response status, etc.
func newRPCCall(service, areaID string, request interface{}, process func(interface{}) error) (*rpcCalls, error) {
	r := rpcCalls{
		service: service,
		request: request,
		results: make(chan *rpcAdapterStatus, rpcChanSize),
		listIF:  make(map[string]*rpcAdapterStatus),
		process: process,
		errs:    make([]error, 0),
	}
	r.statusList(areaID)
	return &r, nil
}

// ResponseProcesser is the interface for
type ResponseProcesser interface {
	Process(interface{}) error
}

type rpcCalls struct {
	service   string
	request   interface{}
	results   chan *rpcAdapterStatus
	processes int
	listIF    map[string]*rpcAdapterStatus
	process   func(interface{}) error
	errs      []error
}

func (r *rpcCalls) run() error {
	// Send all the RPC calls in go routines.
	var startTime time.Time
	if showRunTimes {
		startTime = time.Now()
	}
	err := r.start()
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
				r.listIF[answer.adapter.Name] = answer
				r.processes--
				if answer.err != nil {
					r.errs = append(r.errs, answer.err)
					log.Errorf("RPC call to: %q failed - %s", answer.adapter.Name, answer.err)
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
		log.Info("RPC Call: %q took: %s", r.service, time.Since(startTime))
	}
	return nil
}

func (r *rpcCalls) start() error {
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
			log.Info("Skipping: %s - not connected!", v.adapter.Name)
		}
	}

	fmt.Printf("After run():\n%s\n", r)
	return nil
}

//
func (r *rpcCalls) statusList(areaID string) error {
	var al []*Adapter
	if strings.ToLower(areaID) == "all" {
		al = adapters.Adapters
		log.Debug("using all")
	} else {
		if len(adapters.areaAdapters) == 0 {
			return fmt.Errorf("Area %q is not supported on this Gateway", areaID)
		}
		al = adapters.areaAdapters[areaID]
		log.Debug("using list for: %s", areaID)
	}
	for _, adp := range al {
		rs, err := newAdapterStatus(adp, r.service)
		if err != nil {
			return fmt.Errorf("Error creating Adapter list - %s", err)
		}
		r.listIF[adp.Name] = rs
	}
	return nil
}

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

	rs, err := makeResponseStruct(service)
	if err != nil {
		return nil, fmt.Errorf("Cannot create AdapterStatus - %s", err)
	}
	aStat.response = rs
	return aStat, nil
}

// ==============================================================================================================================
//                                      STRINGS
// ==============================================================================================================================

func (r rpcCalls) String() string {
	ls := new(common.LogString)
	ls.AddS("RPC Call\n")
	ls.AddF("Service: %s\n", r.service)
	ls.AddF("Request: (%p)  results chan: (%p)\n", r.request, r.results)
	ls.AddF("Request: (%p):\n%[1]s\n", r.request)
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
	s := fmt.Sprintf("%-30s     %5t %5t   (%s)%p", fmt.Sprintf("%s (%s @%s)", r.adapter.Name, r.adapter.Type, r.adapter.Address), r.sent, r.replied, reflect.TypeOf(r.response), r.response)
	// s += spew.Sdump(r.response) + "\n"
	return s
}

// ==============================================================================================================================
//                                      MISC
// ==============================================================================================================================

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
