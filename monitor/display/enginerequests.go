package display

import (
	"fmt"
	"time"

	"github.com/open311-gateway/monitor/telemetry"
)

type engRequestType struct {
	id          string
	rType       string
	status      string
	areaID      string
	start       time.Time
	startSet    bool
	complete    time.Time
	completeSet bool
}

func newEngRequest(m telemetry.Message) (dataInterface, error) {
	engRequest := new(engRequestType)
	err := engRequest.update(m)
	if err != nil {
		return nil, err
	}
	return dataInterface(engRequest), nil
}

func (r engRequestType) display() string {
	var dur time.Duration
	if r.startSet && r.completeSet {
		dur = r.complete.Sub(r.start) / time.Millisecond
	} else {
		dur = time.Since(r.start) / time.Millisecond
	}
	return fmt.Sprintf("%-10s  %-25s  %-12s %6dms", r.id, fmt.Sprintf("%s (%s)", r.rType, r.areaID), r.status, dur)
}

func (r *engRequestType) update(m telemetry.Message) error {
	s, err := telemetry.UnmarshalEngRequestMsg(m)
	if err != nil {
		return err
	}
	r.id = s.ID
	r.rType = s.Rtype
	r.status = s.Status
	r.areaID = s.AreaID
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

func (r *engRequestType) key() string {
	return r.id
}

func (r *engRequestType) getLastUpdate() time.Time {
	return time.Now()
}

func (r *engRequestType) setStatus(status string) {
	r.status = status
}
