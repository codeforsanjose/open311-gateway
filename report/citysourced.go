package report

import (
	"Gateway311/common"
	"encoding/xml"
)

// ==============================================================================================================================
//                                      CSReport
// ==============================================================================================================================

// CSReport represents the XML payload for a report request to CitySourced.
type CSReport struct {
	XMLName           xml.Name `xml:"CsRequest"`
	APIAuthKey        string   `json:"ApiAuthKey" xml:"ApiAuthKey"`
	APIRequestType    string   `json:"ApiRequestType" xml:"ApiRequestType"`
	APIRequestVersion string   `json:"ApiRequestVersion" xml:"ApiRequestVersion"`
	DeviceType        string   `json:"DeviceType" xml:"DeviceType"`
	DeviceModel       string   `json:"DeviceModel" xml:"DeviceModel"`
	DeviceID          string   `json:"DeviceId" xml:"DeviceId"`
	RequestType       string   `json:"RequestType" xml:"RequestType"`
	RequestTypeID     int64    `json:"RequestTypeId" xml:"RequestTypeId"`
	ImageURL          string   `json:"ImageUrl" xml:"ImageUrl"`
	ImageURLXl        string   `json:"ImageUrlXl" xml:"ImageUrlXl"`
	ImageURLLg        string   `json:"ImageUrlLg" xml:"ImageUrlLg"`
	ImageURLMd        string   `json:"ImageUrlMd" xml:"ImageUrlMd"`
	ImageURLSm        string   `json:"ImageUrlSm" xml:"ImageUrlSm"`
	ImageURLXs        string   `json:"ImageUrlXs" xml:"ImageUrlXs"`
	Latitude          float64  `json:"Latitude" xml:"Latitude"`
	Longitude         float64  `json:"Longitude" xml:"Longitude"`
	Directionality    string   `json:"Directionality" xml:"Directionality"`
	Description       string   `json:"Description" xml:"Description"`
	AuthorNameFirst   string   `json:"AuthorNameFirst" xml:"AuthorNameFirst"`
	AuthorNameLast    string   `json:"AuthorNameLast" xml:"AuthorNameLast"`
	AuthorEmail       string   `json:"AuthorEmail" xml:"AuthorEmail"`
	AuthorTelephone   string   `json:"AuthorTelephone" xml:"AuthorTelephone"`
	AuthorIsAnonymous bool     `json:"AuthorIsAnonymous" xml:"AuthorIsAnonymous"`
	URLDetail         string   `json:"UrlDetail" xml:"UrlDetail"`
	URLShortened      string   `json:"UrlShortened" xml:"UrlShortened"`
}

// Displays the contents of the Spec_Type custom type.
func (u CSReport) String() string {
	ls := new(common.LogString)
	ls.AddS("CSReport\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", u.DeviceType, u.DeviceModel, u.DeviceID)
	ls.AddF("Request - type: %q  id: %d\n", u.RequestType, u.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v\n", u.Latitude, u.Longitude)
	ls.AddF("Description: %q\n", u.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", u.AuthorIsAnonymous, u.AuthorNameFirst, u.AuthorNameLast, u.AuthorEmail, u.AuthorTelephone)
	return ls.Box(80)
}

// RequestResp represents the XML payload for a report request to CitySourced.
type RequestResp struct {
	Message      string `json:"Message" xml:"Message"`
	RepordID     string `json:"ReportId" xml:"ReportId"`
	AuthorID     string `json:"AuthorId" xml:"AuthorId"`
	ResponseTime string `json:"ResponseTime" xml:"ResponseTime"`
}

// ==============================================================================================================================
//                                      CSReportResp
// ==============================================================================================================================

// CSReportResp is the response to creating or updating a report.
type CSReportResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}
