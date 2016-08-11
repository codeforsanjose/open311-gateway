package request

import (
	"github.com/open311-gateway/adapters/citysourced/logs"
	"github.com/open311-gateway/adapters/citysourced/telemetry"
	"time"
)

var (
	log = logs.Log
)

// Report is the RPC container struct for the Report.Create service.  This service creates
// a new 311 report.
type Report struct{}

// processer is the interface used to run all the common request processing steps (see runRequest()).
type processer interface {
	convertRequest() error
	process() error
	convertResponse() (int, error)
	fail(err error) error
	getIDS() string
	getRoute() string
	String() string
}

// runRequest runs all of the common request processing operations.
func runRequest(r processer) error {
	id := r.getIDS()
	telemetry.SendRPC(id, "open", r.getRoute(), "", 0, time.Now())

	if err := r.convertRequest(); err != nil {
		telemetry.SendRPC(id, "error", r.getRoute(), "", 0, time.Now())
		return r.fail(err)
	}
	if err := r.process(); err != nil {
		telemetry.SendRPC(id, "error", r.getRoute(), "", 0, time.Now())
		return r.fail(err)
	}
	resultCount, err := r.convertResponse()
	if err != nil {
		telemetry.SendRPC(id, "error", r.getRoute(), "", 0, time.Now())
		return r.fail(err)
	}
	// telemetry.SendRPC(r.getID(), data.AdapterName(), "", "", "adp-send")
	telemetry.SendRPC(id, "done", r.getRoute(), "", resultCount, time.Now())
	log.Debugf("Request COMPLETED:%s\n", r.String())
	return nil
}
