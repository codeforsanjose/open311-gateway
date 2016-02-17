package display

import (
	"_sketches/spew"
	"fmt"
	"net"
)

// ==============================================================================================================================
//                                      DATA
// ==============================================================================================================================

// StartReceiver starts the UDP receive process.  The bytes received are parsed
// by the separator character "|" into a slice of strings, and put on the msgChan.
func StartReceiver(addr string, msgChan chan<- message, done <-chan bool) error {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("error resolving address - %s", err.Error())
	}
	fmt.Printf("Address: %s\n", spew.Sdump(a))
	MonitorConn, _ := net.ListenUDP("udp", a)
	defer MonitorConn.Close()
	MonitorConn.SetReadBuffer(1048576)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := MonitorConn.ReadFromUDP(buf)
			if err != nil {
				log.Errorf("Error receiving UDP: %s", err.Error())
			}
			select {
			case <-done:
				return
			default:
			}

			if n > 0 {
				msg, err := newMessage(buf, n)
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
