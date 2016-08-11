package telemetry

import (
	"net"
	"time"

	"github.com/open311-gateway/adapters/email/data"

	log "github.com/jeffizhungry/logrus"
)

var (
	chTQue chan msgSender
)

// SendRPC queues an RPC status message onto the send channel.
func SendRPC(id, status, route, url string, results int, at time.Time) {
	statusMsg := AdpRPCMsgType{
		AdpID:   data.AdapterName(),
		ID:      id,
		Status:  status,
		Route:   route,
		URL:     url,
		Results: results,
		At:      at,
	}
	chTQue <- msgSender(statusMsg)

}

// Shutdown should be called to gracefully stop the telemetry processes.
func Shutdown() {
	close(chTQue)
}

// Init initializes the system monitoring service.
func Init(addr string) {
	chTQue = make(chan msgSender, 100)

	tlmtryServer, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Errorf("Cannot start telemetry - %s", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, tlmtryServer)
	if err != nil {
		log.Errorf("Cannot start telemetry - %s", err)
		return
	}

	go func() {
		log.Debugf("Telemetry sender starting on: %v", addr)
		finish := func() {
			log.Debug("Closing telemetry connection...")
			_ = conn.Close()
		}
		defer finish()
		for m := range chTQue {
			msg, err := m.Marshal()
			if err != nil {
				log.Warning("unable to send message - %s", err)
				continue
			}
			log.Debug(string(msg))
			if _, err := conn.Write(msg); err != nil {
				log.Warning(err.Error())
			}
		}
	}()
}
