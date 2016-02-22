package telemetry

import "fmt"

const (
	debugMaxLen = 30
)

var debugList []string

func init() {
	debugList = make([]string, 0)
}

func DebugMsg(f string, msg ...interface{}) {
	s := fmt.Sprintf(f, msg...)
	// log.Debug(s)
	if len(debugList) >= debugMaxLen {
		debugList = append(debugList[:debugMaxLen], s)
	}
	debugList = append(debugList, s)
}

func DebugListLast(n int) []string {
	if len(debugList) <= n {
		return debugList
	}
	return debugList[(len(debugList) - n):]
}

func DebugClear() {
	debugList = make([]string, 0)
}
