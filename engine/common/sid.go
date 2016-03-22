package common

import (
	"sync/atomic"
)

var (
	requestID SerialID = 100
	rpcID     SerialID = 1
)

// SerialID creates a serial ID type (int64).
type SerialID int64

// Get retrieves the next Serial ID
func (r *SerialID) Get() int64 {
	return atomic.AddInt64((*int64)(r), 1)
}

// RequestID retrieves the next Request ID.
func RequestID() int64 {
	return requestID.Get()
}

// RPCID retrieves the next RPC ID.
func RPCID() int64 {
	return rpcID.Get()
}
