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
}

func (cx *cType) init(p cIface, r *rest.Request) error {
	cx.self = p

	if err := r.DecodeJsonPayload(cx.self); err != nil {
		return fmt.Errorf("Unable to process request: %s", err)
	}
	if err := cx.self.parseQP(r); err != nil {
		return fmt.Errorf("Unable to process request: %s", err)
	}

	cx.self.validate()

	return nil
}

type cIface interface {
	parseQP(r *rest.Request) error
	validate()
}
