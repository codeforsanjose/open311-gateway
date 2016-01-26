package request

import (
	"Gateway311/adapters/citysourced/data"
	"Gateway311/adapters/citysourced/structs"
	"Gateway311/engine/common"
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

// ================================================================================================
//                                      CREATE
// ================================================================================================

// Create is the RPC container struct for the Create service.  This service creates
// a new 311 report.
type Create struct{}

// Run mashals and sends the Create request to the proper back-end, and returns
// the response in Native format.
func (c *Create) Run(rqst *structs.NCreateRequest, resp *structs.NCreateResponse) error {
	fmt.Printf("resp: %p\n", resp)
	fmt.Println(rqst)
	irqst, err := c.makeI(rqst)
	fmt.Printf("rqst: %s\n", *irqst)
	r, err := irqst.Process()
	r.makeN(resp)
	fmt.Printf("resp: %p\n%s\n", resp, *resp)
	return err
}

func (c *Create) makeI(rqst *structs.NCreateRequest) (*ICreateReq, error) {
	// sp, err := router.ServiceProvider(c.TypeIDV)
	// if err != nil {
	// 	return nil, fmt.Errorf("Unable to retrieve Service Provider for Service Type: %v", c.TypeIDV)
	// }

	provider, err := data.MIDProvider(rqst.MID)
	if err != nil {
		return nil, err
	}
	irqst := ICreateReq{
		APIAuthKey:        provider.Key,
		APIRequestType:    "CreateThreeOneOne",
		APIRequestVersion: provider.APIVersion,
		DeviceType:        rqst.DeviceType,
		DeviceModel:       rqst.DeviceModel,
		DeviceID:          rqst.DeviceID,
		RequestType:       rqst.Type,
		RequestTypeID:     rqst.MID.ID,
		Latitude:          rqst.Latitude,
		Longitude:         rqst.Longitude,
		Description:       rqst.Description,
		AuthorNameFirst:   rqst.FirstName,
		AuthorNameLast:    rqst.LastName,
		AuthorEmail:       rqst.Email,
		AuthorTelephone:   rqst.Phone,
		AuthorIsAnonymous: rqst.IsAnonymous,
	}
	return &irqst, nil
}

// ================================================================================================
//                                      CitySourced
// ================================================================================================

// ICreateReq represents the XML payload for a report request to CitySourced.
type ICreateReq struct {
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
func (r *ICreateReq) Process() (*ICreateReqResp, error) {
	// log.Printf("%s\n", r)
	fail := func(err error) (*ICreateReqResp, error) {
		response := ICreateReqResp{
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

	var response ICreateReqResp
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fail(err)
	}

	return &response, nil
}

// Displays the contents of the Spec_Type custom type.
func (r ICreateReq) String() string {
	ls := new(common.LogString)
	ls.AddS("City Sourced - Create\n")
	ls.AddF("API - auth: %q  RequestType: %q  Version: %q\n", r.APIAuthKey, r.APIRequestType, r.APIRequestVersion)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", r.DeviceType, r.DeviceModel, r.DeviceID)
	ls.AddF("Request - type: %q  id: %d\n", r.RequestType, r.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v\n", r.Latitude, r.Longitude)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", r.AuthorIsAnonymous, r.AuthorNameFirst, r.AuthorNameLast, r.AuthorEmail, r.AuthorTelephone)
	return ls.Box(80)
}

// ------------------------------------------------------------------------------------------------

// ICreateReqResp is the response to creating or updating a report.
type ICreateReqResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

// makeN converts the CitySourced data response for a Create request into
// the Native structs.NCreateResponse format.
func (u ICreateReqResp) makeN(resp *structs.NCreateResponse) {
	resp.Message = u.Message
	resp.ID = u.ID
	resp.AuthorID = u.AuthorID
}

// Displays the contents of the Spec_Type custom type.
func (u ICreateReqResp) String() string {
	ls := new(common.LogString)
	ls.AddS("ICreateReqResp\n")
	ls.AddF("Message: %v\n", u.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", u.ID, u.AuthorID)
	return ls.Box(80)
}
