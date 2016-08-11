package request

import (
	"time"

	"github.com/open311-gateway/adapters/citysourced/common"
	"github.com/open311-gateway/adapters/citysourced/data"
	"github.com/open311-gateway/adapters/citysourced/search"
	"github.com/open311-gateway/adapters/citysourced/structs"
	"github.com/open311-gateway/adapters/citysourced/telemetry"
)

const (
	dfltMaxResults     int = 20
	dfltIncludeDetails     = true
	dfltDateRangeStart     = ""
	dfltDateRangeEnd       = ""
)

// ================================================================================================
//                                      SEARCH LL
// ================================================================================================

// SearchLL fully processes a "Search by Location" request.
func (r *Report) SearchLL(rqst *structs.NSearchRequestLL, resp *structs.NSearchResponse) error {
	log.Debugf("SearchLL - request: %p  resp: %p\n", rqst, resp)
	// Make the Search Manager
	cm := &searchLLMgr{
		nreq:  rqst,
		nresp: resp,
	}
	log.Debugf("searchLLMgr: %#v\n", *cm)
	log.Debug(cm.nreq.String())

	return runRequest(processer(cm))
}

// searchLLMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Search for reports by Location.
//  1. Validates and converts the request from the Normal form to the CitySourced native XML form.
//  2. Calls the appropriate CitySourced REST interface with proper credentials.
//  3. Converts the CitySourced reply back to Normal form.
//  4. Returns the Normal Response, and any errors.
type searchLLMgr struct {
	nreq  *structs.NSearchRequestLL
	req   *search.RequestLL
	url   string
	resp  *search.Response
	nresp *structs.NSearchResponse
}

func (c *searchLLMgr) convertRequest() error {
	provider, err := data.RouteProvider(c.nreq.Route)
	if err != nil {
		return err
	}
	c.url = provider.URL
	c.req = &search.RequestLL{
		APIAuthKey:        provider.Key,
		APIRequestType:    "GetReportsByLatLng",
		APIRequestVersion: provider.APIVersion,
		Latitude:          c.nreq.Latitude,
		Longitude:         c.nreq.Longitude,
		Radius:            c.nreq.Radius,
		MaxResults:        dfltMaxResults,
		IncludeDetails:    dfltIncludeDetails,
		DateRangeStart:    dfltDateRangeStart,
		DateRangeEnd:      dfltDateRangeEnd,
	}
	telemetry.SendRPC(c.nreq.GetIDS(), "open", "", c.url, 0, time.Now())
	return nil
}

// Process executes the request to search for reports by location.
func (c *searchLLMgr) process() error {
	resp, err := c.req.Process(c.url)
	c.resp = resp
	return err
}

func (c *searchLLMgr) convertResponse() (resultCount int, err error) {
	log.Debugf("Resp: %s", c.nresp)
	route := c.nreq.GetRoute()
	c.nresp.SetIDF(c.nreq.GetID)
	c.nresp.SetRoute(route)
	c.nresp.Message = c.resp.Message
	c.nresp.ResponseTime = c.resp.ResponseTime
	c.nresp.Reports = make([]structs.NSearchResponseReport, 0)

	for _, rr := range c.resp.Reports.Reports {
		c.nresp.Reports = append(c.nresp.Reports, structs.NSearchResponseReport{
			RID:               structs.NewRID(route, rr.ID),
			DateCreated:       rr.DateCreated,
			DateUpdated:       rr.DateUpdated,
			DeviceType:        rr.DeviceType,
			DeviceModel:       rr.DeviceModel,
			DeviceID:          rr.DeviceID,
			RequestType:       rr.RequestType,
			RequestTypeID:     rr.RequestTypeID,
			MediaURL:          rr.ImageURL,
			City:              rr.City,
			State:             rr.State,
			ZipCode:           rr.ZipCode,
			Latitude:          rr.Latitude,
			Longitude:         rr.Longitude,
			Directionality:    rr.Directionality,
			Description:       rr.Description,
			AuthorNameFirst:   rr.AuthorNameFirst,
			AuthorNameLast:    rr.AuthorNameLast,
			AuthorEmail:       rr.AuthorEmail,
			AuthorTelephone:   rr.AuthorTelephone,
			AuthorIsAnonymous: rr.AuthorIsAnonymous,
			URLDetail:         rr.URLDetail,
			URLShortened:      rr.URLShortened,
			Votes:             rr.Votes,
			StatusType:        rr.StatusType,
			TicketSLA:         rr.TicketSLA,
		})
	}
	return len(c.nresp.Reports), nil
}

func (c *searchLLMgr) fail(err error) error {
	c.nresp.Message = "Failed - " + err.Error()
	return err
}

func (c *searchLLMgr) getIDS() string {
	return c.nreq.GetIDS()
}

