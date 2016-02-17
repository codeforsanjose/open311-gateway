package display

import (
	"Gateway311/monitor/comm"
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

func newEngRequestType(m comm.Message) (dataInterface, error) {
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

func (r *engRequestType) update(m comm.Message) error {
	s, err := comm.UnmarshalEngRequestMsg(m)
	if err != nil {
		return err
	}
	r.id = s.ID
	r.rType = s.Rtype
	r.status = s.Status
	r.at = s.At
	r.areaID = s.AreaID
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
