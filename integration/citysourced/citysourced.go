package main

// considered harmful

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"Gateway311/integration/citysourced/data"
	"Gateway311/integration/citysourced/request"

	// "github.com/davecgh/go-spew/spew"
)

func main() {

	fmt.Println(data.ShowProviderData())

	rpc.Register(&request.Create{})

	rpc.Register(&request.Service{})

	arith := new(Arith)
	rpc.Register(arith)

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
