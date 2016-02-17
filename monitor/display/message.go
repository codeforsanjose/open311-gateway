package display

import (
	"fmt"
	"strings"
)

// ==============================================================================================================================
//                                      MESSAGE DATA
// ==============================================================================================================================

// Message types
const (
	msgTypeES = "ES"
	msgTypeER = "ER"
	msgTypeEA = "EA"

	msgTypeIndex = 0
)

// -------------------------------------------- message --------------------------------------------------------------------

type message struct {
	mType string
	key   string
	data  []string
}

func newMessage(b []byte, n int) (message, error) {
	if n <= 0 {
		return message{}, fmt.Errorf("message has no contents")
	}
	m := strings.Split(string(b[0:n]), "|")
	return message{
		mType: m[0],
		key:   m[msgKeys[m[0]]],
		data:  m,
	}, nil
}

func newMessageTest(msg []string) message {
	return message{
		mType: msg[0],
		key:   msg[1],
		data:  msg,
	}
}

func (r *message) valid() bool {
	return len(r.data) == msgLen[r.mType]
}

func (r message) String() string {
	return fmt.Sprintf("%1s:%-12s  [%v]", r.mType, r.key, r.data)
}

var msgKeys map[string]int

func initMsgKeys() {
	msgKeys = make(map[string]int)
	msgKeys[msgTypeES] = esName
	msgKeys[msgTypeER] = erID
	msgKeys[msgTypeEA] = eaID
}

var msgLen map[string]int

func initMsgLen() {
	msgLen = make(map[string]int)
	msgLen[msgTypeES] = esLength
	msgLen[msgTypeER] = erLength
	msgLen[msgTypeEA] = eaLength
}
