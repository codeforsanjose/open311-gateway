package display

import (
	"fmt"
	"time"

	"github.com/open311-gateway/monitor/telemetry"
)

type engStatusType struct {
	name       string
	status     string
	adapters   string
	addr       string
	rqstCount  int64
	lastUpdate time.Time
}

func newEngStatus(m telemetry.Message) (dataInterface, error) {
	engStatus := new(engStatusType)
	err := engStatus.update(m)
	if err != nil {
		return nil, err
	}
	return dataInterface(engStatus), nil
}

func (r engStatusType) display() string {
	return fmt.Sprintf("%-10s  %10s %8.1f  %s", r.name, r.status, time.Since(r.lastUpdate).Seconds(), r.addr)
}

func (r *engStatusType) update(m telemetry.Message) error {
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

func (r *engStatusType) key() string {
	return r.name
}

func (r *engStatusType) getLastUpdate() time.Time {
	return r.lastUpdate
}

func (r *engStatusType) setStatus(status string) {
	r.status = status
}
