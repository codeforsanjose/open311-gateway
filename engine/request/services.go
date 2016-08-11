package request

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/open311-gateway/engine/common"
	"github.com/open311-gateway/engine/geo"
	"github.com/open311-gateway/engine/router"
	"github.com/open311-gateway/engine/services"
	"github.com/open311-gateway/engine/structs"
	"github.com/open311-gateway/engine/telemetry"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/jeffizhungry/logrus"
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

	// routes structs.NRoutes

	resp *ServicesResp
}

func processServices(rqst *rest.Request) (fresp interface{}, ferr error) {
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
		log.Info("processServices failed - " + err.Error())
		return mgr.resp, fmt.Errorf(err.Error())
	}

	if err := mgr.rqst.DecodeJsonPayload(mgr.req); err != nil {
		if err.Error() != greEmpty {
			return fail(err)
		}
	}

	if err := mgr.validate(); err != nil {
		return fail(err)
	}

	if err := mgr.run(); err != nil {
		log.Error("processServices.callRPC() failed - " + err.Error())
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
		return errors.New(msg)
	}

	v := r.valid
	v.Set("qryParms", "Query parms parsed and loaded ok", false)
	v.Set("inputs", "Type conversion of inputs is OK", false)
	v.Set("areaID", "We have a valid AreaID", false)

	// Load Query Parms.
	if err := r.parseQP(); err != nil {
		return fail("", err)
	}
	v.Set("qryParms", "", true)

	// Convert all string inputs.
	if err := r.req.convert(); err != nil {
		return fail("", err)
	}
	v.Set("inputs", "", true)

	// Location
	switch {
	case common.ValidateLatLng(r.req.LatitudeV, r.req.LongitudeV):
		r.req.City, _ = geo.CityForLatLng(r.req.LatitudeV, r.req.LongitudeV)
		fallthrough

	case len(r.req.FullAddress) > 0:
		addr, err := common.ParseAddress(r.req.FullAddress, true)
		log.Debugf("Parsed full address - addr: %+v", addr)
		if err == nil {
			r.req.Address = addr.Addr
			r.req.City = addr.City
			r.req.State = addr.State
			r.req.Zip = addr.Zip
		} else {
			log.Info("ParseAddress failed - " + err.Error())
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
	r.req.Longitude = r.rqst.URL.Query().Get("lng")
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
	log.Debug(r.req.String())
	services, err := services.GetArea(r.req.areaID)
	if err != nil {
		return fmt.Errorf("Cannot find services for %v - %v", r.req.City, err.Error())
	}
	// log.Debugf("***Services:\n%v", services.String())
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
	Longitude   string  `json:"lng" xml:"lng"`
	LongitudeV  float64 //
	FullAddress string  `json:"address_string" xml:"address_string"`
	Address     string  `json:"addr" xml:"addr"`
	City        string  `json:"city" xml:"city"`
	State       string  `json:"state" xml:"state"`
	Zip         string  `json:"zip" xml:"zip"`
	areaID      string
}

func (r *ServicesReq) validate() error {
	log.Debug("Starting validation:\n" + r.String())
	if x, err := strconv.ParseFloat(r.Latitude, 64); err == nil {
		r.LatitudeV = x
	}
	if x, err := strconv.ParseFloat(r.Longitude, 64); err == nil {
		r.LongitudeV = x
	}
	if len(r.FullAddress) > 0 {
		addr, err := common.ParseAddress(r.FullAddress, true)
		if err == nil {
			r.Address = addr.Addr
			r.City = addr.City
			r.State = addr.State
			r.Zip = addr.Zip
		}
	}
	return nil
}

func (r *ServicesReq) convert() error {
	c := common.NewConversion()
	r.LatitudeV = c.Float("Latitude", r.Latitude)
	r.LongitudeV = c.Float("Longitude", r.Longitude)
	log.Debug("After convert: " + c.String() + "\n" + r.String())
	return nil
}

// =======================================================================================
//                                      ServicesResp
// =======================================================================================

// ServicesResp represents a list of services.
type ServicesResp []*ServicesRespS

func newServiceResp(msg string, ns structs.NServices) (*ServicesResp, error) {
	newSR := ServicesResp{}

	for _, v := range ns {
		newSR = append(newSR, newServicesRespS(v))
	}
	// log.Debugf("***ServiceResp:\n%v", newSR)

	return &newSR, nil
}

// ServicesRespS represents a service in a service list.
type ServicesRespS struct {
	ID          string  `json:"service_code" xml:"service_code"`
	Name        *string `json:"service_name" xml:"service_name"`
	Description *string `json:"description" xml:"description"`
	Metadata    bool    `json:"metadata" xml:"metadata"`
	Stype       *string `json:"type" xml:"type"`
	Keywords    *string `json:"keywords" xml:"keywords"`
	Group       *string `json:"group" xml:"group"`
}

func newServicesRespS(s structs.NService) (sr *ServicesRespS) {
	name := s.Name
	description := s.Name
	metadata := false
	stype := s.ResponseType
	keywords := strings.Join(s.Keywords, ",")
	group := s.Group

	sr = &ServicesRespS{
		ID:          s.ServiceID.MID(),
		Name:        &name,
		Description: &description,
		Metadata:    metadata,
		Stype:       &stype,
		Keywords:    &keywords,
		Group:       &group,
	}
	sr.emptyToNil()
	return sr
}

func (r *ServicesRespS) emptyToNil() {
	if r.Name != nil && *r.Name == "" {
		r.Name = nil
	}
	if r.Description != nil && *r.Description == "" {
		r.Description = nil
	}
	if r.Stype != nil && *r.Stype == "" {
		r.Stype = nil
	}
	if r.Keywords != nil && *r.Keywords == "" {
		r.Keywords = nil
	}
	if r.Group != nil && *r.Group == "" {
		r.Group = nil
	}
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
		ls.AddS(v.String())
	}

	return ls.Box(80)
}

// Displays the contents of the ServicesRespS custom type.
func (r ServicesRespS) String() string {
	var name, stype, group, keywords string
	if r.Name != nil {
		name = *r.Name
	}
	if r.Stype != nil {
		stype = *r.Stype
	}
	if r.Group != nil {
		group = *r.Group
	}
	if r.Keywords != nil {
		keywords = *r.Keywords
	}
	return fmt.Sprintf("%-18s %-30s %-10s [%s] %s\n", r.ID, name, stype, group, keywords)
	// ls.AddF("%-18s %T-%v %T-%v %T-%v %T-%v\n", r.ID, r.Name, r.Stype, r.Group, r.Keywords)
}

// StringD displays the contents of the ServicesRespS custom type, including pointer values.
func (r ServicesRespS) StringD() string {
	var name, stype, group, keywords string
	if r.Name != nil {
		name = *r.Name
	}
	if r.Stype != nil {
		stype = *r.Stype
	}
	if r.Group != nil {
		group = *r.Group
	}
	if r.Keywords != nil {
		keywords = *r.Keywords
	}
	return fmt.Sprintf("%-18s (%p)%-30s (%p)%-10s (%p)[%s] (%p)%s\n", r.ID, r.Name, name, r.Stype, stype, r.Group, group, r.Keywords, keywords)
	// ls.AddF("%-18s %T-%v %T-%v %T-%v %T-%v\n", r.ID, r.Name, r.Stype, r.Group, r.Keywords)
}
