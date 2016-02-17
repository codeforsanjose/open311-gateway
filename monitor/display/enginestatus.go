package display

import (
	"Gateway311/monitor/comm"
	"fmt"
	"time"
)

type engStatusType struct {
	name       string
	lastUpdate time.Time
	status     string
	adapters   string
	addr       string
	rqstCount  int64
}

func newEngStatusType(m comm.Message) (dataInterface, error) {
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

func (r *engStatusType) update(m comm.Message) error {
	s, err := comm.UnmarshalEngStatusMsg(m)
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
