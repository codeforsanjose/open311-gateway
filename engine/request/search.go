package request

import (
	"fmt"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/geo"
	"Gateway311/engine/router"
	"Gateway311/engine/structs"

	"github.com/ant0ine/go-json-rest/rest"
)

const (
	searchRadiusMin  int = 50
	searchRadiusMax  int = 250
	searchRadiusDflt int = 100
)

const (
	srchtUnknown = iota
	srchtReportID
	srchtDeviceID
	stchtLatLng
)

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

	rqst *rest.Request

	req  *SearchRequest
	nreq interface{}

	valid    Validation
	srchType int
	routes   structs.NRoutes

	nresp *structs.NSearchResponse
	resp  *SearchResponse
}

func processSearch(rqst *rest.Request) (fresp interface{}, ferr error) {
	log.Debug("starting processSearch()")
	mgr := searchMgr{
		rqst:  rqst,
		id:    router.GetSID(),
		start: time.Now(),
		req:   &SearchRequest{},
		valid: newValidation(),
		nresp: &structs.NSearchResponse{
			Reports: make([]structs.NSearchResponseReport, 0),
		},
		resp: &SearchResponse{Message: "Request failed"},
	}

	sendTelemetry(mgr.id, "Search", "open")
	defer func() {
		if ferr != nil {
			sendTelemetry(mgr.id, "Search", "error")
		} else {
			sendTelemetry(mgr.id, "Search", "done")
		}
	}()

	fail := func(err error) (interface{}, error) {
		log.Errorf("processSearch failed - %s", err)
		return mgr.resp, fmt.Errorf("Search request failed - %s", err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != "JSON payload is empty" {
			return fail(err)
		}
	}
	log.Debug(mgr.req.String())

	if err := mgr.parseQP(); err != nil {
		log.Errorf("processCreate.parseQP() failed - %s", err)
		return fail(err)
	}
	log.Debug(mgr.req.String())

	if err := mgr.validate(); err != nil {
		log.Errorf("processSearch.validate() failed - %s", err)
		return fail(err)
	}
	log.Debug(mgr.req.String())

	if err := mgr.setRoute(); err != nil {
		log.Errorf("processSearch.route() failed - %s", err)
		return fail(err)
	}

	if err := mgr.callRPC(); err != nil {
		log.Errorf("processSearch.callRPC() failed - %s", err)
		return fail(err)
	}

	mgr.convertResponse()

	return mgr.resp, nil
}

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

func (r *searchMgr) validate() error {
	log.Debug("Starting validate()")
	v := r.valid

	v.Set("convert", "Type conversion of inputs is OK", false)
	v.Set("RID", "Has a Report ID", false)
	v.Set("DID", "Has a Device ID", false)
	v.Set("geo", "Location coordinates are within the continental US", false)
	v.Set("city", "We have a serviced city", false)
	v.Set("route", "Has a viable route", false)

	if err := r.req.convert(); err != nil {
		log.Error(err.Error())
		return err
	}
	v.Set("convert", "", true)

	// ReportID (RID)
	if router.ValidateRID(r.req.RID) {
		v.Set("RID", "", true)
		v.Set("route", "", true)
	}

	// DeviceID
	if len(r.req.DeviceType) > 2 && len(r.req.DeviceID) > 2 {
		v.Set("DID", "", true)
	}

	// Location
	v.Set("geo", "", validateLatLng(r.req.LatitudeV, r.req.LongitudeV))

	// Do we have a valid request?  We must have a ReportID, DeviceID, OR a valid location.
	// If none of those are present, then the request is invalid
	log.Debug("RID: %t", v.IsOK("RID"))
	log.Debug("geo: %t", v.IsOK("geo"))
	log.Debug("DID: %t", v.IsOK("DID"))

	if !(v.IsOK("RID") || v.IsOK("geo") || v.IsOK("DID")) {
		return fmt.Errorf("invalid Search request")
	}

	// Is the Request routable?
	if !v.IsOK("route") && v.IsOK("geo") {
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
		v.Set("route", "", true)
	}
	log.Debug(v.String())
	if !v.IsOK("route") {
		return fmt.Errorf("unable to determine the route for the request")
	}

	// Range-check the search radius.
	switch {
	case r.req.RadiusV < searchRadiusMin:
		r.req.RadiusV = searchRadiusMin
	case r.req.RadiusV > searchRadiusMax:
		r.req.RadiusV = searchRadiusMax
	}

	if err := r.setSearchType(); err != nil {
		return err
	}
	return nil
}

// setRoute gets the route(s) to process the request.
func (r *searchMgr) setRoute() error {
	return nil
}

// callRPC runs the calls to the Adapter(s).
func (r *searchMgr) callRPC() error {
	rpcCall, err := r.prepRPC()
	if err != nil {
		return err
	}

	if err := rpcCall.Run(); err != nil {
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

func (r *searchMgr) convertRequest() {
}

func (r *searchMgr) setSearchType() error {
	v := r.valid

	switch {
	case v.IsOK("RID") && v.IsOK("route"):
		r.srchType = srchtReportID
		r.nreq = r.setRID()
		return nil

	case v.IsOK("DID") && v.IsOK("route"):
		r.srchType = srchtDeviceID
		r.nreq = r.setDID()
		return nil

	case v.IsOK("geo") && v.IsOK("route"):
		r.srchType = stchtLatLng
		r.nreq = r.setLL()
		return nil

	default:
		r.srchType = srchtUnknown
		return fmt.Errorf("invalid query parameters for Search request")
	}
}

func (r *searchMgr) prepRPC() (*router.RPCCall, error) {
	setRPC := func(rpcName string) (*router.RPCCall, error) {
		rpcCall, err := router.NewRPCCall(rpcName, r.nreq, r.adapterReply)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return rpcCall, nil
	}

	switch r.srchType {
	case srchtReportID:
		return setRPC("Report.SearchRID")

	case srchtDeviceID:
		return setRPC("Report.SearchDID")

	case stchtLatLng:
		return setRPC("Report.SearchLL")

	default:
		return nil, fmt.Errorf("cannot prep RPC call - unknown search type: %d", r.srchType)
	}
}

// adapterReply processes the reply returned from the RPC call, by placing a
// pointer to the response in CreateReq.response.
func (r *searchMgr) adapterReply(ndata interface{}) error {
	reply, ok := ndata.(*structs.NSearchResponse)
	log.Debug("reply: %p  ok: %t  size: %v", reply, ok, len(reply.Reports))
	if !ok {
		return fmt.Errorf("invalid interface received: %T", ndata)
	}
	log.Debug("r.nresp.Reports: %T", r.nresp.Reports)
	r.nresp.Reports = append(r.nresp.Reports, reply.Reports...)
	return nil
}

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
				RqstID: r.req.id,
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
				RqstID: r.req.id,
			},
			Rtype: structs.NRTSearchRID,
		},
		RID:    r.req.RID,
		AreaID: r.req.AreaID,
	}
}

