package request

import (
	"Gateway311/adapters/citysourced/create"
	"Gateway311/adapters/citysourced/data"
	"Gateway311/adapters/citysourced/structs"
	"Gateway311/engine/common"
)

// ================================================================================================
//                                      CREATE
// ================================================================================================

// Report is the RPC container struct for the Report.Create service.  This service creates
// a new 311 report.
type Report struct{}

// Create fully processes a request to Create a Report.
func (r *Report) Create(rqst *structs.NCreateRequest, resp *structs.NCreateResponse) error {
	log.Debug("Create - request: %p  resp: %p\n", rqst, resp)
	// Make the Create Manager
	c := &createMgr{
		nreq:  rqst,
		nresp: resp,
	}
	log.Debug("createMgr: %#v\n", *c)

	return runRequest(processer(c))
}

// createMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Create a Report.
//  1. Validates and converts the request from the Normal form to the CitySourced native XML form.
//  2. Calls the appropriate CitySourced REST interface with proper credentials.
//  3. Converts the CitySourced reply back to Normal form.
//  4. Returns the Normal Response, and any errors.
type createMgr struct {
	nreq  *structs.NCreateRequest
	req   *create.Request
	url   string
	resp  *create.Response
	nresp *structs.NCreateResponse
}

func (c *createMgr) convertRequest() error {
	provider, err := data.MIDProvider(c.nreq.MID)
	if err != nil {
		return err
	}
	c.url = provider.URL
	c.req = &create.Request{
		APIAuthKey:        provider.Key,
		APIRequestType:    "CreateThreeOneOne",
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

// Process executes the request to create a new report.
func (c *createMgr) process() error {
	resp, err := c.req.Process(c.url)
	c.resp = resp
	return err
}

func (c *createMgr) convertResponse() error {
	c.nresp.Message = c.resp.Message
	c.nresp.ID = c.resp.ID
	c.nresp.AuthorID = c.resp.AuthorID
	return nil
}

func (c *createMgr) fail(err error) error {
	c.nresp.Message = "Failed - " + err.Error()
	c.nresp.ID = ""
	c.nresp.AuthorID = ""
	return err
}

func (c *createMgr) String() string {
	ls := new(common.LogString)
	ls.AddS("Create\n")
	ls.AddS(c.nreq.String())
	ls.AddS(c.req.String())
	ls.AddS(c.resp.String())
	ls.AddS(c.nresp.String())
	return ls.Box(90)
}
