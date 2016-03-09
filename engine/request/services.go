package request

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/geo"
	"Gateway311/engine/router"
	"Gateway311/engine/services"
	"Gateway311/engine/structs"

	"github.com/ant0ine/go-json-rest/rest"
)

// servicesMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to retrieve a list of Services for an Area.
//  1. Loads all input payload and query parms.
//  2. Validates all input.
//  3. Retrieves the list of services from the Service Cache
//  4. Returns the response.
type serviceMgr struct {
	id    int64
	start time.Time

	reqType structs.NRequestType
	rqst    *rest.Request
	req     *ServicesReq

	valid Validation

	routes structs.NRoutes

	resp *ServicesResp
}

func processServices(rqst *rest.Request) (fresp interface{}, ferr error) {
	log.Debug("starting processServices()")
	mgr := serviceMgr{
		id:      router.GetSID(),
		start:   time.Now(),
		reqType: structs.NRTServicesArea,
		rqst:    rqst,
		req:     &ServicesReq{},
		valid:   newValidation(),
		resp:    &ServicesResp{Message: "Request failed"},
	}
	sendTelemetry(mgr.id, "Services", "open")
	defer func() {
		if ferr != nil {
			sendTelemetry(mgr.id, "Services", "error")
		} else {
			sendTelemetry(mgr.id, "Services", "done")
		}
	}()

	fail := func(err error) (interface{}, error) {
		log.Errorf("processServices failed - %s", err)
		return mgr.resp, fmt.Errorf("Create request failed - %s", err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != "JSON payload is empty" {
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		log.Errorf("processServices.validate() failed - %s", err)
		return fail(err)
	}

	if err := mgr.run(); err != nil {
		log.Errorf("processServices.callRPC() failed - %s", err)
		return fail(err)
	}

	return mgr.resp, nil
}

// -------------------------------------------------------------------------------
//                        VALIDATION
// -------------------------------------------------------------------------------
// validate the unmarshaled data.
func (r *serviceMgr) validate() error {
	log.Debug("Starting validate()")
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
	v.Set("areaID", "We have a valid AreaID", false)

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

	// Location

	switch {
	case validateLatLng(r.req.LatitudeV, r.req.LongitudeV):
		r.req.City, _ = geo.CityForLatLng(r.req.LatitudeV, r.req.LongitudeV)
		fallthrough

	case len(r.req.City) > 2:
		areaID, err := router.GetAreaID(r.req.City)
		if err != nil {
			return fail(fmt.Sprintf("Cannot find services for %v", r.req.City), nil)
		}
		r.req.areaID = areaID
		v.Set("areaID", "", true)
	}

	log.Debug(r.valid.String())
	if !r.valid.Ok() {
		return r.valid
	}
	return nil
}

// parseQP unloads any query parms in the request.
func (r *serviceMgr) parseQP() error {
	r.req.Latitude = r.rqst.URL.Query().Get("lat")
	r.req.Longitude = r.rqst.URL.Query().Get("lng")
	r.req.City = r.rqst.URL.Query().Get("city")
	return nil
}

// -------------------------------------------------------------------------------
//                        RUN
// -------------------------------------------------------------------------------

func (r *serviceMgr) run() error {
	log.Debug("%s", r.req.String())
	services, err := services.GetArea(r.req.areaID)
	if err != nil {
		return fmt.Errorf("Cannot find services for %v - %v", r.req.City, err.Error())
	}
	resp, err := newServiceResp("OK", services)
	if err != nil {
		return err
	}
	r.resp = resp
	return err
}

// =======================================================================================
//                                      ServicesReq
// =======================================================================================

// ServicesReq represents a request to .
type ServicesReq struct {
	// cType              //
	// cIface             //
	Latitude   string  `json:"Latitude" xml:"Latitude"`
	LatitudeV  float64 //
	Longitude  string  `json:"Longitude" xml:"Longitude"`
	LongitudeV float64 //
	City       string  `json:"city" xml:"city"`
	areaID     string
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

func (r *ServicesReq) convert() error {
	c := newConversion()
	r.LatitudeV = c.float("Latitude", r.Latitude)
	r.LongitudeV = c.float("Longitude", r.Longitude)
	log.Debug("After convert: %s\n%s", c.String(), r.String())
	return nil
}

// =======================================================================================
//                                      ServicesResp
// =======================================================================================

// newServiceResp translates structs.NServices to ServicesResp and ServicesRespS.
func newServiceResp(msg string, ns structs.NServices) (*ServicesResp, error) {
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

	return &newSR, nil
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
	ls.AddF("ServicesReq\n")
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
