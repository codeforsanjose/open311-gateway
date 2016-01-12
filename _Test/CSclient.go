package main

// considered harmful

import (
	"Gateway311/gateway/common"
	"_sketches/spew"
	"errors"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Print("dialing:", err)
	}
	// Synchronous call
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Print("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	// Synchronous call
	args = &Args{99, 11}
	var quo Quotient
	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		log.Print("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%v rmd: %v\n", args.A, args.B, quo.Quo, quo.Rem)

	// Synchronous call
	args = &Args{99, 10}
	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		log.Print("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%v rmd: %v\n", args.A, args.B, quo.Quo, quo.Rem)

	// Synchronous call
	args = &Args{99, 0}
	err = client.Call("Arith.Divide", args, &quo)
	if err != nil {
		log.Print("arith error:", err)
	}
	fmt.Printf("Arith: %d/%d=%v rmd: %v\n", args.A, args.B, quo.Quo, quo.Rem)

	// creq := NServiceRequest{
	// 	City: "San Jose",
	// }
	// fmt.Printf("%+v\n", creq)
	// var cresp NServicesResponse
	// err = client.Call("Service.ServicesForCity", &creq, &cresp)
	// if err != nil {
	// 	log.Print("[Services] error: ", err)
	// }
	// fmt.Println(spew.Sdump(cresp))
	{
		creq := NServiceRequest{
			City: "San Jose",
		}
		fmt.Printf("%+v\n", creq)
		var cresp NServicesResponse
		replyCall := client.Go("Service.City", &creq, &cresp, nil)
		answer := <-replyCall.Done
		if replyCall.Error != nil {
			log.Print("[Create] error: ", err)
		}
		fmt.Println(spew.Sdump(answer))
		if answer.Error != nil {
			fmt.Printf("Error on API request: %s\n", answer.Error)
		} else {
			fmt.Printf("Return: %v\n", answer.Reply.(*NServicesResponse))
		}
	}

	{
		creq := NServiceRequest{
			City: "all",
		}
		fmt.Printf("%+v\n", creq)
		var cresp NServicesResponse
		replyCall := client.Go("Service.All", &creq, &cresp, nil)
		answer := <-replyCall.Done
		if replyCall.Error != nil {
			log.Print("[Create] error: ", err)
		}
		fmt.Println(spew.Sdump(answer))
		if answer.Error != nil {
			fmt.Printf("Error on API request: %s\n", answer.Error)
		} else {
			fmt.Printf("Return: %v\n", answer.Reply.(*NServicesResponse))
		}
	}

	// creq := NCreateRequest{
	// 	API: API{
	// 		APIAuthKey:        "a01234567890z",
	// 		APIRequestType:    "CreateThreeOneOne",
	// 		APIRequestVersion: "1.0",
	// 	},
	// 	Type:        "Illegal Dumping / Trash",
	// 	TypeID:      22,
	// 	DeviceType:  "iPhone",
	// 	DeviceModel: "6",
	// 	DeviceID:    "test1",
	// 	Latitude:    37.151198,
	// 	Longitude:   -121.602594,
	// 	FirstName:   "James",
	// 	LastName:    "Haskell",
	// 	Email:       "jameskhaskell@gmail.com",
	// 	Phone:       "4445556666",
	// 	IsAnonymous: false,
	// 	Description: "There's an tattered zebra stuck on the side of the  building.",
	// }
	// fmt.Printf("%+v\n", creq)
	// var cresp NCreateResponse
	// replyCall := client.Go("Create.Run", &creq, &cresp, nil)
	// answer := <-replyCall.Done
	// fmt.Println(spew.Sdump(replyCall))
	// if replyCall.Error != nil {
	// 	log.Print("[Create] error: ", err)
	// }
	// fmt.Printf("Message: %s\n%v\n", answer.Reply.(*NCreateResponse).Message, answer.Reply.(*NCreateResponse))
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

// =======================================================================================
//                                      SERVICES
// =======================================================================================

// NServiceRequest is used to get list of services available to the user.
type NServiceRequest struct {
	City string
}

// Displays the contents of the Spec_Type custom type.
func (c NServiceRequest) String() string {
	ls := new(common.LogString)
	ls.AddS("Services Request\n")
	ls.AddF("Location - city: %v\n", c.City)
	return ls.Box(80)
}

// ------------------------------- Services -------------------------------

// NServicesResponse is the returned struct for a Services request.
type NServicesResponse struct {
	Message  string
	Services NServices
}

// NServices contains a list of Services.
type NServices []NService

// Displays the contents of the Spec_Type custom type.
func (c NServices) String() string {
	ls := new(common.LogString)
	ls.AddS("Services Response\n")
	for _, s := range c {
		ls.AddF("%s\n", s)
	}
	return ls.Box(80)
}

// ------------------------------- Service -------------------------------

// NService represents a Service.  The ID is a combination of the BackEnd Type (IFID),
// the AreaID (i.e. the City id), ProviderID (in case the provider has multiple interfaces),
// and the Service ID.
type NService struct {
	ServiceID  `json:"id"`
	Name       string   `json:"name"`
	Categories []string `json:"catg"`
}

func (s NService) String() string {
	r := fmt.Sprintf("   %s-%s-%d-%d  %-40s  %v", s.IFID, s.AreaID, s.ProviderID, s.ID, s.Name, s.Categories)
	return r
}

// ------------------------------- ServiceID -------------------------------

// ServiceID provides the JSON marshalling conversion between the JSON "ID" and
// the Backend Interface Type, AreaID (City id), ProviderID, and Service ID.
type ServiceID struct {
	IFID       string
	AreaID     string
	ProviderID int
	ID         int
}

// MID creates the string
func (s ServiceID) MID() string {
	return fmt.Sprintf("%s-%s-%d-%d", s.IFID, s.AreaID, s.ProviderID, s.ID)
}
