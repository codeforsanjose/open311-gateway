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
	engRequest := new(engRequestType)
	err := engRequest.update(m)
	if err != nil {
		return nil, err
	}
	return dataInterface(engRequest), nil
}

func (r engRequestType) display() string {
	return fmt.Sprintf("%-10s  %-25s  %-12s %8.1f", r.id, fmt.Sprintf("%s (%s)", r.rType, r.areaID), r.status, time.Since(r.at).Seconds())
}

func (r *engRequestType) update(m message) error {
	s, err := unmarshalEngRequestMsg(m)
	if err != nil {
		return err
	}
	r.id = s.id
	r.rType = s.rType
	r.status = s.status
	r.at = s.at
	r.areaID = s.areaID
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

// -------------------------------------------- message --------------------------------------------------------------------

type engRequestMsgType struct {
	id     string
	rType  string
	status string
	at     time.Time
	areaID string
}

func unmarshalEngRequestMsg(m message) (*engRequestMsgType, error) {
	if m.mType != msgTypeER {
		return &engRequestMsgType{}, fmt.Errorf("invalid message type: %q sent to EngineRequest - message: %v", m.mType, m)
	}
	if !m.valid() {
		return &engRequestMsgType{}, fmt.Errorf("invalid message: %#v", m)
	}

	s := engRequestMsgType{
		id:     m.data[erID],
		rType:  m.data[erRqstType],
		status: m.data[erStatus],
		areaID: m.data[erAreaID],
	}
	if at, err := time.Parse(time.RFC3339Nano, m.data[erAt]); err == nil {
		s.at = at
	} else {
		s.at = time.Now()
	}
	return &s, nil
}

func (r engRequestMsgType) marshal() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s", msgTypeER, msgDelimiter, r.id, msgDelimiter, r.rType, msgDelimiter, r.status, msgDelimiter, r.at.Format(time.RFC3339), msgDelimiter, r.areaID)), nil
}