// ---------------------------- SearchRequest --------------------------------------------

// SearchRequest represents the Search request (Normal form).
type SearchRequest struct {
	cType
	cIface
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
	srchType    int
	response    struct {
		cRType
		*structs.NSearchResponse
	}
}

// convert the unmarshaled data.
func (r *SearchRequest) convert() error {
	log.Debug("starting convert()")
	c := newConversion()
	r.LatitudeV = c.float("Latitude", r.Latitude)
	r.LongitudeV = c.float("Longitude", r.Longitude)
	r.RadiusV = c.int("Radius", r.Radius)
	r.MaxResultsV = c.int("MaxResults", r.MaxResults)
	log.Debug("After convert: %s\n%s", c.String(), r.String())
	if !c.Ok() {
		return c
	}
	return nil
}

// String displays the contents of the SearchRequest custom type.
func (r SearchRequest) String() string {
	ls := new(common.LogString)
	ls.AddF("SearchRequest - %d\n", r.id)
	ls.AddF("RID: %s\n", r.RID)
	ls.AddF("Device Type: %q    ID: %q\n", r.DeviceType, r.DeviceID)
	ls.AddF("Lat: %v (%f)  Lng: %v (%f)\n", r.Latitude, r.LatitudeV, r.Longitude, r.LongitudeV)
	ls.AddF("Radius: %v (%d) AreaID: %q\n", r.Radius, r.RadiusV, r.AreaID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	return ls.Box(80)
}

// SearchResponse contains the search results.
type SearchResponse struct {
	Message      string                 `json:"Message" xml:"Message"`
	ReportCount  int                    `json:"ReportCount" xml:"ReportCount"`
	ResponseTime string                 `json:"ResponseTime" xml:"ResponseTime"`
	Reports      []SearchResponseReport `json:"Reports,omitempty" xml:"Reports,omitempty"`
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
