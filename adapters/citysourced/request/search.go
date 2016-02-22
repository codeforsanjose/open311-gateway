package request

import (
	"time"

	"Gateway311/adapters/citysourced/data"
	"Gateway311/adapters/citysourced/search"
	"Gateway311/adapters/citysourced/structs"
	"Gateway311/adapters/citysourced/telemetry"
	"Gateway311/engine/common"
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
	log.Debug("Search - request: %p  resp: %p\n", rqst, resp)
	// Make the Search Manager
	cm := &searchLLMgr{
		nreq:  rqst,
		nresp: resp,
	}
	log.Debug("searchLLMgr: %#v\n", *cm)
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
			ImageURL:          rr.ImageURL,
			ImageURLXl:        rr.ImageURLXl,
			ImageURLLg:        rr.ImageURLLg,
			ImageURLMd:        rr.ImageURLMd,
			ImageURLSm:        rr.ImageURLSm,
			ImageURLXs:        rr.ImageURLXs,
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
	return c.nreq.GetRoute().SString()
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

/*
// ================================================================================================
//                                      SEARCH DID
// ================================================================================================

// Search fully processes the Search request.
func (r *Report) SearchDID(rqst *structs.NSearchRequestLL, resp *structs.NSearchResponse) error {
	log.Debug("Search - request: %p  resp: %p\n", rqst, resp)
	// Make the Search Manager
	cm := &searchDIDMgr{
		nreq:  rqst,
		nresp: resp,
	}
	log.Debug("searchDIDMgr: %#v\n", *cm)

	return runRequest(processer(cm))
}

// searchDIDMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Search for reports by Location.
//  1. Validates and converts the request from the Normal form to the CitySourced native XML form.
//  2. Calls the appropriate CitySourced REST interface with proper credentials.
//  3. Converts the CitySourced reply back to Normal form.
//  4. Returns the Normal Response, and any errors.
type searchDIDMgr struct {
	nreq  *structs.NSearchRequestLL
	req   *search.RequestLL
	url   string
	resp  *search.Response
	nresp *structs.NSearchResponse
}

func (c *searchDIDMgr) convertRequest() error {
	provider, err := data.MIDProvider(c.nreq.MID)
	if err != nil {
		return err
	}
	c.url = provider.URL
	c.req = &search.RequestLL{
		APIAuthKey:        provider.Key,
		APIRequestType:    "GetReportsByLatLng",
		APIRequestVersion: provider.APIVersion,
		DeviceType:        c.nreq.DeviceType,
		DeviceModel:       c.nreq.DeviceModel,
		DeviceID:          c.nreq.DeviceID,
		RequestType:       c.nreq.Type,
		RequestTypeID:     c.nreq.MID.ID,
		Latitude:          c.nreq.Latitude,
		Longitude:         c.nreq.Longitude,
		Description:       c.nreq.Description,
		AuthorNameFirst:   c.nreq.FirstName,
		AuthorNameLast:    c.nreq.LastName,
		AuthorEmail:       c.nreq.Email,
		AuthorTelephone:   c.nreq.Phone,
		AuthorIsAnonymous: c.nreq.IsAnonymous,
	}
	return nil
}

// Process executes the request to search for reports by location.
func (c *searchDIDMgr) process() error {
	resp, err := c.req.Process(c.url)
	c.resp = resp
	return err
}

func (c *searchDIDMgr) convertResponse() error {
	c.nresp.Message = c.resp.Message
	c.nresp.ID = c.resp.ID
	c.nresp.AuthorID = c.resp.AuthorID
	return nil
}

func (c *searchDIDMgr) fail(err error) error {
	c.nresp.Message = "Failed - " + err.Error()
	c.nresp.ID = ""
	c.nresp.AuthorID = ""
	return err
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
*/

// ================================================================================================
//                                      STRINGS
// ================================================================================================
