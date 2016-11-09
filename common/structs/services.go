package structs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

// UnmarshalJSON implements the conversion from the JSON "ID" to the ServiceID struct.
func (srv *NService) UnmarshalJSON(value []byte) error {
	type T struct {
		ID            int
		Name          string
		Description   string
		Metadata      bool
		Group         string
		Keywords      []string
		ServiceNotice string `json:"service_notice"`
	}
	var t T
	err := json.Unmarshal(value, &t)
	if err != nil {
		return err
	}
	srv.ID = t.ID
	srv.Name = t.Name
	srv.Description = t.Description
	srv.Metadata = t.Metadata
	srv.ServiceNotice = t.ServiceNotice
	srv.Keywords = t.Keywords
	srv.Group = t.Group
	return nil
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

// UnmarshalJSON implements the conversion from the JSON "ID" to the ServiceID struct.
func (s *ServiceID) UnmarshalJSON(value []byte) error {
	cnvInt := func(x string) int {
		y, _ := strconv.ParseInt(x, 10, 64)
		return int(y)
	}
	parts := strings.Split(strings.Trim(string(value), "\" "), "-")
	// log.Debug("[UnmarshalJSON] parts: %+v\n", parts)
	s.AdpID = parts[0]
	s.AreaID = parts[1]
	s.ProviderID = cnvInt(parts[2])
	s.ID = cnvInt(parts[3])
	// log.Debug("[UnmarshalJSON] AdpID: %#v  AreaID: %#v  ProviderID: %#v  ID: %#v\n", s.AdpID, s.AreaID, s.ProviderID, s.ID)
	return nil
}

// MarshalJSON implements the conversion from the ServiceID struct to the JSON "ID".
func (s ServiceID) MarshalJSON() ([]byte, error) {
	// fmt.Printf("  Marshaling s: %#v\n", s)
	return []byte(fmt.Sprintf("\"%s\"", s.MID())), nil
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
