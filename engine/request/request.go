package request

import (
	"net/http"
	"time"

	"Gateway311/engine/logs"
	"Gateway311/engine/router"
	"Gateway311/engine/telemetry"

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
	defer func() {
		if rcvr := recover(); rcvr != nil {
			rest.Error(w, rcvr.(error).Error(), http.StatusInternalServerError)
		}
	}()
	rqstID := router.GetSID()
	sendTelemetry(rqstID, "Services", "open")
	response, err := processServices(r, rqstID)
	if err != nil {
		sendTelemetry(rqstID, "Services", "error")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendTelemetry(rqstID, "Services", "done")
	w.WriteJson(&response)
}

// Create creates a new report.
func Create(w rest.ResponseWriter, r *rest.Request) {
	defer func() {
		if rcvr := recover(); rcvr != nil {
			rest.Error(w, rcvr.(error).Error(), http.StatusInternalServerError)
		}
	}()
	response, err := processCreate(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}

// Search searches for Reports.
func Search(w rest.ResponseWriter, r *rest.Request) {
	// defer func() {
	// 	if rcvr := recover(); rcvr != nil {
	// 		rest.Error(w, rcvr.(error).Error(), http.StatusInternalServerError)
	// 	}
	// }()
	response, err := processSearch(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&response)
}

// =======================================================================================
//                                      TELEMETRY
// =======================================================================================

func sendTelemetry(rqstID int64, op, status string) {
	telemetry.SendRequest(rqstID, op, status, "", time.Now())
}
