package telemetry

import (
	"fmt"
	"net"

	"github.com/open311-gateway/monitor/logs"
)

var (
	log         = logs.Log
	msgChan     chan Message
	done        chan bool
	monitorConn *net.UDPConn
	monitorAddr string
)

func init() {
	msgChan = make(chan Message, 1000)
	done = make(chan bool)
}

// Start starts the telemetry system.
func Start() {

	if err := StartReceiver(monitorAddr, msgChan, done); err != nil {
		Shutdown()
		log.Fatalf("Unable to start receiver - %s", err.Error())
	}
}

// Shutdown ensures an orderly finalization of active processes.
func Shutdown() {
	done <- true
	monitorConn.Close()
}

// SetAddr sets the receiver address.
func SetAddr(addr string) {
	monitorAddr = addr
}

// GetMsgChan returns the message queue channel.
func GetMsgChan() chan Message {
	return msgChan
}

// ==============================================================================================================================
//                                      RECEIVER
// ==============================================================================================================================

// StartReceiver starts the UDP receive process.  The bytes received are parsed
// by the separator character "|" into a slice of strings, and put on the msgChan.
func StartReceiver(addr string, msgChan chan Message, done <-chan bool) error {
	// log.Debug("Address: %v", addr)
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("error resolving address - %s", err.Error())
	}
	monitorConn, _ = net.ListenUDP("udp", a)
	monitorConn.SetReadBuffer(1048576)

	// Receive UDP
	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := monitorConn.ReadFromUDP(buf)
			if err != nil {
				log.Fatalf("Error receiving UDP: %s", err.Error())
			}
			// log.Debug("Received %d bytes: %q", n, string(buf[0:n]))
			select {
			case <-done:
				return
			default:
			}

			if n > 0 {
				msg, err := NewMessage(buf, n)
				if err == nil {
					msgChan <- msg
				} else {
					log.Errorf("invalid message: [[%v]]", buf)
				}
			}
		}
	}()

	return nil
}
