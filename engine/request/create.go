package request

import (
	"errors"
	"fmt"
	"math"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/router"
	"Gateway311/engine/structs"
	"Gateway311/engine/telemetry"

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
	id    int64
	start time.Time

	reqType structs.NRequestType
	rqst    *rest.Request
	req     *CreateRequest
	nreq    *structs.NCreateRequest

	valid common.Validation

	routes structs.NRoutes
	rpc    *router.RPCCallMgr

	nresp *structs.NCreateResponse
	resp  *CreateResponse
}

func processCreate(rqst *rest.Request) (fresp interface{}, ferr error) {
	mgr := createMgr{
		id:      common.RequestID(),
		start:   time.Now(),
		reqType: structs.NRTCreate,
		rqst:    rqst,
		req:     &CreateRequest{},
		valid:   common.NewValidation(),
		resp:    &CreateResponse{Message: "Request failed"},
	}
	telemetry.SendTelemetry(mgr.id, "Create", "open")
	defer func() {
		if ferr != nil {
			telemetry.SendTelemetry(mgr.id, "Create", "error")
		} else {
			telemetry.SendTelemetry(mgr.id, "Create", "done")
		}
	}()

	fail := func(err error) (interface{}, error) {
		log.Errorf("processCreate failed - %s", err)
		return mgr.resp, fmt.Errorf("Create request failed - %s", err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != greEmpty {
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		log.Errorf("processCreate.validate() failed - %s", err)
		return fail(err)
	}

	mgr.convertRequest()

	if err := mgr.callRPC(); err != nil {
		log.Errorf("processCreate.callRPC() failed - %s", err)
		return fail(err)
	}

	mgr.convertResponse()

	return mgr.resp, nil
}

// -------------------------------------------------------------------------------
//                        ROUTER.REQUESTER INTERFACE
// -------------------------------------------------------------------------------
func (r *createMgr) RType() structs.NRequestType {
	return r.reqType
}

func (r *createMgr) Routes() structs.NRoutes {
	return r.routes
}

func (r *createMgr) Data() interface{} {
	return r.nreq
}

func (r *createMgr) Processer() func(ndata interface{}) error {
	return r.processReply
}

// -------------------------------------------------------------------------------
//                        VALIDATION
// -------------------------------------------------------------------------------
// validate the unmarshaled data.
func (r *createMgr) validate() error {
	log.Debug("Starting validate()")
	fail := func(msg string, err error) error {
		if err != nil {
			msg = msg + " - " + err.Error()
		}
		log.Errorf("Validation failed: %s", msg)
		return errors.New(msg)
	}

	v := r.valid
	v.Set("QP", "Query parms parsed and loaded ok", false)
	v.Set("convert", "Type conversion of inputs is OK", false)
	v.Set("geo", "Location coordinates are within the continental US", false)

	// Load Query Parms.
	if err := r.parseQP(); err != nil {
		return fail("", err)
	}
	v.Set("QP", "", true)

	// Convert all string inputs.
	if err := r.req.convert(); err != nil {
		return fail("", err)
	}
	v.Set("convert", "", true)

	// Location
	r.valid.Set("geo", "", common.ValidateLatLng(r.req.LatitudeV, r.req.LongitudeV))

	// Is the Request routable?
	if err := r.setRoute(); err != nil {
		return fail("", err)
	}
	log.Debug("After setRoute() - %s", v.String())

	log.Debug(r.valid.String())
	if !r.valid.Ok() {
		return r.valid
	}
	return nil
}

// parseQP unloads any query parms in the request.
func (r *createMgr) parseQP() error {
	return nil
}

// setRoute gets the route(s) to process the request.
func (r *createMgr) setRoute() error {
	routes, err := router.RoutesMID(r.req.MID)
	log.Debug("Routes: %v", routes.String())
	if err != nil {
		return fmt.Errorf("no routes found - %s", err.Error())
	}
	r.routes = routes
	return nil
}

// -------------------------------------------------------------------------------
//                        RPC
// -------------------------------------------------------------------------------

// callRPC runs the calls to the Adapter(s).
func (r *createMgr) callRPC() (err error) {
	r.rpc, err = router.NewRPCCallMgr(r)
	if err != nil {
		return err
	}

	log.Debug("Before RPC\n%s", r.String())
	if err = r.rpc.Run(); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (r *createMgr) processReply(ndata interface{}) error {
	r.nresp = ndata.(*structs.NCreateResponse)
	return nil
}

// ------------------------------ String -------------------------------------------------

// String displays the contents of the SearchRequest custom type.
func (r createMgr) String() string {
	ls := new(common.LogString)
	ls.AddF("searchMgr - %d\n", r.id)
	ls.AddF("Request type: %v\n", r.reqType.String())
	ls.AddS(r.req.String())
	if r.routes != nil {
		ls.AddS(r.routes.String())
	}
	if r.rpc != nil {
		ls.AddS(r.rpc.String())
	} else {
		ls.AddS("*****RPC uninitialized*****\n")
	}
	ls.AddS(r.nreq.String())
	ls.AddS(r.valid.String())
	if r.nresp != nil {
		ls.AddS(r.nresp.String())
	}
	if r.resp != nil {
		ls.AddS(r.resp.String())
	}
	return ls.Box(120) + "\n\n"
}

// -------------------------------------------------------------------------------
//                        REQUEST
// -------------------------------------------------------------------------------

// CreateRequest represents a new Report
type CreateRequest struct {
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
}

// convert the unmarshaled data.
func (r *CreateRequest) convert() error {
	log.Debug("starting convert()")
	c := common.NewConversion()
	r.LatitudeV = c.Float("Latitude", r.Latitude)
	r.LongitudeV = c.Float("Longitude", r.Longitude)
	r.isAnonymous = c.Bool("IsAnonymous", r.IsAnonymous)
	log.Debug("After convert: %s\n%s", c.String(), r.String())
	if !c.Ok() {
		return c
	}
	return nil
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

// String displays the contents of the CreateRequest type.
func (r CreateRequest) String() string {
	ls := new(common.LogString)
	ls.AddF("CreateRequest\n")
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

// -------------------------------------------------------------------------------
//                        RESPONSE
// -------------------------------------------------------------------------------

// CreateResponse is the response to creating or updating a report.
type CreateResponse struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

// convertResponse converts the NCreateResponse{} to a CreateResponse{}
func (r *createMgr) convertResponse() {
	r.resp = &CreateResponse{
		Message:  r.nresp.Message,
		ID:       r.nresp.RID.RID(),
		AuthorID: r.nresp.AuthorID,
	}
}

// String displays the contents of the CreateRequest type.
func (r CreateResponse) String() string {
	ls := new(common.LogString)
	ls.AddF("CreateResponse - %d\n", r.ID)
	ls.AddF("Message: %s\n", r.Message)
	ls.AddF("AuthorID: %s\n", r.AuthorID)
	return ls.Box(80)
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================
