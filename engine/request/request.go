package request

import (
	"net/http"
	"sync/atomic"
	"time"

	"Gateway311/engine/logs"
	"Gateway311/engine/telemetry"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	log    = logs.Log
	rqstID sidType
)

// Services looks up the service providers and services for the specified location.
// The URL nust contain query parameters of either:
// LatitudeV and LongitudeV, or a city name.
//
// Examples:
//  http;//xyz.com/api/services?lat=34.236144&lon=-118.604794
//  http;//xyz.com/api/services?city=san+jose
func Services(w rest.ResponseWriter, r *rest.Request) {
	rqstID := rqstID.Get()
	SendTelemetry(rqstID, "Services", "open")
	response, err := processServices(r, rqstID)
	if err != nil {
		SendTelemetry(rqstID, "Services", "error")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SendTelemetry(rqstID, "Services", "complete")
	w.WriteJson(&response)
}

// Create creates a new report.
func Create(w rest.ResponseWriter, r *rest.Request) {
	rqstID := rqstID.Get()
	SendTelemetry(rqstID, "Create", "open")
	response, err := processCreate(r, rqstID)
	if err != nil {
		SendTelemetry(rqstID, "Create", "error")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SendTelemetry(rqstID, "Create", "complete")
	w.WriteJson(&response)
}

// Search searches for Reports.
func Search(w rest.ResponseWriter, r *rest.Request) {
	rqstID := rqstID.Get()
	SendTelemetry(rqstID, "Search", "open")
	response, err := processSearch(r, rqstID)
	if err != nil {
		SendTelemetry(rqstID, "Search", "error")
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SendTelemetry(rqstID, "Search", "complete")
	w.WriteJson(&response)
}

func SendTelemetry(rqstID int64, op, status string) {
	telemetry.SendEngRequest(rqstID, op, status, "", time.Now())
}

// =======================================================================================
//                                      MESSAGE ID
// =======================================================================================

type sidType int64

func (r *sidType) Get() int64 {
	return atomic.AddInt64((*int64)(r), 1)
}

func init() {
	rqstID = 1000
}
