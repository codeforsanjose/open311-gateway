package router

import (
	"fmt"
	"strings"

	"Gateway311/engine/structs"

	// "github.com/davecgh/go-spew/spew"
)

// =======================================================================================
//                                      RPC ROUTE MAP
// =======================================================================================

var (
	routeMap map[string]*routeMapMethods
)

type routeMapType map[string]*routeMapMethods

func init() {
	routeMap = make(map[string]*routeMapMethods)

	routeMap["Services.All"] = &routeMapMethods{}
	routeMap["Services.Area"] = &routeMapMethods{}
	routeMap["Report.Create"] = &routeMapMethods{}
	routeMap["Report.SearchDeviceID"] = &routeMapMethods{}
	routeMap["Report.SearchLocation"] = &routeMapMethods{}

	if err := initResponseStructs(); err != nil {
		log.Critical("Unable to initialize the routeMap - %s", err)
		return
	}

	if err := initRPCList(); err != nil {
		log.Critical("Unable to initialize the routeMap - %s", err)
		return
	}

	// log.Debug("---------- routeMap -------------\n%s\n", spew.Sdump(routeMap))
	return
}

func initResponseStructs() error {

	routeMap["Services.All"].newResponse = func() interface{} { return new(structs.NServicesResponse) }
	routeMap["Services.Area"].newResponse = func() interface{} { return new(structs.NServicesResponse) }
	routeMap["Report.Create"].newResponse = func() interface{} { return new(structs.NCreateResponse) }
	routeMap["Report.SearchDeviceID"].newResponse = func() interface{} { return new(structs.SearchResp) }
	routeMap["Report.SearchLocation"].newResponse = func() interface{} { return new(structs.SearchResp) }

	return nil
}

func initRPCList() error {
	adapter := func(rt structs.NRouter, service string) (map[string]*rpcAdapterStatus, error) {
		m := make(map[string]*rpcAdapterStatus)
		adp, err := GetAdapter(rt.Route().AdpID)
		// log.Debug("adp: %s", adp)
		rs, err := newAdapterStatus(adp, service)
		// log.Debug("rs: %s", rs)
		if err != nil {
			return nil, fmt.Errorf("Error creating Adapter list - %s", err)
		}
		m[adp.ID] = rs
		// log.Debug("adapters: %s", m)
		return m, err
	}

	// statusList populates r.adpList with pointers to Adapters that service the specified
	// Area.
	// area := func(areaID, service string) (map[string]*rpcAdapterStatus, error) {
	area := func(rt structs.NRouter, service string) (map[string]*rpcAdapterStatus, error) {
		var al []*Adapter
		m := make(map[string]*rpcAdapterStatus)
		areaID := rt.Route().AreaID
		if strings.ToLower(areaID) == "all" {
			// log.Debug("Using ALL adapters")
			for _, v := range adapters.Adapters {
				al = append(al, v)
			}
		} else {
			// log.Debug("Using only adapters for areaID: %s", areaID)
			var ok bool
			al, ok = adapters.areaAdapters[areaID]
			if !ok {
				return nil, fmt.Errorf("Area %q is not supported on this Gateway", areaID)
			}
		}
		for _, adp := range al {
			rs, err := newAdapterStatus(adp, service)
			if err != nil {
				return nil, fmt.Errorf("Error creating Adapter list - %s", err)
			}
			m[adp.ID] = rs
		}
		return m, nil
	}

	routeMap["Services.All"].buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
		return area(r, "Services.All")
	}
	routeMap["Services.Area"].buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
		return area(r, "Services.Area")
	}
	routeMap["Report.Create"].buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
		return adapter(r, "Report.Create")
	}
	routeMap["Report.SearchDeviceID"].buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
		return area(r, "Report.SearchDeviceID")
	}
	routeMap["Report.SearchLocation"].buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
		return area(r, "Report.SearchLocation")
	}
	return nil
}

// =======================================================================================
//                                      RPC ROUTE METHODS
// =======================================================================================
type routeMapMethods struct {
	newResponse      func() interface{}
	buildAdapterList func(structs.NRouter) (map[string]*rpcAdapterStatus, error)
}
