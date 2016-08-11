package router

import "github.com/open311-gateway/engine/structs"

// "github.com/davecgh/go-spew/spew"

// RoutesRID returns the route list for a ReportID.
func RoutesRID(rid structs.ReportID) (routes structs.NRoutes, err error) {
	return structs.NRoutes{rid.NRoute}, nil
}

// RoutesMID returns the route list for a ReportID.
func RoutesMID(mid structs.ServiceID) (routes structs.NRoutes, err error) {
	return structs.NRoutes{mid.GetRoute()}, nil
}

// RoutesArea returns the route list for an AreaID (i.e. City).
func RoutesArea(areaID string) (routes structs.NRoutes, err error) {
	return nil, nil
}

// RoutesAll returns all routes from all CONFIG'ured adapters.  This call does not
// use the Services cache for this list - it uses the config.json file.
func RoutesAll() (routes structs.NRoutes, err error) {
	for _, adp := range adapters.Adapters {
		routes = append(routes, structs.NRoute{AdpID: adp.ID, AreaID: "all", ProviderID: 0})
	}
	return routes, nil
}
