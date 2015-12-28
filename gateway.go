package main

import (
	"Gateway311/request"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {

	testrequests := request.TestReports{
		Store: map[string]*request.TestReport{},
	}
	// rpt := request.CreateReq{}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		rest.Get("/:jid/testrequests", testrequests.GetAllRequests),
		rest.Post("/:jid/testrequests", testrequests.PostRequest),
		rest.Get("/:jid/testrequests/:id", testrequests.GetRequest),
		rest.Put("/:jid/testrequests/:id", testrequests.PutTestReport),
		rest.Delete("/:jid/testrequests/:id", testrequests.DeleteRequest),

		// rest.Get("/:jid/requests", rpt.GetAll),
		// rest.Get("/:jid/requests/:id", rpt.Get),
		rest.Post("/:jid/requests", request.Create),
		// rest.Put("/:jid/requests/:id", rpt.Update),
		// rest.Delete("/:jid/requests/:id", rpt.Delete),

		rest.Get("/services", request.Services),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
