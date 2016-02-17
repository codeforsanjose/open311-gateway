package display

import "time"

type adpStatusType struct {
	name       string
	status     string
	lastUpdate time.Time
	addr       string
}
