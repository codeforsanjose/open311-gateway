package display

import (
	"fmt"
	"time"

	"github.com/open311-gateway/monitor/telemetry"
)

type adpStatusType struct {
	name       string
	status     string
	addr       string
	rqstCount  int64
	lastUpdate time.Time
}

func newAdpStatusType(m telemetry.Message) (dataInterface, error) {
	adpStatus := new(adpStatusType)
	err := adpStatus.update(m)
	if err != nil {
		return nil, err
	}
	return dataInterface(adpStatus), nil
}

func (r adpStatusType) display() string {
	return fmt.Sprintf("%-10s  %10s %8.1f  %s", r.name, r.status, time.Since(r.lastUpdate).Seconds(), r.addr)
}

func (r *adpStatusType) update(m telemetry.Message) error {
	s, err := telemetry.UnmarshalEngStatusMsg(m)
	if err != nil {
		return err
	}

	r.name = s.Name
	r.lastUpdate = time.Now()
	r.status = s.Status
	if s.Addr > "" {
		r.addr = s.Addr
	}
	return nil
}

func (r *adpStatusType) key() string {
	return r.name
}

func (r *adpStatusType) getLastUpdate() time.Time {
	return r.lastUpdate
}

func (r *adpStatusType) setStatus(status string) {
	r.status = status
}
