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
	cType               //
	cIface              //
	Latitude    string  `json:"Latitude" xml:"Latitude"`
	LatitudeV   float64 //
	Longitude   string  `json:"Longitude" xml:"Longitude"`
	LongitudeV  float64 //
	validLatLng bool    //
	City        string  `json:"city" xml:"city"`
	validCity   bool    //

	bkend string //

}

func (r *ServicesReq) validate() error {
	if x, err := strconv.ParseFloat(r.Latitude, 64); err == nil {
		r.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(r.Longitude, 64); err == nil {
		r.LongitudeV = x
	}
	return nil
}

func (r *ServicesReq) parseQP(rqst *rest.Request) error {
	r.Latitude = rqst.URL.Query().Get("lat")
	r.Longitude = rqst.URL.Query().Get("lng")
	r.City = rqst.URL.Query().Get("city")
	return nil
}

func (r *ServicesReq) init(rqst *rest.Request) error {
	r.load(r, rqst)
	return nil
}

func (r *ServicesReq) run() (interface{}, error) {
	var err error
	fail := func(err string) (*ServicesResp, error) {
		response := ServicesResp{Message: fmt.Sprintf("Failed - %s", err)}
		return &response, fmt.Errorf("%s", err)
	}

	switch {
	case r.LatitudeV > 24.0 && r.LongitudeV >= -180.0 && r.LongitudeV <= -66.0:
		r.City, err = geo.CityForLatLng(r.LatitudeV, r.LongitudeV)
		if err != nil {
			return fail(fmt.Sprintf("Cannot find city for %v:%v - %s", r.Latitude, r.Longitude, err.Error()))
		}
		fallthrough

	case len(r.City) > 2:
		areaID, err := router.GetAreaID(r.City)
		if err != nil {
			return fail(fmt.Sprintf("Cannot find services for %v - %s", r.City, err.Error()))
		}
		services, err := services.GetArea(areaID)
		if err != nil {
			return fail(fmt.Sprintf("Cannot find services for %v - %s", r.City, err.Error()))
		}
		r, err := newServiceResp("OK", services)
		return &r, nil
	}
	return nil, fmt.Errorf("Invalid location - lat: %v lng: %v  city: %v", r.Latitude, r.Longitude, r.City)
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
func (r ServicesReq) String() string {
	ls := new(common.LogString)
	ls.AddS("Services\n")
	ls.AddF("Location - lat: %v  lon: %v  city: %v\n", r.LatitudeV, r.LongitudeV, r.City)
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
