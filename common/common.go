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

type LogString struct {
	raw string
	fmt string
}

func (l *LogString) AddF(format string, args ...interface{}) {
	l.raw = l.raw + fmt.Sprintf(format, args...)
}

func (l *LogString) AddS(s string) {
	l.raw = l.raw + s
}

func (l *LogString) AddSR(s string) {
	l.raw = l.raw + s + "\n"
}

func (l *LogString) Box(w int) string {
	out := "\n"
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

func (l *LogString) BoxC(w int) string {
	out := ""
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

func (l *LogString) Raw() string {
	return l.raw
}

// ==============================================================================================================================
//                                      TIMING
// ==============================================================================================================================
var programStartTime = time.Now()

func ProgramElapsedTime() float64 {
	return ToFixed(time.Since(programStartTime).Seconds(), 2)
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
