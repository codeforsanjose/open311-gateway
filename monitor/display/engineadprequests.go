package display

import (
	"fmt"
	"time"

	"Gateway311/monitor/telemetry"
)

type engAdpRequestType struct {
	id     string
	status string
	route  string
	at     time.Time
}

func newEngAdpRequestType(m telemetry.Message) (dataInterface, error) {
	engAdpRequest := new(engAdpRequestType)
	err := engAdpRequest.update(m)
	if err != nil {
		return nil, err
	}
	return dataInterface(engAdpRequest), nil
}

func (r engAdpRequestType) display() string {
	return fmt.Sprintf("%-10s  %-25s  %-12s %8.1f", r.id, r.route, r.status, time.Since(r.at).Seconds())
}

func (r *engAdpRequestType) update(m telemetry.Message) error {
	s, err := telemetry.UnmarshalAdpEngRequestMsg(m)
	if err != nil {
		return err
	}
	r.id = s.ID
	r.status = s.Status
	r.route = s.Route
	r.at = s.At
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
