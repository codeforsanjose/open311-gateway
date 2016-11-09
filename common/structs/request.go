package structs

import (
	"fmt"

	"github.com/codeforsanjose/open311-gateway/common"
)

// =======================================================================================
//                                      REQUEST
// =======================================================================================

// NID represents the full ID for any Normalized Request or Response.
type NID struct {
	RqstID int64
	RPCID  int64
}

// SetNID sets the NID.
func (r *NID) SetNID(rqstID, rpcID int64) {
	if rqstID > 0 {
		r.RqstID = rqstID
	}
	if rpcID > 0 {
		r.RPCID = rpcID
	}
}

// GetNID gets the NID.
func (r NID) GetNID() (int64, int64) {
	return r.RqstID, r.RPCID
}

// String returns the string representation NID.
func (r NID) String() string {
	return fmt.Sprintf("%d-%d", r.RqstID, r.RPCID)
}

// =======================================================================================
//                                      REQUEST
// =======================================================================================

// NRequestCommon represents properties common to all requests.
type NRequestCommon struct {
	ID    NID
	Route NRoute
	Rtype NRequestType
	NRouter
	NRequester
}

// GetID returns the Request ID
func (r NRequestCommon) GetID() (int64, int64) {
	return r.ID.GetNID()
}

// GetIDS returns the Request ID as a string
func (r NRequestCommon) GetIDS() string {
	x, y := r.GetID()
	return fmt.Sprintf("%v-%v", x, y)
}

// SetID sets the Request ID
func (r *NRequestCommon) SetID(rqstID, rpcID int64) {
	r.ID.SetNID(rqstID, rpcID)
}

// GetType returns the Request Type as a string.
func (r NRequestCommon) GetType() NRequestType {
	return r.Rtype
}

// GetTypeS returns the Request Type as a string.
func (r NRequestCommon) GetTypeS() string {
	return r.Rtype.String()
}

// GetRoute returns NRequestCommon.Route
func (r NRequestCommon) GetRoute() NRoute {
	return r.Route
}

// SetRoute sets the route in NRequestCommon.
func (r *NRequestCommon) SetRoute(route NRoute) {
	r.Route = route
}

// -----------------------------------NRequester --------------------------------------

// NRequester defines the behavior of a Request Package.
type NRequester interface {
	GetID() (int64, int64)
	GetIDS() string
	SetID(int64, int64)
	GetRoute() NRoute
	SetRoute(route NRoute)
	RouteType() NRouteType
	GetType() NRequestType
	GetTypeS() string
}

// -----------------------------------NRequestType --------------------------------------

//go:generate stringer -type=NRequestType

// NRequestType enumerates the valid request types.
type NRequestType int

// NRT* are constants enumerating the valid request types.
const (
	NRTUnknown NRequestType = iota
	NRTServicesAll
	NRTServicesArea
	NRTCreate
	NRTSearchLL
	NRTSearchDID
	NRTSearchRID
)

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// Displays the NRequestCommon custom type.
func (r NRequestCommon) String() string {
	ls := new(common.FmtBoxer)
	ls.AddF("Type: %s\n", r.Rtype.String())
	ls.AddF("ID: %v\n", r.ID)
	ls.AddF("Route: %s\n", r.Route.String())
	return ls.Box(40)
}
