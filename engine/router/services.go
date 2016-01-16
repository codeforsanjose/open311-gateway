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
)

var (
	servicesData ServicesData
)

const (
	rpcChanSize = 50
	rpcTimeout  = time.Second * 3 // 3 seconds
)

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

type rpcStatus struct {
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
func callRPCs(serviceMethod string, request interface{}, reply []interface{}) (chan *rpc.Call, map[string]*rpcStatus, int, error) {
	var (
		done      chan *rpc.Call
		errs      string
		processes int
	)
	done = make(chan *rpc.Call, rpcChanSize)
	listIF := make(map[string]*rpcStatus)

	for i, a := range adapters.Adapters {
		if a.Connected {
			listIF[a.Name] = &rpcStatus{true, false}
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
