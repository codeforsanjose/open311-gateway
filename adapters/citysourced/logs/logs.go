package logs

import (
	"os"

	"github.com/op/go-logging"
)

// ==============================================================================================================================
//                                      LOGS
// ==============================================================================================================================
var (
	modulename = "citysourced"
	Log        = logging.MustGetLogger(modulename)
	// LogPrinter *logPrinter
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

	// LogPrinter = newLogPrinter()
	// go LogPrinter.run()
}

// Password is go-logging type for redacting password in logs.
type Password string

// Redacted is used to hide sensitive information, such as password.
func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}
