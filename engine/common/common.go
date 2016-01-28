package common

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fatih/color"

	// "github.com/davecgh/go-spew/spew"
)

// ==============================================================================================================================
//                                      CONSOLE
// ==============================================================================================================================

// NewLogString creates a new LogString, and initializes color printing.
func NewLogString() *LogString {
	ls := new(LogString)
	ls.color = make(map[string]func(...interface{}) string)
	ls.color["red"] = color.New(color.FgRed).SprintFunc()
	ls.color["green"] = color.New(color.FgGreen).SprintFunc()
	ls.color["blue"] = color.New(color.FgBlue).SprintFunc()
	ls.color["yellow"] = color.New(color.FgYellow).SprintFunc()
	return ls
}

// LogString is used to "box" object representations.
type LogString struct {
	raw   string
	fmt   string
	color map[string]func(...interface{}) string
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

// Color applies the specified color to the string.
func (l *LogString) Color(s, color string) string {
	f, ok := l.color[color]
	if !ok {
		return s
	}
	return f(s)
}

// Color applies the specified color to the string.
func (l *LogString) ColorBool(v bool, strue, sfalse, ctrue, cfalse string) string {
	if v {
		return l.Color(strue, ctrue)
	} else {
		return l.Color(sfalse, cfalse)
	}
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
