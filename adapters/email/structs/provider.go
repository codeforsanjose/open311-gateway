package structs

//go:generate stringer -type=PayloadType

// PayloadType enumerates the valid payload types.
type PayloadType int

// NRT* are constants enumerating the valid request types.
const (
	PTUnknown PayloadType = iota
	PTJSON
	PTXML
	PTByte
	PTString
)

// ------------------------------------ Interfaces ------------------------------------

// Payloader is an interface to a Payload object to retrieve the Payload contents.
type Payloader interface {
	Get() (ptype PayloadType, subject string, contents interface{})
}

// Addresser retrieves the Addressing information for a call to a Provider.
type Addresser interface {
	Get() (to []string, from []string)
}

// ------------------------------------ Payload ------------------------------------

// Payload represents the payload for a call to a Provider.
type Payload struct {
	ptype   PayloadType
	subject string
	content interface{}
}

// NewPayloadString returns a new Payload struct for the specified body of type string.
func NewPayloadString(subject string, body *string) *Payload {
	return &Payload{
		ptype:   PTString,
		subject: subject,
		content: body,
	}
}

// NewPayloadByte returns a new Payload struct for the specified body of type string.
func NewPayloadByte(subject string, body []byte) *Payload {
	return &Payload{
		ptype:   PTByte,
		subject: subject,
		content: body,
	}
}

// Get returns the type. subject and contents of a Payload.
func (r Payload) Get() (ptype PayloadType, subject string, contents interface{}) {
	return r.ptype, r.subject, r.content
}

// ------------------------------------ Address ------------------------------------

// Address represents the addressing information for a call to a Provider.
type Address struct {
	To   []string
	From []string
}

// Get returns the "To, From and Subject address info for a call to a Provider
func (r Address) Get() (to []string, from []string) {
	return r.To, r.From
}
