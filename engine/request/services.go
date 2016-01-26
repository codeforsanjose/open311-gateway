package request

import (
	"fmt"
	"strconv"
	"strings"

	"Gateway311/engine/common"
	"Gateway311/engine/geo"
	"Gateway311/engine/router"
	"Gateway311/engine/services"
	"Gateway311/engine/structs"

	"github.com/ant0ine/go-json-rest/rest"
)

// =======================================================================================
//                                      REQUEST
// =======================================================================================
func processServices(r *rest.Request) (interface{}, error) {
	op := ServicesReq{}
	if err := op.init(r); err != nil {
		return nil, err
	}
	return op.run()
}

// ServicesReq represents a request to .
type ServicesReq struct {
	cType  //
	cIface //
	// JID         int     `json:"jid" xml:"jid"`
	Latitude    string  `json:"Latitude" xml:"Latitude"`
	LatitudeV   float64 //
	Longitude   string  `json:"Longitude" xml:"Longitude"`
	LongitudeV  float64 //
	validLatLng bool    //
	City        string  `json:"city" xml:"city"`
	validCity   bool    //

	bkend string //

}

func (c *ServicesReq) validate() {
	if x, err := strconv.ParseFloat(c.Latitude, 64); err == nil {
		c.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(c.Longitude, 64); err == nil {
		c.LongitudeV = x
	}
	log.Debug("%s\n", c)
	return
}

func (c *ServicesReq) parseQP(r *rest.Request) error {
	c.Latitude = r.URL.Query().Get("lat")
	c.Longitude = r.URL.Query().Get("lng")
	c.City = r.URL.Query().Get("city")
	return nil
}

func (c *ServicesReq) init(r *rest.Request) error {
	c.load(c, r)
	return nil
}

func (c *ServicesReq) run() (interface{}, error) {
	var err error
	fail := func(err string) (*ServicesResp, error) {
		response := ServicesResp{Message: fmt.Sprintf("Failed - %s", err)}
		return &response, fmt.Errorf("%s", err)
	}

	switch {
	case c.LatitudeV > 24.0 && c.LongitudeV >= -180.0 && c.LongitudeV <= -66.0:
		c.City, err = geo.CityForLatLng(c.LatitudeV, c.LongitudeV)
		if err != nil {
			return fail(fmt.Sprintf("Cannot find city for %v:%v - %s", c.Latitude, c.Longitude, err.Error()))
		}
		fallthrough

	case len(c.City) > 2:
		areaID, err := router.GetAreaID(c.City)
		fmt.Printf("AreaID: %q\n", areaID)
		if err != nil {
			return fail(fmt.Sprintf("Cannot find services for %v - %s", c.City, err.Error()))
		}
		services, err := services.GetArea(areaID)
		if err != nil {
			return fail(fmt.Sprintf("Cannot find services for %v - %s", c.City, err.Error()))
		}
		r, err := newServiceResp("OK", services)
		return &r, nil
	}
	return nil, fmt.Errorf("Invalid location - lat: %v lng: %v  city: %v", c.Latitude, c.Longitude, c.City)
}

// =======================================================================================
//                                      Response
// =======================================================================================

// newServiceResp translates structs.NServices to ServicesResp and ServicesRespS.
func newServiceResp(msg string, ns structs.NServices) (ServicesResp, error) {
	newSR := ServicesResp{
		Message:  msg,
		Services: make(map[string]ServicesRespS),
	}

	for _, v := range ns {
		newSR.Services[v.ServiceID.MID()] = ServicesRespS{
			Name:       v.Name,
			Categories: v.Categories,
		}
	}

	return newSR, nil
}

// ServicesResp represents a list of services.
type ServicesResp struct {
	Message  string                   `json:"message" xml:"Message"`
	Services map[string]ServicesRespS `json:"services" xml:"Services"`
}

// ServicesRespS represents a service in a service list.
type ServicesRespS struct {
	Name       string   `json:"name"`
	Categories []string `json:"catg"`
}

// =======================================================================================
//                                      STRINGS
// =======================================================================================

// Displays the contents of the Spec_Type custom type.
func (c ServicesReq) String() string {
	ls := new(common.LogString)
	ls.AddS("Services\n")
	// ls.AddF("JID: %v\n", c.JID)
	ls.AddF("Location - lat: %v  lon: %v  city: %v\n", c.LatitudeV, c.LongitudeV, c.City)
	return ls.Box(80)
}

// Displays the contents of the Spec_Type custom type.
func (r ServicesResp) String() string {
	ls := new(common.LogString)
	ls.AddS("Services Response\n")
	ls.AddF("Message: %v\n", r.Message)
	for k, v := range r.Services {
		ls.AddF("%-18s %-30s [%s]\n", k, v.Name, strings.Join(v.Categories, ", "))
	}

	return ls.Box(80)
}
