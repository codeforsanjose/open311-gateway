package router

import (
	"sync/atomic"
)

var rqstID sidType

func init() {
	rqstID = 1000
}

// GetSID returns the next serial RequestID.
func GetSID() int64 {
	return rqstID.get()
}

type sidType int64

func (r *sidType) get() int64 {
	return atomic.AddInt64((*int64)(r), 1)
}
