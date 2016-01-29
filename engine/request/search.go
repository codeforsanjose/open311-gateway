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
	sreq := SearchReq{}
	if err := sreq.init(r); err != nil {
		log.Errorf("processSearch failed - %s", err)
		log.Errorf("SearchReq: %s", spew.Sdump(sreq))
		return nil, err
	}
	log.Debug("After init:\n%s\n", sreq)
	return sreq.run()
}

// SearchReq represents the Search request (Normal form).
type SearchReq struct {
	cType
	cIface
	bkend string
	structs.SearchReqBase
}

func (c *SearchReq) validate() error {
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

	return nil
}

func (c *SearchReq) parseQP(r *rest.Request) error {
	c.DeviceType = r.URL.Query().Get("dtype")
	c.DeviceID = r.URL.Query().Get("did")
	c.Latitude = r.URL.Query().Get("lat")
	c.Longitude = r.URL.Query().Get("lng")
	c.City = r.URL.Query().Get("city")
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
	c.City = city
	c.AreaID, err = router.GetAreaID(city)
	if err != nil {
		return nil, fmt.Errorf("The city: %q is not serviced by this gateway", c.City)
	}
	log.Debug("%s", c)

	r, err := router.NewRPCCall("Report.SearchLocation", c, c.adapterReply)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = r.Run()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return c.Response, err
}

// adapterReply processes the reply returned from the RPC call, by placing a
// pointer to the response in CreateReq.response.
func (c *SearchReq) adapterReply(ndata interface{}) error {
	c.Response = ndata.(*structs.SearchResp)
	return nil
}

func (c *SearchReq) processCS() (interface{}, error) {
	log.Debug("[processCS] src: %s", spew.Sdump(c))
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
	ls.AddS(c.SearchReqBase.String())
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
