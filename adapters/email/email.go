package main

// considered harmful

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"

	"github.com/open311-gateway/adapters/email/data"
	"github.com/open311-gateway/adapters/email/request"
	"github.com/open311-gateway/adapters/email/telemetry"

	log "github.com/jeffizhungry/logrus"
	// "github.com/davecgh/go-spew/spew"
)

var (
	// Debug switches on some debugging statements.
	Debug      = false
	configFile string
)

func main() {

	log.Setup(false, log.DebugLevel)
	log.Debugf("Command line settings - debug: %t\nConfig file: %q", Debug, configFile)

	if err := data.Init(configFile); err != nil {
		log.Fatal("Unable to start - data initilization failed.\n")
	}
	telemetry.Init(data.GetMonitorAddress())

	rpc.Register(&request.Report{})

	rpc.Register(&request.Services{})

	rpc.HandleHTTP()
	_, _, addr := data.Adapter()
	log.Infof("Listening at: %s\n", addr)

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.Serve(l, nil)
}

func init() {
	flag.BoolVar(&Debug, "debug", false, "Activates debug logging. It is active if either this or the value in 'config.json' are set.")
	flag.StringVar(&configFile, "config", "config.json", "Config file. This is a full or relative path.")
	flag.Parse()

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
