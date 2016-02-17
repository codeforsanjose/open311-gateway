package display

import "time"

type adpRequestType struct {
	id     string
	status string
	recvAt time.Time
	route  string
}
