package router

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/structs"
)

const (
	startupDelay = time.Second * 5
)

var (
	servicesData ServicesData
)

func RefreshServicesList() {
	servicesData.Refresh()
}

// func ServiceList(areaID string)

// ==============================================================================================================================
//                                      SERVICE CACHE
// ==============================================================================================================================

// ServicesData is the cache for services data.  The "list" is the active list.  The
// active list is a simple copy of either list0 or list1 (i.e. it's effectively a
// pointer to the underlying data in list0 or list1).  If the activeList is 0, then list1
// is cleared and is available for loading.  Vice versa for activeList1.
type ServicesData struct {
	list       [2]map[string]structs.NServices // Index: AreaID
	activeList int
	update     chan bool // Update request queue
	sync.RWMutex
}

func (sd *ServicesData) init() {
	sd.Lock()
	defer sd.Unlock()
	sd.list[0] = make(map[string]structs.NServices)
	sd.list[1] = make(map[string]structs.NServices)
	sd.update = make(chan bool, 1)
	sd.activeList = 0

	go func() {
		for {
			_, ok := <-sd.update
			if !ok {
				break
			} else {
				log.Debug("Running refresh...")
				sd.refresh()
			}
		}
	}()

	go func() {
		// Pause a few seconds to make sure everything else is up and running.
		time.Sleep(startupDelay)
		sd.Refresh()
	}()
}

// Shutdown should be called at system shutdown.  It will terminate the update channel, and
// permform any other necessary cleanup.
func (sd *ServicesData) Shutdown() {
	close(sd.update)
}

// Refresh initiates a service list update - that is, it requests the current service lists
// from all Adapters, merges the data back into the standby service list cache, and makes
// the freshly loaded cache the active cache.
func (sd *ServicesData) Refresh() {
	sd.update <- true
}

func (sd *ServicesData) refresh() {
	rqst := &structs.NServiceRequest{"all"}
	r, err := newRPCCall("Service.All", "all", rqst, servicesData.merge)
	if err != nil {
		log.Error(err.Error())
		return
	}
	r.run()
	if err != nil {
		log.Error(err.Error())
		return
	}
	sd.switchList()
}

func (sd *ServicesData) switchList() {
	sd.Lock()
	defer sd.Unlock()
	if sd.activeList == 0 {
		log.Debug("Switched from list 0 to 1")
		sd.activeList = 1
		sd.list[0] = make(map[string]structs.NServices)
	} else {
		log.Debug("Switched from list 1 to 0")
		sd.activeList = 0
		sd.list[1] = make(map[string]structs.NServices)
	}
}

func (sd *ServicesData) merge(ndata interface{}) error {
	data := (ndata.(*structs.NServicesResponse)).Services
	sd.Lock()
	defer sd.Unlock()
	var loadList int
	switch sd.activeList {
	case 0:
		loadList = 1
	case 1:
		loadList = 0
	default:
		msg := fmt.Sprintf("Invalid ServiceData activeList: %v", sd.activeList)
		log.Fatal(msg)
		return errors.New(msg)
	}
	for _, ns := range data {
		if _, ok := sd.list[loadList][ns.AreaID]; !ok {
			log.Debug("Created City: %q", ns.AreaID)
			sd.list[loadList][ns.AreaID] = make(structs.NServices, 0)
		}
		sd.list[loadList][ns.AreaID] = append(sd.list[loadList][ns.AreaID], ns)
		// log.Debug("   Appending: %s - %s", ns.MID(), ns.Name)
	}
	return nil
}

// ==============================================================================================================================
//                                      STRINGS
// ==============================================================================================================================

// Displays the contents of the Spec_Type custom type.
func (sd ServicesData) String() string {
	ls := new(common.LogString)
	ls.AddF("ServicesData [%d]\n", sd.activeList)
	for k, v := range sd.list[sd.activeList] {
		ls.AddF("<<<<<City: %s >>>>>%s", k, v)
	}
	return ls.Box(90)
}

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
	servicesData.init()
}
