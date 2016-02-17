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

func (r *engAdpRequestType) update(m message) error {
	s, err := unmarshalAdpEngRequestMsg(m)
	if err != nil {
		return err
	}
	r.id = s.id
	r.status = s.status
	r.route = s.route
	r.at = s.at
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

// -------------------------------------------- message --------------------------------------------------------------------

type engAdpRequestMsgType struct {
	id     string
	status string
	route  string
	at     time.Time
}

func unmarshalAdpEngRequestMsg(m message) (*engAdpRequestMsgType, error) {
	if m.mType != msgTypeEA {
		return &engAdpRequestMsgType{}, fmt.Errorf("invalid message type: %q sent to EngineRequest - message: %v", m.mType, m)
	}
	if !m.valid() {
		return &engAdpRequestMsgType{}, fmt.Errorf("invalid message: %#v", m)
	}

	s := engAdpRequestMsgType{
		id:     m.data[eaID],
		status: m.data[eaStatus],
		route:  m.data[eaRoute],
	}
	if at, err := time.Parse(time.RFC3339Nano, m.data[eaAt]); err == nil {
		s.at = at
	} else {
		s.at = time.Now()
	}
	return &s, nil

}

func (r engAdpRequestMsgType) marshal() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s", msgTypeEA, msgDelimiter, r.id, msgDelimiter, r.status, msgDelimiter, r.route, msgDelimiter, r.at.Format(time.RFC3339))), nil
}
