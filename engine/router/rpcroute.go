package router

import (
	"fmt"

	"Gateway311/engine/structs"

	"github.com/davecgh/go-spew/spew"

	// "github.com/davecgh/go-spew/spew"
)

// =======================================================================================
//                                      RPC ROUTE MAP
// =======================================================================================

var (
	serviceMap map[string]*serviceMapMethods
)

func init() {
	serviceMap = make(map[string]*serviceMapMethods)

	serviceMap["Services.All"] = &serviceMapMethods{}
	serviceMap["Services.Area"] = &serviceMapMethods{}
	serviceMap["Report.Create"] = &serviceMapMethods{}
	serviceMap["Report.SearchDID"] = &serviceMapMethods{}
	serviceMap["Report.SearchLL"] = &serviceMapMethods{}

	if err := initResponseStructs(); err != nil {
		log.Critical("Unable to initialize the serviceMap - %s", err)
		return
	}

	if err := initRPCList(); err != nil {
		log.Critical("Unable to initialize the serviceMap - %s", err)
		return
	}

	// log.Debug("---------- serviceMap -------------\n%s\n", spew.Sdump(serviceMap))
	return
}

func initResponseStructs() error {

	serviceMap["Services.All"].newResponse = func() interface{} { return new(structs.NServicesResponse) }
	serviceMap["Services.Area"].newResponse = func() interface{} { return new(structs.NServicesResponse) }
	serviceMap["Report.Create"].newResponse = func() interface{} { return new(structs.NCreateResponse) }
	serviceMap["Report.SearchDID"].newResponse = func() interface{} { return new(structs.NSearchResponse) }
	serviceMap["Report.SearchLL"].newResponse = func() interface{} { return new(structs.NSearchResponse) }

	return nil
}

func initRPCList() error {
	adapter := func(rt structs.NRouter, service string) (adapterRouteList, error) {
		// log.Debug("[serviceMap: adapter] service: %q\nroutes: %s\n", service, rt.GetRoutes())
		adpStatList := newAdapterRouteList()
		for i, nroute := range rt.GetRoutes() {
			adp, err := GetAdapter(nroute.AdpID)
			// log.Debug("adp: %s", adp)
			if err != nil {
				return nil, fmt.Errorf("Error creating Adapter list - %s", err)
			}
			rs, err := newAdapterStatus(adp, service, nroute, i+1)
			// log.Debug("rs: %s", rs)
			if err != nil {
				return nil, fmt.Errorf("Error creating Adapter list - %s", err)
			}
			adpStatList[nroute] = rs
			// log.Debug("adapters: %s", adpStatList)
		}
		log.Debug(adpStatList.String())
		return adpStatList, nil
	}

	// statusList populates r.adpList with pointers to Adapters that service the specified
	// Area.
	// area := func(areaID, service string) (adapterRouteList, error) {
	area := func(rt structs.NRouter, service string) (adapterRouteList, error) {
		// log.Debug("[serviceMap: area] service: %q\nroutes: %s\n", service, rt.GetRoutes())
		adpStatList := newAdapterRouteList()
		log.Debug("Routes: %+v", rt.GetRoutes())
		for i, nroute := range rt.GetRoutes() {
			switch nroute.RouteType() {
			case structs.NRtTypAllAdapters:
				log.Debug("Using ALL adapters")
				for _, adp := range adapters.Adapters {
					route := structs.NRoute{AdpID: adp.ID, AreaID: "all", ProviderID: 0}
					adpStat, err := newAdapterStatus(adp, service, route, i+1)
					if err != nil {
						return nil, fmt.Errorf("Error creating route list - %s", err)
					}
					adpStatList[route] = adpStat
				}
			case structs.NRtTypArea:
				log.Debug("Area route: %q", nroute.SString())
				routes, err := GetAreaRoutes(nroute.AreaID)
				if err != nil {
					return nil, fmt.Errorf("Cannot create the Adapter List - no routes for area: %s", nroute.AreaID)
				}
				log.Debug("Routes: %s\n", spew.Sdump(routes))
				for _, route := range routes {
					adp, err := GetAdapter(route.AdpID)
					if err != nil {
						return nil, fmt.Errorf("Invalid adpater id in route: %s", route)
					}
					adpStat, err := newAdapterStatus(adp, service, route, i+1)
					if err != nil {
						return nil, fmt.Errorf("Error creating route list - %s", err)
					}
					log.Debug("adpStat: %s", spew.Sdump(adpStat))
					adpStatList[route] = adpStat
				}
			case structs.NRtTypFull:
				log.Debug("Full route: %q", nroute)
				adp, err := GetAdapter(nroute.AdpID)
				if err != nil {
					return nil, fmt.Errorf("Invalid adpater id in route: %s", nroute)
				}
				adpStat, err := newAdapterStatus(adp, service, nroute, i+1)
				if err != nil {
					return nil, fmt.Errorf("Error creating route list - %s", err)
				}
				adpStatList[nroute] = adpStat

			default:
				// log.Debug("Using only adapters for areaID: %s", areaID)
				return nil, fmt.Errorf("Cannot create the Adapter List - invalid route: %s", nroute)
			}
		}
		log.Debug(adpStatList.String())
		return adpStatList, nil
	}

	serviceMap["Services.All"].buildAdapterList = func(r structs.NRouter) (adapterRouteList, error) {
		return area(r, "Services.All")
	}
	serviceMap["Services.Area"].buildAdapterList = func(r structs.NRouter) (adapterRouteList, error) {
		return area(r, "Services.Area")
	}
	serviceMap["Report.Create"].buildAdapterList = func(r structs.NRouter) (adapterRouteList, error) {
		return adapter(r, "Report.Create")
	}
	serviceMap["Report.SearchDID"].buildAdapterList = func(r structs.NRouter) (adapterRouteList, error) {
		return area(r, "Report.SearchDID")
	}
	serviceMap["Report.SearchLL"].buildAdapterList = func(r structs.NRouter) (adapterRouteList, error) {
		return area(r, "Report.SearchLL")
	}
	return nil
}

// =======================================================================================
//                                      RPC ROUTE METHODS
// =======================================================================================
type serviceMapMethods struct {
	newResponse      func() interface{}
	buildAdapterList func(structs.NRouter) (adapterRouteList, error)
}
