package request

import (
	"math"
	"reflect"
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

// CreateReqBase is used to create a report.  It is an anonymous field in the CreateReq
// struct.
type CreateReqBase struct {
	structs.API
	MID         structs.ServiceID `json:"srvId" xml:"srvId"`
	SrvName     string            `json:"srvName" xml:"srvName"`
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

// CreateReq creates a new Report.
type CreateReq struct {
	cType  //
	cIface //
	// JID    int    `json:"jid" xml:"jid"`
	bkend string //
	CreateReqBase
}

func (c *CreateReq) newNCreate() (structs.NCreateRequest, error) {
	n := structs.NCreateRequest{
		MID:         c.MID,
		Type:        c.SrvName,
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

func (c *CreateReq) validate() {
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
	adp, err := router.GetAdapter(c.MID.IFID)
	if err != nil {
		log.Warning("Unable to get adapter for id: %s", c.MID.IFID)
		return err
	}
	c.bkend = adp.ID
	return nil
}

func (c *CreateReq) run() (interface{}, error) {
	rqst, err := c.newNCreate()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	r, err := router.NewRPCCall("Create.Run", c.bkend, "", rqst, c.rpcReply)
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

func (c *CreateReq) rpcReply(ndata interface{}) error {
	log.Debug("ndata type: %s", reflect.TypeOf(ndata))
	c.response = ndata.(*structs.NCreateResponse)
	return nil
}

// =======================================================================================
//                                      REQUEST
// =======================================================================================

// String displays the contents of the CreateReqBase type.
func (c CreateReqBase) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - CreateReqBase\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", c.DeviceType, c.DeviceModel, c.DeviceID)
	ls.AddF("Request - id: %q  name: %q\n", c.MID.MID(), c.SrvName)
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
