package structs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"Gateway311/engine/common"
)

// =======================================================================================
//                                      REQUEST
// =======================================================================================

// NID represents the full ID for any Normalized Request or Response.
type NID struct {
	rqstID int64
	rpcID  int64
}

// Set sets the NID.
func (r *NID) Set(rqstID, rpcID int64) {
	r.rqstID = rqstID
	r.rpcID = rpcID
}

// Get gets the NID.
func (r NID) Get() (int64, int64) {
	return r.rqstID, r.rpcID
}

// String returns the string representation NID.
func (r NID) String() string {
	return fmt.Sprintf("%d-%d", r.rqstID, r.rpcID)
}

// =======================================================================================
//                                      REQUEST
// =======================================================================================

// NRequestCommon represents properties common to all requests.
type NRequestCommon struct {
	ID    int64
	Route NRoute
	Rtype NRequestType
	NRouter
	NRequester
}

// GetID returns the Request ID
func (r NRequestCommon) GetID() int64 {
	return r.ID
}

// SetID sets the Request ID
func (r *NRequestCommon) SetID(id int64) {
	r.ID = id
}

// GetType returns the Request Type as a string.
func (r NRequestCommon) GetType() NRequestType {
	return r.Rtype
}

// GetTypeS returns the Request Type as a string.
func (r NRequestCommon) GetTypeS() string {
	fmt.Println("[NRequestCommon: GetTypeS()] start")
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
	GetID() int64
	SetID(int64)
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
	NRtTypInvalid NRouteType = iota
	NRtTypFull
	NRtTypArea
	NRtTypAllAreas
	NRtTypAllAdapters
)

