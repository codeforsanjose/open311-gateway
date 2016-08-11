package search

import (
	"bytes"
	"encoding/xml"
	"net/http"

	"github.com/open311-gateway/adapters/citysourced/common"

	log "github.com/jeffizhungry/logrus"
)

// ================================================================================================
//                                      SEARCH LL
// ================================================================================================

// RequestLL represents the CitySourced XML payload for a search by Lat/Lng.
type RequestLL struct {
	XMLName           xml.Name `xml:"CsRequest"`
	APIAuthKey        string   `json:"ApiAuthKey" xml:"ApiAuthKey"`
	APIRequestType    string   `json:"ApiRequestType" xml:"ApiRequestType"`
	APIRequestVersion string   `json:"ApiRequestVersion" xml:"ApiRequestVersion"`
	Latitude          float64  `json:"Latitude" xml:"Latitude"`
	Longitude         float64  `json:"Longitude" xml:"Longitude"`
	Radius            int      `json:"Radius" xml:"Radius"`
	MaxResults        int      `json:"MaxResults" xml:"MaxResults"`
	IncludeDetails    bool     `json:"IncludeDetails" xml:"IncludeDetails"`
	DateRangeStart    string   `json:"DateRangeStart" xml:"DateRangeStart"`
	DateRangeEnd      string   `json:"DateRangeEnd" xml:"DateRangeEnd"`
}

// Process executes the request to create a new report.
func (r *RequestLL) Process(url string) (*Response, error) {
	// log.Printf("%s\n", r)
	fail := func(err error) (*Response, error) {
		response := Response{
			Message: "Failed",
		}
		return &response, err
	}

	var payload = new(bytes.Buffer)
	{
		enc := xml.NewEncoder(payload)
		enc.Indent("  ", "    ")
		_ = enc.Encode(r)
	}
	// log.Printf("Payload:\n%v\n", payload.String())

	client := http.Client{Timeout: common.HttpClientTimeout}
	resp, err := client.Post(url, "application/xml", payload)
	if err != nil {
		return fail(err)
	}

	var response Response
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fail(err)
	}

	return &response, nil
}

// ================================================================================================
//                                      SEARCH DID
// ================================================================================================

// RequestDID represents the XML payload for a report request to CitySourced.
type RequestDID struct {
	XMLName           xml.Name `xml:"CsRequest"`
	APIAuthKey        string   `json:"ApiAuthKey" xml:"ApiAuthKey"`
	APIRequestType    string   `json:"ApiRequestType" xml:"ApiRequestType"`
	APIRequestVersion string   `json:"ApiRequestVersion" xml:"ApiRequestVersion"`
	DeviceType        string   `json:"DeviceType" xml:"DeviceType"`
	DeviceID          string   `json:"DeviceId" xml:"DeviceId"`
	MaxResults        int      `json:"MaxResults" xml:"MaxResults"`
	IncludeDetails    bool     `json:"IncludeDetails" xml:"IncludeDetails"`
	DateRangeStart    string   `json:"DateRangeStart" xml:"DateRangeStart"`
	DateRangeEnd      string   `json:"DateRangeEnd" xml:"DateRangeEnd"`
}

// Process executes the request to create a new report.
func (r *RequestDID) Process(url string) (*Response, error) {
	// log.Printf("%s\n", r)
	fail := func(err error) (*Response, error) {
		response := Response{
			Message: "Failed",
		}
		return &response, err
	}

	var payload = new(bytes.Buffer)
	{
		enc := xml.NewEncoder(payload)
		enc.Indent("  ", "    ")
		_ = enc.Encode(r)
	}
	log.Debugf("Payload:\n%v\n", payload.String())

	client := http.Client{Timeout: common.HttpClientTimeout}
	resp, err := client.Post(url, "application/xml", payload)
	if err != nil {
		return fail(err)
	}

	var response Response
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fail(err)
	}

	return &response, nil
}

// ================================================================================================
//                                      SEARCH RID
// ================================================================================================

