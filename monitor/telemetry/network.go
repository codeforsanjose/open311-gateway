package telemetry

import (
	"fmt"
	"net"

	"Gateway311/monitor/logs"

	"github.com/davecgh/go-spew/spew"
)

var (
	log         = logs.Log
	msgChan     chan Message
	done        chan bool
	monitorConn *net.UDPConn
)

const (
	monitorAddr = "127.0.0.1:5051"
)

func init() {
	msgChan = make(chan Message, 1000)
	done = make(chan bool)

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

// ==============================================================================================================================
//                                      DATA
// ==============================================================================================================================

// StartReceiver starts the UDP receive process.  The bytes received are parsed
// by the separator character "|" into a slice of strings, and put on the msgChan.
func StartReceiver(addr string, msgChan chan<- Message, done <-chan bool) error {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("error resolving address - %s", err.Error())
	}
	fmt.Printf("Address: %s\n", spew.Sdump(a))
	monitorConn, _ = net.ListenUDP("udp", a)
	monitorConn.SetReadBuffer(1048576)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := monitorConn.ReadFromUDP(buf)
			if err != nil {
				log.Fatalf("Error receiving UDP: %s", err.Error())
			}
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
