package logs_test

import (
	"ACS/iDirectStat/logs"

	"time"

	"testing"
)

// ==============================================================================================================================
//                                      Ship
// ==============================================================================================================================

func TestConsole(t *testing.T) {
	l := new(logs.LogString)

	logs.Init(true, true)

	s1 := "hi there"
	s2 := "bye now"
	l.AddS("TTitle\n")
	l.AddF("in quotes %q  not in quotes: %s\n", s1, s2)
	l.BCon(60)

	for i := 1001; i < 1021; i++ {
		l = new(logs.LogString)
		l.AddF("TTitle %d\n", i-1000)
		l.AddF("And here we are at line: %d\n", i)
		l.AddS("line 2\n")
		l.BCon(60)
	}

	n1 := "James"
	n2 := "Romo"
	logs.LogDebug("This is a test by %q and %q\n", n1, n2)

	time.Sleep(3 * time.Second)
}