func (c *searchLLMgr) getRoute() string {
	return c.nreq.GetRoute().String()
}

func (c *searchLLMgr) String() string {
	ls := new(common.LogString)
	ls.AddS("SearchLL\n")
	ls.AddS(c.nreq.String())
	ls.AddS(c.req.String())
	ls.AddS(c.resp.String())
	ls.AddS(c.nresp.String())
	return ls.Box(90)
}

// ================================================================================================
//                                      SEARCH RID
// ================================================================================================

// SearchRID fully processes the Search by DeviceID request.
func (r *Report) SearchRID(rqst *structs.NSearchRequestRID, resp *structs.NSearchResponse) error {
	log.Debugf("Search - request: %p  resp: %p\n", rqst, resp)
	// Make the Search Manager
	cm := &searchRIDMgr{
		nreq:  rqst,
		nresp: resp,
	}
	log.Debugf("searchRIDMgr: %#v\n", *cm)

	return runRequest(processer(cm))
}

// searchRIDMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Search for reports by Location.
//  1. Validates and converts the request from the Normal form to the CitySourced native XML form.
//  2. Calls the appropriate CitySourced REST interface with proper credentials.
//  3. Converts the CitySourced reply back to Normal form.
//  4. Returns the Normal Response, and any errors.
type searchRIDMgr struct {
	nreq  *structs.NSearchRequestRID
	req   *search.RequestRID
	url   string
	resp  *search.Response
	nresp *structs.NSearchResponse
}

func (c *searchRIDMgr) convertRequest() error {
	provider, err := data.RouteProvider(c.nreq.Route)
	if err != nil {
		return err
	}
	c.url = provider.URL
	c.req = &search.RequestRID{
		APIAuthKey:        provider.Key,
		APIRequestType:    "GetReport",
		APIRequestVersion: provider.APIVersion,
		ReportID:          c.nreq.RID.ID,
		MaxResults:        dfltMaxResults,
		IncludeDetails:    dfltIncludeDetails,
		DateRangeStart:    dfltDateRangeStart,
		DateRangeEnd:      dfltDateRangeEnd,
	}
	telemetry.SendRPC(c.nreq.GetIDS(), "open", "", c.url, 0, time.Now())
	return nil
}

// Process executes the request to search for reports by location.
func (c *searchRIDMgr) process() error {
	resp, err := c.req.Process(c.url)
	c.resp = resp
	return err
}

func (c *searchRIDMgr) convertResponse() (resultCount int, err error) {
	log.Debugf("Resp: %s", c.nresp)
	route := c.nreq.GetRoute()
	c.nresp.SetIDF(c.nreq.GetID)
	c.nresp.SetRoute(route)
	c.nresp.Message = c.resp.Message
	c.nresp.ResponseTime = c.resp.ResponseTime
	c.nresp.Reports = make([]structs.NSearchResponseReport, 0)

	for _, rr := range c.resp.Reports.Reports {
		c.nresp.Reports = append(c.nresp.Reports, structs.NSearchResponseReport{
			RID:               structs.NewRID(route, rr.ID),
			DateCreated:       rr.DateCreated,
			DateUpdated:       rr.DateUpdated,
			DeviceType:        rr.DeviceType,
			DeviceModel:       rr.DeviceModel,
			DeviceID:          rr.DeviceID,
			RequestType:       rr.RequestType,
			RequestTypeID:     rr.RequestTypeID,
			MediaURL:          rr.ImageURL,
			City:              rr.City,
			State:             rr.State,
			ZipCode:           rr.ZipCode,
			Latitude:          rr.Latitude,
			Longitude:         rr.Longitude,
			Directionality:    rr.Directionality,
			Description:       rr.Description,
			AuthorNameFirst:   rr.AuthorNameFirst,
			AuthorNameLast:    rr.AuthorNameLast,
			AuthorEmail:       rr.AuthorEmail,
			AuthorTelephone:   rr.AuthorTelephone,
			AuthorIsAnonymous: rr.AuthorIsAnonymous,
			URLDetail:         rr.URLDetail,
			URLShortened:      rr.URLShortened,
			Votes:             rr.Votes,
			StatusType:        rr.StatusType,
			TicketSLA:         rr.TicketSLA,
		})
	}
	return len(c.nresp.Reports), nil
}

func (c *searchRIDMgr) fail(err error) error {
	c.nresp.Message = "Failed - " + err.Error()
	return err
}

func (c *searchRIDMgr) getIDS() string {
	return c.nreq.GetIDS()
}

func (c *searchRIDMgr) getRoute() string {
	return c.nreq.GetRoute().String()
}

func (c *searchRIDMgr) String() string {
	ls := new(common.LogString)
	ls.AddS("SearchRID\n")
	ls.AddS(c.nreq.String())
	ls.AddS(c.req.String())
	ls.AddS(c.resp.String())
	ls.AddS(c.nresp.String())
	return ls.Box(90)
}