// RequestRID represents the XML payload for a report request to CitySourced.
type RequestRID struct {
	XMLName           xml.Name `xml:"CsRequest"`
	APIAuthKey        string   `json:"ApiAuthKey" xml:"ApiAuthKey"`
	APIRequestType    string   `json:"ApiRequestType" xml:"ApiRequestType"`
	APIRequestVersion string   `json:"ApiRequestVersion" xml:"ApiRequestVersion"`
	ReportID          string   `json:"ReportId" xml:"ReportId"`
	MaxResults        int      `json:"MaxResults" xml:"MaxResults"`
	IncludeDetails    bool     `json:"IncludeDetails" xml:"IncludeDetails"`
	DateRangeStart    string   `json:"DateRangeStart" xml:"DateRangeStart"`
	DateRangeEnd      string   `json:"DateRangeEnd" xml:"DateRangeEnd"`
}

// Process executes the request to create a new report.
func (r *RequestRID) Process(url string) (*Response, error) {
	// log.Printf("%s\n", r)
	fail := func(err error) (*Response, error) {
		response := Response{
			Message: "Failed",
		}
		return &response, err
	}

	var payload = new(bytes.Buffer)
	{
		enc := xml.NewEncoder(payload)
		enc.Indent("  ", "    ")
		_ = enc.Encode(r)
	}
	log.Debugf("Payload:\n%v\n", payload.String())

	client := http.Client{Timeout: common.HttpClientTimeout}
	resp, err := client.Post(url, "application/xml", payload)
	if err != nil {
		return fail(err)
	}

	var response Response
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fail(err)
	}

	return &response, nil
}

// ================================================================================================
//                                      RESPONSE
// ================================================================================================

// Response contains the search results.
type Response struct {
	XMLName      xml.Name `xml:"CsResponse"`
	Message      string   `xml:"Message"`
	ResponseTime string   `xml:"ResponseTime"`
	Reports      Reports  `xml:"Reports"`
}

// Reports is the <Reports> sub-element in the CitySourced XML response.  It contains
// a list of the reports meeting the search criteria.
type Reports struct {
	ReportCount int       `xml:"ReportCount"`
	Reports     []*Report `xml:"Report"`
}

// Report is the <Report> sub-element in the CitySourced XML response.
type Report struct {
	XMLName           xml.Name `xml:"Report" json:"Report"`
	ID                string   `json:"Id" xml:"Id"`
	DateCreated       string   `json:"DateCreated" xml:"DateCreated"`
	DateUpdated       string   `json:"DateUpdated" xml:"DateUpdated"`
	DeviceType        string   `json:"DeviceType" xml:"DeviceType"`
	DeviceModel       string   `json:"DeviceModel" xml:"DeviceModel"`
	DeviceID          string   `json:"DeviceId" xml:"DeviceId"`
	RequestType       string   `json:"RequestType" xml:"RequestType"`
	RequestTypeID     string   `json:"RequestTypeId" xml:"RequestTypeId"`
	ImageURL          string   `json:"ImageUrl" xml:"ImageUrl"`
	ImageURLXl        string   `json:"ImageUrlXl" xml:"ImageUrlXl"`
	ImageURLLg        string   `json:"ImageUrlLg" xml:"ImageUrlLg"`
	ImageURLMd        string   `json:"ImageUrlMd" xml:"ImageUrlMd"`
	ImageURLSm        string   `json:"ImageUrlSm" xml:"ImageUrlSm"`
	ImageURLXs        string   `json:"ImageUrlXs" xml:"ImageUrlXs"`
	City              string   `json:"City" xml:"City"`
	State             string   `json:"State" xml:"State"`
	ZipCode           string   `json:"ZipCode" xml:"ZipCode"`
	Latitude          string   `xml:"Latitude" json:"Latitude"`
	Longitude         string   `xml:"Longitude" json:"Longitude"`
	Directionality    string   `json:"Directionality" xml:"Directionality"`
	Description       string   `json:"Description" xml:"Description"`
	AuthorNameFirst   string   `json:"AuthorNameFirst" xml:"AuthorNameFirst"`
	AuthorNameLast    string   `json:"AuthorNameLast" xml:"AuthorNameLast"`
	AuthorEmail       string   `json:"AuthorEmail" xml:"AuthorEmail"`
	AuthorTelephone   string   `json:"AuthorTelephone" xml:"AuthorTelephone"`
	AuthorIsAnonymous string   `xml:"AuthorIsAnonymous" json:"AuthorIsAnonymous"`
	URLDetail         string   `json:"UrlDetail" xml:"UrlDetail"`
	URLShortened      string   `json:"UrlShortened" xml:"UrlShortened"`
	Votes             string   `json:"Votes" xml:"Votes"`
	StatusType        string   `json:"StatusType" xml:"StatusType"`
	TicketSLA         string   `json:"TicketSla" xml:"TicketSla"`
}

