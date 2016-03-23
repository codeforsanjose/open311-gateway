package request

import (
	"errors"
	"net/http"

	"Gateway311/engine/logs"

	"github.com/ant0ine/go-json-rest/rest"
)

const (
	debugRecover = false

	greEmpty = "JSON payload is empty"
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
	runRequest(w, r, processServices)
}

// Create creates a new report.
func Create(w rest.ResponseWriter, r *rest.Request) {
	runRequest(w, r, processCreate)
}

// Search searches for Reports.
func Search(w rest.ResponseWriter, r *rest.Request) {
	runRequest(w, r, processSearch)
}

func runRequest(w rest.ResponseWriter, r *rest.Request, f func(*rest.Request) (interface{}, error)) {
	if debugRecover {
		defer func() {
			if rcvr := recover(); rcvr != nil {
				rest.Error(w, rcvr.(error).Error(), http.StatusInternalServerError)
			}
		}()
	}
	response, err := f(r)
	if err != nil {
		errorResp(w, newErrorsResponseJ().errorJ(400, err.Error()), http.StatusBadRequest)
		return
	}
	if err := w.WriteJson(&response); err != nil {
		log.Error(err.Error())
	}
}

func errorResp(w rest.ResponseWriter, errResp ErrorsResponseJ, code int) {
	w.WriteHeader(code)
	err := w.WriteJson(errResp)
	if err != nil {
		panic(errors.New("invalid error response"))
	}
}
