package request

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/open311-gateway/adapters/email/common"
	"github.com/open311-gateway/adapters/email/create"
	"github.com/open311-gateway/adapters/email/data"
	"github.com/open311-gateway/adapters/email/structs"
	"github.com/open311-gateway/adapters/email/telemetry"

	log "github.com/jeffizhungry/logrus"
)

var successMessages map[string]bool

// ================================================================================================
//                                      CREATE
// ================================================================================================

// Create fully processes the Create request.
func (r *Report) Create(rqst *structs.NCreateRequest, resp *structs.NCreateResponse) (err error) {
	log.Debugf("Create - request: %p  resp: %p\n", rqst, resp)
	// Make the Create Manager
	cm := &createMgr{
		nreq:  rqst,
		nresp: resp,
	}

	if cm.nsrv, err = data.ServiceFromID(cm.nreq.MID); err != nil {
		return fmt.Errorf("unable to process Create request - %s", err.Error())
	}

	log.Debugf("createMgr: %#v\n", *cm)

	return runRequest(processer(cm))
}

// createMgr conglomerates the Normal and Native structs and supervisor logic
// for processing a request to Create a Report.
//  1. Validates and converts the request from the Normal form to the CitySourced native XML form.
//  2. Calls the appropriate CitySourced REST interface with proper credentials.
//  3. Converts the CitySourced reply back to Normal form.
//  4. Returns the Normal Response, and any errors.
type createMgr struct {
	nreq *structs.NCreateRequest
	req  *create.Request

	nsrv structs.NService

	resp  *create.Response
	nresp *structs.NCreateResponse
}

func (c *createMgr) convertRequest() error {
	fail := func(err string) error {
		return fmt.Errorf("Unable to create the email - %s", err)
	}
	telemetry.SendRPC(c.nreq.GetIDS(), "open", "", "", 0, time.Now())

	// Fill in the Service Name
	c.nreq.ServiceName = c.nsrv.Name

	// Get the EmailSender interface.
	provider, err := data.MIDProvider(c.nreq.MID)
	if err != nil {
		return fail(fmt.Sprintf("unable to determine route/sender for the Create request - %s", err))
	}
	sender := provider.Email

	// Execute the template
	body, err := c.createBody(sender.Template())
	if err != nil {
		fail(err.Error())
	}

	c.req = &create.Request{
		Sender: sender,
		Body:   structs.NewPayloadString(&body),
	}

	return nil
}

// Process executes the request to create a new report.
func (c *createMgr) process() error {
	resp, err := c.req.Process()
	c.resp = resp
	return err
}

func (c *createMgr) convertResponse() (int, error) {
	rspOK := func(msg string) bool {
		_, ok := successMessages[strings.ToLower(msg)]
		log.Debugf("successMessages: %#v", successMessages)
		log.Debugf("Msg: %q / %q result: %t", msg, strings.ToLower(msg), ok)
		return ok
	}
	route := c.nreq.GetRoute()
	c.nresp.SetIDF(c.nreq.GetID)
	c.nresp.SetRoute(route)

	if !rspOK(c.resp.Message) {
		log.Errorf("Returning response message: %s", c.resp.Message)
		c.nresp.Message = c.resp.Message
		return 1, nil
	}

	log.Debugf("Service: %s", c.nsrv.SString())

	c.nresp.Message = fmt.Sprintf(c.nsrv.ServiceNotice, c.nsrv.Name)
	return 1, nil
}

func (c *createMgr) fail(err error) error {
	c.nresp.Message = "Failed - " + err.Error()
	c.nresp.RID = structs.ReportID{}
	c.nresp.AccountID = ""
	return err
}

func (c *createMgr) getIDS() string {
	return c.nreq.GetIDS()
}

func (c *createMgr) getRoute() string {
	return c.nreq.GetRoute().String()
}

// createEmail creates an email message from the request using the specified template
func (c *createMgr) createBody(tmpl *template.Template) (string, error) {
	var doc bytes.Buffer
	// Apply the values we have initialized in our struct context to the template.
	log.Debugf("Executing template: %p", tmpl)
	if err := tmpl.Execute(&doc, c.nreq); err != nil {
		log.Errorf("error trying to execute email template - %s", err.Error())
		return "", err
	}
	log.Debugf("Doc:\n%s", doc.String())
	return doc.String(), nil
}

func (c *createMgr) String() string {
	ls := new(common.LogString)
	ls.AddS("Create\n")
	ls.AddS(c.nreq.String())
	ls.AddS(c.req.String())
	ls.AddS(c.nsrv.String())
	ls.AddS(c.resp.String())
	ls.AddS(c.nresp.String())
	return ls.Box(90)
}

func init() {
	successMessages = map[string]bool{
		"success": true,
		"ok":      true,
	}
}
