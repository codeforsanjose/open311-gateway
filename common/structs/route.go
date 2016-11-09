package structs

import (
	"fmt"

	"github.com/codeforsanjose/open311-gateway/common"
	"github.com/fatih/color"
)

// =======================================================================================
//                                      ROUTE
// =======================================================================================

// NRouter is the interface to retrieve the routing data (AdpID, AreaID) from
// any N*Request.
type NRouter interface {
	GetRoutes() NRoutes
}

// NRoutes represents a list of Routes for a request.
type NRoutes []NRoute

// NewNRoutes returns a new instance of NRoutes.
func NewNRoutes() NRoutes {
	return make([]NRoute, 0)
}

// NRoutes represents a list of Routes for a request.
func (r NRoutes) add(nr NRoute) NRoutes {
	r = append(r, nr)
	return r
}

// NRoute represents the data needed to route requests to Adapters.
type NRoute struct {
	AdpID      string
	AreaID     string
	ProviderID int
}

// NRouteType enumerates the valid route types.
type NRouteType int

// NRT* are constants enumerating the valid request types.
const (
	NRtTypEmpty NRouteType = iota
	NRtTypInvalid
	NRtTypFull
	NRtTypArea
	NRtTypAllAreas
	NRtTypAllAdapters
)

// RouteType returns the validity and type of the NRoute.
func (r NRoute) RouteType() NRouteType {
	switch {
	case r.AdpID == "" && r.AreaID == "" && r.ProviderID == 0:
		return NRtTypEmpty
	case r.AdpID > "" && r.AreaID > "" && r.ProviderID > 0:
		return NRtTypFull
	case r.AdpID > "" && r.AreaID == "all" && r.ProviderID == 0:
		return NRtTypAllAreas
	case r.AdpID == "" && r.AreaID > "" && r.ProviderID == 0:
		return NRtTypArea
	case r.AdpID == "all" && r.AreaID == "" && r.ProviderID == 0:
		return NRtTypAllAdapters
	default:
		return NRtTypInvalid
	}
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// String returns a short representation of Routes.
func (r NRoutes) String() string {
	ls := new(common.FmtBoxer)
	ls.AddS("NRoutes\n")
	for _, r := range r {
		ls.AddF("%s\n", r)
	}
	return ls.Box(40)
}

// String returns a short representation of a Route.
func (r NRoute) String() string {
	return fmt.Sprintf("%s-%s-%d", r.AdpID, r.AreaID, r.ProviderID)
}

func (r NRouteType) String() string {
	switch r {
	case NRtTypEmpty:
		return color.BlueString("empty")
	case NRtTypFull:
		return color.GreenString("Full")
	case NRtTypAllAdapters:
		return color.YellowString("AllAdps")
	case NRtTypAllAreas:
		return color.YellowString("AllAreas")
	case NRtTypArea:
		return color.YellowString("Area")
	default:
		return color.RedString("Invalid")
	}
}

// SString displays a Route.
func (r NRoute) SString() string {
	// fmtEmpty := color.New(color.BgRed, color.FgWhite, color.Bold).SprintFunc()
	// empty := fmtEmpty("\u2205")
	// empty := color.RedString("\u2205")
	if r.AdpID == "" && r.AreaID == "" && r.ProviderID == 0 {
		return fmt.Sprintf("[%s] %s", r.RouteType(), color.RedString("\u2205\u2205\u2205"))
	}
	AdpID, AreaID, ProviderID := r.AdpID, r.AreaID, r.ProviderID
	if r.AdpID == "" {
		AdpID = color.RedString("\u2205")
	}
	if r.AreaID == "" {
		// r.AreaID = "\u00F8"
		AreaID = color.RedString("\u2205")
	}
	return fmt.Sprintf("[%s] %s-%s-%d", r.RouteType(), AdpID, AreaID, ProviderID)
}
