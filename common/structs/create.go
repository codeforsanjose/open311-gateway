package structs

import "github.com/codeforsanjose/open311-gateway/common"

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
//                                      STRINGS
// =======================================================================================

// Displays the NCreateRequest custom type.
func (r NCreateRequest) String() string {
	ls := new(common.FmtBoxer)
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
	ls := new(common.FmtBoxer)
	ls.AddS("NCreateResponse\n")
	ls.AddS(r.NResponseCommon.String())
	ls.AddF("Message: %s\n", r.Message)
	ls.AddF("ID: %v  AccountID: %v\n", r.ID, r.AccountID)
	return ls.Box(80)
}
