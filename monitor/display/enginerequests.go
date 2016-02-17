package display

import (
	"fmt"
	"time"
)

type engRequestType struct {
	id     string
	rType  string
	status string
	at     time.Time
	areaID string
}

const (
	erID int = 1 + iota
	erRqstType
	erStatus
	erAt
	erAreaID
	erLength
)

func newEngRequestType(m message) (dataInterface, error) {
	if m.mType != msgTypeER {
		return nil, fmt.Errorf("invalid message type: %q sent to Engine Request - message: %v", m.mType, m)
	}
	if !m.valid() {
		return nil, fmt.Errorf("invalid message: %#v", m)
	}
	id := m.data[erID]
	status := m.data[erStatus]
	at, err := time.Parse(time.RFC3339Nano, m.data[erAt])
	if err != nil {
		at = time.Now()
	}
	rType := m.data[erRqstType]
	areaID := m.data[erAreaID]
	// log.Debug("Name: %q  status: %q  adapters: %q  addr: %q", name, status, adapters, addr)

	d := &engRequestType{
		id:     id,
		status: status,
		at:     at,
		rType:  rType,
		areaID: areaID,
	}
	// log.Debug("d: %#v", d)
	return dataInterface(d), nil
}

func (r engRequestType) display() string {
	return fmt.Sprintf("%-10s  %-25s  %-12s %8.1f", r.id, fmt.Sprintf("%s (%s)", r.rType, r.areaID), r.status, time.Since(r.at).Seconds())
}

func (r *engRequestType) update(m message) error {
	// log.Debug(m.String())
	if m.mType != msgTypeES {
		return fmt.Errorf("invalid message type: %q sent to System Update - message: %v", m.mType, m)
	}
	if !m.valid() {
		return fmt.Errorf("invalid message: %#v", m)
	}

	r.id = m.data[erID]
	r.rType = m.data[erRqstType]
	r.status = m.data[erStatus]
	at, err := time.Parse(time.RFC3339Nano, m.data[erAt])
	if err != nil {
		at = time.Now()
	}
	r.at = at
	r.areaID = m.data[erAreaID]
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
