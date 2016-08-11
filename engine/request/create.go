package request

import (
	"errors"
	"fmt"
	"time"

	"github.com/open311-gateway/engine/common"
	"github.com/open311-gateway/engine/router"
	"github.com/open311-gateway/engine/services"
	"github.com/open311-gateway/engine/structs"
	"github.com/open311-gateway/engine/telemetry"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/jeffizhungry/logrus"
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
		// resp:    &CreateResponse{Message: "Request failed"},
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
		log.Warn("processCreate failed - " + err.Error())
		return mgr.resp, fmt.Errorf("Create request failed - %s", err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != greEmpty {
			log.Error("Decode failed")
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		log.Warn("processCreate.validate() failed - " + err.Error())
		return fail(err)
	}

	mgr.convertRequest()

	if err := mgr.callRPC(); err != nil {
		log.Warn("processCreate.callRPC() failed - " + err.Error())
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

/* ----- NOTES ------

Route Validation

1. If there is a full address string:
  a. Parse it.
  b. Get lat/lng.
  c. Overwrite input lat/lng.
2. If there is a lat/lng:
  a. Validate it is in the US.
  b. Overwrite the address




--------------------- */

// validate the unmarshaled data.
func (r *createMgr) validate() error {
	log.Debug("Starting validate()")
	fail := func(msg string, err error) error {
		if err != nil {
			msg = msg + " - " + err.Error()
		}
		log.Warn("Validation failed: " + msg)
		return errors.New(msg)
	}

	v := r.valid
	v.Set("qryParms", "Query parms parsed and loaded ok", false)
	v.Set("inputs", "Type conversion of inputs is OK", false)
	v.Set("SrvID", "The ServiceID is valid", false)
	v.Set("geo", "Location coordinates are within the continental US", false)
	v.Set("city", "We have a city", false)
	v.Set("SrvArea", "The Service ID corresponds to the location", false)

	// Load Query Parms.
	if err := r.parseQP(); err != nil {
		return fail("", err)
	}
	v.Set("qryParms", "", true)

	// Check the ServiceID
	if err := r.req.validateServiceID(); err != nil {
		return fail("", err)
	}

	// Convert all string inputs.
	if err := r.req.convert(); err != nil {
		return fail("", err)
	}
	v.Set("inputs", "", true)

	// Is it anonymous:
	if r.req.Email == "" && r.req.FirstName == "" && r.req.LastName == "" {
		r.req.isAnonymous = true
	}

	// A request currently MUST have a location.
	if err := r.req.validateLocation(); err != nil {
		return fail("", err)
	}
	v.Set("city", "", true)

	// Verify the ServiceID (MID) matches the location.
	if err := r.req.validateLocationMID(); err == nil {
		v.Set("SrvArea", "", true)
	}
	v.Set("SrvID", "", true)

	// Location - the AreaID in the MID must match the location specified by the address
	// or the lat/lng.
	r.valid.Set("geo", "", common.ValidateLatLng(r.req.LatitudeV, r.req.LongitudeV))
	log.Debug("After ValidateLatLng: " + v.String() + r.req.String())

	// Is the Request routable?
	if err := r.setRoute(); err != nil {
		return fail("", err)
	}

	log.Debug(r.valid.String())

	if !r.valid.Ok() {
		return r.valid
	}
	return nil
}

// parseQP unloads any query parms in the request.
func (r *createMgr) parseQP() error {

	for key, values := range r.rqst.URL.Query() {
		for i, value := range values {
			if i > 0 {
				return fmt.Errorf("Invalid query parms")
			}
			switch key {
			case "lat":
				r.req.Latitude = value
			case "lng":
				r.req.Longitude = value
			case "address_string":
				r.req.FullAddress = value

			case "email":
				r.req.Email = value
			case "device_id":
				r.req.DeviceID = value
			case "account_id":
				r.req.AccountID = value
			case "first_name":
				r.req.FirstName = value
			case "last_name":
				r.req.LastName = value
			case "phone":
				r.req.Phone = value
			case "description":
				r.req.Description = value
			case "media_url":
				r.req.MediaURL = value

			case "addr":
				r.req.Address = value
			case "city":
				r.req.City = value
			case "state":
				r.req.State = value
			case "zip":
				r.req.Zip = value
			}
		}
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
		Latitude:    r.req.LatitudeV,
		Longitude:   r.req.LongitudeV,
		FullAddress: r.req.FullAddress,
		Email:       r.req.Email,
		DeviceID:    r.req.DeviceID,
		FirstName:   r.req.FirstName,
		LastName:    r.req.LastName,
		Phone:       r.req.Phone,
		Description: r.req.Description,
		MediaURL:    r.req.MediaURL,

		DeviceType:  r.req.DeviceType,
		DeviceModel: r.req.DeviceModel,
		Address:     r.req.Address,
		Area:        r.req.City,
		State:       r.req.State,
		Zip:         r.req.Zip,
		IsAnonymous: r.req.isAnonymous,
	}
}

// convertResponse converts the NCreateResponse{} to a CreateResponse{}
func (r *createMgr) convertResponse() {
	rid := r.nresp.RID.RID()
	r.resp = &CreateResponse{
		ID:        &rid,
		Notice:    &r.nresp.Message,
		AccountID: &r.nresp.AccountID,
	}
	r.resp.emptyToNil()
}

// setRoute gets the route(s) to process the request.
func (r *createMgr) setRoute() error {
	routes, err := router.RoutesMID(r.req.MID)
	log.Debug("Routes: " + routes.String())
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

	log.Debug("Before RPC\n" + r.String())
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
	MID            structs.ServiceID `json:"service_code" xml:"service_code"`
	APIKey         string            `json:"api_key" xml:"api_key"`
	JurisdictionID string            `json:"jurisdiction_id" xml:"jurisdiction_id"`
	Latitude       string            `json:"lat" xml:"lat"`
	Longitude      string            `json:"lng" xml:"lng"`
	FullAddress    string            `json:"address_string" xml:"address_string"`
	Email          string            `json:"email" xml:"email"`
	DeviceID       string            `json:"device_id" xml:"device_id"`
	AccountID      string            `json:"account_id" xml:"account_id"`
	FirstName      string            `json:"first_name" xml:"first_name"`
	LastName       string            `json:"last_name" xml:"last_name"`
	Phone          string            `json:"phone" xml:"phone"`
	Description    string            `json:"description" xml:"description"`
	MediaURL       string            `json:"media_url" xml:"media_url"`

	LatitudeV  float64 //
	LongitudeV float64 //
	AreaID     string  //

	Address string `json:"address" xml:"address"`
	City    string `json:"city" xml:"city"`
	State   string `json:"state" xml:"state"`
	Zip     string `json:"zip" xml:"zip"`

	isAnonymous bool //

	DeviceType  string `json:"device_type" xml:"device_type"`
	DeviceModel string `json:"device_model" xml:"device_model"`
}

// convert the unmarshaled data.
func (r *CreateRequest) convert() error {
	c := common.NewConversion()
	r.LatitudeV = c.Float("Latitude", r.Latitude)
	r.LongitudeV = c.Float("Longitude", r.Longitude)
	return nil
}

func (r *CreateRequest) validateServiceID() (err error) {
	if !services.ValidateServiceID(r.MID) {
		return fmt.Errorf("The requested ServiceID: %s is invalid", r.MID.MID())
	}
	return nil
}

// validateLocation does the following:
// 1. If there is a non-blank full address, attempt to use it by calling validateAddress()
// 2. If validateAddress() is successful, set the Lat/Long to the address' location and return.
// 3. If validateAddress() fails, then try to find the location using the LongitudeV and LatitudeV.
// 4. If LongitudeV and LatitudeV are invalid, return error.
// 5. If the lodation can be found using LongitudeV and LatitudeV, then set the address and return.
func (r *CreateRequest) validateLocation() (err error) {
	var addr common.Address
	success := func() error {

		r.LatitudeV = addr.Lat
		r.LongitudeV = addr.Long
		r.FullAddress = addr.FullAddr()
		r.Address = addr.Addr
		r.City = addr.City
		r.State = addr.State
		r.Zip = addr.Zip
		log.Debug("validateLocation SUCCESS - " + r.String() + "\n" + addr.String())
		return nil
	}
	fail := func(e string) error {
		log.Debug("validateLocation FAIL: " + e + "\n" + r.String())
		return fmt.Errorf(e)
	}

	// Try the FullAddress first.
	if len(r.FullAddress) > 0 {
		log.Debug("Trying FullAddress...")
		addr, err = common.NewAddr(r.FullAddress, true)
		if err == nil {
			return success()
		}
	}

	// Try the address parts next.
	log.Debug("Trying AddressParts...")
	addr, err = common.NewAddrP(r.Address, r.City, r.State, r.Zip, true)
	if err == nil {
		return success()
	}

	// Finally, try reversing the Lat/Long to an address.
	log.Debug("Validate LatLong...")
	if !common.ValidateLatLng(r.LatitudeV, r.LongitudeV) {
		return fail("unable to determine the request location")
	}

	log.Debug("Getting Address for Lat/Long...")
	addr, err = common.AddrForLatLng(r.LatitudeV, r.LongitudeV)
	if err == nil {
		return success()
	}

	return nil
}

func (r *CreateRequest) validateLocationMID() error {
	if err := r.setAreaID(); err != nil {
		return err
	}

	if r.AreaID != r.MID.AreaID {
		log.WithFields(log.Fields{
			"areaID_address":   r.AreaID,
			"areaID_ServiceID": r.MID.AreaID,
		}).Warn("Address area ID does not match the ServiceID area ID ")
		return fmt.Errorf("the ServiceID: %s is outside the specified location: %s", r.MID.MID(), r.City)
	}

	return nil
}

func (r *CreateRequest) setAreaID() error {
	areaID, err := router.GetAreaID(r.City)
	if err != nil {
		return err
	}
	r.AreaID = areaID
	return nil
}

// String displays the contents of the CreateRequest type.
func (r CreateRequest) String() string {
	ls := new(common.LogString)
	ls.AddF("CreateRequest\n")
	ls.AddF("Request - id: %q   JurisdictionID: %q\n", r.MID.MID(), r.JurisdictionID)
	ls.AddF("Device - ID: %s  type %s  model: %s\n", r.DeviceID, r.DeviceType, r.DeviceModel)
	ls.AddF("Location - lat: %v(%q)  lon: %v(%q)  AreaID: %q\n", r.LatitudeV, r.Latitude, r.LongitudeV, r.Longitude, r.AreaID)
	ls.AddF("          %s\n", r.FullAddress)
	ls.AddF("          %s, %s   %s\n", r.City, r.State, r.Zip)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Author (anon: %t) %s %s  Email: %s  Phone: %s  AcctID: %s\n", r.isAnonymous, r.FirstName, r.LastName, r.Email, r.Phone, r.AccountID)
	ls.AddF("MediaURL: %s\n", r.MediaURL)
	return ls.Box(80)
}

// -------------------------------------------------------------------------------
//                        RESPONSE
// -------------------------------------------------------------------------------

// CreateResponse is the response to creating or updating a report.
type CreateResponse struct {
	ID        *string `json:"service_request_id" xml:"service_request_id"`
	Notice    *string `json:"service_notice" xml:"service_notice"`
	AccountID *string `json:"account_id" xml:"account_id"`
}

func (r *CreateResponse) emptyToNil() {
	if r.ID != nil && *r.ID == "" {
		r.ID = nil
	}
	if r.Notice != nil && *r.Notice == "" {
		r.Notice = nil
	}
	if r.AccountID != nil && *r.AccountID == "" {
		r.AccountID = nil
	}
}

// String displays the contents of the CreateRequest type.
func (r CreateResponse) String() string {
	ls := new(common.LogString)
	ls.AddF("CreateResponse - %d\n", r.ID)
	ls.AddF("AccountID: %s\n", r.AccountID)
	return ls.Box(80)
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================
