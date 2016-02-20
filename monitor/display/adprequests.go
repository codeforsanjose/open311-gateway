package display

import (
	"fmt"
	"time"

	"Gateway311/monitor/telemetry"
)

type adpRequestType struct {
	id          string
	status      string
	route       string
	url         string
	start       time.Time
	startSet    bool
	complete    time.Time
	completeSet bool
}

func newAdpRPC(m telemetry.Message) (dataInterface, error) {
	adpRequest := new(adpRequestType)
	err := adpRequest.update(m)
	if err != nil {
		return nil, err
	}
	return dataInterface(adpRequest), nil
}

func (r adpRequestType) display() string {
	var dur time.Duration
	if r.startSet && r.completeSet {
		dur = r.complete.Sub(r.start) / time.Millisecond
	} else {
		dur = time.Since(r.start) / time.Millisecond
	}
	return fmt.Sprintf("%-10s  %-60s %-10s %6dms", r.id, fmt.Sprintf("%s-->%s", r.route, r.url), r.status, dur)
}

func (r *adpRequestType) update(m telemetry.Message) error {
	s, err := telemetry.UnmarshalAdpRPCMsg(m)
	if err != nil {
		return err
	}
	r.id = s.ID
	r.status = s.Status
	r.route = s.Route
	if s.URL > "" {
		r.url = s.URL
	}
	switch {
	case s.Status == "open" && s.At.Year() > 2000 && !r.startSet:
		r.start = s.At
		r.startSet = true
	case (s.Status == "complete" || s.Status == "error") && s.At.Year() > 2000:
		r.complete = s.At
		r.completeSet = true
	}
	return nil
}

func (r *adpRequestType) key() string {
	return r.id
}

func (r *adpRequestType) getLastUpdate() time.Time {
	return time.Now()
}

func (r *adpRequestType) setStatus(status string) {
	r.status = status
}
