package display

import (
	"fmt"
	"time"

	"github.com/open311-gateway/monitor/telemetry"
)

type engAdpRequestType struct {
	id          string
	status      string
	route       string
	start       time.Time
	startSet    bool
	complete    time.Time
	completeSet bool
}

func newEngRPC(m telemetry.Message) (dataInterface, error) {
	engAdpRequest := new(engAdpRequestType)
	err := engAdpRequest.update(m)
	if err != nil {
		return nil, err
	}
	return dataInterface(engAdpRequest), nil
}

func (r engAdpRequestType) display() string {
	var dur time.Duration
	if r.startSet && r.completeSet {
		dur = r.complete.Sub(r.start) / time.Millisecond
	} else {
		dur = time.Since(r.start) / time.Millisecond
	}
	return fmt.Sprintf("%-10s  %-25s  %-12s %6dms", r.id, r.route, r.status, dur)
}

func (r *engAdpRequestType) update(m telemetry.Message) error {
	s, err := telemetry.UnmarshalEngRPCMsg(m)
	if err != nil {
		return err
	}
	r.id = s.ID
	r.status = s.Status
	if s.Route > "" {
		r.route = s.Route
	}
	switch {
	case s.Status == "open" && s.At.Year() > 2000 && !r.startSet:
		r.start = s.At
		r.startSet = true
	case (s.Status == "done" || s.Status == "error") && s.At.Year() > 2000:
		r.complete = s.At
		r.completeSet = true
	}
	return nil
}

func (r *engAdpRequestType) key() string {
	return r.id
}

func (r *engAdpRequestType) getLastUpdate() time.Time {
	return time.Now()
}

func (r *engAdpRequestType) setStatus(status string) {
	r.status = status
}
