package structs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"github.com/open311-gateway/engine/common"
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
//                                      SERVICES
// =======================================================================================

// NServiceRequest is used to get list of services available to the user.
type NServiceRequest struct {
	NRequestCommon
	Area string
}

// GetRoutes returns the routing data.
func (r NServiceRequest) GetRoutes() NRoutes {
	if r.Area == "all" {
		return NewNRoutes().add(NRoute{"all", "", 0})
	}
	return NewNRoutes().add(NRoute{"", r.Area, 0})
}

// NServicesResponse is the returned struct for a Services request.
type NServicesResponse struct {
	NResponseCommon
	AdpID    string
	Message  string
	Services NServices
}

// ------------------------------- Services -------------------------------

// NServices contains a list of Services.
type NServices []NService

// ------------------------------- Service -------------------------------

// NService represents a Service.  The ID is a combination of the BackEnd Type (AdpID),
// the AreaID (i.e. the Area id), ProviderID (in case the provider has multiple interfaces),
// and the Service ID.
type NService struct {
	ServiceID
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Metadata      bool     `json:"metadata"`
	ResponseType  string   `json:"responseType"`
	ServiceNotice string   `json:"service_notice"`
	Keywords      []string `json:"keywords"`
	Group         string   `json:"group"`
}

// ------------------------------- ServiceID -------------------------------

// ServiceID provides the JSON marshalling conversion between the JSON "ID" and
// the Backend Interface Type, AreaID (Area id), ProviderID, and Service ID.
type ServiceID struct {
	AdpID      string
	AreaID     string
	ProviderID int
	ID         int
}

// GetRoute returns the NRoute for the ServiceID.
func (r ServiceID) GetRoute() NRoute {
	return NRoute{
		AdpID:      r.AdpID,
		AreaID:     r.AreaID,
		ProviderID: r.ProviderID,
	}
}

// ------------------------------- ReportID -------------------------------

// ReportID adds routing information to a ReportID returned by a call to a
// Service Provider.
type ReportID struct {
	NRoute
	ID string
}

// NewRID creates a ReportID by concatenating a Route (NRoute string) with a message
// response ID.
func NewRID(route NRoute, id string) ReportID {
	return ReportID{
		NRoute: route,
		ID:     id,
	}

}

// =======================================================================================
//                                      CREATE
// =======================================================================================

// NCreateRequest is used to create a new Report.  It is the "native" format of the
// data, and is used by the Engine and all backend Adapters.
type NCreateRequest struct {
	NRequestCommon
	MID         ServiceID
	ServiceName string
	DeviceType  string
	DeviceModel string
	DeviceID    string
	Latitude    float64
	Longitude   float64
	FullAddress string
	Address     string
	Area        string
	State       string
	Zip         string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	IsAnonymous bool
	Description string
	MediaURL    string
}

// GetRoutes returns the routing data.
func (r NCreateRequest) GetRoutes() NRoutes {
	return NewNRoutes().add(NRoute{r.MID.AdpID, r.MID.AreaID, r.MID.ProviderID})
}

// NCreateResponse is the response to creating or updating a report.
type NCreateResponse struct {
	NResponseCommon `json:"-"`
	Message         string   `json:"Message" xml:"Message"`
	RID             ReportID `json:"ReportId" xml:"ReportId"`
	// ID              string `json:"ReportId" xml:"ReportId"`
	AccountID string `json:"AuthorId" xml:"AuthorId"`
}

// =======================================================================================
//                                      SEARCH
// =======================================================================================

// NSearchRequestLL represents the Normal struct for a location based search request.
type NSearchRequestLL struct {
	NRequestCommon
	Latitude   float64
	Longitude  float64
	Radius     int // in meters
	AreaID     string
	MaxResults int
}

// GetRoutes returns the routing data.
func (r NSearchRequestLL) GetRoutes() NRoutes {
	return NewNRoutes().add(NRoute{"", r.AreaID, 0})
}

// NSearchRequestDID represents the Normal struct for a request to search for all Reports
// authored by the specified Device ID.
type NSearchRequestDID struct {
	NRequestCommon
	DeviceType string
	DeviceID   string
	MaxResults int
	RouteList  NRoutes
	AreaID     string
}

// GetRoutes returns the routing data.
func (r NSearchRequestDID) GetRoutes() NRoutes {
	if len(r.RouteList) > 0 {
		return r.RouteList
	}
	return NewNRoutes().add(NRoute{"", r.AreaID, 0})
}

