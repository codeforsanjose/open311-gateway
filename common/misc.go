package common

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"time"

	// "github.com/davecgh/go-spew/spew"
)

// ==============================================================================================================================
//                                      TIMING
// ==============================================================================================================================
var programStartTime = time.Now()

// ProgramElapsedTime returns the time elapsed since the program started.
func ProgramElapsedTime() float64 {
	return SetPrecision(time.Since(programStartTime).Seconds(), 2)
}

// TimeoutChan returns a "chan bool" that will receive a "true" value after the
// specifed dur`ation t.
func TimeoutChan(t time.Duration) chan bool {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(t)
		timeout <- true
	}()
	return timeout
}

// ==============================================================================================================================
//                                      MATH
// ==============================================================================================================================

// Round returns an int of the rounded value of the float input.
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// SetPrecision effectively rounds a floating point number to the specified precision.
func SetPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

// ==============================================================================================================================
//                                      STRINGS
// ==============================================================================================================================

// BoxIt creates a box around the input string with the specified width.
func BoxIt(str []string, width int) string {
	out := "\n"
	for i, ln := range str {
		if i == 0 {
			x := ((width - len(ln)) / 2) - 1
			out += fmt.Sprintf("\u2554%s %s %s\n", strings.Repeat("\u2550", x), ln, strings.Repeat("\u2550", x))
		} else {
			out += fmt.Sprintf("\u2551%s\n", strings.Replace(ln, "\n", "\n\u2551", -1))
		}
	}
	out += fmt.Sprintf("\u255A%s\n\n", strings.Repeat("\u2550", width))
	return out
}

// ==============================================================================================================================
//                                      TIMESTAMP
// ==============================================================================================================================

// UnixTimestampType reqresents a Unix timestamp.
type UnixTimestampType int64

// SetCurrentTime sets the UnixTimestampType variable to the current time.
func (r *UnixTimestampType) SetCurrentTime() {
	*r = UnixTimestampType(time.Now().Unix())
}

// String returns a string represenation of a UnixTimestampType custom type.
func (r UnixTimestampType) String() string {
	if r == 0 {
		return "***"
	}
	return fmt.Sprintf("%v", time.Unix(int64(r), 0))
}

// ==============================================================================================================================
//                                      BYTES
// ==============================================================================================================================

// ByteToString converts a byte array to a string.
func ByteToString(input []byte, length int) string {
	// fmt.Printf("length: %v\n", length)
	if length <= 0 {
		// Scan for a zero terminator
		zeroPos := bytes.IndexByte(input, 0)
		if zeroPos > 0 {
			length = zeroPos
		} else {
			length = len(input)
		}
	}
	// fmt.Printf("adjusted length: %v\n", length)
	if length == 0 {
		return ""
	}
	return string(input[:length])
}
