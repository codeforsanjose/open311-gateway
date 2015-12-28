package request

import (
	"Gateway311/common"
	"Gateway311/integration"
	"Gateway311/router"
	"_sketches/spew"

	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// =======================================================================================
//                                      Request
// =======================================================================================
func processCreate(r *rest.Request) (interface{}, error) {
	report := CreateReq{cType: cType{inputBody: true}}
	if err := report.init(r); err != nil {
		return nil, err
	}
	return report.run()
}

// CreateReq is used to create a report.
type CreateReq struct {
	cType         //
	cIface        //
	JID    int    `json:"jid" xml:"jid"`
	bkend  string //

	ID          string  `json:"id" xml:"id"`
	Type        string  `json:"type" xml:"type"`
	TypeID      string  `json:"typeId" xml:"typeId"`
	typeID      int     //
	DeviceType  string  `json:"deviceType" xml:"deviceType"`
	DeviceModel string  `json:"deviceModel" xml:"deviceModel"`
	DeviceID    string  `json:"deviceID" xml:"deviceID"`
	Latitude    string  `json:"latitude" xml:"latitude"`
	latitude    float64 //
	Longitude   string  `json:"longitude" xml:"longitude"`
	longitude   float64 //
	Address     string  `json:"address" xml:"address"`
	City        string  `json:"city" xml:"city"`
	State       string  `json:"state" xml:"state"`
	Zip         string  `json:"zip" xml:"zip"`
	FirstName   string  `json:"firstName" xml:"firstName"`
	LastName    string  `json:"lastName" xml:"lastName"`
	Email       string  `json:"email" xml:"email"`
	Phone       string  `json:"phone" xml:"phone"`
	IsAnonymous string  `json:"isAnonymous" xml:"isAnonymous"`
	isAnonymous bool    //
	Description string  `json:"Description" xml:"Description"`
}

func (c *CreateReq) validate() {
	if x, err := strconv.ParseInt(c.TypeID, 10, 64); err == nil {
		c.typeID = int(x)
	}

	if x, err := strconv.ParseFloat(c.Latitude, 64); err == nil {
		c.latitude = x
	}

	if x, err := strconv.ParseFloat(c.Longitude, 64); err == nil {
		c.longitude = x
	}

	if x, err := strconv.ParseBool(c.IsAnonymous); err == nil {
		c.isAnonymous = x
	}
	return
}

func (c *CreateReq) parseQP(r *rest.Request) error {
	return nil
}

func (c *CreateReq) init(r *rest.Request) error {
	c.cType.init(c, r)
	itype, err := router.ServiceProviderInterface(c.typeID)
	if err != nil {
		return fmt.Errorf("Cannot determine Service Provider Interface for type: %v", c.typeID)
	}
	c.bkend = itype
	// log.Printf("After init: \n%s\n", c)
	return nil
}

func (c *CreateReq) run() (interface{}, error) {
	switch c.bkend {
	case "CitySourced":
		return c.processCS()
	}
	return nil, fmt.Errorf("Unsupported backend: %q", c.bkend)
}

func (c *CreateReq) processCS() (interface{}, error) {
	log.Printf("[processCS] src: %s", spew.Sdump(c))
	rqst, _ := c.toCreateCS()
	resp, _ := rqst.Process(c.JID)
	ourResp, _ := fromCreateCS(resp)

	return ourResp, nil
}

// Displays the contents of the Spec_Type custom type.
func (c CreateReq) String() string {
	ls := new(common.LogString)
	ls.AddS("Report\n")
	ls.AddF("Parse - body: %t  qp: %t  bkend: %s\n", c.inputBody, c.inputQP, c.bkend)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", c.DeviceType, c.DeviceModel, c.DeviceID)
	ls.AddF("Request - type: %q  id: %q\n", c.Type, c.TypeID)
	if math.Abs(c.latitude) > 1 {
		ls.AddF("Location - lat: %v  lon: %v\n", c.latitude, c.longitude)
	}
	if len(c.City) > 1 {
		ls.AddF("          %s, %s   %s\n", c.City, c.State, c.Zip)
	}
	ls.AddF("Description: %q\n", c.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", c.isAnonymous, c.FirstName, c.LastName, c.Email, c.Phone)
	return ls.Box(80)
}

// --------------------------- Integrations ----------------------------------------------

func (c *CreateReq) toCreateCS() (*integration.CSReport, error) {
	sp, err := router.ServiceProvider(c.typeID)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Service Provider for Service Type: %v", c.typeID)
	}

	rqst := integration.CSReport{
		APIAuthKey:        sp.Key,
		APIRequestType:    "CreateThreeOneOne",
		APIRequestVersion: sp.APIVersion,
		DeviceType:        c.DeviceType,
		DeviceModel:       c.DeviceModel,
		DeviceID:          c.DeviceID,
		RequestType:       c.Type,
		RequestTypeID:     c.typeID,
		Latitude:          c.latitude,
		Longitude:         c.longitude,
		Description:       c.Description,
		AuthorNameFirst:   c.FirstName,
		AuthorNameLast:    c.LastName,
		AuthorEmail:       c.Email,
		AuthorTelephone:   c.Phone,
		AuthorIsAnonymous: c.isAnonymous,
	}
	return &rqst, nil
}

// =======================================================================================
//                                      Response
// =======================================================================================

// CreateResp is the response to creating or updating a report.
type CreateResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

func fromCreateCS(src *integration.CSReportResp) (*CreateResp, error) {
	resp := CreateResp{
		Message:  src.Message,
		ID:       src.ID,
		AuthorID: src.AuthorID,
	}
	return &resp, nil
}
