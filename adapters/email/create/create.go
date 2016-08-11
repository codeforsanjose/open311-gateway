package create

import (
	"fmt"
	"strings"

	"github.com/open311-gateway/adapters/email/common"
	"github.com/open311-gateway/adapters/email/data"
	"github.com/open311-gateway/adapters/email/mail"
	"github.com/open311-gateway/adapters/email/structs"
)

// ================================================================================================
//                                      CREATE
// ================================================================================================

// Request represents the Sender and Body for the email.
type Request struct {
	Sender data.EmailSender
	Body   *structs.Payload
}

// Process executes the request to create a new report.
func (r *Request) Process() (*Response, error) {
	fail := func(err error) (*Response, error) {
		response := Response{
			Message: fmt.Sprintf("unable to send email - %s", err),
		}
		return &response, err
	}

	if err := mail.Send(r.Sender, r.Body); err != nil {
		fail(err)
	}

	return &Response{"Success"}, nil
}

// ------------------------------------------------------------------------------------------------

// Response is the response to creating or updating a report.
type Response struct {
	Message string `json:"Message" xml:"Message"`
}

// ================================================================================================
//                                      TEMPLATES
// ================================================================================================

/*
// createEmail creates an email message from the request using the specified template
func (r *Request) createEmail(tmpl *template.Template) (string, error) {
	var doc bytes.Buffer
	// Apply the values we have initialized in our struct context to the template.
	if err := tmpl.Execute(&doc, r); err != nil {
		log.Error("error trying to execute email template ", err)
		return "", err
	}
	log.Debug("Doc:\n%s", doc.String())
	return doc.String(), nil
}

*/

// ================================================================================================
//                                      STRINGS
// ================================================================================================

// String displays a Request
func (r Request) String() string {
	ls := new(common.LogString)
	ls.AddS("create.Request\n")
	to, from, subject := r.Sender.Address()
	ls.AddF("Sender - to: %#v  from: %#v\n", strings.Join(to, ", "), strings.Join(from, ", "))
	ls.AddF("Subject: %q\n", subject)
	ls.AddF("Message:\n%s\n", r.Body)
	return ls.Box(80)
}

// String displays a Response
func (r Response) String() string {
	ls := new(common.LogString)
	ls.AddS("create.Response\n")
	ls.AddF("Message: %v\n", r.Message)
	return ls.Box(80)
}
