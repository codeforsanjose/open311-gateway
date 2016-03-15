package mail

import (
	"fmt"

	"Gateway311/adapters/email/common"
	"Gateway311/adapters/email/data"
	"Gateway311/adapters/email/logs"
	"Gateway311/adapters/email/structs"

	"gopkg.in/gomail.v2"
)

var (
	log    = logs.Log
	auth   data.EmailAuthData
	dialer *gomail.Dialer
)

func Init() {
	logs.Init(true)
	auth = data.GetEmailAuth()
	log.Debug("Auth: %v", auth)

	dialer = gomail.NewDialer(
		auth.Server,
		auth.Port,
		auth.Account,
		auth.Password,
	)
}

func Send(a structs.Addresser, p structs.Payloader) error {
	if dialer == nil {
		Init()
	}
	var msg string
	recipients, from := a.Get()
	log.Debug("recipients: %#v  from: %#v", recipients, from)
	ptype, subject, c := p.Get()
	log.Debug("ptype: %v  subject: %v  c: %v (%T)", ptype, subject, c, c)
	// log.Debug("dialer:\n%s\n", spew.Sdump(dialer))

	fail := func() error {
		return fmt.Errorf("invalid payload (type: %T) received by Send() - must be either string or []byte", c)
	}

	// Validate the Payload - if it's []byte, convert it to a string.
	switch tcontents := c.(type) {
	case *string:
		if ptype != structs.PTString {
			return fail()
		}
		msg = *tcontents

	case string:
		if ptype != structs.PTString {
			return fail()
		}
		msg = tcontents

	case []byte:
		if ptype != structs.PTByte {
			return fail()
		}
		msg = common.ByteToString(tcontents, 0)
	default:
		return fmt.Errorf("invalid Payload received by Send() - must be either string or []byte")

	}

	if len(msg) == 0 {
		return fmt.Errorf("no message")
	}

	fmt.Printf("recipients: %v  subject: %v  body: %v\n", recipients, subject, msg)
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from[0], from[1])
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", msg)

	if err := dialer.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}
