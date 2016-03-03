package request

import (
	"fmt"
	"math"
	"strconv"

	"Gateway311/engine/common"
	"Gateway311/engine/router"
	"Gateway311/engine/structs"

	"github.com/ant0ine/go-json-rest/rest"
)

// createMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Create a Report.
//  1. Loads all input payload and query parms.
//  2. Validates all input.
//  3. Determines the route(s).  Returns error if no valid route(s) is found.
//  4. Converts the input to the Normal form.
//  5. Call RPC Router to process the request.
//  6. Validates and merges Normal form RPC response(s).
//  7. Converts Normal form to response.
//  8. Returns response.
type createMgr struct {
	id int64

	rqst *rest.Request

	req  *CreateReq
	nreq *structs.NCreateRequest

	valid  Validation
	routes structs.NRoutes

	nresp *structs.NCreateResponse
	resp  *CreateResp
}

func processCreate(rqst *rest.Request, rqstID int64) (interface{}, error) {
	log.Debug("starting processCreate()")
	mgr := createMgr{
		rqst:  rqst,
		id:    rqstID,
		req:   &CreateReq{},
		valid: newValidation(),
		resp:  &CreateResp{Message: "Request failed"},
	}

	fail := func(err error) (interface{}, error) {
		log.Errorf("processCreate failed - %s", err)
		return mgr.resp, fmt.Errorf("Create request failed - %s", err.Error())
	}

	if rqstID == 0 {
		mgr.id = router.GetSID()
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != "JSON payload is empty" {
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		log.Errorf("processCreate.validate() failed - %s", err)
		return fail(err)
	}

	mgr.convertRequest()

	if err := mgr.setRoute(); err != nil {
		log.Errorf("processCreate.route() failed - %s", err)
		return fail(err)
	}

	if err := mgr.callRPC(); err != nil {
		log.Errorf("processCreate.callRPC() failed - %s", err)
		return fail(err)
	}

	mgr.convertResponse()

	return mgr.resp, nil
}

// parseQP unloads any query parms in the request.
func (r *createMgr) parseQP(rqst *rest.Request) error {
	return nil
}

// validate the unmarshaled data.
func (r *createMgr) validate() error {
	log.Debug("Starting validate()")
	if err := r.req.convert(); err != nil {
		return err
	}
	r.valid.Set("location", "Location must be in continental US", validateLatLng(r.req.LatitudeV, r.req.LongitudeV))

	if x, err := strconv.ParseBool(r.req.IsAnonymous); err == nil {
		r.req.isAnonymous = x
	}
	log.Debug(r.valid.String())
	if !r.valid.Ok() {
		return r.valid
	}
	return nil
}

// setRoute gets the route(s) to process the request.
func (r *createMgr) setRoute() error {
	var router structs.NRouter = r.nreq
	r.routes = router.GetRoutes()
	if len(r.routes) == 0 {
		return fmt.Errorf("no routes found")
	}
	return nil
}

// callRPC runs the calls to the Adapter(s).
func (r *createMgr) callRPC() error {
	loadData := func(ndata interface{}) error {
		r.nresp = ndata.(*structs.NCreateResponse)
		return nil
	}

	rpcCall, err := router.NewRPCCall("Report.Create", r.nreq, loadData)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = rpcCall.Run()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return err
}

// convertResponse converts the NCreateResponse{} to a CreateResp{}
func (r *createMgr) convertResponse() {
	r.resp = &CreateResp{
		Message:  r.nresp.Message,
		ID:       r.nresp.RID.RID(),
		AuthorID: r.nresp.AuthorID,
	}
}

func (r *createMgr) convertRequest() {
	r.nreq = &structs.NCreateRequest{
		NRequestCommon: structs.NRequestCommon{
			ID: structs.NID{
				RqstID: r.id,
			},
			Rtype: structs.NRTCreate,
		},
		MID:         r.req.MID,
		Type:        r.req.ServiceName,
		DeviceType:  r.req.DeviceType,
		DeviceModel: r.req.DeviceModel,
		DeviceID:    r.req.DeviceID,
		Latitude:    r.req.LatitudeV,
		Longitude:   r.req.LongitudeV,
		Address:     r.req.Address,
		State:       r.req.State,
		Zip:         r.req.Zip,
		FirstName:   r.req.FirstName,
		LastName:    r.req.LastName,
		Email:       r.req.Email,
		Phone:       r.req.Phone,
		IsAnonymous: r.req.isAnonymous,
		Description: r.req.Description,
	}
}

// CreateReq represents a new Report.  The struct is composed of three anonymous structs:
//
// CreateReqBase - represents a new report.  Contains all fields to unmarshal a
// a request to create a new report.
// cType - defined in common.go, responsible for unmarshaling and parsing the request.
// cIFace - defined in common.go, an interface type for parseQP() and validate() methods.
type CreateReq struct {
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

// convert the unmarshaled data.
func (r *CreateReq) convert() error {
	log.Debug("starting convert()")
	c := newConversion()
	r.LatitudeV = c.float("Latitude", r.Latitude)
	r.LongitudeV = c.float("Longitude", r.Longitude)
	r.isAnonymous = c.bool("IsAnonymous", r.IsAnonymous)
	log.Debug("After convert: %s\n%s", c.String(), r.String())
	if !c.Ok() {
		return c
	}
	return nil
}

// CreateResp is the response to creating or updating a report.
type CreateResp struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// String displays the contents of the CreateReqBase type.
func (r CreateReq) String() string {
	ls := new(common.LogString)
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
