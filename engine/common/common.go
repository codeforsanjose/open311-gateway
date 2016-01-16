package common

import (
	"fmt"
	"math"
	"strings"
	"time"

	// "github.com/davecgh/go-spew/spew"
)

// ==============================================================================================================================
//                                      CONSOLE
// ==============================================================================================================================

// LogString is used to "box" object representations.
type LogString struct {
	raw string
	fmt string
}

// AddF adds a formated line of text, like Printf().
func (l *LogString) AddF(format string, args ...interface{}) {
	l.raw = l.raw + fmt.Sprintf(format, args...)
}

// AddS adds a single line of text, with no terminating line return.
func (l *LogString) AddS(s string) {
	l.raw = l.raw + s
}

// AddSR adds a single line of text (with line return), like Println().
func (l *LogString) AddSR(s string) {
	l.raw = l.raw + s + "\n"
}

// Box draws a box around the LogString with the specified line width, with a leading line return.
func (l *LogString) Box(w int) string {
	return l.box(w, true)
}

// BoxC draws a box around the LogString with the specified line width, without a leading line return.
func (l *LogString) BoxC(w int) string {
	return l.box(w, false)
}

// box draws a box around the LogString with the specified line width, and leading line return.
func (l *LogString) box(w int, lr bool) string {
	var out string
	if lr {
		out = "\n"
	}
	ss := strings.Split(l.raw, "\n")
	ls := len(ss)
	for i, ln := range ss {
		if i == 0 {
			x := ((w - len(ln)) / 2) - 1
			out += fmt.Sprintf("\u2554%s %s %s\n", strings.Repeat("\u2550", x), ln, strings.Repeat("\u2550", x))
		} else if i == (ls-1) && len(ln) == 0 {
			continue
		} else {
			out += fmt.Sprintf("\u2551%s\n", strings.Replace(ln, "\n", "\n\u2551", -1))
		}
	}
	out += fmt.Sprintf("\u255A%s\n", strings.Repeat("\u2550", w))
	l.fmt = out
	return l.fmt
}

/*
// Raw retrieves the unprocessed LogString.  This is all of the strings to be printerd,
// separated by "\n".
func (l *LogString) Raw() string {
	return l.raw
}

// BCon sends a LogString to the log printer queue (see l.Con() and l.run() below).
func (l *LogString) BCon(w int) {
	LogPrinter.con(l.Box(w))
}

// logPrinter is a string channel that LogStrings can be sent to using the logPrinter.con() method.
// The LogPrinter go routine will receive the strings and print them.
type logPrinter struct {
	todo chan string
}

// newLogPrinter creates a new logPrinter and the associated job channel.  If a queued
// log print is to be used, this must be called to create the logPrinter, followed by
// a call to logPrinter.run() to start the printing go routine.
func newLogPrinter() *logPrinter {
	// Log.Debug("newLogPrinter()... ")
	l := new(logPrinter)
	l.todo = make(chan string, 100)
	return l
}

// con sends a string to the logPrinter.
func (l *logPrinter) con(s string) {
	l.todo <- s
}

// run must be called
func (l *logPrinter) run() {
	// Log.Debug("logPrinter.run()... ")
	for msg := range l.todo {
		fmt.Println(msg)
	}
}
*/

// ==============================================================================================================================
//                                      TIMING
// ==============================================================================================================================
var programStartTime = time.Now()

func ProgramElapsedTime() float64 {
	return ToFixed(time.Since(programStartTime).Seconds(), 2)
}

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
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

// ==============================================================================================================================
//                                      STRINGS
// ==============================================================================================================================
func VarNameToGo(varName, EorP string) string {
	ss := strings.Split(varName, "_")
	for i := range ss {
		if !(EorP == "private" && i == 0) {
			ss[i] = strings.Title(ss[i])
		}
	}
	newVarName := strings.Join(ss, "")
	return newVarName
}

func BoxIt(str []string, w int) string {
	out := "\n"
	for i, ln := range str {
		if i == 0 {
			x := ((w - len(ln)) / 2) - 1
			out += fmt.Sprintf("\u2554%s %s %s\n", strings.Repeat("\u2550", x), ln, strings.Repeat("\u2550", x))
		} else {
			out += fmt.Sprintf("\u2551%s\n", strings.Replace(ln, "\n", "\n\u2551", -1))
		}
	}
	out += fmt.Sprintf("\u255A%s\n\n", strings.Repeat("\u2550", w))
	return out
}

// ==============================================================================================================================
//                                      INTERFACE TYPE
// ==============================================================================================================================
func GetInterfacePtrType(interfacePtr interface{}) string {
	if _, ok := interfacePtr.(*int); ok {
		return "int"
	}
	if _, ok := interfacePtr.(*int32); ok {
		return "int32"
	}
	if _, ok := interfacePtr.(*int64); ok {
		return "int64"
	}

	if _, ok := interfacePtr.(*float32); ok {
		return "float32"
	}
	if _, ok := interfacePtr.(*float64); ok {
		return "float64"
	}

	if _, ok := interfacePtr.(*string); ok {
		return "string"
	}
	return "unknown"
}

// ==============================================================================================================================
//                                      TIMESTAMP
// ==============================================================================================================================
type UnixTimestamp_type int64

func (this *UnixTimestamp_type) SetCurrentTime() {
	*this = UnixTimestamp_type(time.Now().Unix())
}

func (this UnixTimestamp_type) String() string {
	if this == 0 {
		return "***"
	} else {
		return fmt.Sprintf("%v", time.Unix(int64(this), 0))
	}
}
