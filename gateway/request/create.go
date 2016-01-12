package request

import (
	"Gateway311/gateway/common"
	"Gateway311/gateway/router"
	"Gateway311/gateway/structs"
	"_sketches/spew"
	"math"

	"fmt"
	"log"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// =======================================================================================
//                                      Request
// =======================================================================================
func processCreate(r *rest.Request) (interface{}, error) {
	op := CreateReq{}
	if err := op.init(r); err != nil {
		return nil, err
	}
	return op.run()
}

// CreateReq is used to create a report.
type CreateReqBase struct {
	structs.API
	ID          string  `json:"id" xml:"id"`
	Type        string  `json:"type" xml:"type"`
	TypeID      string  `json:"typeId" xml:"typeId"`
	TypeIDV     int     //
	DeviceType  string  `json:"deviceType" xml:"deviceType"`
	DeviceModel string  `json:"deviceModel" xml:"deviceModel"`
	DeviceID    string  `json:"deviceID" xml:"deviceID"`
	Latitude    string  `json:"LatitudeV" xml:"LatitudeV"`
	LatitudeV    float64 //
	Longitude   string  `json:"LongitudeV" xml:"LongitudeV"`
	LongitudeV   float64 //
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

// Displays the contents of the Spec_Type custom type.
func (c CreateReqBase) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - Create\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", c.DeviceType, c.DeviceModel, c.DeviceID)
	ls.AddF("Request - type: %q  id: %q(%v)\n", c.Type, c.TypeID, c.TypeIDV)
	ls.AddF("Location - lat: %v(%q)  lon: %v(%q)\n", c.LatitudeV, c.Latitude, c.LongitudeV, c.Longitude)
	ls.AddF("          %s, %s   %s\n", c.City, c.State, c.Zip)
	if math.Abs(c.LatitudeV) > 1 {
		ls.AddF("Location - lat: %v(%q)  lon: %v(%q)\n", c.LatitudeV, c.Latitude, c.LongitudeV, c.Longitude)
	}
	if len(c.City) > 1 {
		ls.AddF("          %s, %s   %s\n", c.City, c.State, c.Zip)
	}
	ls.AddF("Description: %q\n", c.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", c.isAnonymous, c.FirstName, c.LastName, c.Email, c.Phone)
	return ls.Box(80)
}

type CreateReq struct {
	cType  //
	cIface //
	// JID    int    `json:"jid" xml:"jid"`
	bkend string //
	CreateReqBase
}

func (c *CreateReq) validate() {
	if x, err := strconv.ParseInt(c.TypeID, 10, 64); err == nil {
		c.TypeIDV = int(x)
	}
	if x, err := strconv.ParseFloat(c.Latitude, 64); err == nil {
		c.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(c.Longitude, 64); err == nil {
		c.LongitudeV = x
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
	c.load(c, r)
	itype, err := router.ServiceProviderInterface(c.TypeIDV)
	if err != nil {
		return fmt.Errorf("Cannot determine Service Provider Interface for type: %v", c.TypeIDV)
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
	resp, _ := rqst.Process()
	ourResp, _ := fromCreateCS(resp)

	return ourResp, nil
}

// --------------------------- Integrations ----------------------------------------------

func (c *CreateReq) toCreateCS() (*integration.CSReport, error) {
	sp, err := router.ServiceProvider(c.TypeIDV)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Service Provider for Service Type: %v", c.TypeIDV)
	}

	rqst := integration.CSReport{
		APIAuthKey:        sp.Key,
		APIRequestType:    "CreateThreeOneOne",
		APIRequestVersion: sp.APIVersion,
		DeviceType:        c.DeviceType,
		DeviceModel:       c.DeviceModel,
		DeviceID:          c.DeviceID,
		RequestType:       c.Type,
		RequestTypeID:     c.TypeIDV,
		Latitude:          c.LatitudeV,
		Longitude:         c.LongitudeV,
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

// Displays the contents of the Spec_Type custom type.
func (c CreateResp) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - Resp\n")
	ls.AddF("Message: %s\n", c.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", c.ID, c.AuthorID)
	return ls.Box(80)
}

func fromCreateCS(src *integration.CSReportResp) (*CreateResp, error) {
	resp := CreateResp{
		Message:  src.Message,
		ID:       src.ID,
		AuthorID: src.AuthorID,
	}
	return &resp, nil
}
