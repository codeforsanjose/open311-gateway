package structs

import (
	"fmt"

	"github.com/codeforsanjose/open311-gateway/common"
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
//                                      STRINGS
// =======================================================================================

// Displays the NServices custom type.
func (r NServices) String() string {
	ls := new(common.FmtBoxer)
	ls.AddS("NServices\n")
	for _, s := range r {
		ls.AddF("%s\n", s)
	}
	return ls.Box(80)
}

// Displays the NService custom type.
func (s NService) String() string {
	r := fmt.Sprintf("  %-20s  %-40s  %-14s %v", s.MID(), s.Name, s.Group, s.Keywords)
	return r
}

// SString should be used to display one NService struct.
func (s NService) SString() string {
	ls := new(common.FmtBoxer)
	ls.AddS("NService\n")
	ls.AddF("MID: %s    Name: %s\n", s.MID(), s.Name)
	ls.AddF("Group: %s   Keywords: %v\n", s.Group, s.Keywords)
	ls.AddF("ServiceNotice: %v\n", s.ServiceNotice)
	return ls.Box(70)
}

// Displays the NServiceRequest custom type.
func (r NServiceRequest) String() string {
	ls := new(common.FmtBoxer)
	ls.AddF("NServiceRequest\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("Location - area: %v\n", r.Area)
	return ls.Box(80)
}

// Displays the NServicesResponse custom type.
func (r NServicesResponse) String() string {
	ls := new(common.FmtBoxer)
	ls.AddS("NServicesResponse\n")
	ls.AddS(r.NResponseCommon.String())
	ls.AddF("Message: %q%s", r.Message, r.Services)
	return ls.Box(90)
}
