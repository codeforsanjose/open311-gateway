package display

import (
	"fmt"
	"time"
)

type engStatusType struct {
	name       string
	lastUpdate time.Time
	status     string
	addr       string
	adapters   string
	rqstCount  int64
}

const (
	esName int = 1 + iota
	esStatus
	esAdapters
	esAddr
	esLength
)

func newEngStatusType(m message) (dataInterface, error) {
	if m.mType != msgTypeES {
		return nil, fmt.Errorf("invalid message type: %q sent to System Update - message: %v", m.mType, m)
	}
	if !m.valid() {
		return nil, fmt.Errorf("invalid message: %#v", m)
	}
	name := m.data[esName]
	status := m.data[esStatus]
	adapters := m.data[esAdapters]
	addr := m.data[esAddr]
	// log.Debug("Name: %q  status: %q  adapters: %q  addr: %q", name, status, adapters, addr)

	d := &engStatusType{
		name:       name,
		lastUpdate: time.Now(),
		status:     status,
		adapters:   adapters,
		addr:       addr,
	}
	// log.Debug("d: %#v", d)
	return dataInterface(d), nil
}

func (r engStatusType) display() string {
	return fmt.Sprintf("%-10s  %10s %8.1f  %s", r.name, r.status, time.Since(r.lastUpdate).Seconds(), r.addr)
}

func (r *engStatusType) update(m message) error {
	// log.Debug(m.String())
	if m.mType != msgTypeES {
		return fmt.Errorf("invalid message type: %q sent to System Update - message: %v", m.mType, m)
	}
	if !m.valid() {
		return fmt.Errorf("invalid message: %#v", m)
	}

	r.name = m.data[esName]
	r.lastUpdate = time.Now()
	r.status = m.data[esStatus]
	if r.addr > "" {
		r.addr = m.data[esAddr]
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
