package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/logs"
	"Gateway311/engine/router"
	"Gateway311/engine/structs"
)

const (
	startupDelay = time.Second * 5
)

var (
	servicesData cache
	log          = logs.Log
)

// GetArea returns a list of Services for the specified Area.
func GetArea(areaID string) (structs.NServices, error) {
	return servicesData.getArea(areaID)
}

// Refresh initiates a refresh of the Services Cache.
func Refresh() {
	servicesData.refresh()
}

// Shutdown should be called at system shutdown.  It will terminate the update channel, and
// permform any other necessary cleanup.
func Shutdown() {
	servicesData.shutdown()
}

// ==============================================================================================================================
//                                      SERVICE CACHE
// ==============================================================================================================================

// cache is the cache for Services data.  The "list" is the active list.  The
// active list is a simple copy of either list0 or list1 (i.e. it's effectively a
// pointer to the underlying data in list0 or list1).  If the activeSet is 0, then list1
// is cleared and is available for loading.  Vice versa for activeSet1.
type cache struct {
	list      [2]map[string]structs.NServices // Index: AreaID
	activeSet int
	update    chan bool // Update request queue
	sync.RWMutex
}

// getArea retrieves the ServiceList for the specified area.
func (r *cache) getArea(areaID string) (structs.NServices, error) {
	r.RLock()
	defer r.RUnlock()
	l, ok := r.list[r.activeSet][areaID]
	if !ok {
		return nil, fmt.Errorf("The requested AreaID: %q is not serviced by this gateway.", areaID)
	}
	return l, nil
}

func (r *cache) processRefresh() {
	rqst := &structs.NServiceRequest{
		NRequestCommon: structs.NRequestCommon{
			Rtype: structs.NRTServicesAll,
		},
		Area: "all",
	}
	rpcCall, err := router.NewRPCCall("Services.All", rqst, servicesData.merge)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = rpcCall.Run()
	if err != nil {
		log.Error(err.Error())
		return
	}
	r.indexAreaAdapters()
	r.sendRoutes()
	r.switchSet()
	// log.Debug("Cache%s", r)
}

// sendRoutes builds a unique list of all NRoutes and posts it to the
// router.GetChRouteUpd() channel.
func (r *cache) sendRoutes() error {
	ls := r.loadSet()

	// Build a unique list of NRoutes.
	allRoutes := make(map[structs.NRoute]bool)
	for _, nservices := range r.list[ls] {
		for _, nservice := range nservices {
			allRoutes[nservice.GetRoute()] = true
		}
	}
	log.Debug("Sending ALL route data")
	router.GetChRouteUpd() <- allRoutes
	return nil
}

// indexAreaAdapters indexes all Adapters by the AreaID, and sends to router via
// the router.GetChAreaAdp() channel
func (r *cache) indexAreaAdapters() error {
	ls := r.loadSet()
	type mk struct {
		areaID, adp string
	}

	// Build a unique list of [AreaID, AdpID] combinations.
	areaAdp := make(map[mk]bool)
	for _, nservices := range r.list[ls] {
		for _, nservice := range nservices {
			areaAdp[mk{nservice.AreaID, nservice.AdpID}] = true
		}
	}

	// Convert the above list to a map of a list of AdpID for each AreaID.
	alist := make(map[string][]string)
	for k := range areaAdp {
		if _, ok := alist[k.areaID]; !ok {
			alist[k.areaID] = make([]string, 0)
		}
		alist[k.areaID] = append(alist[k.areaID], k.adp)
	}

	router.GetChAreaAdp() <- alist

	return nil
}

func (r *cache) loadSet() (ls int) {
	if r.activeSet == 0 {
		ls = 1
	}
	return
}

func (r *cache) clearLoadSet(ds int) {
	r.list[r.loadSet()] = make(map[string]structs.NServices)
}

// refresh initiates a service list update - that is, it requests the current service lists
// from all Adapters, merges the data back into the standby service list cache, and makes
// the freshly loaded cache the active cache.
func (r *cache) refresh() {
	r.update <- true
}

func (r *cache) switchSet() {
	r.Lock()
	defer r.Unlock()
	if r.activeSet == 0 {
		log.Debug("[ServicesCache] Switched from list 0 to 1")
		r.activeSet = 1
		r.clearLoadSet(0)
	} else {
		log.Debug("[ServicesCache] Switched from list 1 to 0")
		r.activeSet = 0
		r.clearLoadSet(1)
	}
}

func (r *cache) merge(ndata interface{}) error {
	data := (ndata.(*structs.NServicesResponse)).Services
	r.Lock()
	defer r.Unlock()
	var loadSet int
	switch r.activeSet {
	case 0:
		loadSet = 1
	case 1:
		loadSet = 0
	default:
		msg := fmt.Sprintf("Invalid ServiceData activeSet: %v", r.activeSet)
		log.Fatal(msg)
		return errors.New(msg)
	}
	for _, ns := range data {
		if _, ok := r.list[loadSet][ns.AreaID]; !ok {
			log.Info("Created Area: %q", ns.AreaID)
			r.list[loadSet][ns.AreaID] = make(structs.NServices, 0)
		}
		r.list[loadSet][ns.AreaID] = append(r.list[loadSet][ns.AreaID], ns)
		// log.Debug("   Appending: %s - %s", ns.MID(), ns.Name)
	}
	return nil
}

func (r *cache) init() {
	r.Lock()
	defer r.Unlock()
	r.list[0] = make(map[string]structs.NServices)
	r.list[1] = make(map[string]structs.NServices)

	r.update = make(chan bool, 1)
	r.activeSet = 0

	go func() {
		for {
			_, ok := <-r.update
			if !ok {
				log.Info("Terminating cache refresh queue...")
				break
			} else {
				log.Debug("Running Services refresh...")
				r.processRefresh()
			}
		}
	}()
}

// Shutdown should be called at system shutdown.  It will terminate the update channel, and
// permform any other necessary cleanup.
func (r *cache) shutdown() {
	close(r.update)
}

// ==============================================================================================================================
//                                      STRINGS
// ==============================================================================================================================

// Displays the contents of the Spec_Type custom type.
func (r cache) String() string {
	ls := new(common.LogString)
	ls.AddF("cache [%d]\n", r.activeSet)
	ls.AddS("------- Service List --------\n")
	for k, v := range r.list[r.activeSet] {
		ls.AddF("<<<<<Area: %s >>>>>%s", k, v)
	}
	return ls.Box(90)
}

// ==============================================================================================================================
//                                      MISC
// ==============================================================================================================================

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

func init() {
	servicesData.init()
}
