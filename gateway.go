package main

import (
	"Gateway311/report"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {

	testrequests := report.TestReports{
		Store: map[string]*report.TestReport{},
	}
	rpt := report.CreateReport{}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		rest.Get("/:jid/testrequests", testrequests.GetAllRequests),
		rest.Post("/:jid/testrequests", testrequests.PostRequest),
		rest.Get("/:jid/testrequests/:id", testrequests.GetRequest),
		rest.Put("/:jid/testrequests/:id", testrequests.PutTestReport),
		rest.Delete("/:jid/testrequests/:id", testrequests.DeleteRequest),

		rest.Get("/:jid/requests", rpt.GetAll),
		rest.Get("/:jid/requests/:id", rpt.Get),
		rest.Post("/:jid/requests", report.Create),
		rest.Put("/:jid/requests/:id", rpt.Update),
		rest.Delete("/:jid/requests/:id", rpt.Delete),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
