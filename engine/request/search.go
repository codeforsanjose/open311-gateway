package request

import (
	"errors"
	"fmt"
	"time"

	"github.com/open311-gateway/engine/common"
	"github.com/open311-gateway/engine/geo"
	"github.com/open311-gateway/engine/router"
	"github.com/open311-gateway/engine/structs"
	"github.com/open311-gateway/engine/telemetry"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/jeffizhungry/logrus"
)

var (
	searchRadiusMin int
	searchRadiusMax int
)

// =======================================================================================
//                                      SEARCH MANAGER
// =======================================================================================

// searchMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to initiate a Search.
//  1. Loads all input payload and query parms.
//  2. Validates all input.
//  3. Determines the route(s).  Returns error if no valid route(s) is found.
//  4. Converts the input to the Normal form.
//  5. Call RPC Router to process the request.
//  6. Validates and merges Normal form RPC response(s).
//  7. Converts Normal form to response.
//  8. Returns response.
type searchMgr struct {
	id    int64
	start time.Time

	reqType structs.NRequestType
	rqst    *rest.Request
	req     *SearchRequest
	nreq    interface{}

	valid common.Validation

	routes structs.NRoutes
	rpc    *router.RPCCallMgr

	nresp *structs.NSearchResponse
	resp  SearchResponse
}

