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
	"Gateway311/engine/telemetry"

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

	valid common.Validation

	routes structs.NRoutes

	resp *ServicesResp
}

func processServices(rqst *rest.Request) (fresp interface{}, ferr error) {
	log.Debug("starting processServices()")
	mgr := serviceMgr{
		id:      common.RequestID(),
		start:   time.Now(),
		reqType: structs.NRTServicesArea,
		rqst:    rqst,
		req:     &ServicesReq{},
		valid:   common.NewValidation(),
	}
	telemetry.SendTelemetry(mgr.id, "Services", "open")
	defer func() {
		if ferr != nil {
			telemetry.SendTelemetry(mgr.id, "Services", "error")
		} else {
			telemetry.SendTelemetry(mgr.id, "Services", "done")
		}
	}()

	fail := func(err error) (interface{}, error) {
		log.Warningf("processServices failed - %s", err)
		return mgr.resp, fmt.Errorf(err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != greEmpty {
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		log.Warningf("processServices.validate() failed - %s", err)
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
		log.Warningf("Validation failed: %s", msg)
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
	case common.ValidateLatLng(r.req.LatitudeV, r.req.LongitudeV):
		r.req.City, _ = geo.CityForLatLng(r.req.LatitudeV, r.req.LongitudeV)
		fallthrough

	case len(r.req.FullAddress) > 0:
		addr, city, state, zip, err := common.ParseAddress(r.req.FullAddress, true)
		log.Debug("Parsed full address - addr: %q  city: %q  state: %q  zip: %q", addr, city, state, zip)
		if err == nil {
			r.req.Address = addr
			r.req.City = city
			r.req.State = state
			r.req.Zip = zip
		} else {
			log.Warningf("ParseAddress failed - %s", err.Error())
		}
	}

	if len(r.req.City) > 2 {
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
	r.req.Longitude = r.rqst.URL.Query().Get("long")
	r.req.FullAddress = r.rqst.URL.Query().Get("address_string")
	r.req.Address = r.rqst.URL.Query().Get("addr")
	r.req.City = r.rqst.URL.Query().Get("city")
	r.req.State = r.rqst.URL.Query().Get("state")
	r.req.Zip = r.rqst.URL.Query().Get("zip")
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
	Latitude    string  `json:"lat" xml:"lat"`
	LatitudeV   float64 //
	Longitude   string  `json:"long" xml:"long"`
	LongitudeV  float64 //
	FullAddress string  `json:"address_string" xml:"address_string"`
	Address     string  `json:"addr" xml:"addr"`
	City        string  `json:"city" xml:"city"`
	State       string  `json:"state" xml:"state"`
	Zip         string  `json:"zip" xml:"zip"`
	areaID      string
}

func (r *ServicesReq) validate() error {
	log.Debug("Starting validation:\n%s", r.String())
	if x, err := strconv.ParseFloat(r.Latitude, 64); err == nil {
		r.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(r.Longitude, 64); err == nil {
		r.LongitudeV = x
	}
	if len(r.FullAddress) > 0 {
		addr, city, state, zip, err := common.ParseAddress(r.FullAddress, true)
		if err == nil {
			r.Address = addr
			r.City = city
			r.State = state
			r.Zip = zip
		}
	}
	return nil
}

func (r *ServicesReq) convert() error {
	c := common.NewConversion()
	r.LatitudeV = c.Float("Latitude", r.Latitude)
	r.LongitudeV = c.Float("Longitude", r.Longitude)
	log.Debug("After convert: %s\n%s", c.String(), r.String())
	return nil
}

// =======================================================================================
//                                      ServicesResp
// =======================================================================================

// ServicesResp represents a list of services.
type ServicesResp []*ServicesRespS

// ServicesRespS represents a service in a service list.
type ServicesRespS struct {
	ID          string `json:"service_code" xml:"service_code"`
	Name        string `json:"service_name" xml:"service_name"`
	Description string `json:"description" xml:"description"`
	Metadata    bool   `json:"metadata" xml:"metadata"`
	Stype       string `json:"type" xml:"type"`
	Keywords    string `json:"keywords" xml:"keywords"`
	Group       string `json:"group" xml:"group"`
}

// newServiceResp translates structs.NServices to ServicesResp and ServicesRespS.
func newServiceResp(msg string, ns structs.NServices) (*ServicesResp, error) {
	newSR := ServicesResp{}

	for _, v := range ns {
		newSR = append(newSR, &ServicesRespS{
			ID:          v.ServiceID.MID(),
			Name:        v.Name,
			Description: v.Name,
			Metadata:    false,
			Stype:       "realtime",
			Keywords:    strings.Join(v.Keywords, ","),
			Group:       v.Group,
		})
	}

	return &newSR, nil
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
	for _, v := range r {
		ls.AddF("%-18s %-30s %-10s [%s] %s\n", v.ID, v.Name, v.Stype, v.Group, v.Keywords)
	}

	return ls.Box(80)
}
