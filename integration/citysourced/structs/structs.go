package structs

import (
	"Gateway311/gateway/common"
	"fmt"
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
//                                      SERVICES
// =======================================================================================

// NServiceRequest is used to get list of services available to the user.
type NServiceRequest struct {
	City string
}

// Displays the contents of the Spec_Type custom type.
func (c NServiceRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("Services Request\n")
	ls.AddF("Location - city: %v\n", c.City)
	return ls.Box(80)
}

// ------------------------------- Services -------------------------------

// NServicesResponse is the returned struct for a Services request.
type NServicesResponse struct {
	Message  string
	Services NServices
}

// Displays the contents of the Spec_Type custom type.
func (c NServicesResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("Services Response\n")
	ls.AddF("Message: %q%s", c.Message, c.Services)
	return ls.Box(90)
}

// NServices contains a list of Services.
type NServices []NService

// Displays the contents of the Spec_Type custom type.
func (c NServices) String() string {
	ls := new(common.LogString)
	ls.AddS("NServices\n")
	for _, s := range c {
		ls.AddF("%s\n", s)
	}
	return ls.Box(80)
}

// ------------------------------- Service -------------------------------

// NService represents a Service.  The ID is a combination of the BackEnd Type (IFID),
// the AreaID (i.e. the City id), ProviderID (in case the provider has multiple interfaces),
// and the Service ID.
type NService struct {
	ServiceID  `json:"id"`
	Name       string   `json:"name"`
	Categories []string `json:"catg"`
}

func (s NService) String() string {
	// r := fmt.Sprintf("  %s-%s-%d-%d  %-40s  %v", s.IFID, s.AreaID, s.ProviderID, s.ID, s.Name, s.Categories)
	r := fmt.Sprintf("  %-20s  %-40s  %v", s.MID(), s.Name, s.Categories)
	return r
}

// ------------------------------- ServiceID -------------------------------

// ServiceID provides the JSON marshalling conversion between the JSON "ID" and
// the Backend Interface Type, AreaID (City id), ProviderID, and Service ID.
type ServiceID struct {
	IFID       string
	AreaID     string
	ProviderID int
	ID         int
}

// MID creates the Master ID string for the Service.
func (s ServiceID) MID() string {
	return fmt.Sprintf("%s-%s-%d-%d", s.IFID, s.AreaID, s.ProviderID, s.ID)
}

// =======================================================================================
//                                      CREATE
// =======================================================================================

// NCreateRequest is used to create a new Report.  It is the "native" format of the
// data, and is used by the Engine and all backend Adapters.
type NCreateRequest struct {
	API
	TypeID      int
	Type        string
	DeviceType  string
	DeviceModel string
	DeviceID    string
	Latitude    float64
	Longitude   float64
	Address     string
	City        string
	State       string
	Zip         string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	IsAnonymous bool
	Description string
}

// Displays the contents of the Spec_Type custom type.
func (c NCreateRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - Create\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", c.DeviceType, c.DeviceModel, c.DeviceID)
	ls.AddF("Request - type: (%v) %q\n", c.TypeID, c.Type)
	ls.AddF("Location - lat: %v lon: %v\n", c.Latitude, c.Longitude)
	ls.AddF("          %s, %s   %s\n", c.City, c.State, c.Zip)
	// if math.Abs(c.latitude) > 1 {
	// 	ls.AddF("Location - lat: %v(%q)  lon: %v(%q)\n", c.latitude, c.Latitude, c.longitude, c.Longitude)
	// }
	// if len(c.City) > 1 {
	// 	ls.AddF("          %s, %s   %s\n", c.City, c.State, c.Zip)
	// }
	ls.AddF("Description: %q\n", c.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", c.IsAnonymous, c.FirstName, c.LastName, c.Email, c.Phone)
	return ls.Box(80)
}

// NCreateResponse is the response to creating or updating a report.
type NCreateResponse struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

// Displays the contents of the Spec_Type custom type.
func (c NCreateResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - Resp\n")
	ls.AddF("Message: %s\n", c.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", c.ID, c.AuthorID)
	return ls.Box(80)
}

// =======================================================================================
//                                      SEARCH
// =======================================================================================

// SearchReqBase is used to create a report.
type SearchReqBase struct {
	API
	DeviceType string  `json:"deviceType" xml:"deviceType"`
	DeviceID   string  `json:"deviceId" xml:"deviceId"`
	Latitude   string  `json:"latitude" xml:"latitude"`
	latitude   float64 //
	Longitude  string  `json:"longitude" xml:"longitude"`
	longitude  float64 //
	Radius     string  `json:"radius" xml:"radius"`
	radius     int     // in meters
	Address    string  `json:"address" xml:"address"`
	City       string  `json:"city" xml:"city"`
	State      string  `json:"state" xml:"state"`
	Zip        string  `json:"zip" xml:"zip"`
	MaxResults string  `json:"maxResults" xml:"maxResults"`
	maxResults int     //
	SearchType string  //
}

// SearchResp is the response to creating or updating a report.
type SearchResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}
