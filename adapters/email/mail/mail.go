package mail

import (
	"CitySourcedAPI/logs"
	"fmt"

	"github.com/open311-gateway/adapters/email/common"
	"github.com/open311-gateway/adapters/email/data"
	"github.com/open311-gateway/adapters/email/structs"

	log "github.com/jeffizhungry/logrus"
	"gopkg.in/gomail.v2"
)

var (
	auth   data.EmailAuthData
	dialer *gomail.Dialer
)

// Init should be called at program startup to initialize
func Init() {
	logs.Init(true)
	auth = data.GetEmailAuth()
	log.Debugf("Auth: %v", auth)

	dialer = gomail.NewDialer(
		auth.Server,
		auth.Port,
		auth.Account,
		auth.Password,
	)
}

// Send sends an email.
func Send(a data.EmailSender, p structs.Payloader) error {
	if dialer == nil {
		Init()
	}
	var msg string
	to, from, subject := a.Address()
	log.Debugf("to: %#v  from: %#v  subject: %q", to, from, subject)
	ptype, content := p.Get()
	log.Debugf("ptype: %v  content: %v (%[2]T)", ptype, content)
	// log.Debugf("dialer:\n%s\n", spew.Sdump(dialer))

	fail := func() error {
		return fmt.Errorf("invalid payload (type: %T) received by Send() - must be either string or []byte", content)
	}

	// Validate the Payload - if it's []byte, convert it to a string.
	switch content := content.(type) {
	case *string:
		if ptype != structs.PTString {
			return fail()
		}
		msg = *content

	case string:
		if ptype != structs.PTString {
			return fail()
		}
		msg = content

	case []byte:
		if ptype != structs.PTByte {
			return fail()
		}
		msg = common.ByteToString(content, 0)
	default:
		return fmt.Errorf("invalid Payload received by Send() - must be either string or []byte")

	}

	if len(msg) == 0 {
		return fmt.Errorf("no message")
	}

	fmt.Printf("to: %v  subject: %v  body: %v\n", to, subject, msg)
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from[0], from[1])
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", msg)

	go func() {
		if err := dialer.DialAndSend(m); err != nil {
			log.Error(err)
		}
	}()

	return nil
}
