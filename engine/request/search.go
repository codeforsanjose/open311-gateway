package request

import (
	"fmt"
	"strconv"

	"Gateway311/engine/common"
	"Gateway311/engine/geo"
	"Gateway311/engine/router"
	"Gateway311/engine/structs"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/davecgh/go-spew/spew"
)

const (
	searchRadiusMin  int = 50
	searchRadiusMax  int = 250
	searchRadiusDflt int = 100
)

// =======================================================================================
//                                      Request
// =======================================================================================

func processSearch(rqst *rest.Request, rqstID int64) (interface{}, error) {
	sreq := SearchRequest{}
	if err := sreq.init(rqst, rqstID); err != nil {
		log.Errorf("processSearch failed - %s", err)
		log.Errorf("SearchRequest: %s", spew.Sdump(sreq))
		return nil, err
	}
	log.Debug("After init:\n%s\n", sreq)
	return sreq.run()
}

const (
	srchtUnknown = iota
	srchtReportID
	srchtDeviceID
	stchtLatLng
)

// SearchRequest represents the Search request (Normal form).
type SearchRequest struct {
	cType
	cIface
	RID         string  `json:"reportID" xml:"reportID"`
	DeviceType  string  `json:"deviceType" xml:"deviceType"`
	DeviceID    string  `json:"deviceId" xml:"deviceId"`
	Latitude    string  `json:"latitude" xml:"latitude"`
	LatitudeV   float64 //
	Longitude   string  `json:"longitude" xml:"longitude"`
	LongitudeV  float64 //
	Radius      string  `json:"radius" xml:"radius"`
	RadiusV     int     // in meters
	Address     string  `json:"address" xml:"address"`
	City        string  `json:"city" xml:"city"`
	AreaID      string  //
	State       string  `json:"state" xml:"state"`
	Zip         string  `json:"zip" xml:"zip"`
	MaxResults  string  `json:"MaxResultsV" xml:"MaxResultsV"`
	MaxResultsV int     //
	validRID    bool
	validDID    bool
	validLatLng bool
	validRoute  bool
	srchType    int
	response    struct {
		cRType
		*structs.NSearchResponse
	}
}

func (r *SearchRequest) setSearchType() error {
	switch {
	case r.validRID && r.validRoute:
		r.srchType = srchtReportID
	case r.validDID && r.validRoute:
		r.srchType = srchtDeviceID
	case r.validLatLng && r.validRoute:
		r.srchType = stchtLatLng
	default:
		r.srchType = srchtUnknown
		return fmt.Errorf("invalid query parameters for Search request")
	}
	return nil
}

