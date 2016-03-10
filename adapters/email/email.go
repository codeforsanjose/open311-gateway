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

	"Gateway311/adapters/email/data"
	"Gateway311/adapters/email/logs"
	"Gateway311/adapters/email/request"

	// "github.com/davecgh/go-spew/spew"
)

var (
	log = logs.Log
	// Debug switches on some debugging statements.
	Debug = false
)

func main() {

	rpc.Register(&request.Report{})

	rpc.Register(&request.Services{})

	rpc.HandleHTTP()
	_, _, addr := data.Adapter()
	log.Info("Listening at: %s\n", addr)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

func init() {
	var configFile string
	flag.BoolVar(&Debug, "debug", false, "Activates debug logging. It is active if either this or the value in 'config.json' are set.")
	flag.StringVar(&configFile, "config", "data/config.json", "Config file. This is a full or relative path.")
	flag.Parse()

	logs.Init(Debug)

	if err := data.Init(configFile); err != nil {
		log.Fatal("Unable to start - data initilization failed.\n")
	}

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
