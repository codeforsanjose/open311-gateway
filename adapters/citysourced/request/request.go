package request

import (
	"Gateway311/adapters/citysourced/logs"
	"Gateway311/adapters/citysourced/telemetry"
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
	convertResponse() error
	fail(err error) error
	getIDS() string
	getRoute() string
	String() string
}

// runRequest runs all of the common request processing operations.
func runRequest(r processer) error {
	telemetry.SendRPC(r.getIDS(), "open", r.getRoute(), "", time.Now())

	if err := r.convertRequest(); err != nil {
		telemetry.SendRPC(r.getIDS(), "error", r.getRoute(), "", time.Now())
		return r.fail(err)
	}
	if err := r.process(); err != nil {
		telemetry.SendRPC(r.getIDS(), "error", r.getRoute(), "", time.Now())
		return r.fail(err)
	}
	if err := r.convertResponse(); err != nil {
		telemetry.SendRPC(r.getIDS(), "error", r.getRoute(), "", time.Now())
		return r.fail(err)
	}
	// telemetry.SendRPC(r.getID(), data.AdapterName(), "", "", "adp-send")
	telemetry.SendRPC(r.getIDS(), "complete", r.getRoute(), "", time.Now())
	log.Debug("COMPLETED:%s\n", r.String())
	return nil
}