// NSearchRequestRID represents the Normal struct for a request to find a single,
// specific Report.
type NSearchRequestRID struct {
	NRequestCommon
	RID       ReportID
	RouteList NRoutes
	AreaID    string
}

// GetRoutes returns the routing data.
func (r NSearchRequestRID) GetRoutes() NRoutes {
	if len(r.RouteList) > 0 {
		return r.RouteList
	}
	return NewNRoutes().add(NRoute{r.RID.AdpID, r.RID.AreaID, r.RID.ProviderID})
}

// NSearchResponse contains the search results.
type NSearchResponse struct {
	NResponseCommon `json:"-"`
	Message         string
	ReportCount     int
	ResponseTime    string
	Reports         []NSearchResponseReport
}

// NSearchResponseReport represents a report.
type NSearchResponseReport struct {
	RID               ReportID `xml:"ID" json:"ID"`
	DateCreated       string
	DateUpdated       string
	DeviceType        string
	DeviceModel       string
	DeviceID          string
	RequestType       string
	RequestTypeID     string
	MediaURL          string
	City              string
	State             string
	ZipCode           string
	Latitude          string
	Longitude         string
	Directionality    string
	Description       string
	AuthorNameFirst   string
	AuthorNameLast    string
	AuthorEmail       string
	AuthorTelephone   string
	AuthorIsAnonymous string
	URLDetail         string
	URLShortened      string
	Votes             string
	StatusType        string
	TicketSLA         string
}

// FullAddress returns the standard formatted full address.
func (r *NSearchResponseReport) FullAddress() string {
	if len(r.City+r.State+r.ZipCode) == 0 {
		return ""
	}
	return fmt.Sprintf("%s, %s %s", r.City, r.State, r.ZipCode)
}

// =======================================================================================
//                                      MID
// =======================================================================================

// SplitRMID breaks down an MID or RID, and returns all subfields.
func SplitRMID(mid string) (string, string, int, int, error) {
	fail := func() (string, string, int, int, error) {
		return "", "", 0, 0, fmt.Errorf("Invalid RMID: %q", mid)
	}
	parts := strings.Split(mid, "-")
	if len(parts) != 4 {
		fail()
	}
	pid, err := strconv.Atoi(parts[2])
	if err != nil {
		fail()
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		fail()
	}
	return parts[0], parts[1], pid, id, nil
}

const emInvalidMid = "Invalid MID: %q"

// SplitMID breaks down an MID, and returns all subfields.
func SplitMID(mid string) (string, string, int, int, error) {
	adpID, areaID, provID, id, err := SplitRMID(mid)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf(emInvalidMid, mid)
	}
	return adpID, areaID, provID, id, nil
}

// MidAdpID breaks down a MID, and returns the AdpID.
func MidAdpID(mid string) (string, error) {
	adpID, _, _, _, err := SplitRMID(mid)
	if err != nil {
		return "", fmt.Errorf(emInvalidMid, mid)
	}
	return adpID, nil
}

// MidAreaID breaks down a MID, and returns the AreaID.
func MidAreaID(mid string) (string, error) {
	_, areaID, _, _, err := SplitRMID(mid)
	if err != nil {
		return "", fmt.Errorf(emInvalidMid, mid)
	}
	return areaID, nil
}

// MidProviderID breaks down an MID, and returns the ProviderID.
func MidProviderID(mid string) (int, error) {
	_, _, provID, _, err := SplitRMID(mid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidMid, mid)
	}
	return provID, nil
}

// MidID breaks down an MID, and returns the Service ID.
func MidID(mid string) (int, error) {
	_, _, _, id, err := SplitRMID(mid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidMid, mid)
	}
	return id, nil
}

// =======================================================================================
//                                      RID
// =======================================================================================

const emInvalidRid = "Invalid RID: %q"

// RIDFromString converts a reportID string to a new ReportID struct.
func RIDFromString(rids string) (ReportID, NRoute, error) {
	if rids == "" {
		return ReportID{}, NRoute{}, fmt.Errorf("empty RID: %q", rids)
	}
	adpID, areaID, providerID, reportID, err := SplitRID(rids)
	if err != nil {
		return ReportID{}, NRoute{}, fmt.Errorf(emInvalidRid, rids)
	}
	nr := NRoute{
		AdpID:      adpID,
		AreaID:     areaID,
		ProviderID: providerID,
	}

	return ReportID{
		NRoute: nr,
		ID:     fmt.Sprintf("%v", reportID),
	}, nr, nil
}