// RouteType returns the validity and type of the NRoute.
func (r NRoute) RouteType() NRouteType {
	switch {
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
	ID    int64
	Route NRoute
	Rtype NResponseType
	NResponseer
}

// GetType returns the Response Type as a string.
func (r NResponseCommon) GetType() NResponseType {
	return r.Rtype
}

// GetTypeS returns the Response Type as a string.
func (r NResponseCommon) GetTypeS() string {
	fmt.Println("[NResponseCommon: GetTypeS()] start")
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

// NResponseer defines the behavior of a Response Package.
type NResponseer interface {
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
	ServiceID  `json:"id"`
	Name       string   `json:"name"`
	Categories []string `json:"catg"`
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

// =======================================================================================
//                                      CREATE
// =======================================================================================

// NCreateRequest is used to create a new Report.  It is the "native" format of the
// data, and is used by the Engine and all backend Adapters.
type NCreateRequest struct {
	NRequestCommon
	MID         ServiceID
	Type        string
	DeviceType  string
	DeviceModel string
	DeviceID    string
	Latitude    float64
	Longitude   float64
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
}

// GetRoutes returns the routing data.
func (r NCreateRequest) GetRoutes() NRoutes {
	return NewNRoutes().add(NRoute{r.MID.AdpID, r.MID.AreaID, r.MID.ProviderID})
}

// NCreateResponse is the response to creating or updating a report.
type NCreateResponse struct {
	NResponseCommon `json:"-"`
	Message         string `json:"Message" xml:"Message"`
	ID              string `json:"ReportId" xml:"ReportId"`
	AuthorID        string `json:"AuthorId" xml:"AuthorId"`
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

// NSearchRequestDID represents the Normal struct for a request to search for all reports
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
	return r.RouteList
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
	ID                int64
	DateCreated       string
	DateUpdated       string
	DeviceType        string
	DeviceModel       string
	DeviceID          string
	RequestType       string
	RequestTypeID     string
	ImageURL          string
	ImageURLXl        string
	ImageURLLg        string
	ImageURLMd        string
	ImageURLSm        string
	ImageURLXs        string
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

// =======================================================================================
//                                      MISC
// =======================================================================================

// SplitMID breaks down an MID, and returns all subfields.
func SplitMID(mid string) (string, string, int, int, error) {
	fail := func() (string, string, int, int, error) {
		return "", "", 0, 0, fmt.Errorf("Invalid MID: %s", mid)
	}
	parts := strings.Split(mid, "-")
	fmt.Printf("MID: %+v\n", parts)
	if len(parts) != 4 {
		fail()
	}
	pid, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		fail()
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		fail()
	}
	return parts[0], parts[1], int(pid), int(id), nil
}

// MidAdpID breaks down a MID, and returns the AdpID.
func MidAdpID(mid string) (string, error) {
	parts := strings.Split(mid, "-")
	fmt.Printf("MID: %+v\n", parts)
	if len(parts) != 4 {
		return "", fmt.Errorf("Invalid MID: %s", mid)
	}
	return parts[0], nil
}

// MidAreaID breaks down a MID, and returns the AreaID.
func MidAreaID(mid string) (string, error) {
	parts := strings.Split(mid, "-")
	fmt.Printf("MID: %+v\n", parts)
	if len(parts) != 4 {
		return "", fmt.Errorf("Invalid MID: %s", mid)
	}
	return parts[1], nil
}

// MidProviderID breaks down an MID, and returns the ProviderID.
func MidProviderID(mid string) (int, error) {
	fail := func() (int, error) {
		return 0, fmt.Errorf("Invalid MID: %s", mid)
	}
	parts := strings.Split(mid, "-")
	fmt.Printf("MID: %+v\n", parts)
	if len(parts) != 4 {
		return 0, fmt.Errorf("Invalid MID: %s", mid)
	}
	pid, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		fail()
	}
	return int(pid), nil
}

// MidID breaks down an MID, and returns the Service ID.
func MidID(mid string) (int, error) {
	fail := func() (int, error) {
		return 0, fmt.Errorf("Invalid MID: %s", mid)
	}
	parts := strings.Split(mid, "-")
	fmt.Printf("MID: %+v\n", parts)
	if len(parts) != 4 {
		return 0, fmt.Errorf("Invalid MID: %s", mid)
	}
	id, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		fail()
	}
	return int(id), nil
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
	r := fmt.Sprintf("  %-20s  %-40s  %v", s.MID(), s.Name, s.Categories)
	return r
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
	return fmt.Sprintf("%s-%s-%d-%d", r.AdpID, r.AreaID, r.ProviderID, r.ID)
}

func (r NRouteType) String() string {
	switch r {
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

func (r NRoute) SString() string {
	return fmt.Sprintf("%s-%s-%d", r.AdpID, r.AreaID, r.ProviderID)
}

// String displays the type.
func (r NRoute) String() string {
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
	ls.AddF("Device - type %s  model: %s  ID: %s\n", r.DeviceType, r.DeviceModel, r.DeviceID)
	ls.AddF("Request - %s:  %s\n", r.MID.MID(), r.Type)
	ls.AddF("Location - lat: %v lon: %v\n", r.Latitude, r.Longitude)
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
	ls.AddF("ID: %v  AuthorID: %v\n", r.ID, r.AuthorID)
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

// Displays the the NSearchRequestDID custom type.
func (r NSearchResponseReport) String() string {
	ls := new(common.LogString)
	ls.AddF("NSearchResponseReport %d\n", r.ID)
	ls.AddF("DateCreated \"%v\"\n", r.DateCreated)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", r.DeviceType, r.DeviceModel, r.DeviceID)
	ls.AddF("Request - type: %q  id: %q\n", r.RequestType, r.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v  directionality: %q\n", r.Latitude, r.Longitude, r.Directionality)
	ls.AddF("          %s, %s   %s\n", r.City, r.State, r.ZipCode)
	ls.AddF("Votes: %d\n", r.Votes)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Images - std: %s\n", r.ImageURL)
	if len(r.ImageURLXs) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLXs)
	}
	if len(r.ImageURLSm) > 0 {
		ls.AddF("          SM: %s\n", r.ImageURLSm)
	}
	if len(r.ImageURLMd) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLMd)
	}
	if len(r.ImageURLLg) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLLg)
	}
	if len(r.ImageURLXl) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLXl)
	}
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
