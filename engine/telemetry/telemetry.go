package telemetry

import (
	"fmt"
	"net"
	"time"

	log "github.com/jeffizhungry/logrus"
)

var (
	chTQue      chan msgSender
	monitorAddr = "127.0.0.1:5051"
)

// SendTelemetry sends a telemetry message.
func SendTelemetry(rqstID int64, op, status string) {
	SendRequest(rqstID, op, status, "", time.Now())
}

// SendRequest queues an Engine REST Request message onto the send channel.
func SendRequest(msgID int64, rType, status, areaID string, at time.Time) {
	statusMsg := EngRequestMsgType{
		ID:     fmt.Sprintf("%d", msgID),
		Rtype:  rType,
		Status: status,
		At:     at,
		AreaID: areaID,
	}
	chTQue <- msgSender(statusMsg)

}

// SendRPC sends an Adapter RPC status message to the monitor.
func SendRPC(id, status, route string, at time.Time) {
	statusMsg := EngRPCMsgType{
		ID:     id,
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
		// log.Debug("Telemetry sender starting...")
		defer conn.Close()
		for m := range chTQue {
			msg, err := m.Marshal()
			if err != nil {
				log.Warning("unable to send message - %s", err.Error())
				continue
			}
			// log.Debug(string(msg))
			if _, err := conn.Write(msg); err != nil {
				log.Warning(err.Error())
			}
		}
	}()
}
