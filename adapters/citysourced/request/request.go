package request

import (
	"Gateway311/engine/logs"
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
	String() string
}

// runRequest runs all of the common request processing operations.
func runRequest(r processer) error {
	if err := r.convertRequest(); err != nil {
		return r.fail(err)
	}
	if err := r.process(); err != nil {
		return r.fail(err)
	}
	if err := r.convertResponse(); err != nil {
		return r.fail(err)
	}
	log.Debug("COMPLETED:%s\n", r.String())
	return nil
}
