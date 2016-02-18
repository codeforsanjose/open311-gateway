package telemetry

import (
	"fmt"
	"net"
	"time"

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
	chTQue      chan msgSender
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
func SendEngRequest(msgID int64, rType, status, areaID string, at time.Time) {
	statusMsg := EngRequestMsgType{
		ID:     fmt.Sprintf("%d", msgID),
		Rtype:  rType,
		Status: status,
		At:     at,
		AreaID: areaID,
	}
	chTQue <- msgSender(statusMsg)

}

// SendEngRPC sends an Adapter RPC status message to the monitor.
func SendEngRPC(RqstID, RPCID int64, status, route string, at time.Time) {
	statusMsg := EngRPCMsgType{
		ID:     fmt.Sprintf("%d-%d", RqstID, RPCID),
		Status: status,
		Route:  route,
		At:     at,
	}
	chTQue <- msgSender(statusMsg)

}

// Shutdown should be called to gracefully stop the telemetry processes.
func Shutdown() {
	close(chTQue)
}

func init() {
	chTQue = make(chan msgSender, 100)

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
			msg, err := m.Marshal()
			if err != nil {
				log.Warning("unable to send message - %s", err.Error())
				continue
			}
			log.Debug(string(msg))
			if _, err := conn.Write(msg); err != nil {
				log.Warning(err.Error())
			}
		}
	}()

	//
}
