package telemetry

import (
	"fmt"
	"net"

	"Gateway311/engine/logs"
)

const (
	MIMsgID int = iota
	MISysID
	MIOp
	MIDest
	MIStatus
)

var (
	log         = logs.Log
	chTQue      chan status
	monitorAddr = "127.0.0.1:5051"
)

type status struct {
	MsgID  int64
	SysID  string
	Op     string
	Dest   string
	Status string
}

// Send queues a status message onto the send channel.
func Send(MsgID int64, SysID, Op, Dest, Status string) {
	statusMsg := status{
		MsgID:  MsgID,
		SysID:  SysID,
		Op:     Op,
		Dest:   Dest,
		Status: Status,
	}
	chTQue <- statusMsg

}

// Shutdown should be called to gracefully stop the telemetry processes.
func Shutdown() {
	close(chTQue)
}

func init() {
	chTQue = make(chan status, 100)

	tlmtryServer, err := net.ResolveUDPAddr("udp", monitorAddr)
	if err != nil {
		log.Errorf("Cannot start telemetry - %s", err.Error())
		return
	}

	conn, err := net.DialUDP("udp", nil, tlmtryServer)
	if err != nil {
		log.Errorf("Cannot start telemetry - %s", err.Error())
		return
	}

	go func() {
		log.Debug("Telemetry sender starting...")
		defer conn.Close()
		for m := range chTQue {
			msg := fmt.Sprintf("%d|%s|%s|%s|%s", m.MsgID, m.SysID, m.Op, m.Dest, m.Status)
			log.Debug(msg)
			buf := []byte(msg)
			_, err := conn.Write(buf)
			if err != nil {
				log.Warning(err.Error())
			}
		}
	}()

	//
}
