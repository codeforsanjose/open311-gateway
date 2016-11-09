package structs

import (
	"fmt"

	"github.com/codeforsanjose/open311-gateway/common"
)

// =======================================================================================
//                                      RESPONSE
// =======================================================================================

// NResponseCommon represents properties common to all requests.
type NResponseCommon struct {
	ID    NID
	Route NRoute
	Rtype NResponseType
	NResponser
}

// GetID returns the Request ID
func (r NResponseCommon) GetID() (int64, int64) {
	return r.ID.GetNID()
}

// GetIDS returns the Request ID as a string
func (r NResponseCommon) GetIDS() string {
	x, y := r.GetID()
	return fmt.Sprintf("%v-%v", x, y)
}

// SetID sets the Request ID
func (r *NResponseCommon) SetID(rqstID, rpcID int64) {
	r.ID.SetNID(rqstID, rpcID)
}

// SetIDF sets the Request ID using the specified function.
func (r *NResponseCommon) SetIDF(f func() (int64, int64)) {
	x, y := f()
	r.ID.SetNID(x, y)
}

// GetType returns the Response Type as a string.
func (r NResponseCommon) GetType() NResponseType {
	return r.Rtype
}

// GetTypeS returns the Response Type as a string.
func (r NResponseCommon) GetTypeS() string {
	return r.Rtype.String()
}

// GetRoute returns NResponseCommon.Route
func (r NResponseCommon) GetRoute() NRoute {
	return r.Route
}

// SetRoute sets the route in NResponseCommon.
func (r *NResponseCommon) SetRoute(route NRoute) {
	r.Route = route
}

// -----------------------------------NResponseer --------------------------------------

// NResponser defines the behavior of a Response Package.
type NResponser interface {
	GetID() (int64, int64)
	GetIDS() string
	SetID(int64, int64)
	GetType() NResponseType
	GetTypeS() string
	GetRoute() NRoute
	SetRoute(route NRoute)
}

// -----------------------------------NResponseType --------------------------------------

//go:generate stringer -type=NResponseType

// NResponseType enumerates the valid request types.
type NResponseType int

// NRspT* are constants enumerating the valid request types.
const (
	NRspTUnknown NResponseType = iota
	NRspTServices
	NRspTServicesArea
	NRspTCreate
	NRspTSearchLL
	NRspTSearchDID
	NRspTSearchRID
)

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// Displays the NResponseCommon custom type.
func (r NResponseCommon) String() string {
	ls := new(common.FmtBoxer)
	ls.AddF("Type: %s\n", r.Rtype.String())
	ls.AddF("Route: %s\n", r.Route.String())
	return ls.Box(40)
}
