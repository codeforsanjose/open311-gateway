package request

import (
	"Gateway311/common"
	"Gateway311/integration"

	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// Services looks up the service providers and services for the specified location.
func Services(w rest.ResponseWriter, r *rest.Request) {
	m, _ := url.ParseQuery(r.URL.RawQuery)
	for k, v := range m {
		fmt.Printf("%s: %#v\n", k, v)
	}
	response := "OK!"

	w.WriteJson(&response)
}

// todo: this is a to do item.
// Create creates a new Report and adds it to Reports.
func Create(w rest.ResponseWriter, r *rest.Request) {
	jid, err := strconv.ParseInt(r.PathParam("jid"), 10, 64)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("[Create] - jid: (%T)%d\n", jid, jid)
	report := CreateReport{}
	if err := r.DecodeJsonPayload(&report); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := processCS(&report)
	if err != nil {
		fmt.Println("!! PrepOut failed.")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("   response: %#v\n", response)

	w.WriteJson(&response)
}

// ==============================================================================================================================
//                                      CreateReport
// ==============================================================================================================================

// CreateReport is used to create a report.
type CreateReport struct {
	JID         int     `json:"jid" xml:"jid"`
	ID          string  `json:"id" xml:"id"`
	Type        string  `json:"type" xml:"type"`
	TypeID      string  `json:"typeId" xml:"typeId"`
	DeviceType  string  `json:"deviceType" xml:"deviceType"`
	DeviceModel string  `json:"deviceModel" xml:"deviceModel"`
	DeviceID    string  `json:"deviceID" xml:"deviceID"`
	Latitude    float64 `json:"latitude" xml:"latitude"`
	Longitude   float64 `json:"longitude" xml:"longitude"`
	Address     string  `json:"address" xml:"address"`
	City        string  `json:"city" xml:"city"`
	State       string  `json:"state" xml:"state"`
	Zip         string  `json:"zip" xml:"zip"`
	FirstName   string  `json:"firstName" xml:"firstName"`
	LastName    string  `json:"lastName" xml:"lastName"`
	Email       string  `json:"email" xml:"email"`
	Phone       string  `json:"phone" xml:"phone"`
	IsAnonymous bool    `json:"isAnonymous" xml:"isAnonymous"`
	Description string  `json:"Description" xml:"Description"`
}

// GetAll retrieves all reports, with no filtering.
func (u *CreateReport) GetAll(w rest.ResponseWriter, r *rest.Request) {
	fmt.Printf("[GetAll]\n")
	jid := r.PathParam("jid")
	r.ParseForm()
	fmt.Printf("  jid: (%T)%q Form: %#v\n", jid, jid, r.Form)

	response := CreateReportResp{
		Message:  "success",
		ID:       "12345",
		AuthorID: "99999",
	}
	w.WriteJson(&response)
}

// Get retrieves a single report specified by it's ID.
func (u *CreateReport) Get(w rest.ResponseWriter, r *rest.Request) {
	fmt.Printf("[Get]\n")
	jid := r.PathParam("jid")
	id := r.PathParam("id")
	fmt.Printf("  jid: %s id: %s\n", jid, id)
	response := CreateReportResp{
		Message:  "success",
		ID:       "12345",
		AuthorID: "99999",
	}
	w.WriteJson(&response)
}

// Update updates the report specified by ID.
func (u *CreateReport) Update(w rest.ResponseWriter, r *rest.Request) {
	// id := r.PathParam("id")
	response := CreateReportResp{
		Message:  "success",
		ID:       "12345",
		AuthorID: "99999",
	}
	w.WriteJson(&response)
}

// Delete deletes the report specified by ID.
func (u *CreateReport) Delete(w rest.ResponseWriter, r *rest.Request) {
	// id := r.PathParam("id")
	response := CreateReportResp{
		Message:  "success",
		ID:       "12345",
		AuthorID: "99999",
	}
	w.WriteJson(&response)
}

func toCreateCS(src *CreateReport) (*integration.CSReport, error) {
	requestTypeID, err := strconv.ParseInt(src.TypeID, 10, 64)
	if err != nil {
		fmt.Printf("Unable to parse request type id: %q\n", src.TypeID)
		return nil, fmt.Errorf("Unable to parse request type id: %q", src.TypeID)
	}
	rqst := integration.CSReport{
		APIAuthKey:        "a01234567890z",
		APIRequestType:    "CreateThreeOneOne",
		APIRequestVersion: "1",
		DeviceType:        src.DeviceType,
		DeviceModel:       src.DeviceModel,
		DeviceID:          src.DeviceID,
		RequestType:       src.Type,
		RequestTypeID:     requestTypeID,
		Latitude:          src.Latitude,
		Longitude:         src.Longitude,
		Description:       src.Description,
		AuthorNameFirst:   src.FirstName,
		AuthorNameLast:    src.LastName,
		AuthorEmail:       src.Email,
		AuthorTelephone:   src.Phone,
		AuthorIsAnonymous: src.IsAnonymous,
	}
	return &rqst, nil
}

func fromCreateCS(src *integration.CSReportResp) (*CreateReportResp, error) {
	resp := CreateReportResp{
		Message:  src.Message,
		ID:       src.ID,
		AuthorID: src.AuthorID,
	}
	return &resp, nil
}

func processCS(src *CreateReport) (interface{}, error) {
	rqst, _ := toCreateCS(src)
	resp, _ := rqst.Process(1)
	ourResp, _ := fromCreateCS(resp)

	return ourResp, nil
}

// Displays the contents of the Spec_Type custom type.
func (u CreateReport) String() string {
	ls := new(common.LogString)
	ls.AddS("Report\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", u.DeviceType, u.DeviceModel, u.DeviceID)
	ls.AddF("Request - type: %q  id: %q\n", u.Type, u.TypeID)
	ls.AddF("Location - lat: %v  lon: %v\n", u.Latitude, u.Longitude)
	ls.AddF("          %s, %s   %s\n", u.City, u.State, u.Zip)
	ls.AddF("Description: %q\n", u.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", u.IsAnonymous, u.FirstName, u.LastName, u.Email, u.Phone)
	return ls.Box(80)
}

// ==============================================================================================================================
//                                      CreateReportResp
// ==============================================================================================================================

// CreateReportResp is the response to creating or updating a report.
type CreateReportResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}
