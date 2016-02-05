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

func processSearch(r *rest.Request) (interface{}, error) {
	sreq := SearchRequest{}
	if err := sreq.init(r); err != nil {
		log.Errorf("processSearch failed - %s", err)
		log.Errorf("SearchRequest: %s", spew.Sdump(sreq))
		return nil, err
	}
	log.Debug("After init:\n%s\n", sreq)
	return sreq.run()
}

// SearchRequest represents the Search request (Normal form).
type SearchRequest struct {
	cType
	cIface
	DeviceType  string  `json:"deviceType" xml:"deviceType"`
	DeviceID    string  `json:"deviceId" xml:"deviceId"`
	Latitude    string  `json:"LatitudeV" xml:"LatitudeV"`
	LatitudeV   float64 //
	Longitude   string  `json:"LongitudeV" xml:"LongitudeV"`
	LongitudeV  float64 //
	Radius      string  `json:"RadiusV" xml:"RadiusV"`
	RadiusV     int     // in meters
	Address     string  `json:"address" xml:"address"`
	City        string  `json:"city" xml:"city"`
	AreaID      string  //
	State       string  `json:"state" xml:"state"`
	Zip         string  `json:"zip" xml:"zip"`
	MaxResults  string  `json:"MaxResultsV" xml:"MaxResultsV"`
	MaxResultsV int     //
	Response    *structs.NSearchResponse
}

func (r *SearchRequest) newNSearchLL() (structs.NSearchRequestLL, error) {
	n := structs.NSearchRequestLL{
		NRequestCommon: structs.NRequestCommon{
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

func (r *SearchRequest) validate() error {
	if x, err := strconv.ParseFloat(r.Latitude, 64); err == nil {
		r.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(r.Longitude, 64); err == nil {
		r.LongitudeV = x
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

	return nil
}

func (r *SearchRequest) parseQP(rqst *rest.Request) error {
	r.DeviceType = rqst.URL.Query().Get("dtype")
	r.DeviceID = rqst.URL.Query().Get("did")
	r.Latitude = rqst.URL.Query().Get("lat")
	r.Longitude = rqst.URL.Query().Get("lng")
	r.Radius = rqst.URL.Query().Get("radius")
	r.City = rqst.URL.Query().Get("city")
	return nil
}

func (r *SearchRequest) init(rqst *rest.Request) error {
	r.load(r, rqst)
	return nil
}

func (r *SearchRequest) run() (interface{}, error) {
	city, err := geo.CityForLatLng(r.LatitudeV, r.LongitudeV)
	if err != nil {
		return nil, fmt.Errorf("The lat/lng: %v:%v is not in a city", r.LatitudeV, r.LongitudeV)
	}
	r.City = city
	r.AreaID, err = router.GetAreaID(city)
	if err != nil {
		return nil, fmt.Errorf("The city: %q is not serviced by this gateway", r.City)
	}
	log.Debug("%s", r)

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
	err = rpcCall.Run()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return r.Response, err
}

// adapterReply processes the reply returned from the RPC call, by placing a
// pointer to the response in CreateReq.response.
func (r *SearchRequest) adapterReply(ndata interface{}) error {
	r.Response = ndata.(*structs.NSearchResponse)
	return nil
}

func (r *SearchRequest) convertN() (interface{}, error) {
	var err error
	switch {
	case r.LatitudeV > 24.0 && r.LongitudeV >= -180.0 && r.LongitudeV <= -66.0:
		r.City, err = geo.CityForLatLng(r.LatitudeV, r.LongitudeV)
		if err != nil {
			return nil, fmt.Errorf("Cannot find city for %v:%v - %s", r.Latitude, r.Longitude, err.Error())
		}
		areaID, err := router.GetAreaID(r.City)
		if err != nil {
			return nil, err
		}
		r.AreaID = areaID
		return r.convertLL()

	case len(r.City) > 2:
		areaID, err := router.GetAreaID(r.City)
		if err != nil {
			return nil, err
		}
		r.AreaID = areaID
		return r.convertLL()

	case len(r.DeviceType) > 2 && len(r.DeviceID) > 2:
		return r.convertDID()

	}
	return nil, fmt.Errorf("Invalid search request.")
}

func (r *SearchRequest) convertLL() (interface{}, error) {
	return structs.NSearchRequestLL{
		// NSearchRequest: structs.NSearchRequest{
		// 	SearchType: structs.NSTLocation,
		// },
		Latitude:   r.LatitudeV,
		Longitude:  r.LongitudeV,
		AreaID:     r.AreaID,
		Radius:     r.RadiusV,
		MaxResults: r.MaxResultsV,
	}, nil
}

func (r *SearchRequest) convertDID() (interface{}, error) {
	return structs.NSearchRequestDID{
		// NSearchRequest: structs.NSearchRequest{
		// 	SearchType: structs.NSTDeviceID,
		// },
		DeviceType: r.DeviceType,
		DeviceID:   r.DeviceID,
		MaxResults: r.MaxResultsV,
	}, nil
}

// String displays the contents of the SearchRequest custom type.
func (r SearchRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("SearchRequest\n")
	ls.AddF("Device Type: %s ID: %s\n", r.DeviceType, r.DeviceID)
	ls.AddF("Lat: %v (%f)  Lng: %v (%f)\n", r.Latitude, r.LatitudeV, r.Longitude, r.LongitudeV)
	ls.AddF("Radius: %v (%d) AreaID: %q\n", r.Radius, r.RadiusV, r.AreaID)
	ls.AddF("MaxResults: %v\n", r.MaxResults)
	return ls.Box(80)
}
