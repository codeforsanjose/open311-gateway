package request

import (
	"Gateway311/geo"
	"Gateway311/router"

	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ant0ine/go-json-rest/rest"
)

// Services looks up the service providers and services for the specified location.
// The URL nust contain query parameters of either:
// latitude and longitude, or a city name.
//
// Examples:
//  http;//xyz.com/api/services?lat=34.236144&lon=-118.604794
//  http;//xyz.com/api/services?city=san+jose
func Services(w rest.ResponseWriter, r *rest.Request) {
	response := "Error!"
	// m, _ := url.ParseQuery(r.URL.RawQuery)
	m := r.URL.Query()
	for k, v := range m {
		fmt.Printf("%s: %#v\n", k, v)
	}

	if _, ok := m["lat"]; ok {
		log.Printf("   QueryParms: Lat/Lng...\n")
		lat, err1 := strconv.ParseFloat(m["lat"][0], 64)
		lng, err2 := strconv.ParseFloat(m["lng"][0], 64)
		if err1 != nil || err2 != nil {
			msg := fmt.Sprintf("Invalid lat/lng: %s:%s", m["lat"][0], m["lng"][0])
			log.Printf(msg)
			rest.Error(w, msg, http.StatusInternalServerError)
			return
		}
		log.Printf("Lat: %v lng: %v\n", lat, lng)
		city, err := geo.GetCity(lat, lng)
		if err != nil {
			log.Printf("Can't find city for: %v:%v\n", lat, lng)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r := RespServices{}
		r.JID, r.Services, err = router.Services(city)
		if err != nil {
			log.Printf("%s", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteJson(&r)
		return

	} else if _, ok := m["city"]; ok {
		err := errors.New("")
		city := m["city"][0]
		log.Printf("   QueryParms: Address - city: %s\n", city)
		r := RespServices{}
		r.JID, r.Services, err = router.Services(city)
		if err != nil {
			log.Printf("%s", err)
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteJson(&r)
		return
	}

	// ToDo: build correct response
	response = "OK!"
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

// ==============================================================================================================================
//                                      SERVICES
// ==============================================================================================================================

// RespServices is used to return a service list.
type RespServices struct {
	JID      int               `json:"jid" xml:"jid"`
	Services []*router.Service `json:"services" xml:"services"`
}
