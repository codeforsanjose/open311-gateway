package logs

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/op/go-logging"
)

// ==============================================================================================================================
//                                      LOGS
// ==============================================================================================================================
var (
	modulename  string
	Log         = logging.MustGetLogger(modulename)
	LogPrinter  *logPrinter
	initialized bool
)

// Init configures the logging system.
func Init(debug bool) {
	if initialized {
		return
	}
	initialized = true
	var syslogfmtstr, logfmtstr string
	if debug {
		syslogfmtstr = "[%{shortpkg}: %{shortfile}: %{shortfunc}] %{message}"
		logfmtstr = "%{color}%{time:15:04:05} [%{shortpkg}: %{shortfile}: %{shortfunc}()] ▶ %{level:.4s} ◀  %{color:reset} %{message}"
	} else {
		syslogfmtstr = "[%{shortpkg}: %{shortfile}: %{shortfunc}] %{message}"
		logfmtstr = "%{color}%{time:15:04:05} [%{shortpkg}] ▶ %{level:.4s} ◀  %{color:reset} %{message}"
	}
	syslogformat := logging.MustStringFormatter(syslogfmtstr)
	syslog, _ := logging.NewSyslogBackend(modulename)
	syslogF := logging.NewBackendFormatter(syslog, syslogformat)
	syslogL := logging.AddModuleLevel(syslogF)
	syslogL.SetLevel(logging.WARNING, "")

	logformat := logging.MustStringFormatter(logfmtstr)
	console := logging.NewLogBackend(os.Stderr, "", 0)
	consoleF := logging.NewBackendFormatter(console, logformat)
	consoleFLev := logging.AddModuleLevel(consoleF)
	if debug {
		consoleFLev.SetLevel(logging.DEBUG, modulename)
	} else {
		consoleFLev.SetLevel(logging.INFO, modulename)
	}

	logging.SetBackend(syslogL, consoleFLev)

	LogPrinter = newLogPrinter()
	go LogPrinter.run()
}

// Password is go-logging type for redacting password in logs.
type Password string

// Redacted is used to hide sensitive information, such as password.
func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

// ==============================================================================================================================
//                                      CONSOLE
// ==============================================================================================================================

// NewFmtBoxer creates a new FmtBoxer, and initializes color printing.
func NewFmtBoxer() *FmtBoxer {
	ls := new(FmtBoxer)
	ls.color = make(map[string]func(...interface{}) string)
	ls.color["red"] = color.New(color.FgRed).SprintFunc()
	ls.color["green"] = color.New(color.FgGreen).SprintFunc()
	ls.color["blue"] = color.New(color.FgBlue).SprintFunc()
	ls.color["yellow"] = color.New(color.FgYellow).SprintFunc()
	return ls
}

// FmtBoxer is used to "box" object representations.
type FmtBoxer struct {
	raw   string
	fmt   string
	color map[string]func(...interface{}) string
}

// Color applies the specified color to the string.
func (l *FmtBoxer) Color(color, s string) string {
	f, ok := l.color[color]
	if !ok {
		return s
	}
	return f(s)
}

// AddF adds a formated line of text, like Printf().
func (l *FmtBoxer) AddF(format string, args ...interface{}) {
	l.raw = l.raw + fmt.Sprintf(format, args...)
}

// AddS adds a single line of text, with no terminating line return.
func (l *FmtBoxer) AddS(s string) {
	l.raw = l.raw + s
}

// AddSR adds a single line of text (with line return), like Println().
func (l *FmtBoxer) AddSR(s string) {
	l.raw = l.raw + s + "\n"
}

// Box draws a box around the FmtBoxer with the specified line width, with a leading line return.
func (l *FmtBoxer) Box(w int) string {
	return l.box(w, true)
}

// BoxC draws a box around the FmtBoxer with the specified line width, without a leading line return.
func (l *FmtBoxer) BoxC(w int) string {
	return l.box(w, false)
}

// box draws a box around the FmtBoxer with the specified line width, and leading line return.
func (l *FmtBoxer) box(w int, lr bool) string {
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

// BCon sends a FmtBoxer to the log printer queue (see l.Con() and l.run() below).
func (l *FmtBoxer) BCon(w int) {
	LogPrinter.con(l.Box(w))
}

// Raw retrieves the unprocessed FmtBoxer.  This is all of the strings to be printerd,
// separated by "\n".
func (l *FmtBoxer) Raw() string {
	return l.raw
}

// logPrinter is a string channel that FmtBoxers can be sent to using the logPrinter.con() method.
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
