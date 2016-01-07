package main

// considered harmful

import (
	"Gateway311/gateway/common"
	"errors"
	"fmt"
	"log"
	"net/rpc"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	// Synchronous call
	args = &Args{99, 11}
	var quo Quotient
	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%v rmd: %v\n", args.A, args.B, quo.Quo, quo.Rem)

	// Synchronous call
	args = &Args{99, 10}
	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%v rmd: %v\n", args.A, args.B, quo.Quo, quo.Rem)

	creq := NCreateRequest{
		API: API{
			APIAuthKey:        "a01234567890z",
			APIRequestType:    "CreateThreeOneOne",
			APIRequestVersion: "1.0",
		},
		Type:        "Illegal Dumping / Trash",
		TypeID:      22,
		DeviceType:  "iPhone",
		DeviceModel: "6",
		DeviceID:    "test1",
		Latitude:    37.151198,
		Longitude:   -121.602594,
		FirstName:   "James",
		LastName:    "Haskell",
		Email:       "jameskhaskell@gmail.com",
		Phone:       "4445556666",
		IsAnonymous: false,
		Description: "There's an tattered zebra stuck on the side of the  building.",
	}
	fmt.Printf("%+v\n", creq)
	var cresp NCreateResponse
	replyCall := client.Go("Create.Run", &creq, &cresp, nil)
	answer := <-replyCall.Done
	fmt.Println(spew.Sdump(replyCall))
	if replyCall.Error != nil {
		log.Fatal("[Create] error: ", err)
	}
	fmt.Printf("Message: %s\n%v\n", answer.Reply.(*NCreateResponse).Message, answer.Reply.(*NCreateResponse))
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

// API contains the information required by the Backend to process a transation - e.g. the
// API authorization key, API call, etc.
type API struct {
	APIAuthKey        string `json:"ApiAuthKey" xml:"ApiAuthKey"`
	APIRequestType    string `json:"ApiRequestType" xml:"ApiRequestType"`
	APIRequestVersion string `json:"ApiRequestVersion" xml:"ApiRequestVersion"`
}

// NCreateRequest is used to create a new Report.  It is the "native" format of the
// data, and is used by the Engine and all backend Adapters.
type NCreateRequest struct {
	API
	TypeID      int
	Type        string
	DeviceType  string
	DeviceModel string
	DeviceID    string
	Latitude    float64
	Longitude   float64
	Address     string
	City        string
	State       string
	Zip         string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	IsAnonymous bool
	Description string
}

// Displays the contents of the Spec_Type custom type.
func (c NCreateRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - Create\n")
	ls.AddF("Device - type %s  model: %s  ID: %s\n", c.DeviceType, c.DeviceModel, c.DeviceID)
	ls.AddF("Request - type: (%v) %q\n", c.TypeID, c.Type)
	ls.AddF("Location - lat: %v lon: %v\n", c.Latitude, c.Longitude)
	ls.AddF("          %s, %s   %s\n", c.City, c.State, c.Zip)
	ls.AddF("Description: %q\n", c.Description)
	ls.AddF("Author(anon: %t) %s %s  Email: %s  Tel: %s\n", c.IsAnonymous, c.FirstName, c.LastName, c.Email, c.Phone)
	return ls.Box(80)
}

// NCreateResponse is the response to creating or updating a report.
type NCreateResponse struct {
	Message  string
	ID       string
	AuthorID string
}

// Displays the contents of the Spec_Type custom type.
func (c NCreateResponse) String() string {
	ls := new(common.LogString)
	ls.AddS("Report - Resp\n")
	ls.AddF("Message: %s\n", c.Message)
	ls.AddF("ID: %v  AuthorID: %v\n", c.ID, c.AuthorID)
	return ls.Box(80)
}
