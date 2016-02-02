package structs

import (
	"fmt"
	"strconv"
	"strings"

	"Gateway311/engine/common"
)

// =======================================================================================
//                                      API
// =======================================================================================

// API contains the information required by the Backend to process a transation - e.g. the
// API authorization key, API call, etc.
type API struct {
	APIAuthKey        string `json:"ApiAuthKey" xml:"ApiAuthKey"`
	APIRequestType    string `json:"ApiRequestType" xml:"ApiRequestType"`
	APIRequestVersion string `json:"ApiRequestVersion" xml:"ApiRequestVersion"`
}

// =======================================================================================
//                                      ROUTE
// =======================================================================================

// NRouter is the interface to retrieve the routing data (AdpID, AreaID) from
// any N*Request.
type NRouter interface {
	Route() NRoute
}

// NRoute represents the data needed to route requests to Adapters.
type NRoute struct {
	AdpID      string
	AreaID     string
	ProviderID int
}

// =======================================================================================
//                                      SERVICES
// =======================================================================================

// NServiceRequest is used to get list of services available to the user.
type NServiceRequest struct {
	NRouter
	Area string
}

// Route returns the routing data.
func (n NServiceRequest) Route() NRoute {
	return NRoute{"", n.Area, 0}
}

// NServicesResponse is the returned struct for a Services request.
type NServicesResponse struct {
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

// =======================================================================================
//                                      CREATE
// =======================================================================================

// NCreateRequest is used to create a new Report.  It is the "native" format of the
// data, and is used by the Engine and all backend Adapters.
type NCreateRequest struct {
	NRouter
	API
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

// Route returns the routing data.
func (ncr NCreateRequest) Route() NRoute {
	return NRoute{ncr.MID.AdpID, ncr.MID.AreaID, ncr.MID.ProviderID}
}

// NCreateResponse is the response to creating or updating a report.
type NCreateResponse struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

// =======================================================================================
//                                      SEARCH
// =======================================================================================

// NSearchRequestLL represents the Normal struct for a location based search request.
type NSearchRequestLL struct {
	NRouter
	Routing NRoute
	API
	SearchType NSearchType
	Latitude   float64
	Longitude  float64
	AreaID     string
	Radius     int // in meters
	MaxResults int
	Response   *NSearchResponse
}

// Route returns the routing data.
func (r NSearchRequestLL) Route() NRoute {
	return r.Routing
}

// NSearchRequestDID represents the Normal struct for a request to search for all reports
// authored by the specified Device ID.
type NSearchRequestDID struct {
	NRouter
	Routing NRoute
	API
	SearchType NSearchType
	DeviceType string
	DeviceID   string
	MaxResults int
	Response   *NSearchResponse
}

// Route returns the routing data.
func (r NSearchRequestDID) Route() NRoute {
	return r.Routing
}

//go:generate stringer -type=NSearchType

// NSearchType enumerates the valid search types.
type NSearchType int

// NSearchType definitions.
const (
	NSTUnknown NSearchType = iota
	NSTLocation
	NSTDeviceID
)

// NSearchResponse contains the search results.
type NSearchResponse struct {
	Message      string
	ReportCount  int
	ResponseTime string
	Reports      []NSearchResponseReport
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

// Displays the NServiceRequest custom type.
func (n NServiceRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("NServiceRequest\n")
	ls.AddF("Location - area: %v\n", n.Area)
	return ls.Box(80)
}

// Displays the NServicesResponse custom type.
func (c NServicesResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NServicesResponse\n")
	ls.AddF("Message: %q%s", c.Message, c.Services)
	return ls.Box(90)
}

// Displays the NService custom type.
func (s NService) String() string {
	// r := fmt.Sprintf("  %s-%s-%d-%d  %-40s  %v", s.AdpID, s.AreaID, s.ProviderID, s.ID, s.Name, s.Categories)
	r := fmt.Sprintf("  %-20s  %-40s  %v", s.MID(), s.Name, s.Categories)
	return r
}

// Displays the NServices custom type.
func (c NServices) String() string {
	ls := new(common.LogString)
	ls.AddS("NServices\n")
	for _, s := range c {
		ls.AddF("%s\n", s)
	}
	return ls.Box(80)
}

// MID creates the Master ID string for the Service.
func (s ServiceID) MID() string {
	return fmt.Sprintf("%s-%s-%d-%d", s.AdpID, s.AreaID, s.ProviderID, s.ID)
}

// Displays the NCreateRequest custom type.
func (ncr NCreateRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("NCreateRequest\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", ncr.DeviceType, ncr.DeviceModel, ncr.DeviceID)
	ls.AddF("Request - %s:  %s\n", ncr.MID.MID(), ncr.Type)
	ls.AddF("Location - lat: %v lon: %v\n", ncr.Latitude, ncr.Longitude)
	ls.AddF("          %s, %s   %s\n", ncr.Area, ncr.State, ncr.Zip)
	ls.AddF("Description: %q\n", ncr.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", ncr.IsAnonymous, ncr.FirstName, ncr.LastName, ncr.Email, ncr.Phone)
	return ls.Box(80)
}

// Displays the NCreateResponse custom type.
func (c NCreateResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NCreateResponse\n")
	ls.AddF("Message: %s\n", c.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", c.ID, c.AuthorID)
	return ls.Box(80)
}

// Displays the NSearchRequestLL custom type.
func (r NSearchRequestLL) String() string {
	ls := new(common.LogString)
	ls.AddS("NSearchRequestLL\n")
	ls.AddF("SearchType: %s\n", r.SearchType)
	ls.AddF("Lat: %v  Lng: %v   Radius: %v AreaID: %q\n", r.Latitude, r.Longitude, r.Radius, r.AreaID)
	ls.AddF("MaxResults: %v  Routing: %s\n", r.MaxResults, r.Routing)
	return ls.Box(80)
}

// Displays the NSearchRequestDID custom type.
func (r NSearchRequestDID) String() string {
	ls := new(common.LogString)
	ls.AddS("NSearchRequestDID\n")
	ls.AddF("SearchType: %s\n", r.SearchType)
	ls.AddF("Device type: %v  ID: %v\n", r.DeviceType, r.DeviceID)
	ls.AddF("MaxResults: %v  Routing: %s\n", r.MaxResults, r.Routing)
	return ls.Box(80)
}

// Displays the NSearchResponse custom type.
func (r NSearchResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NSearchResponse\n")
	ls.AddF("Count: %v RspTime: %v Message: %v\n", r.ReportCount, r.ResponseTime, r.Message)
	for _, x := range r.Reports {
		ls.AddS(x.String())
	}
	return ls.Box(90)
}

// Displays the the NSearchRequestDID custom type.
func (r NSearchResponseReport) String() string {
	ls := new(common.LogString)
	ls.AddF("Report %d\n", r.ID)
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
