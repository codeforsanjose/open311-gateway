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
	Route() []NRoute
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
func (n NServiceRequest) Route() []NRoute {
	return []NRoute{{"", n.Area, 0}}
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
func (ncr NCreateRequest) Route() []NRoute {
	return []NRoute{{ncr.MID.AdpID, ncr.MID.AreaID, ncr.MID.ProviderID}}
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

// SearchReqBase is used to create a report.
type SearchReqBase struct {
	NRouter
	API
	DeviceType  string  `json:"deviceType" xml:"deviceType"`
	DeviceID    string  `json:"deviceId" xml:"deviceId"`
	Latitude    string  `json:"LatitudeV" xml:"LatitudeV"`
	LatitudeV   float64 //
	Longitude   string  `json:"LongitudeV" xml:"LongitudeV"`
	LongitudeV  float64 //
	Radius      string  `json:"RadiusV" xml:"RadiusV"`
	RadiusV     int     // in meters
	Address     string  `json:"address" xml:"address"`
	City        string  `json:"city" xml:"city"`
	AreaID      string
	State       string `json:"state" xml:"state"`
	Zip         string `json:"zip" xml:"zip"`
	MaxResults  string `json:"MaxResultsV" xml:"MaxResultsV"`
	MaxResultsV int    //
	SearchType  string //
	Response    *SearchResp
}

// Route returns the routing data.
func (r SearchReqBase) Route() []NRoute {
	return []NRoute{{"", r.AreaID, 0}}
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

// NSearchReqLL represents the Normal struct for a location based search request.
type NSearchReqLL struct {
	NRouter
	API
	SearchType NSearchType
	Latitude   float64
	Longitude  float64
	AreaID     string
	Radius     int // in meters
	MaxResults int
	Response   *SearchResp
}

// Route returns the routing data.
func (r NSearchReqLL) Route() []NRoute {
	return []NRoute{{"", r.AreaID, 0}}
}

// NSearchReqDID represents the Normal struct for a request to search for all reports
// authored by the specified Device ID.
type NSearchReqDID struct {
	NRouter
	API
	SearchType NSearchType
	DeviceType string
	DeviceID   string
	MaxResults int
	Routes     []NRoute
	Response   *SearchResp
}

// Route returns the routing data.
func (r NSearchReqDID) Route() []NRoute {
	return r.Routes
}

// SearchResp is the response to creating or updating a report.
type SearchResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
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

// Displays the contents of the Spec_Type custom type.
func (n NServiceRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("NServiceRequest\n")
	ls.AddF("Location - area: %v\n", n.Area)
	return ls.Box(80)
}

// Displays the contents of the Spec_Type custom type.
func (c NServicesResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NServicesResponse\n")
	ls.AddF("Message: %q%s", c.Message, c.Services)
	return ls.Box(90)
}

func (s NService) String() string {
	// r := fmt.Sprintf("  %s-%s-%d-%d  %-40s  %v", s.AdpID, s.AreaID, s.ProviderID, s.ID, s.Name, s.Categories)
	r := fmt.Sprintf("  %-20s  %-40s  %v", s.MID(), s.Name, s.Categories)
	return r
}

// Displays the contents of the Spec_Type custom type.
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

// Displays the contents of the Spec_Type custom type.
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

// Displays the contents of the Spec_Type custom type.
func (c NCreateResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("NCreateResponse\n")
	ls.AddF("Message: %s\n", c.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", c.ID, c.AuthorID)
	return ls.Box(80)
}

// Displays the contents of the Spec_Type custom type.
func (r SearchReqBase) String() string {
	ls := new(common.LogString)
	ls.AddS("SearchReqBase\n")
	ls.AddF("Search type: %s\n", r.SearchType)
	ls.AddF("Device - type %s  ID: %s\n", r.DeviceType, r.DeviceID)
	ls.AddF("GeoLoc - lat: %v (%f)  lon: %v (%f)  Radius: %v (%d)\n", r.Latitude, r.LatitudeV, r.Longitude, r.LongitudeV, r.Radius, r.RadiusV)
	ls.AddF("Address: %s, %s   %s\n", r.City, r.State, r.Zip)
	ls.AddF("Max results: %v (%d)\n", r.MaxResults, r.MaxResultsV)
	return ls.Box(80)
}
