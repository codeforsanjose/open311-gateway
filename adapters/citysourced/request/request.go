package request

import (
	"Gateway311/engine/logs"
)

var (
	log = logs.Log
)

type processer interface {
	convertRequest() error
	process() error
	convertResponse() error
	fail(err error) error
	String() string
}

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
