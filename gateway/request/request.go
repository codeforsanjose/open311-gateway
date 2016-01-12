package request

import (
	"fmt"
	"net/http"
	"net/rpc"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/davecgh/go-spew/spew"
)

// Services looks up the service providers and services for the specified location.
// The URL nust contain query parameters of either:
// LatitudeV and LongitudeV, or a city name.
//
// Examples:
//  http;//xyz.com/api/services?lat=34.236144&lon=-118.604794
//  http;//xyz.com/api/services?city=san+jose
func Services(w rest.ResponseWriter, r *rest.Request) {
	response, err := processServices(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}

// Create creates a new Report and adds it to Reports.
func Create(w rest.ResponseWriter, r *rest.Request) {
	response, err := processCreate(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}

/*
creq := NServiceRequest{
	City: "San Jose",
}
fmt.Printf("%+v\n", creq)
var cresp NServicesResponse
replyCall := client.Go("Service.ServicesForCity", &creq, &cresp, nil)
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
*/

var adapters map[string]struct{name string, port int}{
	"CS": struct{name string, port int}{"rpc1", 5001},
	"CS2": struct{name string, port int}{"rpc2", 5002},
}

func init() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Print("dialing:", err)
	}
}

func CallAdapter(apiCall string, request, response interface{}) *rpc.Call {
	fmt.Println(spew.Sdump(request))
	return client.Go(apiCall, request, response, nil)
}
