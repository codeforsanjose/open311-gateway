package request

import (
	"Gateway311/gateway/common"
	"Gateway311/gateway/geo"
	"Gateway311/gateway/structs"
	"_sketches/spew"

	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
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
	op := SearchReq{}
	if err := op.init(r); err != nil {
		return nil, err
	}
	return op.run()
}

type SearchReq struct {
	cType  //
	cIface //
	// JID    int    `json:"jid" xml:"jid"`
	bkend string //
	structs.SearchReqBase
}

func (c *SearchReq) validate() {
	if x, err := strconv.ParseFloat(c.Latitude, 64); err == nil {
		c.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(c.Longitude, 64); err == nil {
		c.LongitudeV = x
	}
	if x, err := strconv.ParseInt(c.Radius, 10, 64); err == nil {
		switch {
		case int(x) < searchRadiusMin:
			c.RadiusV = searchRadiusMin
		case int(x) > searchRadiusMax:
			c.RadiusV = searchRadiusMax
		default:
			c.RadiusV = int(x)
		}
	}
	if x, err := strconv.ParseInt(c.MaxResults, 0, 64); err == nil {
		c.MaxResultsV = int(x)
	}
	return
}

func (c *SearchReq) parseQP(r *rest.Request) error {
	return nil
}

func (c *SearchReq) init(r *rest.Request) error {
	c.load(c, r)
	return nil
}

func (c *SearchReq) run() (interface{}, error) {
	city, err := geo.CityForLatLng(c.LatitudeV, c.LongitudeV)
	if err != nil {
		return nil, fmt.Errorf("The lat/lng: %v:%v is not in a city", c.LatitudeV, c.LongitudeV)
	}
	fmt.Printf("[toCSSearchLL] city: %q\n", city)

	switch c.bkend {
	case "CitySourced":
		return c.processCS()
	}
	return nil, fmt.Errorf("Unsupported backend: %q", c.bkend)
}

func (c *SearchReq) processCS() (interface{}, error) {
	log.Printf("[processCS] src: %s", spew.Sdump(c))
	// rqst, _ := c.toCSSearchLL()
	// resp, _ := rqst.Process()
	// ourResp, _ := fromSearchCS(resp)

	// return ourResp, nil

	return nil, nil
}

// Displays the contents of the Spec_Type custom type.
func (c SearchReq) String() string {
	ls := new(common.LogString)
	ls.AddS("Search\n")
	ls.AddF("Bkend: %s\n", c.bkend)
	ls.AddF("Device ID: %s\n", c.DeviceID)
	ls.AddF("Location\n")
	if math.Abs(c.LatitudeV) > 1 {
		ls.AddF("   lat: %v  lon: %v\n", c.LatitudeV, c.LongitudeV)
	}
	if len(c.City) > 1 {
		ls.AddF("   \n", c.Address)
		ls.AddF("   %s, %s   %s\n", c.City, c.State, c.Zip)
	}
	return ls.Box(80)
}

//
// // --------------------------- Integrations ----------------------------------------------
//
// func (c *SearchReq) toCSSearchLL() (*integration.CSSearchLLReq, error) {
//
// 	rqst := integration.CSSearchLLReq{
// 	// APIAuthKey:        sp.Key,
// 	// APIRequestType:    "SearchThreeOneOne",
// 	// APIRequestVersion: sp.APIVersion,
// 	// DeviceType:        c.DeviceType,
// 	// DeviceModel:       c.DeviceModel,
// 	// DeviceID:          c.DeviceID,
// 	// RequestType:       c.Type,
// 	// RequestTypeID:     c.TypeIDV,
// 	// Latitude:          c.LatitudeV,
// 	// Longitude:         c.LongitudeV,
// 	// Description:       c.Description,
// 	// AuthorNameFirst:   c.FirstName,
// 	// AuthorNameLast:    c.LastName,
// 	// AuthorEmail:       c.Email,
// 	// AuthorTelephone:   c.Phone,
// 	// AuthorIsAnonymous: c.isAnonymous,
// 	}
// 	return &rqst, nil
// }
//
// // =======================================================================================
// //                                      Response
// // =======================================================================================
//
// func fromSearchCS(src *integration.CSSearchResp) (*SearchResp, error) {
// 	resp := SearchResp{
// 		Message: src.Message,
// 	}
// 	return &resp, nil
// }
