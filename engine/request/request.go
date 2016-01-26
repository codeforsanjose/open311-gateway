package request

import (
	"net/http"

	"Gateway311/engine/logs"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	log = logs.Log
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

// Create creates a new report.
func Create(w rest.ResponseWriter, r *rest.Request) {
	response, err := processCreate(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}
