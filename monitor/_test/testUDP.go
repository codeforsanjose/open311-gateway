package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	Debug = true
	addr  string
)

func main() {
	checkErr := func(err error) {
		if err != nil {
			fmt.Printf("Shutting down - %s\n", err.Error())
			stop()
			os.Exit(-1)
		}
	}

	fmt.Printf("Starting test...\n")

	ServerAddr, err := net.ResolveUDPAddr("udp", addr)
	checkErr(err)

	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	checkErr(err)

	defer Conn.Close()

	i := 0
	for {
		msg := strconv.Itoa(i)
		fmt.Printf("%d >>\n", i)
		i++
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}

func init() {
	flag.BoolVar(&Debug, "debug", false, "Activates debug logging.")
	flag.StringVar(&addr, "addr", "127.0.0.1:5050", "The address the monitor will listen on.")
	flag.Parse()

	fmt.Printf("Debug: %v\n", Debug)
	fmt.Printf("Address: %v\n", addr)

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