// ================================================================================================
//                                      STRINGS
// ================================================================================================

// Displays the contents of the Spec_Type custom type.
func (r RequestLL) String() string {
	ls := new(common.LogString)
	ls.AddS("RequestLL\n")
	ls.AddF("Location - lat: %v  lon: %v\n", r.Latitude, r.Longitude)
	ls.AddF("MaxResults: %v  IncludeDetails: %v\n", r.MaxResults, r.IncludeDetails)
	ls.AddF("Date Range - start: %v  end: %v\n", r.DateRangeStart, r.DateRangeEnd)
	return ls.Box(80)
}

// Displays the contents of the Spec_Type custom type.
func (r RequestDID) String() string {
	ls := new(common.LogString)
	ls.AddS("RequestDID\n")
	ls.AddF("Device - type %s  ID: %s\n", r.DeviceType, r.DeviceID)
	ls.AddF("MaxResults: %v  IncludeDetails: %v\n", r.MaxResults, r.IncludeDetails)
	ls.AddF("Date Range - start: %v  end: %v\n", r.DateRangeStart, r.DateRangeEnd)
	return ls.Box(80)
}

// Displays the contents of the Spec_Type custom type.
func (r RequestRID) String() string {
	ls := new(common.LogString)
	ls.AddS("RequestRID\n")
	ls.AddF("ReportID: %s\n", r.ReportID)
	ls.AddF("MaxResults: %v  IncludeDetails: %v\n", r.MaxResults, r.IncludeDetails)
	ls.AddF("Date Range - start: %v  end: %v\n", r.DateRangeStart, r.DateRangeEnd)
	return ls.Box(80)
}

// Displays the contents of the Spec_Type custom type.
func (r Response) String() string {
	ls := new(common.LogString)
	ls.AddS("Response\n")
	ls.AddF("Count: %v RspTime: %v Message: %v\n", r.Reports.ReportCount, r.ResponseTime, r.Message)
	for _, x := range r.Reports.Reports {
		ls.AddS(x.String())
	}
	return ls.Box(90)
}

// Displays the NSearchRequestDID custom type.
func (s Report) String() string {
	ls := new(common.LogString)
	ls.AddF("Report %v\n", s.ID)
	ls.AddF("DateCreated \"%v\"\n", s.DateCreated)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", s.DeviceType, s.DeviceModel, s.DeviceID)
	ls.AddF("Request - type: %q  id: %q\n", s.RequestType, s.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v  directionality: %q\n", s.Latitude, s.Longitude, s.Directionality)
	ls.AddF("          %s, %s   %s\n", s.City, s.State, s.ZipCode)
	ls.AddF("Votes: %v\n", s.Votes)
	ls.AddF("Description: %q\n", s.Description)
	ls.AddF("Images - std: %s\n", s.ImageURL)
	if len(s.ImageURLXs) > 0 {
		ls.AddF("          XS: %s\n", s.ImageURLXs)
	}
	if len(s.ImageURLSm) > 0 {
		ls.AddF("          SM: %s\n", s.ImageURLSm)
	}
	if len(s.ImageURLMd) > 0 {
		ls.AddF("          MD: %s\n", s.ImageURLMd)
	}
	if len(s.ImageURLLg) > 0 {
		ls.AddF("          LG: %s\n", s.ImageURLLg)
	}
	if len(s.ImageURLXl) > 0 {
		ls.AddF("          XL: %s\n", s.ImageURLXl)
	}
	ls.AddF("Author(anon: %v) %s %s  Email: %s  Tel: %s\n", s.AuthorIsAnonymous, s.AuthorNameFirst, s.AuthorNameLast, s.AuthorEmail, s.AuthorTelephone)
	ls.AddF("SLA: %s\n", s.TicketSLA)
	return ls.Box(80)
}
