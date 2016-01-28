package router

import (
	"_sketches/spew"
	"fmt"
	"strings"

	"Gateway311/engine/structs"
)

var (
	routeMap routeMapType
)

func init() {
	routeMap.init()
	fmt.Printf("---------- routeMap -------------\n%s\n", spew.Sdump(routeMap))
}

// =======================================================================================
//                                      RPC ROUTE MAP
// =======================================================================================
type routeMapType map[string]*routeMapMethods

func (rm routeMapType) init() error {
	rm = make(map[string]*routeMapMethods)

	rm["Services.All"] = &routeMapMethods{}
	rm["Services.Area"] = &routeMapMethods{}
	rm["Report.Create"] = &routeMapMethods{}
	rm["Report.SearchDeviceID"] = &routeMapMethods{}
	rm["Report.SearchLocation"] = &routeMapMethods{}

	if err := rm.initResponseStructs(); err != nil {
		return err
	}

	if err := rm.initRPCList(); err != nil {
		return err
	}

	return nil
}

func (rm routeMapType) initResponseStructs() error {
	for service, methods := range rm {
		switch service {
		case "Services.All":
			methods.newResponse = func() interface{} { return new(structs.NServicesResponse) }
		case "Services.Area":
			methods.newResponse = func() interface{} { return new(structs.NServicesResponse) }
		case "Report.Create":
			methods.newResponse = func() interface{} { return new(structs.NCreateResponse) }
		case "Report.SearchDeviceID":
			methods.newResponse = func() interface{} { return new(structs.SearchResp) }
		case "Report.SearchLocation":
			methods.newResponse = func() interface{} { return new(structs.SearchResp) }
		default:
			return fmt.Errorf("Unknown request: %q - unable to initialize Route Map", service)
		}
	}
	return nil
}

func (rm routeMapType) initRPCList() error {
	adapter := func(rt structs.NRouter, service string) (map[string]*rpcAdapterStatus, error) {
		m := make(map[string]*rpcAdapterStatus)
		adp, err := GetAdapter(rt.Route().IFID)
		log.Debug("adp: %s", adp)
		rs, err := newAdapterStatus(adp, service)
		log.Debug("rs: %s", rs)
		if err != nil {
			return nil, fmt.Errorf("Error creating Adapter list - %s", err)
		}
		m[adp.ID] = rs
		log.Debug("adapters: %s", m)
		return m, err
	}

	// statusList populates r.listIF with pointers to Adapters that service the specified
	// Area.
	// area := func(areaID, service string) (map[string]*rpcAdapterStatus, error) {
	area := func(rt structs.NRouter, service string) (map[string]*rpcAdapterStatus, error) {
		var al []*Adapter
		m := make(map[string]*rpcAdapterStatus)
		areaID := rt.Route().AreaID
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

	for service, methods := range rm {
		switch service {
		case "Services.All":
			methods.buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
				return area(r, service)
			}
		case "Services.Area":
			methods.buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
				return area(r, service)
			}
		case "Report.Create":
			methods.buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
				return adapter(r, service)
			}
		case "Report.SearchDeviceID":
			methods.buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
				return area(r, service)
			}
		case "Report.SearchLocation":
			methods.buildAdapterList = func(r structs.NRouter) (map[string]*rpcAdapterStatus, error) {
				return area(r, service)
			}
		default:
			return fmt.Errorf("Unknown request: %q - unable to initialize Route Map", service)
		}
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
