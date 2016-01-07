package request

import (
	"CitySourcedAPI/logs"
	"Gateway311/gateway/common"
	"bytes"
	"encoding/xml"
	"net/http"
)

// ================================================================================================
//                                      Search - Lat/Lng
// ================================================================================================

// CSSearchLLReq represents the CitySourced XML payload for a search by Lat/Lng.
type CSSearchLLReq struct {
	XMLName           xml.Name `xml:"CsSearch"`
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
func (r *CSSearchLLReq) Process() (*CSSearchResp, error) {
	// log.Printf("%s\n", r)
	fail := func(err error) (*CSSearchResp, error) {
		response := CSSearchResp{
			Message: "Failed",
		}
		return &response, err
	}

	var payload = new(bytes.Buffer)
	{
		enc := xml.NewEncoder(payload)
		enc.Indent("  ", "    ")
		enc.Encode(r)
	}
	// log.Printf("Payload:\n%v\n", payload.String())

	url := "http://localhost:5050/api/"
	resp, err := http.Post(url, "application/xml", payload)
	if err != nil {
		return fail(err)
	}

	var response CSSearchResp
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fail(err)
	}

	return &response, nil
}

// Displays the contents of the Spec_Type custom type.
func (r CSSearchLLReq) String() string {
	ls := new(common.LogString)
	ls.AddS("City Sourced - Search LL\n")
	ls.AddF("Location - lat: %v  lon: %v\n", r.Latitude, r.Longitude)
	ls.AddF("MaxResults: %v  IncludeDetails: %v\n", r.MaxResults, r.IncludeDetails)
	ls.AddF("Date Range - start: %v  end: %v\n", r.DateRangeStart, r.DateRangeEnd)
	return ls.Box(80)
}

// ================================================================================================
//                                      Search - Device ID
// ================================================================================================

// CSSearchDIDReq represents the XML payload for a report request to CitySourced.
type CSSearchDIDReq struct {
	XMLName           xml.Name `xml:"CsSearch"`
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

// Displays the contents of the Spec_Type custom type.
func (r CSSearchDIDReq) String() string {
	ls := new(common.LogString)
	ls.AddS("City Sourced - Search\n")
	ls.AddF("Device - type %s  ID: %s\n", r.DeviceType, r.DeviceID)
	ls.AddF("MaxResults: %v  IncludeDetails: %v\n", r.MaxResults, r.IncludeDetails)
	ls.AddF("Date Range - start: %v  end: %v\n", r.DateRangeStart, r.DateRangeEnd)
	return ls.Box(80)
}

// ================================================================================================
//                                      Search - Zip
// ================================================================================================

// CSSearchZipReq represents the XML payload for a report request to CitySourced.
type CSSearchZipReq struct {
	XMLName           xml.Name `xml:"CsSearch"`
	APIAuthKey        string   `json:"ApiAuthKey" xml:"ApiAuthKey"`
	APIRequestType    string   `json:"ApiRequestType" xml:"ApiRequestType"`
	APIRequestVersion string   `json:"ApiRequestVersion" xml:"ApiRequestVersion"`
	Zip               string   `json:"ZipCode" xml:"ZipCode"`
	MaxResults        int      `json:"MaxResults" xml:"MaxResults"`
	IncludeDetails    bool     `json:"IncludeDetails" xml:"IncludeDetails"`
	DateRangeStart    string   `json:"DateRangeStart" xml:"DateRangeStart"`
	DateRangeEnd      string   `json:"DateRangeEnd" xml:"DateRangeEnd"`
}

// Displays the contents of the Spec_Type custom type.
func (r CSSearchZipReq) String() string {
	ls := new(common.LogString)
	ls.AddS("City Sourced - Search\n")
	ls.AddF("Location - zip: %v\n", r.Zip)
	ls.AddF("MaxResults: %v  IncludeDetails: %v\n", r.MaxResults, r.IncludeDetails)
	ls.AddF("Date Range - start: %v  end: %v\n", r.DateRangeStart, r.DateRangeEnd)
	return ls.Box(80)
}

// ------------------------------------------------------------------------------------------------

// CSSearchResp contains the search results.
type CSSearchResp struct {
	XMLName      xml.Name            `xml:"CsResponse"`
	Message      string              `xml:"Message"`
	ResponseTime string              `xml:"ResponseTime"`
	Reports      CSSearchRespReports `xml:"Reports"`
}

// Displays the contents of the Spec_Type custom type.
func (r CSSearchResp) String() string {
	ls := new(common.LogString)
	ls.AddS("CSSearchResp\n")
	ls.AddF("Count: %v RspTime: %v Message: %v\n", r.Reports.ReportCount, r.ResponseTime, r.Message)
	for _, x := range r.Reports.Reports {
		ls.AddS(x.String())
	}
	return ls.Box(80)
}

// CSSearchRespReports is the <Reports> sub-element in the CitySourced XML response.  It contains
// a list of the reports meeting the search criteria.
type CSSearchRespReports struct {
	ReportCount int               `xml:"ReportCount"`
	Reports     []*CSSearchReport `xml:"Report"`
}

// CSSearchReport is the <Report> sub-element in the CitySourced XML response.
type CSSearchReport struct {
	XMLName           xml.Name `xml:"Report" json:"Report"`
	ID                int64    `json:"Id" xml:"Id"`
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

func (s CSSearchReport) String() string {
	ls := new(logs.LogString)
	ls.AddF("Report %d\n", s.ID)
	ls.AddF("DateCreated \"%v\"\n", s.DateCreated)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", s.DeviceType, s.DeviceModel, s.DeviceID)
	ls.AddF("Request - type: %q  id: %q\n", s.RequestType, s.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v  directionality: %q\n", s.Latitude, s.Longitude, s.Directionality)
	ls.AddF("          %s, %s   %s\n", s.City, s.State, s.ZipCode)
	ls.AddF("Votes: %d\n", s.Votes)
	ls.AddF("Description: %q\n", s.Description)
	ls.AddF("Images - std: %s\n", s.ImageURL)
	if len(s.ImageURLXs) > 0 {
		ls.AddF("          XS: %s\n", s.ImageURLXs)
	}
	if len(s.ImageURLSm) > 0 {
		ls.AddF("          SM: %s\n", s.ImageURLSm)
	}
	if len(s.ImageURLMd) > 0 {
		ls.AddF("          XS: %s\n", s.ImageURLMd)
	}
	if len(s.ImageURLLg) > 0 {
		ls.AddF("          XS: %s\n", s.ImageURLLg)
	}
	if len(s.ImageURLXl) > 0 {
		ls.AddF("          XS: %s\n", s.ImageURLXl)
	}
	ls.AddF("Author(anon: %v) %s %s  Email: %s  Tel: %s\n", s.AuthorIsAnonymous, s.AuthorNameFirst, s.AuthorNameLast, s.AuthorEmail, s.AuthorTelephone)
	ls.AddF("SLA: %s\n", s.TicketSLA)
	return ls.Box(80)
}