// NRouteFromString converts a reportID string to a new NRoute struct.
func NRouteFromString(rids string) (NRoute, error) {
	adpID, areaID, providerID, _, err := SplitRID(rids)
	if err != nil {
		return NRoute{}, fmt.Errorf("invalid RID: %q", rids)
	}
	return NRoute{
		AdpID:      adpID,
		AreaID:     areaID,
		ProviderID: providerID,
	}, nil

}

// SplitRID breaks down an RID, and returns all subfields.
func SplitRID(rid string) (string, string, int, int, error) {
	adpID, areaID, provID, id, err := SplitRMID(rid)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf(emInvalidRid, rid)
	}
	return adpID, areaID, provID, id, nil
}

// RidAdpID breaks down a RID, and returns the AdpID.
func RidAdpID(rid string) (string, error) {
	adpID, _, _, _, err := SplitRMID(rid)
	if err != nil {
		return "", fmt.Errorf(emInvalidRid, rid)
	}
	return adpID, nil
}

// RidAreaID breaks down a RID, and returns the AreaID.
func RidAreaID(rid string) (string, error) {
	_, areaID, _, _, err := SplitRMID(rid)
	if err != nil {
		return "", fmt.Errorf(emInvalidRid, rid)
	}
	return areaID, nil
}

// RidProviderID breaks down an RID, and returns the ProviderID.
func RidProviderID(rid string) (int, error) {
	_, _, provID, _, err := SplitRMID(rid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidRid, rid)
	}
	return provID, nil
}

