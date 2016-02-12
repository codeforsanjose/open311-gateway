package request

import (
	"Gateway311/adapters/citysourced/data"
	"Gateway311/adapters/citysourced/logs"
	"Gateway311/engine/telemetry"
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
	convertResponse() error
	fail(err error) error
	getID() int64
	String() string
}

// runRequest runs all of the common request processing operations.
func runRequest(r processer) error {
	telemetry.Send(r.getID(), data.AdapterName(), "", "", "adp-recv")

	if err := r.convertRequest(); err != nil {
		return r.fail(err)
	}
	if err := r.process(); err != nil {
		return r.fail(err)
	}
	if err := r.convertResponse(); err != nil {
		return r.fail(err)
	}
	telemetry.Send(r.getID(), data.AdapterName(), "", "", "adp-send")
	log.Debug("COMPLETED:%s\n", r.String())
	return nil
}
