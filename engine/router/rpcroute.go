package router

import (
	"fmt"

	"Gateway311/engine/structs"

	// "github.com/davecgh/go-spew/spew"
)

// =======================================================================================
//                                      RPC ROUTE MAP
// =======================================================================================

var (
	routeMap map[string]*routeMapMethods
)

type routeMapType map[structs.NRoute]*routeMapMethods

func init() {
	routeMap = make(map[string]*routeMapMethods)

	routeMap["Services.All"] = &routeMapMethods{}
	routeMap["Services.Area"] = &routeMapMethods{}
	routeMap["Report.Create"] = &routeMapMethods{}
	routeMap["Report.SearchDeviceID"] = &routeMapMethods{}
	routeMap["Report.SearchLL"] = &routeMapMethods{}

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
	routeMap["Report.SearchDeviceID"].newResponse = func() interface{} { return new(structs.NSearchResponse) }
	routeMap["Report.SearchLL"].newResponse = func() interface{} { return new(structs.NSearchResponse) }

	return nil
}

func initRPCList() error {
	showAdpMap := func(adpMap map[structs.NRoute]*rpcAdapterStatus) {
		for k, v := range adpMap {
			log.Debug("%s%s", k, v)
		}
	}

	adapter := func(rt structs.NRouter, service string) (map[structs.NRoute]*rpcAdapterStatus, error) {
		log.Debug("[RouteMap: adapter] service: %q\nroutes: %s\n", service, rt.GetRoutes())
		m := make(map[structs.NRoute]*rpcAdapterStatus)
		for _, nroute := range rt.GetRoutes() {
			adp, err := GetAdapter(nroute.AdpID)
			// log.Debug("adp: %s", adp)
			if err != nil {
				return nil, fmt.Errorf("Error creating Adapter list - %s", err)
			}
			rs, err := newAdapterStatus(adp, service, nroute)
			// log.Debug("rs: %s", rs)
			if err != nil {
				return nil, fmt.Errorf("Error creating Adapter list - %s", err)
			}
			m[nroute] = rs
			// log.Debug("adapters: %s", m)
		}
		showAdpMap(m)
		return m, nil
	}

	// statusList populates r.adpList with pointers to Adapters that service the specified
	// Area.
	// area := func(areaID, service string) (map[structs.NRoute]*rpcAdapterStatus, error) {
	area := func(rt structs.NRouter, service string) (map[structs.NRoute]*rpcAdapterStatus, error) {
		log.Debug("[RouteMap: area] service: %q\nroutes: %s\n", service, rt.GetRoutes())
		adpStatList := make(map[structs.NRoute]*rpcAdapterStatus)
		for _, nroute := range rt.GetRoutes() {
			switch nroute.RouteType() {
			case structs.NRtTypAllAdapters:
				log.Debug("Using ALL adapters")
				for _, adp := range adapters.Adapters {
					route := structs.NRoute{AdpID: adp.ID, AreaID: "all", ProviderID: 0}
					adpStat, err := newAdapterStatus(adp, service, route)
					if err != nil {
						return nil, fmt.Errorf("Error creating route list - %s", err)
					}
					adpStatList[route] = adpStat
				}
			case structs.NRtTypArea:
				log.Debug("Area route: %q", nroute.String())
				routes, err := GetAreaRoutes(nroute.AreaID)
				if err != nil {
					return nil, fmt.Errorf("Cannot create the Adapter List - no routes for area: %s", nroute.AreaID)
				}
				for _, route := range routes {
					adp, err := GetAdapter(route.AdpID)
					if err != nil {
						return nil, fmt.Errorf("Invalid adpater id in route: %s", route)
					}
					adpStat, err := newAdapterStatus(adp, service, route)
					if err != nil {
						return nil, fmt.Errorf("Error creating route list - %s", err)
					}
					adpStatList[route] = adpStat
				}
			case structs.NRtTypFull:
				log.Debug("Full route: %q", nroute)
				adp, err := GetAdapter(nroute.AdpID)
				if err != nil {
					return nil, fmt.Errorf("Invalid adpater id in route: %s", nroute)
				}
				adpStat, err := newAdapterStatus(adp, service, nroute)
				if err != nil {
					return nil, fmt.Errorf("Error creating route list - %s", err)
				}
				adpStatList[nroute] = adpStat

			default:
				// log.Debug("Using only adapters for areaID: %s", areaID)
				return nil, fmt.Errorf("Cannot create the Adapter List - invalid route: %s", nroute)
			}
		}
		showAdpMap(adpStatList)
		return adpStatList, nil
	}

	routeMap["Services.All"].buildAdapterList = func(r structs.NRouter) (map[structs.NRoute]*rpcAdapterStatus, error) {
		return area(r, "Services.All")
	}
	routeMap["Services.Area"].buildAdapterList = func(r structs.NRouter) (map[structs.NRoute]*rpcAdapterStatus, error) {
		return area(r, "Services.Area")
	}
	routeMap["Report.Create"].buildAdapterList = func(r structs.NRouter) (map[structs.NRoute]*rpcAdapterStatus, error) {
		return adapter(r, "Report.Create")
	}
	routeMap["Report.SearchDeviceID"].buildAdapterList = func(r structs.NRouter) (map[structs.NRoute]*rpcAdapterStatus, error) {
		return area(r, "Report.SearchDeviceID")
	}
	routeMap["Report.SearchLL"].buildAdapterList = func(r structs.NRouter) (map[structs.NRoute]*rpcAdapterStatus, error) {
		return area(r, "Report.SearchLL")
	}
	return nil
}

// =======================================================================================
//                                      RPC ROUTE METHODS
// =======================================================================================
type routeMapMethods struct {
	newResponse      func() interface{}
	buildAdapterList func(structs.NRouter) (map[structs.NRoute]*rpcAdapterStatus, error)
}
