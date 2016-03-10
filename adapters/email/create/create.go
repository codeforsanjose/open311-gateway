package create

import (
	"CitySourcedAPI/logs"
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"text/template"

	"Gateway311/adapters/email/common"
	"Gateway311/adapters/email/mail"
)

var (
	tmpl *template.Template
	log  = logs.Log
)

// ================================================================================================
//                                      CREATE
// ================================================================================================

// Request represents the XML payload for a report request to CitySourced.
type Request struct {
	To                string  //
	From              string  //
	Subject           string  //
	RequestType       string  `json:"RequestType" xml:"RequestType"`
	RequestTypeID     int     `json:"RequestTypeId" xml:"RequestTypeId"`
	ImageURL          string  `json:"ImageUrl" xml:"ImageUrl"`
	Latitude          float64 `json:"Latitude" xml:"Latitude"`
	Longitude         float64 `json:"Longitude" xml:"Longitude"`
	Description       string  `json:"Description" xml:"Description"`
	AuthorNameFirst   string  `json:"AuthorNameFirst" xml:"AuthorNameFirst"`
	AuthorNameLast    string  `json:"AuthorNameLast" xml:"AuthorNameLast"`
	AuthorEmail       string  `json:"AuthorEmail" xml:"AuthorEmail"`
	AuthorTelephone   string  `json:"AuthorTelephone" xml:"AuthorTelephone"`
	AuthorIsAnonymous bool    `json:"AuthorIsAnonymous" xml:"AuthorIsAnonymous"`
}

// Process executes the request to create a new report.
func (r *Request) Process(url string) (*Response, error) {
	fail := func(err error) (*Response, error) {
		response := Response{
			Message:  "Failed",
			ID:       "",
			AuthorID: "",
		}
		return &response, err
	}

	var payload = new(bytes.Buffer)
	{
		enc := xml.NewEncoder(payload)
		enc.Indent("  ", "    ")
		enc.Encode(r)
	}

	client := http.Client{Timeout: common.HttpClientTimeout}
	resp, err := client.Post(url, "application/xml", payload)
	if err != nil {
		return fail(err)
	}

	var response Response
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fail(err)
	}

	return &response, nil
}

// ------------------------------------------------------------------------------------------------

// Response is the response to creating or updating a report.
type Response struct {
	Message  string `json:"Message" xml:"Message"`
	ID       string `json:"ReportId" xml:"ReportId"`
	AuthorID string `json:"AuthorId" xml:"AuthorId"`
}

// ================================================================================================
//                                      TEMPLATES
// ================================================================================================

func (r *Request) SendEmail(recipients []string) error {
	fail := func(err error) error {
		errmsg := "unable to send email - " + err.Error()
		log.Errorf(errmsg)
		return errors.New(errmsg)
	}
	doc, err := r.createEmail()
	if err != nil {
		return fail(err)
	}

	mail.Send(recipients, doc)

	return nil
}

// createEmail creates an email message from the request using the standard email
// template (tmplCreateStd).
func (r *Request) createEmail() (doc bytes.Buffer, err error) {
	// Apply the values we have initialized in our struct context to the template.
	if err = tmpl.Execute(&doc, r); err != nil {
		log.Error("error trying to execute email template ", err)
	}
	log.Debug("Doc:\n%s", doc.String())
	return doc, nil
}

func init() {
	var err error
	// Create a new template for our SMTP message.
	tmpl = template.New("emailTemplate")
	if tmpl, err = tmpl.Parse(tmplCreateStd); err != nil {
		log.Error("error trying to parse mail template ", err)
	}
}

const tmplCreateStd = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

Request Type: {{.RequestType}}

Description: {{.Description}}

Location: {{.Latitude}}, {{.Longitude}}

---Author---
{{.AuthorNameLast}}, {{.AuthorNameFirst}}
Email: {{.AuthorEmail}}
Phone: {{.AuthorTelephone}}
`

// ================================================================================================
//                                      STRINGS
// ================================================================================================

// String displays a Request
func (r Request) String() string {
	ls := new(common.LogString)
	ls.AddS("create.Request\n")
	ls.AddF("Request - type: %q  id: %d\n", r.RequestType, r.RequestTypeID)
	ls.AddF("Location - lat: %v  lon: %v\n", r.Latitude, r.Longitude)
	ls.AddF("Description: %q\n", r.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", r.AuthorIsAnonymous, r.AuthorNameFirst, r.AuthorNameLast, r.AuthorEmail, r.AuthorTelephone)
	return ls.Box(80)
}

// String displays a Response
func (r Response) String() string {
	ls := new(common.LogString)
	ls.AddS("create.Response\n")
	ls.AddF("Message: %v\n", r.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", r.ID, r.AuthorID)
	return ls.Box(80)
}
