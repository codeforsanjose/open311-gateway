package telemetry

import "fmt"

const (
	debugMaxLen = 30
)

var debugList []string

func init() {
	debugList = make([]string, 0)
}

// DebugMsg adds a new Debug message (a string) onto the Debug Message List.
func DebugMsg(f string, msg ...interface{}) {
	s := fmt.Sprintf(f, msg...)
	// log.Debug(s)
	if len(debugList) >= debugMaxLen {
		debugList = append(debugList[:debugMaxLen], s)
	} else {
		debugList = append(debugList, s)
	}
}

// DebugListLast returns the last N messages from the Debug Message List.
func DebugListLast(n int) []string {
	if len(debugList) <= n {
		return debugList
	}
	return debugList[(len(debugList) - n):]
}

// DebugClear clears all the Debug Message List.
func DebugClear() {
	debugList = make([]string, 0)
}
