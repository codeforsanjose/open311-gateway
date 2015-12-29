package integration

import (
	"Gateway311/common"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
The integration package only knows about each of the supported intergrations.  It does NOT
know about the "standard" data formats in the report package.  To avoid circular imports,
it is the responsibility of the report package to:
1. Select the proper backend (CitySourced, SeeClickFix, etc).
2. Translate the standard format into the the proper struct for the backend (Create, Update, Find, etc.).
3. Initiate the backend service.
4. Translate the response back to the standard format.
*/

// runCS communicates with the specified backend to process the request.
func runCS(payload []byte) (body []byte, err error) {
	fmt.Printf("[runCS]\n")

	url := "http://localhost:5050/api/"
	resp, err := http.Post(url, "application/xml", bytes.NewBuffer(payload))

	// url := "http://localhost:5050/api/"
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(xmlPayload))
	// req.Header.Set("Content-Type", "application/xml")
	// client := &http.Client{}
	// resp, err := client.Do(req)

	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	fmt.Println("response\n   Status:", resp.Status)
	fmt.Println("   Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("   Body:", string(body))

	return body, nil
}

// ================================================================================================
//                                      CSReport
// ================================================================================================

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
	RequestTypeID     int      `json:"RequestTypeId" xml:"RequestTypeId"`
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

// Process executes the request to create a new report.
func (r *CSReport) Process() (*CSReportResp, error) {
	// log.Printf("%s\n", r)
	fail := func(err error) (*CSReportResp, error) {
		response := CSReportResp{
			Message:  "Failed",
			ID:       "",
			AuthorID: "",
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

	var response CSReportResp
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fail(err)
	}

	return &response, nil
}

// Displays the contents of the Spec_Type custom type.
func (r CSReport) String() string {
	ls := new(common.LogString)
	ls.AddS("CSReport\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", r.DeviceType, r.DeviceModel, r.DeviceID)
	ls.AddF("Request - type: %q  id: %d\n", r.RequestType, r.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v\n", r.Latitude, r.Longitude)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", r.AuthorIsAnonymous, r.AuthorNameFirst, r.AuthorNameLast, r.AuthorEmail, r.AuthorTelephone)
	return ls.Box(80)
}

// ------------------------------------------------------------------------------------------------

// CSReportResp is the response to creating or updating a report.
type CSReportResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

// Displays the contents of the Spec_Type custom type.
func (u CSReportResp) String() string {
	ls := new(common.LogString)
	ls.AddS("CSReportResp\n")
	ls.AddF("Message: %v\n", u.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", u.ID, u.AuthorID)
	return ls.Box(80)
}