/// ================================================================================================
//                                      SEARCH DID
// ================================================================================================

// SearchDID fully processes the Search by DeviceID request.
func (r *Report) SearchDID(rqst *structs.NSearchRequestDID, resp *structs.NSearchResponse) error {
	log.Debugf("Search - request: %p  resp: %p\n", rqst, resp)
	// Make the Search Manager
	cm := &searchDIDMgr{
		nreq:  rqst,
		nresp: resp,
	}
	log.Debugf("searchDIDMgr: %#v\n", *cm)

	return runRequest(processer(cm))
}

// searchDIDMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Search for reports by Location.
//  1. Validates and converts the request from the Normal form to the CitySourced native XML form.
//  2. Calls the appropriate CitySourced REST interface with proper credentials.
//  3. Converts the CitySourced reply back to Normal form.
//  4. Returns the Normal Response, and any errors.
type searchDIDMgr struct {
	nreq  *structs.NSearchRequestDID
	req   *search.RequestDID
	url   string
	resp  *search.Response
	nresp *structs.NSearchResponse
}

func (c *searchDIDMgr) convertRequest() error {
	provider, err := data.RouteProvider(c.nreq.Route)
	if err != nil {
		return err
	}
	c.url = provider.URL
	c.req = &search.RequestDID{
		APIAuthKey:        provider.Key,
		APIRequestType:    "GetReportsByDeviceId",
		APIRequestVersion: provider.APIVersion,
		DeviceType:        c.nreq.DeviceType,
		DeviceID:          c.nreq.DeviceID,
		MaxResults:        dfltMaxResults,
		IncludeDetails:    dfltIncludeDetails,
		DateRangeStart:    dfltDateRangeStart,
		DateRangeEnd:      dfltDateRangeEnd,
	}
	telemetry.SendRPC(c.nreq.GetIDS(), "open", "", c.url, 0, time.Now())
	return nil
}

// Process executes the request to search for reports by location.
func (c *searchDIDMgr) process() error {
	resp, err := c.req.Process(c.url)
	c.resp = resp
	return err
}

func (c *searchDIDMgr) convertResponse() (resultCount int, err error) {
	log.Debugf("Resp: %s", c.nresp)
	route := c.nreq.GetRoute()
	c.nresp.SetIDF(c.nreq.GetID)
	c.nresp.SetRoute(route)
	c.nresp.Message = c.resp.Message
	c.nresp.ResponseTime = c.resp.ResponseTime
	c.nresp.Reports = make([]structs.NSearchResponseReport, 0)

	for _, rr := range c.resp.Reports.Reports {
		c.nresp.Reports = append(c.nresp.Reports, structs.NSearchResponseReport{
			RID:               structs.NewRID(route, rr.ID),
			DateCreated:       rr.DateCreated,
			DateUpdated:       rr.DateUpdated,
			DeviceType:        rr.DeviceType,
			DeviceModel:       rr.DeviceModel,
			DeviceID:          rr.DeviceID,
			RequestType:       rr.RequestType,
			RequestTypeID:     rr.RequestTypeID,
			MediaURL:          rr.ImageURL,
			City:              rr.City,
			State:             rr.State,
			ZipCode:           rr.ZipCode,
			Latitude:          rr.Latitude,
			Longitude:         rr.Longitude,
			Directionality:    rr.Directionality,
			Description:       rr.Description,
			AuthorNameFirst:   rr.AuthorNameFirst,
			AuthorNameLast:    rr.AuthorNameLast,
			AuthorEmail:       rr.AuthorEmail,
			AuthorTelephone:   rr.AuthorTelephone,
			AuthorIsAnonymous: rr.AuthorIsAnonymous,
			URLDetail:         rr.URLDetail,
			URLShortened:      rr.URLShortened,
			Votes:             rr.Votes,
			StatusType:        rr.StatusType,
			TicketSLA:         rr.TicketSLA,
		})
	}
	return len(c.nresp.Reports), nil
}

func (c *searchDIDMgr) fail(err error) error {
	c.nresp.Message = "Failed - " + err.Error()
	return err
}

func (c *searchDIDMgr) getIDS() string {
	return c.nreq.GetIDS()
}

func (c *searchDIDMgr) getRoute() string {
	return c.nreq.GetRoute().String()
}

func (c *searchDIDMgr) String() string {
	ls := new(common.LogString)
	ls.AddS("SearchDID\n")
	ls.AddS(c.nreq.String())
	ls.AddS(c.req.String())
	ls.AddS(c.resp.String())
	ls.AddS(c.nresp.String())
	return ls.Box(90)
}

// ================================================================================================
//                                      STRINGS
// ================================================================================================
