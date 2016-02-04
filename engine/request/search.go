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
	bkend       string
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

func (sr *SearchRequest) validate() error {
	if x, err := strconv.ParseFloat(sr.Latitude, 64); err == nil {
		sr.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(sr.Longitude, 64); err == nil {
		sr.LongitudeV = x
	}
	if x, err := strconv.ParseInt(sr.Radius, 10, 64); err == nil {
		switch {
		case int(x) < searchRadiusMin:
			sr.RadiusV = searchRadiusMin
		case int(x) > searchRadiusMax:
			sr.RadiusV = searchRadiusMax
		default:
			sr.RadiusV = int(x)
		}
	}
	if x, err := strconv.ParseInt(sr.MaxResults, 0, 64); err == nil {
		sr.MaxResultsV = int(x)
	}

	return nil
}

func (sr *SearchRequest) parseQP(r *rest.Request) error {
	sr.DeviceType = r.URL.Query().Get("dtype")
	sr.DeviceID = r.URL.Query().Get("did")
	sr.Latitude = r.URL.Query().Get("lat")
	sr.Longitude = r.URL.Query().Get("lng")
	sr.City = r.URL.Query().Get("city")
	return nil
}

func (sr *SearchRequest) init(r *rest.Request) error {
	sr.load(sr, r)
	return nil
}

func (sr *SearchRequest) run() (interface{}, error) {
	city, err := geo.CityForLatLng(sr.LatitudeV, sr.LongitudeV)
	if err != nil {
		return nil, fmt.Errorf("The lat/lng: %v:%v is not in a city", sr.LatitudeV, sr.LongitudeV)
	}
	sr.City = city
	sr.AreaID, err = router.GetAreaID(city)
	if err != nil {
		return nil, fmt.Errorf("The city: %q is not serviced by this gateway", sr.City)
	}
	log.Debug("%s", sr)

	r, err := router.NewRPCCall("Report.SearchLocation", sr, sr.adapterReply)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = r.Run()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return sr.Response, err
}

// adapterReply processes the reply returned from the RPC call, by placing a
// pointer to the response in CreateReq.response.
func (sr *SearchRequest) adapterReply(ndata interface{}) error {
	sr.Response = ndata.(*structs.NSearchResponse)
	return nil
}

func (sr *SearchRequest) convertN() (interface{}, error) {
	var err error
	switch {
	case sr.LatitudeV > 24.0 && sr.LongitudeV >= -180.0 && sr.LongitudeV <= -66.0:
		sr.City, err = geo.CityForLatLng(sr.LatitudeV, sr.LongitudeV)
		if err != nil {
			return nil, fmt.Errorf("Cannot find city for %v:%v - %s", sr.Latitude, sr.Longitude, err.Error())
		}
		areaID, err := router.GetAreaID(sr.City)
		if err != nil {
			return nil, err
		}
		sr.AreaID = areaID
		return sr.convertLL()

	case len(sr.City) > 2:
		areaID, err := router.GetAreaID(sr.City)
		if err != nil {
			return nil, err
		}
		sr.AreaID = areaID
		return sr.convertLL()

	case len(sr.DeviceType) > 2 && len(sr.DeviceID) > 2:
		return sr.convertDID()

	}
	return nil, fmt.Errorf("Invalid search request.")
}

func (sr *SearchRequest) convertLL() (interface{}, error) {
	return structs.NSearchRequestLL{
		// NSearchRequest: structs.NSearchRequest{
		// 	SearchType: structs.NSTLocation,
		// },
		Latitude:   sr.LatitudeV,
		Longitude:  sr.LongitudeV,
		AreaID:     sr.AreaID,
		Radius:     sr.RadiusV,
		MaxResults: sr.MaxResultsV,
	}, nil
}

func (sr *SearchRequest) convertDID() (interface{}, error) {
	return structs.NSearchRequestDID{
		// NSearchRequest: structs.NSearchRequest{
		// 	SearchType: structs.NSTDeviceID,
		// },
		DeviceType: sr.DeviceType,
		DeviceID:   sr.DeviceID,
		MaxResults: sr.MaxResultsV,
	}, nil
}

// String displays the contents of the SearchRequest custom type.
func (sr SearchRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("Search\n")
	ls.AddF("Bkend: %s\n", sr.bkend)
	ls.AddF("Device Type: %s ID: %s\n", sr.DeviceType, sr.DeviceID)
	ls.AddF("Lat: %v (%f)  Lng: %v (%f)\n", sr.Latitude, sr.LatitudeV, sr.Longitude, sr.LongitudeV)
	ls.AddF("Radius: %v (%d) AreaID: %q\n", sr.Radius, sr.RadiusV, sr.AreaID)
	ls.AddF("MaxResults: %v\n", sr.MaxResults)
	return ls.Box(80)
}
