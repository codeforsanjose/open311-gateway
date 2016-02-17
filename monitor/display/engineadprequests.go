package display

import (
	"fmt"
	"time"
)

type engAdpRequestType struct {
	id     string
	status string
	route  string
	at     time.Time
}

const (
	eaID int = 1 + iota
	eaStatus
	eaRoute
	eaAt
	eaLength
)

func newEngAdpRequestType(m message) (dataInterface, error) {
	if m.mType != msgTypeEA {
		return nil, fmt.Errorf("invalid message type: %q sent to Eng Adp Request - message: %v", m.mType, m)
	}
	if !m.valid() {
		return nil, fmt.Errorf("invalid message: %#v", m)
	}
	id := m.data[eaID]
	status := m.data[eaStatus]
	at, err := time.Parse(time.RFC3339Nano, m.data[eaAt])
	if err != nil {
		at = time.Now()
	}
	route := m.data[eaRoute]
	// log.Debug("Name: %q  status: %q  adapters: %q  addr: %q", name, status, adapters, addr)

	d := &engAdpRequestType{
		id:     id,
		status: status,
		at:     at,
		route:  route,
	}
	// log.Debug("d: %#v", d)
	return dataInterface(d), nil
}

func (r engAdpRequestType) display() string {
	return fmt.Sprintf("%-10s  %-25s  %-12s %8.1f", r.id, r.route, r.status, time.Since(r.at).Seconds())
}

func (r *engAdpRequestType) update(m message) error {
	// log.Debug(m.String())
	if m.mType != msgTypeES {
		return fmt.Errorf("invalid message type: %q sent to System Update - message: %v", m.mType, m)
	}
	if !m.valid() {
		return fmt.Errorf("invalid message: %#v", m)
	}

	r.id = m.data[eaID]
	r.route = m.data[eaRoute]
	r.status = m.data[eaStatus]
	at, err := time.Parse(time.RFC3339Nano, m.data[eaAt])
	if err != nil {
		at = time.Now()
	}
	r.at = at
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
