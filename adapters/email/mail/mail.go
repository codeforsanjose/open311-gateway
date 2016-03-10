package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"Gateway311/adapters/email/logs"
)

var (
	log = logs.Log
)

func Send(recipients []string, body bytes.Buffer) {

	auth := smtp.PlainAuth(
		"",
		"cfsj.test@gmail.com",
		"Ma%GvphDwKqy74Ct",
		"smtp.gmail.com",
	)

	err = smtp.SendMail("smtp.gmail.com:587",
		auth,
		"cfsj.test@gmail.com",
		[]string{"jameskhaskell@gmail.com"},
		body.Bytes())
	if err != nil {
		log.Error("ERROR: attempting to send a mail ", err)
	}
}

func Send1(recipient string, body string) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		"cfsj.test@gmail.com",
		"Ma%GvphDwKqy74Ct",
		"smtp.google.com",
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{recipient}
	msg := []byte("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail("smtp.gmail.com:465", auth, "cfsj.test@gmail.com", to, msg)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		// log.Fatal(err)
	}
}

func Send2(recipient string, body string) {
	auth := smtp.PlainAuth(
		"",
		"cfsj.test@gmail.com",
		"Ma%GvphDwKqy74Ct",
		"smtp.gmail.com",
	)

	c, err := smtp.Dial("smtp.google.com:465")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	c.Auth(auth)

	c.Mail("cfsj.test+from@gmail.com")
	c.Rcpt(recipient)

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(body)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}

type SmtpTemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

var err error
var doc bytes.Buffer

// Go template
const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,

{{.From}}
`

func Send3() {

	// Authenticate with Gmail (analagous to logging in to your Gmail account in the browser)
	auth := smtp.PlainAuth(
		"",
		"cfsj.test@gmail.com",
		"Ma%GvphDwKqy74Ct",
		"smtp.gmail.com",
	)

	// Set the context for the email template.
	context := &SmtpTemplateData{"Code for San Jose <cfsj.test@gmail.com>",
		"Recipient Person <???>",
		"RethinkDB is so slick!",
		"Hey Recipient, just wanted to drop you a line and let you know how I feel about ReQL..."}

	// Create a new template for our SMTP message.
	t := template.New("emailTemplate")
	if t, err = t.Parse(emailTemplate); err != nil {
		log.Error("error trying to parse mail template ", err)
	}

	// Apply the values we have initialized in our struct context to the template.
	if err = t.Execute(&doc, context); err != nil {
		log.Error("error trying to execute mail template ", err)
	}
	fmt.Print("Doc: ", doc.String())

	// Actually perform the step of sending the email - (3) above
	err = smtp.SendMail("smtp.gmail.com:587",
		auth,
		"cfsj.test@gmail.com",
		[]string{"jameskhaskell@gmail.com"},
		doc.Bytes())
	if err != nil {
		log.Error("ERROR: attempting to send a mail ", err)
	}
}
