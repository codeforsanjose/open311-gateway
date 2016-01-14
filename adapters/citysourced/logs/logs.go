package logs

import (
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"
)

// ==============================================================================================================================
//                                      LOGS
// ==============================================================================================================================
var (
	modulename = "citysourced"
	Log        = logging.MustGetLogger(modulename)
	LogPrinter *logPrinter
)

// Init configures the logging system.
func Init(debug bool) {
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

// BCon sends a LogString to the log printer queue (see l.Con() and l.run() below).
func (l *LogString) BCon(w int) {
	LogPrinter.con(l.Box(w))
}

// Raw retrieves the unprocessed LogString.  This is all of the strings to be printerd,
// separated by "\n".
func (l *LogString) Raw() string {
	return l.raw
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
	Log.Debug("newLogPrinter()... ")
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
	Log.Debug("logPrinter.run()... ")
	for msg := range l.todo {
		fmt.Println(msg)
	}
}
