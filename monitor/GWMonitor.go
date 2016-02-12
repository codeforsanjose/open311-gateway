package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"Gateway311/monitor/display"
	"Gateway311/monitor/logs"

	"github.com/davecgh/go-spew/spew"
)

var (
	log         = logs.Log
	Debug       = true
	monitorAddr string
)

const (
	MIMsgID int = iota
	MISysID
	MIOp
	MIDest
	MIStatus
)

type status struct {
	ID        int64
	CreatedAt time.Time
	SysID     string
	Op        string
	Dest      string
}

type fullStatus struct {
	status
	engSent bool
	engRecv bool
	adpRecv bool
	adpSent bool
}

var sMap map[int64]*fullStatus

func main() {
	// RunOld()
	display.RunTest()
}

func RunOld() {
	log.Debug("Here we go...")

	go Display()
	a, err := net.ResolveUDPAddr("udp", monitorAddr)
	if err != nil {
		fmt.Printf("Error resolving address - %s", err.Error())
		stop()
		os.Exit(-1)
	}
	fmt.Printf("Address: %s\n", spew.Sdump(a))
	MonitorConn, _ := net.ListenUDP("udp", a)
	defer MonitorConn.Close()
	MonitorConn.SetReadBuffer(1048576)

	buf := make([]byte, 1024)

	fmt.Printf("MsgID           Source                Type                 Destination\n")
	for {
		n, _, err := MonitorConn.ReadFromUDP(buf)
		msgStr := strings.Split(string(buf[0:n]), "|")
		// fmt.Printf("%-15s %-20s  %-20s %s\n", msgStr[MIMsgID], fmt.Sprintf("[%s: %s]", msgStr[MISysID], msgStr[MIStatus]), msgStr[MIOp], msgStr[MIDest])
		msgID, err := strconv.ParseInt(msgStr[MIMsgID], 10, 64)
		if err != nil || msgID <= 0 {
			log.Error("Invalid message id: %q", msgStr[MIMsgID])
			continue
		}
		_, found := sMap[msgID]
		if !found {
			sMap[msgID] = &fullStatus{
				status: status{
					ID:        msgID,
					CreatedAt: time.Now(),
					SysID:     msgStr[MISysID],
					Op:        msgStr[MIOp],
					Dest:      msgStr[MIDest],
				},
			}
		}

		switch msgStr[MIStatus] {
		case "eng-send":
			sMap[msgID].engSent = true
		case "eng-recv":
			sMap[msgID].engRecv = true
		case "adp-recv":
			sMap[msgID].adpRecv = true
		case "adp-send":
			sMap[msgID].adpSent = true
		}

		displayItem(sMap[msgID])

	}
}

func Display() {

	for {
		// Cleanup
		for k, v := range sMap {
			if v.engSent && v.engRecv && v.adpRecv && v.adpSent {
				fmt.Printf("Completing: %d\n", k)
				delete(sMap, k)
			}
			if time.Since(v.CreatedAt).Seconds() > 10 {
				fmt.Printf("Expired: %d!\n", k)
				delete(sMap, k)
			}
		}

		// Display
		fmt.Printf("----------------------------------------------------------------\n")
		for _, stat := range sMap {
			displayItem(stat)
		}
		fmt.Println()

		time.Sleep(5 * time.Second)
	}
}

func displayItem(stat *fullStatus) {
	fmt.Printf("%-10d %-15s %-6s %5t %5t %5t %5t\n", stat.ID, time.Since(stat.CreatedAt), fmt.Sprintf("[%s]", stat.SysID), stat.engSent, stat.engRecv, stat.adpRecv, stat.adpSent)
}

func init() {
	flag.BoolVar(&Debug, "debug", false, "Activates debug logging.")
	flag.StringVar(&monitorAddr, "addr", "127.0.0.1:5050", "The address the monitor will listen on.")
	flag.Parse()

	logs.Init(Debug)
	sMap = make(map[int64]*fullStatus)

	log.Debug("Debug: %v", Debug)
	log.Debug("Address: %v", monitorAddr)

	go signalHandler(make(chan os.Signal, 1))
	fmt.Println("Press Ctrl-C to shutdown...")
}

func signalHandler(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	for s := <-c; ; s = <-c {
		switch s {
		case os.Interrupt:
			fmt.Println("Ctrl-C Received!")
			stop()
			os.Exit(0)
		case os.Kill:
			fmt.Println("SIGKILL Received!")
			stop()
			os.Exit(1)
		}
	}
}

func stop() error {
	return nil
}
