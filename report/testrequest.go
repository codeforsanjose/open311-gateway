package report

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

// TestReport is a specific service request.
type TestReport struct {
	JID         string  `json:"jid" xml:"jid"`
	ID          string  `json:"id" xml:"id"`
	Type        string  `json:"type" xml:"type"`
	TypeID      string  `json:"typeId" xml:"typeId"`
	DeviceType  string  `json:"deviceType" xml:"deviceType"`
	DeviceModel string  `json:"deviceModel" xml:"deviceModel"`
	DeviceID    string  `json:"deviceID" xml:"deviceID"`
	Latitude    float64 `json:"latitude" xml:"latitude"`
	Longitude   float64 `json:"longitude" xml:"longitude"`
	Address     string  `json:"address" xml:"address"`
	City        string  `json:"city" xml:"city"`
	State       string  `json:"state" xml:"state"`
	Zip         string  `json:"zip" xml:"zip"`
	FirstName   string  `json:"firstName" xml:"firstName"`
	LastName    string  `json:"lastName" xml:"lastName"`
	Email       string  `json:"email" xml:"email"`
	Phone       string  `json:"phone" xml:"phone"`
	IsAnonymous bool    `json:"isAnonymous" xml:"isAnonymous"`
	Description string  `json:"Description" xml:"Description"`
}

// TestReports is a struct for storing all existing requests.
type TestReports struct {
	sync.RWMutex
	Store  map[string]*TestReport
	lastID int
}

// GetAllRequests retrieves all requests, with no filtering.
func (u *TestReports) GetAllRequests(w rest.ResponseWriter, r *rest.Request) {
	fmt.Printf("[GetAllRequests]\n")
	jid := r.PathParam("jid")
	r.ParseForm()
	fmt.Printf("  jid: (%T)%q Form: %#v\n", jid, jid, r.Form)
	u.RLock()
	var requests []TestReport
	for _, request := range u.Store {
		if request.JID == jid {
			requests = append(requests, *request)
		}
	}
	u.RUnlock()
	w.WriteJson(&requests)
}

// GetRequest retrieves a single request specified by it's ID.
func (u *TestReports) GetRequest(w rest.ResponseWriter, r *rest.Request) {
	fmt.Printf("[GetRequest]\n")
	jid := r.PathParam("jid")
	id := r.PathParam("id")
	fmt.Printf("  jid: %s id: %s\n", jid, id)
	u.RLock()
	var request *TestReport
	if u.Store[id] != nil {
		if u.Store[id].JID == jid {
			request = &TestReport{}
			*request = *u.Store[id]
		}
	}
	u.RUnlock()
	if request == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(request)
}

// PostRequest creates a new TestReport and adds it to TestReports.
func (u *TestReports) PostRequest(w rest.ResponseWriter, r *rest.Request) {
	jid := r.PathParam("jid")
	fmt.Printf("[PostRequest] - jid: (%T)%q\n", jid, jid)
	request := TestReport{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Lock()
	u.lastID++
	request.JID = jid
	id := fmt.Sprintf("%d", u.lastID)
	request.ID = id
	u.Store[id] = &request
	u.Unlock()
	w.Header().Set("Location", fmt.Sprintf("http://localhost:8080/requests/%s", id))
	w.WriteJson(&request)
}

// PutTestReport updates the request specified by ID.
func (u *TestReports) PutTestReport(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	u.Lock()
	if u.Store[id] == nil {
		rest.NotFound(w, r)
		u.Unlock()
		return
	}
	request := TestReport{}
	err := r.DecodeJsonPayload(&request)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		u.Unlock()
		return
	}
	request.ID = id
	u.Store[id] = &request
	u.Unlock()
	w.WriteJson(&request)
}

// DeleteRequest deletes the request specified by ID.
func (u *TestReports) DeleteRequest(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	u.Lock()
	delete(u.Store, id)
	u.Unlock()
	w.WriteHeader(http.StatusOK)
}