// RidID breaks down an RID, and returns the Report ID.
func RidID(rid string) (int, error) {
	_, _, _, id, err := SplitRMID(rid)
	if err != nil {
		return 0, fmt.Errorf(emInvalidRid, rid)
	}
	return id, nil
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// Displays the NRequestCommon custom type.
func (r NRequestCommon) String() string {
	ls := new(common.LogString)
	ls.AddF("Type: %s\n", r.Rtype.String())
	ls.AddF("ID: %v\n", r.ID)
	ls.AddF("Route: %s\n", r.Route.String())
	return ls.Box(40)
}

// Displays the NServiceRequest custom type.
func (r NServiceRequest) String() string {
	ls := new(common.LogString)
	ls.AddF("NServiceRequest\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("Location - area: %v\n", r.Area)
	return ls.Box(80)
}

// Displays the NResponseCommon custom type.
func (r NResponseCommon) String() string {
	ls := new(common.LogString)
	ls.AddF("Type: %s\n", r.Rtype.String())
	ls.AddF("Route: %s\n", r.Route.String())
	return ls.Box(40)
}

// Displays the NServicesResponse custom type.
func (r NServicesResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NServicesResponse\n")
	ls.AddS(r.NResponseCommon.String())
	ls.AddF("Message: %q%s", r.Message, r.Services)
	return ls.Box(90)
}

// Displays the NService custom type.
func (s NService) String() string {
	r := fmt.Sprintf("  %-20s  %-40s  %-14s %v", s.MID(), s.Name, s.Group, s.Keywords)
	return r
}

// SString should be used to display one NService struct.
func (s NService) SString() string {
	ls := new(common.LogString)
	ls.AddS("NService\n")
	ls.AddF("MID: %s    Name: %s\n", s.MID(), s.Name)
	ls.AddF("Group: %s   Keywords: %v\n", s.Group, s.Keywords)
	ls.AddF("ServiceNotice: %v\n", s.ServiceNotice)
	return ls.Box(70)
}

// Displays the NServices custom type.
func (r NServices) String() string {
	ls := new(common.LogString)
	ls.AddS("NServices\n")
	for _, s := range r {
		ls.AddF("%s\n", s)
	}
	return ls.Box(80)
}

// MID creates the Master ID string for the Service.
func (r ServiceID) MID() string {
	if r.AdpID == "" && r.AreaID == "" && r.ProviderID == 0 && r.ID == 0 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%d-%d", r.AdpID, r.AreaID, r.ProviderID, r.ID)
}

// RID creates the Master ID string for the Service.
func (r ReportID) RID() string {
	if r.AdpID == "" && r.AreaID == "" && r.ProviderID == 0 && r.ID == "" {
		return ""
	}
	return fmt.Sprintf("%s-%s-%d-%s", r.AdpID, r.AreaID, r.ProviderID, r.ID)
}

// Display the string represenation of a ReportID.
func (r ReportID) String() string {
	return fmt.Sprintf("%s-%s", r.NRoute, r.ID)
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

// String returns a short representation of a Route.
func (r NRoute) String() string {
	return fmt.Sprintf("%s-%s-%d", r.AdpID, r.AreaID, r.ProviderID)
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

// Displays the NCreateRequest custom type.
func (r NCreateRequest) String() string {
	ls := new(common.LogString)
	ls.AddF("NCreateRequest\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("Request: %s\n", r.ServiceName)
	ls.AddF("Device - ID: %s  type: %s  model: %s\n", r.DeviceID, r.DeviceType, r.DeviceModel)
	ls.AddF("Request - %s\n", r.MID.MID())
	ls.AddF("Location - lat: %v lon: %v\n", r.Latitude, r.Longitude)
	ls.AddF("          %s\n", r.Address)
	ls.AddF("          %s, %s   %s\n", r.Area, r.State, r.Zip)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", r.IsAnonymous, r.FirstName, r.LastName, r.Email, r.Phone)
	return ls.Box(80)
}

// Displays the NCreateResponse custom type.
func (r NCreateResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NCreateResponse\n")
	ls.AddS(r.NResponseCommon.String())
	ls.AddF("Message: %s\n", r.Message)
	ls.AddF("ID: %v  AccountID: %v\n", r.ID, r.AccountID)
	return ls.Box(80)
}

// Displays the NSearchRequestLL custom type.
func (r NSearchRequestLL) String() string {
	ls := new(common.LogString)
	ls.AddF("NSearchRequestLL\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("Lat: %v  Lng: %v   Radius: %v AreaID: %q\n", r.Latitude, r.Longitude, r.Radius, r.AreaID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	return ls.Box(80)
}

// Displays the NSearchRequestDID custom type.
func (r NSearchRequestDID) String() string {
	ls := new(common.LogString)
	ls.AddF("NSearchRequestDID\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("Device type: %v  ID: %v\n", r.DeviceType, r.DeviceID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	return ls.Box(80)
}

// Displays the NSearchRequestRID custom type.
func (r NSearchRequestRID) String() string {
	ls := new(common.LogString)
	ls.AddF("NSearchRequestRID\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("ReportID: %v    AreaID: %q\n", r.RID.String(), r.AreaID)
	ls.AddF("RouteList: %v\n", r.RouteList.String())
	return ls.Box(80)
}

// Displays the NSearchResponse custom type.
func (r NSearchResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NSearchResponse\n")
	ls.AddS(r.NResponseCommon.String())
	ls.AddF("Count: %v RspTime: %v Message: %v\n", r.ReportCount, r.ResponseTime, r.Message)
	for _, x := range r.Reports {
		ls.AddS(x.String())
	}
	return ls.Box(90)
}

// Displays the the NSearchResponseReport custom type.
func (r NSearchResponseReport) String() string {
	ls := new(common.LogString)
	ls.AddF("NSearchResponseReport %s\n", r.RID.RID())
	ls.AddF("DateCreated \"%v\"\n", r.DateCreated)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", r.DeviceType, r.DeviceModel, r.DeviceID)
	ls.AddF("Request - type: %q  id: %q\n", r.RequestType, r.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v  directionality: %q\n", r.Latitude, r.Longitude, r.Directionality)
	ls.AddF("          %s, %s   %s\n", r.City, r.State, r.ZipCode)
	ls.AddF("Votes: %v\n", r.Votes)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Images - std: %s\n", r.MediaURL)
	ls.AddF("Author(anon: %v) %s %s  Email: %s  Tel: %s\n", r.AuthorIsAnonymous, r.AuthorNameFirst, r.AuthorNameLast, r.AuthorEmail, r.AuthorTelephone)
	ls.AddF("SLA: %s\n", r.TicketSLA)
	return ls.Box(80)
}

func (r NRoutes) String() string {
	ls := new(common.LogString)
	ls.AddS("NRoutes\n")
	for _, r := range r {
		ls.AddF("%s\n", r)
	}
	return ls.Box(40)
}
