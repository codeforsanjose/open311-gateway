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
	modulename string = "CSSimulator"
	Log               = logging.MustGetLogger(modulename)
	LogPrinter *logPrinter
)

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

	LogPrinter = NewLogPrinter()
	go LogPrinter.run()
}

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

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

func (l *LogString) BCon(w int) {
	LogPrinter.Con(l.Box(w))
}

func (l *LogString) Raw() string {
	return l.raw
}

type logPrinter struct {
	todo chan string
}

func NewLogPrinter() *logPrinter {
	Log.Debug("NewLogPrinter()... ")
	l := new(logPrinter)
	l.todo = make(chan string, 100)
	return l
}

func (l *logPrinter) Con(s string) {
	l.todo <- s
}

func (l *logPrinter) run() {
	Log.Debug("logPrinter.run()... ")
	for msg := range l.todo {
		fmt.Println(msg)
	}
}
