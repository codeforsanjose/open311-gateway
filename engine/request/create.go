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
func processCreate(rqst *rest.Request, rqstID int64) (interface{}, error) {
	creq := CreateReq{}
	if err := creq.init(rqst, rqstID); err != nil {
		log.Errorf("processCreate failed - %s", err)
		log.Errorf("CreateReq: %s", spew.Sdump(creq))
		return nil, err
	}
	return creq.run()
}

// CreateReq represents a new Report.  The struct is composed of three anonymous structs:
//
// CreateReqBase - represents a new report.  Contains all fields to unmarshal a
// a request to create a new report.
// cType - defined in common.go, responsible for unmarshaling and parsing the request.
// cIFace - defined in common.go, an interface type for parseQP() and validate() methods.
type CreateReq struct {
	cType
	cIface
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
	response    struct {
		cRType
		*structs.NCreateResponse
	}
}

func (r *CreateReq) newNCreate() (structs.NCreateRequest, error) {
	n := structs.NCreateRequest{
		NRequestCommon: structs.NRequestCommon{
			Rtype: structs.NRTCreate,
		},
		MID:         r.MID,
		Type:        r.ServiceName,
		DeviceType:  r.DeviceType,
		DeviceModel: r.DeviceModel,
		DeviceID:    r.DeviceID,
		Latitude:    r.LatitudeV,
		Longitude:   r.LongitudeV,
		Address:     r.Address,
		State:       r.State,
		Zip:         r.Zip,
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		Email:       r.Email,
		Phone:       r.Phone,
		IsAnonymous: r.isAnonymous,
		Description: r.Description,
	}
	return n, nil

}

// validate the unmarshaled data.
func (r *CreateReq) validate() error {
	if x, err := strconv.ParseFloat(r.Latitude, 64); err == nil {
		r.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(r.Longitude, 64); err == nil {
		r.LongitudeV = x
	}
	if x, err := strconv.ParseBool(r.IsAnonymous); err == nil {
		r.isAnonymous = x
	}
	return nil
}

// parseQP unloads any query parms in the request.
func (r *CreateReq) parseQP(rqst *rest.Request) error {
	return nil
}

// init is called to load and initialize the CreateReq.  It calls cType.load():
// 1. Decodes the input payload.
// 2. Calls parseQP to parse and load any query parms into the struct.
// 3. Calls validate() to check all inputs.
func (r *CreateReq) init(rqst *rest.Request, rqstID int64) error {
	r.load(r, rqstID, rqst)
	_, err := router.GetAdapter(r.MID.AdpID)
	if err != nil {
		log.Warning("Unable to get adapter for id: %s", r.MID.AdpID)
		return err
	}
	return nil
}

// run sends the request to the appropriate Adapter, and waits for a reponse.
func (r *CreateReq) run() (interface{}, error) {
	log.Debug(r.String())
	rqst, err := r.newNCreate()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	rpcCall, err := router.NewRPCCall("Report.Create", &rqst, r.adapterReply)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = rpcCall.Run()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return r.response, err
}

// adapterReply processes the reply returned from the RPC call, by placing a
// pointer to the response in CreateReq.response.
func (r *CreateReq) adapterReply(ndata interface{}) error {
	r.response.NCreateResponse = ndata.(*structs.NCreateResponse)
	r.response.id = r.id
	return nil
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// String displays the contents of the CreateReqBase type.
func (r CreateReq) String() string {
	ls := new(common.LogString)
	ls.AddF("CreateReq - %d\n", r.id)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", r.DeviceType, r.DeviceModel, r.DeviceID)
	ls.AddF("Request - id: %q  name: %q\n", r.MID.MID(), r.ServiceName)
	if math.Abs(r.LatitudeV) > 1 {
		ls.AddF("Location - lat: %v(%q)  lon: %v(%q)\n", r.LatitudeV, r.Latitude, r.LongitudeV, r.Longitude)
	}
	if len(r.City) > 1 {
		ls.AddF("          %s, %s   %s\n", r.City, r.State, r.Zip)
	}
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", r.isAnonymous, r.FirstName, r.LastName, r.Email, r.Phone)
	return ls.Box(80)
}
