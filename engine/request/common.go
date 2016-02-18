package request

import (
	"fmt"

	"github.com/ant0ine/go-json-rest/rest"
)

// ==============================================================================================================================
//                                      COMMON
// ==============================================================================================================================
type cType struct {
	self cIface
	id   int64
}

func (r *cType) load(p cIface, rqstID int64, rqst *rest.Request) error {
	r.self = p
	r.id = rqstID

	if err := rqst.DecodeJsonPayload(r.self); err != nil {
		if err.Error() != "JSON payload is empty" {
			return fmt.Errorf("Unable to process request: %s", err)
		}
	}
	if err := r.self.parseQP(rqst); err != nil {
		return fmt.Errorf("Unable to process request: %s", err)
	}

	if err := r.self.validate(); err != nil {
		return fmt.Errorf("Unable to process request: %s", err)
	}

	return nil
}

type cIface interface {
	parseQP(r *rest.Request) error
	validate() error
}

type cRType struct {
	id int64
}

func (r *cRType) SetID(id int64) {
	r.id = id
}
