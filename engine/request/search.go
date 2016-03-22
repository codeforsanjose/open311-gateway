package request

import (
	"errors"
	"fmt"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/geo"
	"Gateway311/engine/router"
	"Gateway311/engine/structs"
	"Gateway311/engine/telemetry"

	"github.com/ant0ine/go-json-rest/rest"
)

const (
	searchRadiusMin int = 50
	searchRadiusMax int = 250
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
	resp  *SearchResponse
}

func processSearch(rqst *rest.Request) (fresp interface{}, ferr error) {
	log.Debug("starting processSearch()")
	mgr := searchMgr{
		rqst:  rqst,
		id:    common.RequestID(),
		start: time.Now(),
		req:   &SearchRequest{},
		valid: common.NewValidation(),
		resp:  &SearchResponse{Message: "Request failed"},
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
		log.Errorf("processSearch failed - %s", err)
		return mgr.resp, fmt.Errorf("Search request failed - %s", err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != greEmpty {
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		log.Errorf("processSearch.validate() failed - %s", err)
		return fail(err)
	}

	log.Debug("Before RPC Call:\n%s", mgr.String())
	if err := mgr.callRPC(); err != nil {
		log.Errorf("processSearch.callRPC() failed - %s", err)
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
		log.Errorf("Validation failed: %s", msg)
		return errors.New(msg)
	}

	v := r.valid
	v.Set("QP", "Query parms parsed and loaded ok", false)
	v.Set("convert", "Type conversion of inputs is OK", false)
	v.Set("RID", "Has a Report ID", false)
	v.Set("DID", "Has a Device ID", false)
	v.Set("geo", "Location coordinates are within the continental US", false)
	v.Set("city", "We have a serviced city", false)
	v.Set("route", "Has a viable route", false)

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
	log.Debug("After setRoute() - %s", v.String())

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
			log.Debug("City: %q", city)
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
	log.Debug("reply: %p  ok: %t  size: %v", reply, ok, len(reply.Reports))
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
	var rpts []SearchResponseReport
	for _, rpt := range r.nresp.Reports {
		rpts = append(rpts, SearchResponseReport{
			RID:               rpt.RID,
			DateCreated:       rpt.DateCreated,
			DateUpdated:       rpt.DateUpdated,
			DeviceType:        rpt.DeviceType,
			DeviceModel:       rpt.DeviceModel,
			DeviceID:          rpt.DeviceID,
			RequestType:       rpt.RequestType,
			RequestTypeID:     rpt.RequestTypeID,
			ImageURL:          rpt.ImageURL,
			ImageURLXl:        rpt.ImageURLXl,
			ImageURLLg:        rpt.ImageURLLg,
			ImageURLMd:        rpt.ImageURLMd,
			ImageURLSm:        rpt.ImageURLSm,
			ImageURLXs:        rpt.ImageURLXs,
			City:              rpt.City,
			State:             rpt.State,
			ZipCode:           rpt.ZipCode,
			Latitude:          rpt.Latitude,
			Longitude:         rpt.Longitude,
			Directionality:    rpt.Directionality,
			Description:       rpt.Description,
			AuthorNameFirst:   rpt.AuthorNameFirst,
			AuthorNameLast:    rpt.AuthorNameLast,
			AuthorEmail:       rpt.AuthorEmail,
			AuthorTelephone:   rpt.AuthorTelephone,
			AuthorIsAnonymous: rpt.AuthorIsAnonymous,
			URLDetail:         rpt.URLDetail,
			URLShortened:      rpt.URLShortened,
			Votes:             rpt.Votes,
			StatusType:        rpt.StatusType,
			TicketSLA:         rpt.TicketSLA,
		})
	}
	r.resp = &SearchResponse{
		Message:      r.nresp.Message,
		ReportCount:  r.nresp.ReportCount,
		ResponseTime: r.nresp.ResponseTime,
		Reports:      rpts,
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
	Latitude    string           `json:"latitude" xml:"latitude"`
	LatitudeV   float64          //
	Longitude   string           `json:"longitude" xml:"longitude"`
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
type SearchResponse struct {
	Message      string                 `json:"Message" xml:"Message"`
	ReportCount  int                    `json:"ReportCount" xml:"ReportCount"`
	ResponseTime string                 `json:"ResponseTime" xml:"ResponseTime"`
	Reports      []SearchResponseReport `json:"Reports,omitempty" xml:"Reports,omitempty"`
}

// Displays the SearchResponse custom type.
func (r SearchResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("SearchResponse\n")
	ls.AddF("Count: %v RspTime: %v Message: %v\n", r.ReportCount, r.ResponseTime, r.Message)
	for _, x := range r.Reports {
		ls.AddS(x.String())
	}
	return ls.Box(90)
}

// SearchResponseReport represents a report.
type SearchResponseReport struct {
	RID               structs.ReportID `xml:"ID" json:"ID"`
	DateCreated       string           `json:"DateCreated" xml:"DateCreated"`
	DateUpdated       string           `json:"DateUpdated" xml:"DateUpdated"`
	DeviceType        string           `json:"DeviceType" xml:"DeviceType"`
	DeviceModel       string           `json:"DeviceModel" xml:"DeviceModel"`
	DeviceID          string           `json:"DeviceID" xml:"DeviceID"`
	RequestType       string           `json:"RequestType" xml:"RequestType"`
	RequestTypeID     string           `json:"RequestTypeID" xml:"RequestTypeID"`
	ImageURL          string           `json:"ImageURL" xml:"ImageURL"`
	ImageURLXl        string           `json:"ImageURLXl" xml:"ImageURLXl"`
	ImageURLLg        string           `json:"ImageURLLg" xml:"ImageURLLg"`
	ImageURLMd        string           `json:"ImageURLMd" xml:"ImageURLMd"`
	ImageURLSm        string           `json:"ImageURLSm" xml:"ImageURLSm"`
	ImageURLXs        string           `json:"ImageURLXs" xml:"ImageURLXs"`
	City              string           `json:"City" xml:"City"`
	State             string           `json:"State" xml:"State"`
	ZipCode           string           `json:"ZipCode" xml:"ZipCode"`
	Latitude          string           `json:"Latitude" xml:"Latitude"`
	Longitude         string           `json:"Longitude" xml:"Longitude"`
	Directionality    string           `json:"Directionality" xml:"Directionality"`
	Description       string           `json:"Description" xml:"Description"`
	AuthorNameFirst   string           `json:"AuthorNameFirst" xml:"AuthorNameFirst"`
	AuthorNameLast    string           `json:"AuthorNameLast" xml:"AuthorNameLast"`
	AuthorEmail       string           `json:"AuthorEmail" xml:"AuthorEmail"`
	AuthorTelephone   string           `json:"AuthorTelephone" xml:"AuthorTelephone"`
	AuthorIsAnonymous string           `json:"AuthorIsAnonymous" xml:"AuthorIsAnonymous"`
	URLDetail         string           `json:"URLDetail" xml:"URLDetail"`
	URLShortened      string           `json:"URLShortened" xml:"URLShortened"`
	Votes             string           `json:"Votes" xml:"Votes"`
	StatusType        string           `json:"StatusType" xml:"StatusType"`
	TicketSLA         string           `json:"TicketSLA" xml:"TicketSLA"`
}

// Displays the the SearchResponseReport custom type.
func (r SearchResponseReport) String() string {
	ls := new(common.LogString)
	ls.AddF("SearchResponseReport %d\n", r.RID.RID())
	ls.AddF("DateCreated \"%v\"\n", r.DateCreated)
	ls.AddF("Device - type %s  model: %s  ID: %s\n", r.DeviceType, r.DeviceModel, r.DeviceID)
	ls.AddF("Request - type: %q  id: %q\n", r.RequestType, r.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v  directionality: %q\n", r.Latitude, r.Longitude, r.Directionality)
	ls.AddF("          %s, %s   %s\n", r.City, r.State, r.ZipCode)
	ls.AddF("Votes: %v\n", r.Votes)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Images - std: %s\n", r.ImageURL)
	if len(r.ImageURLXs) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLXs)
	}
	if len(r.ImageURLSm) > 0 {
		ls.AddF("          SM: %s\n", r.ImageURLSm)
	}
	if len(r.ImageURLMd) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLMd)
	}
	if len(r.ImageURLLg) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLLg)
	}
	if len(r.ImageURLXl) > 0 {
		ls.AddF("          XS: %s\n", r.ImageURLXl)
	}
	ls.AddF("Author(anon: %v) %s %s  Email: %s  Tel: %s\n", r.AuthorIsAnonymous, r.AuthorNameFirst, r.AuthorNameLast, r.AuthorEmail, r.AuthorTelephone)
	ls.AddF("SLA: %s\n", r.TicketSLA)
	return ls.Box(80)
}
