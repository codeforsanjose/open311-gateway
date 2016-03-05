package router

import (
	"Gateway311/engine/structs"
)

// "github.com/davecgh/go-spew/spew"

// RoutesRID returns the route list for a ReportID.
func RoutesRID(rid structs.ReportID) (routes structs.NRoutes, err error) {
	return structs.NRoutes{rid.NRoute}, nil
}

// RoutesArea returns the route list for an AreaID (i.e. City).
func RoutesArea(areaID string) (routes structs.NRoutes, err error) {
	return nil, nil
}

// RoutesAll returns all routes.
func RoutesAll() (routes structs.NRoutes, err error) {
	return GetAllRoutes(), nil
}
