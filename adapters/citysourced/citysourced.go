package main

// considered harmful

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"

	"Gateway311/adapters/citysourced/data"
	"Gateway311/adapters/citysourced/logs"
	"Gateway311/adapters/citysourced/request"

	// "github.com/davecgh/go-spew/spew"
)

var (
	log = logs.Log
	// Debug switches on some debugging statements.
	Debug = false
)

func main() {

	rpc.Register(&request.Create{})

	rpc.Register(&request.Service{})

	arith := new(Arith)
	rpc.Register(arith)

	rpc.HandleHTTP()
	_, _, addr := data.Adapter()
	log.Info("Listening at: %s\n", addr)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

// Args ...
type Args struct {
	A, B int
}

// Quotient ...
type Quotient struct {
	Quo, Rem int
}

// Arith ...
type Arith int

// Multiply ...
func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

// Divide ...
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
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
