package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/open311-gateway/monitor/display"
	"github.com/open311-gateway/monitor/logs"
	"github.com/open311-gateway/monitor/telemetry"
)

var (
	log         = logs.Log
	Debug       = true
	monitorAddr string
)

func main() {
	telemetry.Start()
	display.Start()
}

func init() {
	var monitorAddr string
	flag.BoolVar(&Debug, "debug", false, "Activates debug logging.")
	flag.StringVar(&monitorAddr, "addr", "127.0.0.1:5051", "The address the monitor will listen on.")
	flag.Parse()

	logs.Init(Debug)

	telemetry.SetAddr(monitorAddr)

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
	telemetry.Shutdown()
	return nil
}
