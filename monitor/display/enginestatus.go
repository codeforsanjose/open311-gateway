package display

import (
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

const (
	esName int = 1 + iota
	esStatus
	esAdapters
	esAddr
	esLength
)

func newEngStatusType(m message) (dataInterface, error) {
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

func (r *engStatusType) update(m message) error {
	s, err := unmarshalEngStatusMsg(m)
	if err != nil {
		return err
	}

	r.name = s.name
	r.lastUpdate = time.Now()
	r.status = s.status
	if s.addr > "" {
		r.addr = s.addr
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

// -------------------------------------------- message --------------------------------------------------------------------

type engStatusMsgType struct {
	name     string
	status   string
	adapters string
	addr     string
}

func unmarshalEngStatusMsg(m message) (*engStatusMsgType, error) {
	if m.mType != msgTypeES {
		return &engStatusMsgType{}, fmt.Errorf("invalid message type: %q sent to EngineStatus - message: %v", m.mType, m)
	}
	if !m.valid() {
		return &engStatusMsgType{}, fmt.Errorf("invalid message: %#v", m)
	}

	return &engStatusMsgType{
		name:     m.data[esName],
		status:   m.data[esStatus],
		adapters: m.data[esAdapters],
		addr:     m.data[esAddr],
	}, nil
}

func (r engStatusMsgType) marshal() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s", msgTypeES, msgDelimiter, r.name, msgDelimiter, r.status, msgDelimiter, r.addr, msgDelimiter, r.adapters)), nil
}
