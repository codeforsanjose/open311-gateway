package mail

import (
	"fmt"
	"os"
	"testing"

	"github.com/open311-gateway/adapters/email/data"
	"github.com/open311-gateway/adapters/email/structs"
)

func TestMain(m *testing.M) {
	if err := data.Init("/Users/james/Dropbox/Development/go/src/Gateway311/adapters/email/data/config.json"); err != nil {
		panic("Unable to Init Data")
	}

	Init()
	os.Exit(m.Run())
}

func TestLoad(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestLoad] =============================\n\n")

	fmt.Printf("Auth: %+v\n", auth)
	fmt.Printf("Dialer: %+v\n", dialer)
}

func TestSend(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestSend] =============================\n\n")

	msg := `Request Type: Graffiti

Description: Spray paint!!

Location: 10.0, -20.0

---Author---
Haskell, James
Email: james@cloudyoperations.com
Phone: 1234567890
`
	{
		fmt.Printf("\n------------------ Sending String -------------------------\n")
		payload := structs.NewPayloadString("Open311 Request", &msg)

		address := structs.Address{
			To:   []string{"jameskhaskell@gmail.com"},
			From: []string{"cfsj.test@gmail.com", "Code for San Jose"},
		}

		if err := Send(address, payload); err != nil {
			t.Errorf("Send failed - %s", err)
		}
	}

	{
		fmt.Printf("\n------------------ Sending Byte -------------------------\n")
		payload := structs.NewPayloadByte("Open311 Request", []byte(msg))

		address := structs.Address{
			To:   []string{"jameskhaskell@gmail.com"},
			From: []string{"cfsj.test@gmail.com", "Code for San Jose"},
		}

		if err := Send(address, payload); err != nil {
			t.Errorf("Send failed - %s", err)
		}
	}

}