func (r *SearchRequest) validate() error {
	// ReportID (RID)
	if len(r.RID) >= 10 {
		r.validRID = true
		r.validRoute = true
	}

	// Lat/Lng
	if x, err := strconv.ParseFloat(r.Latitude, 64); err == nil {
		r.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(r.Longitude, 64); err == nil {
		r.LongitudeV = x
	}
	if r.LatitudeV > 24.0 && r.LongitudeV >= -180.0 && r.LongitudeV <= -66.0 {
		log.Debug("r.validLatLng = true")
		r.validLatLng = true
	}

	// DeviceID
	if len(r.DeviceType) > 2 && len(r.DeviceID) > 2 {
		r.validDID = true
	}

	// Do we have a valid request?
	if !r.validRID && !r.validLatLng && !r.validDID {
		return fmt.Errorf("invalid Search request")
	}

	// Is the Request routable?
	if !r.validRoute && r.validLatLng {
		if city, err := geo.CityForLatLng(r.LatitudeV, r.LongitudeV); err == nil {
			log.Debug("City: %q", city)
			r.City = city
			r.validRoute = true
		}
		var err error
		r.AreaID, err = router.GetAreaID(r.City)
		if err != nil {
			return fmt.Errorf("the city: %q is not serviced by this gateway", r.City)
		}
		log.Debug("%s", r)
	}
	if !r.validRoute {
		return fmt.Errorf("unable to determine the route(s) for the request")
	}

	if x, err := strconv.ParseInt(r.Radius, 10, 64); err == nil {
		switch {
		case int(x) < searchRadiusMin:
			r.RadiusV = searchRadiusMin
		case int(x) > searchRadiusMax:
			r.RadiusV = searchRadiusMax
		default:
			r.RadiusV = int(x)
		}
	}

	if x, err := strconv.ParseInt(r.MaxResults, 0, 64); err == nil {
		r.MaxResultsV = int(x)
	}

	if err := r.setSearchType(); err != nil {
		return err
	}
	return nil
}

func (r *SearchRequest) parseQP(rqst *rest.Request) error {
	r.RID = rqst.URL.Query().Get("rid")
	r.DeviceType = rqst.URL.Query().Get("dtype")
	r.DeviceID = rqst.URL.Query().Get("did")
	r.Latitude = rqst.URL.Query().Get("lat")
	r.Longitude = rqst.URL.Query().Get("lng")
	r.Radius = rqst.URL.Query().Get("radius")
	return nil
}

func (r *SearchRequest) init(rqst *rest.Request, rqstID int64) error {
	if err := r.load(r, rqstID, rqst); err != nil {
		return err
	}
	r.response.NSearchResponse = &structs.NSearchResponse{
		Reports: make([]structs.NSearchResponseReport, 0),
	}
	return nil
}

func (r *SearchRequest) run() (interface{}, error) {
	var rpcCall *router.RPCCall
	switch r.srchType {
	case srchtUnknown:
		return nil, fmt.Errorf("unknown search type - invalid query parms")
	case srchtReportID:
		return nil, fmt.Errorf("Search by ReportID not implemented")
	case srchtDeviceID:
		var err error
		rpcCall, err = r.setSearchDID()
		if err != nil {
			return nil, err
		}
	case stchtLatLng:
		var err error
		rpcCall, err = r.setSearchLL()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown search type: %v", r.srchType)
	}

	if err := rpcCall.Run(); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	r.response.NSearchResponse.ReportCount = len(r.response.NSearchResponse.Reports)
	if r.response.NSearchResponse.ReportCount > 0 {
		r.response.NSearchResponse.Message = "OK"
	} else {
		r.response.NSearchResponse.Message = "No reports found"
	}
	return r.response, nil
}

// adapterReply processes the reply returned from the RPC call, by placing a
// pointer to the response in CreateReq.response.
func (r *SearchRequest) adapterReply(ndata interface{}) error {
	reply, ok := ndata.(*structs.NSearchResponse)
	log.Debug("reply: %p  ok: %t", reply, ok)
	if !ok {
		return fmt.Errorf("invalid interface received: %T", ndata)
	}
	r.response.NSearchResponse.Reports = append(r.response.NSearchResponse.Reports, reply.Reports...)
	r.response.id = r.id
	return nil
}

// ---------------------------- DeviceID --------------------------------------------

func (r *SearchRequest) setSearchDID() (*router.RPCCall, error) {
	rqst, err := r.newNSearchDID()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	rpcCall, err := router.NewRPCCall("Report.SearchDID", &rqst, r.adapterReply)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return rpcCall, nil
}

func (r *SearchRequest) newNSearchDID() (structs.NSearchRequestDID, error) {
	n := structs.NSearchRequestDID{
		NRequestCommon: structs.NRequestCommon{
			ID: structs.NID{
				RqstID: r.id,
			},
			Rtype: structs.NRTSearchDID,
		},
		DeviceType: r.DeviceType,
		DeviceID:   r.DeviceID,
		AreaID:     r.AreaID,
		MaxResults: r.MaxResultsV,
	}
	return n, nil

}

// ---------------------------- Lat/Lng --------------------------------------------

func (r *SearchRequest) setSearchLL() (*router.RPCCall, error) {
	rqst, err := r.newNSearchLL()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	rpcCall, err := router.NewRPCCall("Report.SearchLL", &rqst, r.adapterReply)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return rpcCall, nil
}

func (r *SearchRequest) newNSearchLL() (structs.NSearchRequestLL, error) {
	n := structs.NSearchRequestLL{
		NRequestCommon: structs.NRequestCommon{
			ID: structs.NID{
				RqstID: r.id,
			},
			Rtype: structs.NRTSearchLL,
		},
		Latitude:   r.LatitudeV,
		Longitude:  r.LongitudeV,
		Radius:     r.RadiusV,
		AreaID:     r.AreaID,
		MaxResults: r.MaxResultsV,
	}
	return n, nil

}

// ---------------------------- Strings --------------------------------------------

// String displays the contents of the SearchRequest custom type.
func (r SearchRequest) String() string {
	ls := new(common.LogString)
	ls.AddF("SearchRequest - %d\n", r.id)
	ls.AddF("Device Type: %q    ID: %q\n", r.DeviceType, r.DeviceID)
	ls.AddF("Lat: %v (%f)  Lng: %v (%f)\n", r.Latitude, r.LatitudeV, r.Longitude, r.LongitudeV)
	ls.AddF("Radius: %v (%d) AreaID: %q\n", r.Radius, r.RadiusV, r.AreaID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	ls.AddF("Valid - RID: %t  DID: %t  LatLng: %t  Route: %t", r.validRID, r.validDID, r.validLatLng, r.validRoute)
	return ls.Box(80)
}
