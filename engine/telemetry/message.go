package telemetry

import (
	"fmt"
	"strings"
	"time"
)

// ==============================================================================================================================
//                                      MESSAGE DATA
// ==============================================================================================================================

// Message types
const (
	MsgTypeES   = "ES"   // Engine Status
	MsgTypeER   = "ER"   // Engine Request
	MsgTypeERPC = "ERPC" // Engine RPC

	MsgTypeAS   = "AS"   // Adapter Status
	MsgTypeARPC = "ARPC" // Adapter RPC

	MsgTypeIndex = 0
	msgDelimiter = "|"
)

var (
	msgKeys map[string]int
	msgLen  map[string]int
)

func init() {
	initMsgKeys()
	initMsgLen()
}

func initMsgKeys() {
	msgKeys = make(map[string]int)
	msgKeys[MsgTypeES] = esName
	msgKeys[MsgTypeER] = erID
	msgKeys[MsgTypeERPC] = eaID
}

func initMsgLen() {
	msgLen = make(map[string]int)
	msgLen[MsgTypeES] = esLength
	msgLen[MsgTypeER] = erLength
	msgLen[MsgTypeERPC] = eaLength
}

type msgSender interface {
	Marshal() ([]byte, error)
}

// -------------------------------------------- message --------------------------------------------------------------------

// Message represents a Raw tatus Message (i.e. text).
type Message struct {
	mType string
	key   string
	data  []string
}

// NewMessage converts a Raw Message ([]byte) into a Message.
func NewMessage(b []byte, n int) (Message, error) {
	if n <= 0 {
		return Message{}, fmt.Errorf("message has no contents")
	}
	m := strings.Split(string(b[0:n]), msgDelimiter)
	return Message{
		mType: m[0],
		key:   m[msgKeys[m[0]]],
		data:  m,
	}, nil
}

// NewMessageTest returns a Raw Message from an array of strings.  It is for testing purposes only.
func NewMessageTest(msg []string) Message {
	return Message{
		mType: msg[0],
		key:   msg[1],
		data:  msg,
	}
}

func (r *Message) valid() bool {
	return len(r.data) == msgLen[r.mType]
}

// Key returns the message key.
func (r *Message) Key() string {
	return r.key
}

// Mtype returns the message type (mType).
func (r *Message) Mtype() string {
	return r.mType
}

// Data returns data (Raw Message).
func (r *Message) Data() []string {
	return r.data
}

func (r Message) String() string {
	return fmt.Sprintf("%1s:%-12s  [%v]", r.mType, r.key, r.data)
}

// -------------------------------------------- EngStatusMsgType --------------------------------------------------------------------

// EngStatusMsgType represents the Engine Status messages.
type EngStatusMsgType struct {
	Name     string
	Status   string
	Adapters string
	Addr     string
}

const (
	esName int = 1 + iota
	esStatus
	esAdapters
	esAddr
	esLength
)

// UnmarshalEngStatusMsg converts a Raw Message to an EngStatusMsgType instance
func UnmarshalEngStatusMsg(m Message) (*EngStatusMsgType, error) {
	if m.mType != MsgTypeES {
		return &EngStatusMsgType{}, fmt.Errorf("invalid message type: %q sent to EngineStatus - message: %v", m.mType, m)
	}
	if !m.valid() {
		return &EngStatusMsgType{}, fmt.Errorf("invalid message: %#v", m)
	}

	return &EngStatusMsgType{
		Name:     m.data[esName],
		Status:   m.data[esStatus],
		Adapters: m.data[esAdapters],
		Addr:     m.data[esAddr],
	}, nil
}

// Marshal converts a EngStatusMsgType to a Raw Message.
func (r EngStatusMsgType) Marshal() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s", MsgTypeES, msgDelimiter, r.Name, msgDelimiter, r.Status, msgDelimiter, r.Addr, msgDelimiter, r.Adapters)), nil
}

// -------------------------------------------- EngRequestMsgType --------------------------------------------------------------------

// EngRequestMsgType represents the Engine Request messages.
type EngRequestMsgType struct {
	ID     string
	Rtype  string
	Status string
	At     time.Time
	AreaID string
}

const (
	erID int = 1 + iota
	erRqstType
	erStatus
	erAt
	erAreaID
	erLength
)

// UnmarshalEngRequestMsg converts a Raw Message to an EngRequestMsgType instance
func UnmarshalEngRequestMsg(m Message) (*EngRequestMsgType, error) {
	if m.mType != MsgTypeER {
		return &EngRequestMsgType{}, fmt.Errorf("invalid message type: %q sent to EngineRequest - message: %v", m.mType, m)
	}
	if !m.valid() {
		return &EngRequestMsgType{}, fmt.Errorf("invalid message: %#v", m)
	}

	s := EngRequestMsgType{
		ID:     m.data[erID],
		Rtype:  m.data[erRqstType],
		Status: m.data[erStatus],
		AreaID: m.data[erAreaID],
	}
	if at, err := time.Parse(time.RFC3339Nano, m.data[erAt]); err == nil {
		s.At = at
	} else {
		s.At = time.Now()
	}
	return &s, nil
}

// Marshal converts a EngRequestMsgType to a Raw Message.
func (r EngRequestMsgType) Marshal() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s", MsgTypeER, msgDelimiter, r.ID, msgDelimiter, r.Rtype, msgDelimiter, r.Status, msgDelimiter, r.At.Format(time.RFC3339), msgDelimiter, r.AreaID)), nil
}

// -------------------------------------------- EngRPCMsgType --------------------------------------------------------------------

// EngRPCMsgType represents the Engine Adapter Request messages.
type EngRPCMsgType struct {
	ID     string
	Status string
	Route  string
	At     time.Time
}

const (
	eaID int = 1 + iota
	eaStatus
	eaRoute
	eaAt
	eaLength
)

// UnmarshalEngRPCMsg converts a Raw Message to an EngRPCMsgType instance
func UnmarshalEngRPCMsg(m Message) (*EngRPCMsgType, error) {
	if m.mType != MsgTypeERPC {
		return &EngRPCMsgType{}, fmt.Errorf("invalid message type: %q sent to EngineRequest - message: %v", m.mType, m)
	}
	if !m.valid() {
		return &EngRPCMsgType{}, fmt.Errorf("invalid message: %#v", m)
	}

	s := EngRPCMsgType{
		ID:     m.data[eaID],
		Status: m.data[eaStatus],
		Route:  m.data[eaRoute],
	}
	if at, err := time.Parse(time.RFC3339Nano, m.data[eaAt]); err == nil {
		s.At = at
	} else {
		s.At = time.Now()
	}
	return &s, nil

}

// Marshal converts a EngRPCMsgType to a Raw Message.
func (r EngRPCMsgType) Marshal() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s", MsgTypeERPC, msgDelimiter, r.ID, msgDelimiter, r.Status, msgDelimiter, r.Route, msgDelimiter, r.At.Format(time.RFC3339))), nil
}