func processSearch(rqst *rest.Request) (fresp interface{}, ferr error) {
	mgr := searchMgr{
		rqst:  rqst,
		id:    common.RequestID(),
		start: time.Now(),
		req:   &SearchRequest{},
		valid: common.NewValidation(),
		nresp: &structs.NSearchResponse{
			Reports: make([]structs.NSearchResponseReport, 0),
		},
	}

	telemetry.SendTelemetry(mgr.id, "Search", "open")
	defer func() {
		if ferr != nil {
			telemetry.SendTelemetry(mgr.id, "Search", "error")
		} else {
			telemetry.SendTelemetry(mgr.id, "Search", "done")
		}
	}()

	fail := func(err error) (interface{}, error) {
		log.Warn("processSearch failed - " + err.Error())
		return mgr.resp, fmt.Errorf("Search request failed - %s", err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != greEmpty {
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		log.Warn("processSearch.validate() failed - " + err.Error())
		return fail(err)
	}

	log.Debug("Before RPC Call:\n%s", mgr.String())
	if err := mgr.callRPC(); err != nil {
		log.Error("processSearch.callRPC() failed - " + err.Error())
		return fail(err)
	}

	mgr.convertResponse()

	return mgr.resp, nil
}

// -------------------------------------------------------------------------------
//                        ROUTER.REQUESTER INTERFACE
// -------------------------------------------------------------------------------
func (r *searchMgr) RType() structs.NRequestType {
	return r.reqType
}

func (r *searchMgr) Routes() structs.NRoutes {
	return r.routes
}

func (r *searchMgr) Data() interface{} {
	return r.nreq
}

func (r *searchMgr) Processer() func(ndata interface{}) error {
	return r.processReply
}

// -------------------------------------------------------------------------------
//                        VALIDATION
// -------------------------------------------------------------------------------

// validate converts and verifies all input parameters.  It calls:
//    setRoute() - determines if there are viable Adapter routes to process the search.
func (r *searchMgr) validate() error {
	fail := func(msg string, err error) error {
		if err != nil {
			msg = msg + " - " + err.Error()
		}
		log.Warn("Validation failed - " + msg)
		return errors.New(msg)
	}

	v := r.valid
	v.Set("qryParms", "Query parms parsed and loaded ok", false)
	v.Set("inputs", "Type conversion of inputs is OK", false)
	v.Set("RID", "Has a Report ID", false)
	v.Set("DID", "Has a Device ID", false)
	v.Set("geo", "Location coordinates are within the continental US", false)
	v.Set("city", "We have a serviced city", false)
	v.Set("route", "Has a viable route", false)

	// Load Query Parms.
	if err := r.parseQP(); err != nil {
		return fail("", err)
	}
	v.Set("qryParms", "", true)

	// Convert all string inputs.
	if err := r.req.convert(); err != nil {
		return fail("", err)
	}
	v.Set("inputs", "", true)

	// ReportID (RID)
	if router.ValidateRID(r.req.RID) {
		v.Set("RID", "", true)
	}

	// DeviceID
	if len(r.req.DeviceType) > 2 && len(r.req.DeviceID) > 2 {
		v.Set("DID", "", true)
	}

	// Location
	v.Set("geo", "", common.ValidateLatLng(r.req.LatitudeV, r.req.LongitudeV))

	// Range-check the search radius.
	log.Debugf("Search radius min/max: %v-%v", searchRadiusMin, searchRadiusMax)
	switch {
	case r.req.RadiusV < searchRadiusMin:
		r.req.RadiusV = searchRadiusMin
	case r.req.RadiusV > searchRadiusMax:
		r.req.RadiusV = searchRadiusMax
	}

	// Do we have a valid request?  We must have a ReportID, DeviceID, OR a valid location.
	// If none of those are present, then the request is invalid
	if !(v.IsOK("RID") || v.IsOK("geo") || v.IsOK("DID")) {
		return fail("invalid Search request", nil)
	}

	// Is the Request routable?
	if err := r.setRoute(); err != nil {
		return fail("", err)
	}
	log.Debug("After setRoute() - " + v.String())

	if err := r.setSearchType(); err != nil {
		return fail("", err)
	}
	return nil
}

// parseQP parses the query parameters, and loads them into the searchMgr.req struct.
func (r *searchMgr) parseQP() error {
	rid, _, err := structs.RIDFromString(r.rqst.URL.Query().Get("rid"))
	if err == nil {
		r.req.RID = rid
	}
	r.req.DeviceType = r.rqst.URL.Query().Get("dtype")
	r.req.DeviceID = r.rqst.URL.Query().Get("did")
	r.req.Latitude = r.rqst.URL.Query().Get("lat")
	r.req.Longitude = r.rqst.URL.Query().Get("lng")
	r.req.Radius = r.rqst.URL.Query().Get("radius")
	return nil
}

// setRoute gets the route(s) to process the request.
// One of the following, in order, MUST be present to determine a route.
//   1. If a RID is present, it is used.
//   2. If we have a valid Lat/Lng ("geo"), then we use it to get a city.
func (r *searchMgr) setRoute() error {
	v := r.valid

	switch {
	case v.IsOK("RID"):
		r.routes = structs.NRoutes{r.req.RID.NRoute}
		v.Set("route", "", true)
		return nil

	case v.IsOK("geo"):
		if city, err := geo.CityForLatLng(r.req.LatitudeV, r.req.LongitudeV); err == nil {
			log.Debug("City: " + city)
			r.req.City = city
		}
		var err error
		r.req.AreaID, err = router.GetAreaID(r.req.City)
		if err != nil {
			return fmt.Errorf("the city: %q is not serviced by this gateway", r.req.City)
		}
		v.Set("city", "", true)
		routes, err := router.GetAreaRoutes(r.req.AreaID)
		if err != nil {
			return err
		}
		r.routes = routes
		v.Set("route", "", true)
		return nil

	default:
		return fmt.Errorf("can't find a route")
	}
}

func (r *searchMgr) setSearchType() error {
	v := r.valid

	switch {
	case v.IsOK("RID") && v.IsOK("route"):
		r.reqType = structs.NRTSearchRID
		r.nreq = r.setRID()
	case v.IsOK("DID") && v.IsOK("route"):
		r.reqType = structs.NRTSearchDID
		r.nreq = r.setDID()
	case v.IsOK("geo") && v.IsOK("route"):
		r.reqType = structs.NRTSearchLL
		r.nreq = r.setLL()
	default:
		r.reqType = structs.NRTUnknown
		return fmt.Errorf("invalid query parameters for Search request")
	}
	r.nreq.(structs.NRequester).SetID(r.id, 0)
	return nil
}

// -------------------------------------------------------------------------------
//                        RPC
// -------------------------------------------------------------------------------

// callRPC runs the calls to the Adapter(s).
func (r *searchMgr) callRPC() (err error) {
	r.rpc, err = router.NewRPCCallMgr(r)
	if err != nil {
		return err
	}

	if err = r.rpc.Run(); err != nil {
		log.Error(err.Error())
		return err
	}
	r.nresp.ReportCount = len(r.nresp.Reports)
	if r.nresp.ReportCount > 0 {
		r.nresp.Message = "OK"
	} else {
		r.nresp.Message = "No reports found"
	}
	return nil
}

func (r *searchMgr) processReply(ndata interface{}) error {
	reply, ok := ndata.(*structs.NSearchResponse)
	log.WithFields(log.Fields{
		"reply": reply,
		"ok":    ok,
		"size":  len(reply.Reports),
	}).Debug("processReply...")
	if !ok {
		return fmt.Errorf("wrong type of data: %T returned by RPC call", ndata)
	}
	r.nresp.Reports = append(r.nresp.Reports, reply.Reports...)
	return nil
}

// -------------------------------------------------------------------------------
//                        RESPONSE
// -------------------------------------------------------------------------------

func (r *searchMgr) convertResponse() {
	r.resp = SearchResponse{}

	for _, rpt := range r.nresp.Reports {
		fullAddress := rpt.FullAddress()
		newRsp := SearchResponseReport{
			RID: rpt.RID,

			Status:            &rpt.StatusType,
			StatusNotes:       &rpt.TicketSLA,
			ServiceName:       &rpt.RequestType,
			ServiceCode:       &rpt.RequestTypeID,
			Description:       &rpt.Description,
			AgencyResponsible: nil,
			ServiceNotice:     nil,
			RequestedAt:       &rpt.DateCreated,
			UpdatedAt:         &rpt.DateUpdated,
			ExpectedAt:        nil,
			Address:           &fullAddress,
			AddressID:         nil,
			ZipCode:           &rpt.ZipCode,
			Latitude:          &rpt.Latitude,
			Longitude:         &rpt.Longitude,
			MediaURL:          &rpt.MediaURL,
		}
		newRsp.emptyToNil()
		r.resp = append(r.resp, newRsp)
	}
}

// ---------------------------- DeviceID --------------------------------------------
func (r *searchMgr) setDID() *structs.NSearchRequestDID {
	return &structs.NSearchRequestDID{
		NRequestCommon: structs.NRequestCommon{
			ID: structs.NID{
				RqstID: r.id,
			},
			Rtype: structs.NRTSearchDID,
		},
		DeviceType: r.req.DeviceType,
		DeviceID:   r.req.DeviceID,
		AreaID:     r.req.AreaID,
		MaxResults: r.req.MaxResultsV,
	}
}

// ---------------------------- Lat/Lng --------------------------------------------
func (r *searchMgr) setLL() *structs.NSearchRequestLL {
	return &structs.NSearchRequestLL{
		NRequestCommon: structs.NRequestCommon{
			ID: structs.NID{
				RqstID: r.id,
			},
			Rtype: structs.NRTSearchLL,
		},
		Latitude:   r.req.LatitudeV,
		Longitude:  r.req.LongitudeV,
		Radius:     r.req.RadiusV,
		AreaID:     r.req.AreaID,
		MaxResults: r.req.MaxResultsV,
	}
}

// ---------------------------- ReportID --------------------------------------------
func (r *searchMgr) setRID() *structs.NSearchRequestRID {
	return &structs.NSearchRequestRID{
		NRequestCommon: structs.NRequestCommon{
			ID: structs.NID{
				RqstID: r.id,
			},
			Rtype: structs.NRTSearchRID,
		},
		RID:    r.req.RID,
		AreaID: r.req.AreaID,
	}
}

// String displays the contents of the SearchRequest custom type.
func (r searchMgr) String() string {
	ls := new(common.LogString)
	ls.AddF("searchMgr - %d\n", r.id)
	ls.AddF("Request type: %v\n", r.reqType.String())
	ls.AddS(r.routes.String())
	ls.AddS(r.req.String())
	if r.rpc != nil {
		ls.AddS(r.rpc.String())
	} else {
		ls.AddS("*****RPC uninitialized*****\n")
	}
	if s, ok := r.nreq.(fmt.Stringer); ok {
		ls.AddS(s.String())
	}
	ls.AddS(r.valid.String())
	if r.routes != nil {
		ls.AddS(r.routes.String())
	}
	if r.nresp != nil {
		ls.AddS(r.nresp.String())
	}
	if r.resp != nil {
		ls.AddS(r.resp.String())
	}
	return ls.Box(120) + "\n\n"
}

// =======================================================================================
//                                      SEARCH REQUEST
// =======================================================================================

// SearchRequest represents the Search request (Normal form).
type SearchRequest struct {
	RID         structs.ReportID `json:"reportID" xml:"reportID"`
	DeviceType  string           `json:"deviceType" xml:"deviceType"`
	DeviceID    string           `json:"deviceId" xml:"deviceId"`
	Latitude    string           `json:"lat" xml:"lat"`
	LatitudeV   float64          //
	Longitude   string           `json:"lng" xml:"lng"`
	LongitudeV  float64          //
	Radius      string           `json:"radius" xml:"radius"`
	RadiusV     int              // in meters
	Address     string           `json:"address" xml:"address"`
	City        string           `json:"city" xml:"city"`
	AreaID      string           //
	State       string           `json:"state" xml:"state"`
	Zip         string           `json:"zip" xml:"zip"`
	MaxResults  string           `json:"MaxResultsV" xml:"MaxResultsV"`
	MaxResultsV int              //
}

// convert the unmarshaled data.
func (r *SearchRequest) convert() error {
	c := common.NewConversion()
	r.LatitudeV = c.Float("Latitude", r.Latitude)
	r.LongitudeV = c.Float("Longitude", r.Longitude)
	r.RadiusV = c.Int("Radius", r.Radius)
	r.MaxResultsV = c.Int("MaxResults", r.MaxResults)
	if !c.Ok() {
		return c
	}
	return nil
}

// String displays the contents of the SearchRequest custom type.
func (r SearchRequest) String() string {
	ls := new(common.LogString)
	ls.AddF("SearchRequest\n")
	ls.AddF("RID: %s\n", r.RID)
	ls.AddF("Device Type: %q    ID: %q\n", r.DeviceType, r.DeviceID)
	ls.AddF("Lat: %v (%f)  Lng: %v (%f)\n", r.Latitude, r.LatitudeV, r.Longitude, r.LongitudeV)
	ls.AddF("Radius: %v (%d) AreaID: %q\n", r.Radius, r.RadiusV, r.AreaID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	return ls.Box(80)
}

// =======================================================================================
//                                      SEARCH RESPONSE
// =======================================================================================

// SearchResponse contains the search results.
type SearchResponse []SearchResponseReport

// Displays the SearchResponse custom type.
func (r SearchResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("SearchResponse\n")
	for _, x := range r {
		ls.AddS(x.String())
	}
	return ls.Box(90)
}

// SearchResponseReport represents a report.
type SearchResponseReport struct {
	RID               structs.ReportID `xml:"service_request_id" json:"service_request_id"`
	Status            *string          `json:"status" xml:"status"`
	StatusNotes       *string          `json:"status_notes,omitempty" xml:"status_notes,omitempty"`
	ServiceName       *string          `json:"service_name" xml:"service_name"`
	ServiceCode       *string          `json:"service_code" xml:"service_code"`
	Description       *string          `json:"description" xml:"description"`
	AgencyResponsible *string          `json:"agency_responsible,omitempty" xml:"agency_responsible,omitempty"`
	ServiceNotice     *string          `json:"service_notice,omitempty" xml:"service_notice,omitempty"`
	RequestedAt       *string          `json:"requested_datetime" xml:"requested_datetime"`
	UpdatedAt         *string          `json:"updated_datetime" xml:"updated_datetime"`
	ExpectedAt        *string          `json:"expected_datetime,omitempty" xml:"expected_datetime,omitempty"`
	Address           *string          `json:"address" xml:"address"`
	AddressID         *string          `json:"address_id" xml:"address_id"`
	ZipCode           *string          `json:"zipcode" xml:"zipcode"`
	Latitude          *string          `json:"lat" xml:"lat"`
	Longitude         *string          `json:"lng" xml:"lng"`
	MediaURL          *string          `json:"media_url" xml:"media_url"`
}

func (r *SearchResponseReport) emptyToNil() {
	if r.Status != nil && *r.Status == "" {
		r.Status = nil
	}
	if r.StatusNotes != nil && *r.StatusNotes == "" {
		r.StatusNotes = nil
	}
	if r.ServiceName != nil && *r.ServiceName == "" {
		r.ServiceName = nil
	}
	if r.ServiceCode != nil && *r.ServiceCode == "" {
		r.ServiceCode = nil
	}
	if r.Description != nil && *r.Description == "" {
		r.Description = nil
	}
	if r.AgencyResponsible != nil && *r.AgencyResponsible == "" {
		r.AgencyResponsible = nil
	}
	if r.ServiceNotice != nil && *r.ServiceNotice == "" {
		r.ServiceNotice = nil
	}
	if r.RequestedAt != nil && *r.RequestedAt == "" {
		r.RequestedAt = nil
	}
	if r.UpdatedAt != nil && *r.UpdatedAt == "" {
		r.UpdatedAt = nil
	}
	if r.ExpectedAt != nil && *r.ExpectedAt == "" {
		r.ExpectedAt = nil
	}
	if r.Address != nil && *r.Address == "" {
		r.Address = nil
	}
	if r.AddressID != nil && *r.AddressID == "" {
		r.AddressID = nil
	}
	if r.ZipCode != nil && *r.ZipCode == "" {
		r.ZipCode = nil
	}
	if r.Latitude != nil && *r.Latitude == "" {
		r.Latitude = nil
	}
	if r.Longitude != nil && *r.Longitude == "" {
		r.Longitude = nil
	}
	if r.MediaURL != nil && *r.MediaURL == "" {
		r.MediaURL = nil
	}

	return
}

// Displays the the SearchResponseReport custom type.
func (r SearchResponseReport) String() string {
	ls := new(common.LogString)
	ls.AddF("ServiceRequestID: %s\n", r.RID.RID())
	ls.AddF("Service - id: %q   name: %q\n", r.ServiceCode, r.ServiceName)
	ls.AddF("Created: %v   Updated: %v   Expected: %v\n", r.RequestedAt, r.UpdatedAt, r.ExpectedAt)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Agency: %s\n", r.AgencyResponsible)
	ls.AddF("Location - lat: %v  lon: %v \n", r.Latitude, r.Longitude)
	ls.AddF("          %s  zip: %s\n", r.Address, r.ZipCode)
	ls.AddF("MediaURL: %s\n", r.MediaURL)
	return ls.Box(80)
}
