package request

import (
	"math"
	"strconv"

	"Gateway311/engine/common"
	"Gateway311/engine/router"
	"Gateway311/engine/structs"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/davecgh/go-spew/spew"
)

// =======================================================================================
//                                      REQUEST
// =======================================================================================
func processCreate(r *rest.Request) (interface{}, error) {
	op := CreateReq{}
	if err := op.init(r); err != nil {
		log.Debug("op failed: %s", spew.Sdump(op))
		log.Debug("Error: %s", err)
		return nil, err
	}
	log.Debug("op: %s", spew.Sdump(op))
	return op.run()
}

// CreateReqBase represents a new report.  It is an anonymous field in the CreateReq
// struct.
type CreateReqBase struct {
	structs.API
	MID         structs.ServiceID `json:"srvId" xml:"srvId"`
	ServiceName string            `json:"srvName" xml:"srvName"`
	DeviceType  string            `json:"deviceType" xml:"deviceType"`
	DeviceModel string            `json:"deviceModel" xml:"deviceModel"`
	DeviceID    string            `json:"deviceID" xml:"deviceID"`
	Latitude    string            `json:"latitude" xml:"latitude"`
	LatitudeV   float64           //
	Longitude   string            `json:"longitude" xml:"longitude"`
	LongitudeV  float64           //
	Address     string            `json:"address" xml:"address"`
	City        string            `json:"city" xml:"city"`
	State       string            `json:"state" xml:"state"`
	Zip         string            `json:"zip" xml:"zip"`
	FirstName   string            `json:"firstName" xml:"firstName"`
	LastName    string            `json:"lastName" xml:"lastName"`
	Email       string            `json:"email" xml:"email"`
	Phone       string            `json:"phone" xml:"phone"`
	IsAnonymous string            `json:"isAnonymous" xml:"isAnonymous"`
	isAnonymous bool              //
	Description string            `json:"Description" xml:"Description"`
	response    *structs.NCreateResponse
}

// CreateReq represents a new Report.  The struct is composed of three anonymous structs:
//
// CreateReqBase - represents a new report.  Contains all fields to unmarshal a
// a request to create a new report.
// cType - defined in common.go, responsible for unmarshaling and parsing the request.
// cIFace - defined in common.go, an interface type for parseQP() and validate() methods.
type CreateReq struct {
	cType         //
	cIface        //
	bkend  string //
	CreateReqBase
}

func (c *CreateReq) newNCreate() (structs.NCreateRequest, error) {
	n := structs.NCreateRequest{
		MID:         c.MID,
		Type:        c.ServiceName,
		DeviceType:  c.DeviceType,
		DeviceModel: c.DeviceModel,
		DeviceID:    c.DeviceID,
		Latitude:    c.LatitudeV,
		Longitude:   c.LongitudeV,
		Address:     c.Address,
		State:       c.State,
		Zip:         c.Zip,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Email:       c.Email,
		Phone:       c.Phone,
		IsAnonymous: c.isAnonymous,
		Description: c.Description,
	}
	return n, nil

}

// validate the unmarshaled data.
func (c *CreateReq) validate() error {
	if x, err := strconv.ParseFloat(c.Latitude, 64); err == nil {
		c.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(c.Longitude, 64); err == nil {
		c.LongitudeV = x
	}
	if x, err := strconv.ParseBool(c.IsAnonymous); err == nil {
		c.isAnonymous = x
	}
	return nil
}

// parseQP unloads any query parms in the request.
func (c *CreateReq) parseQP(r *rest.Request) error {
	return nil
}

// init is called to load and initialize the CreateReq.  It calls cType.load():
// 1. Decodes the input payload.
// 2. Calls parseQP to parse and load any query parms into the struct.
// 3. Calls validate() to check all inputs.
func (c *CreateReq) init(r *rest.Request) error {
	c.load(c, r)
	adp, err := router.GetAdapter(c.MID.IFID)
	if err != nil {
		log.Warning("Unable to get adapter for id: %s", c.MID.IFID)
		return err
	}
	c.bkend = adp.ID
	return nil
}

// run sends the request to the appropriate Adapter, and waits for a reponse.
func (c *CreateReq) run() (interface{}, error) {
	rqst, err := c.newNCreate()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	r, err := router.NewRPCCall("Report.Create", rqst, c.adapterReply)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = r.Run()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	log.Debug("After run - r: %s", r)
	return c.response, err
}

// adapterReply processes the reply returned from the RPC call, by placing a
// pointer to the response in CreateReq.response.
func (c *CreateReq) adapterReply(ndata interface{}) error {
	c.response = ndata.(*structs.NCreateResponse)
	return nil
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// String displays the contents of the CreateReqBase type.
func (c CreateReqBase) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - CreateReqBase\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", c.DeviceType, c.DeviceModel, c.DeviceID)
	ls.AddF("Request - id: %q  name: %q\n", c.MID.MID(), c.ServiceName)
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

// String displays the contents of the CreateReqBase type.
func (c CreateReq) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - CreateReq\n")
	ls.AddF("Backend: %s\n", c.bkend)
	ls.AddS(c.CreateReqBase.String())
	return ls.Box(90)
}
