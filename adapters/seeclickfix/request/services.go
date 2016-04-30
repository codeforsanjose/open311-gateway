package request

import (
	"time"

	"Gateway311/adapters/seeclickfix/common"
	"Gateway311/adapters/seeclickfix/data"
	"Gateway311/adapters/seeclickfix/structs"
	"Gateway311/adapters/seeclickfix/telemetry"
)

// ================================================================================================
//                                      CREATE
// ================================================================================================

// Service fully processes the Service request.
func (r *Report) Service(rqst *structs.NServiceRequest, resp *structs.NServiceResponse) error {
	log.Debugf("Service - request: %p  resp: %p\n", rqst, resp)
	// Make the Service Manager
	cm := &serviceMgr{
		nreq:  rqst,
		nresp: resp,
	}
	log.Debugf("serviceMgr: %#v\n", *cm)

	return runRequest(processer(cm))
}

// serviceMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Service a Report.
//  1. Validates and converts the request from the Normal form to the CitySourced native XML form.
//  2. Calls the appropriate CitySourced REST interface with proper credentials.
//  3. Converts the CitySourced reply back to Normal form.
//  4. Returns the Normal Response, and any errors.
type serviceMgr struct {
	nreq  *structs.NServiceRequest
	req   *services.Request
	url   string
	resp  *services.Response
	nresp *structs.NServiceResponse
}

func (c *serviceMgr) convertRequest() error {
	provider, err := data.MIDProvider(c.nreq.MID)
	if err != nil {
		return err
	}
	c.url = provider.URL
	c.req = &services.Request{
		APIAuthKey:        provider.Key,
		APIRequestType:    "ServiceThreeOneOne",
		APIRequestVersion: provider.APIVersion,
		DeviceType:        c.nreq.DeviceType,
		DeviceModel:       c.nreq.DeviceModel,
		DeviceID:          c.nreq.DeviceID,
		RequestType:       c.nreq.ServiceName,
		RequestTypeID:     c.nreq.MID.ID,
		ImageURL:          c.nreq.MediaURL,
		Latitude:          c.nreq.Latitude,
		Longitude:         c.nreq.Longitude,
		Description:       c.nreq.Description,
		AuthorNameFirst:   c.nreq.FirstName,
		AuthorNameLast:    c.nreq.LastName,
		AuthorEmail:       c.nreq.Email,
		AuthorTelephone:   c.nreq.Phone,
		AuthorIsAnonymous: c.nreq.IsAnonymous,
	}
	telemetry.SendRPC(c.nreq.GetIDS(), "open", "", c.url, 0, time.Now())
	return nil
}

// Process executes the request to service a new report.
func (c *serviceMgr) process() error {
	resp, err := c.req.Process(c.url)
	c.resp = resp
	return err
}

func (c *serviceMgr) convertResponse() (int, error) {
	route := c.nreq.GetRoute()
	c.nresp.SetIDF(c.nreq.GetID)
	c.nresp.SetRoute(route)
	c.nresp.RID = structs.NewRID(route, c.resp.ID)
	c.nresp.Message = c.resp.Message
	c.nresp.AccountID = c.resp.AuthorID
	return 1, nil
}

func (c *serviceMgr) fail(err error) error {
	c.nresp.Message = "Failed - " + err.Error()
	c.nresp.RID = structs.ReportID{}
	c.nresp.AccountID = ""
	return err
}

func (c *serviceMgr) getIDS() string {
	return c.nreq.GetIDS()
}

func (c *serviceMgr) getRoute() string {
	return c.nreq.GetRoute().String()
}

func (c *serviceMgr) String() string {
	ls := new(common.LogString)
	ls.AddS("Service\n")
	ls.AddS(c.nreq.String())
	ls.AddS(c.req.String())
	ls.AddS(c.resp.String())
	ls.AddS(c.nresp.String())
	return ls.Box(90)
}
