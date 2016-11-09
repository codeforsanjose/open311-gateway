package structs

import (
	"fmt"

	"github.com/codeforsanjose/open311-gateway/common"
)

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
//                                      STRINGS
// =======================================================================================

// Displays the NSearchRequestLL custom type.
func (r NSearchRequestLL) String() string {
	ls := new(common.FmtBoxer)
	ls.AddF("NSearchRequestLL\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("Lat: %v  Lng: %v   Radius: %v AreaID: %q\n", r.Latitude, r.Longitude, r.Radius, r.AreaID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	return ls.Box(80)
}

// Displays the NSearchRequestDID custom type.
func (r NSearchRequestDID) String() string {
	ls := new(common.FmtBoxer)
	ls.AddF("NSearchRequestDID\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("Device type: %v  ID: %v\n", r.DeviceType, r.DeviceID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	return ls.Box(80)
}

// Displays the NSearchRequestRID custom type.
func (r NSearchRequestRID) String() string {
	ls := new(common.FmtBoxer)
	ls.AddF("NSearchRequestRID\n")
	ls.AddS(r.NRequestCommon.String())
	ls.AddF("ReportID: %v    AreaID: %q\n", r.RID.String(), r.AreaID)
	ls.AddF("RouteList: %v\n", r.RouteList.String())
	return ls.Box(80)
}

// Displays the NSearchResponse custom type.
func (r NSearchResponse) String() string {
	ls := new(common.FmtBoxer)
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
	ls := new(common.FmtBoxer)
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
