package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"Gateway311/engine/logs"
	"Gateway311/engine/request"
	"Gateway311/engine/router"
	"Gateway311/engine/services"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	log = logs.Log
	// Debug switches on some debugging statements.
	Debug = false
)

func main() {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(

		rest.Get("/services", request.Services),
		rest.Post("/requests", request.Create),
		// rest.Get("/:jid/requests", rpt.GetAll),
		// rest.Get("/:jid/requests/:id", rpt.Get),
		// rest.Put("/:jid/requests/:id", rpt.Update),
		// rest.Delete("/:jid/requests/:id", rpt.Delete),

	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

func init() {
	var configFile string
	flag.BoolVar(&Debug, "debug", false, "Activates debug logging.")
	flag.StringVar(&configFile, "config", "data/config.json", "Config file. This is a full or relative path.")
	flag.Parse()

	logs.Init(Debug)

	if err := router.Init(configFile); err != nil {
		log.Fatal("Unable to start - data initilization failed.\n")
	}

	go signalHandler(make(chan os.Signal, 1))
	fmt.Println("Press Ctrl-C to shutdown...")

	time.Sleep(time.Second * 2)
	services.Refresh()
}

func signalHandler(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	for s := <-c; ; s = <-c {
		switch s {
		case os.Interrupt:
			fmt.Println("Ctrl-C Received!")
			stop()
			os.Exit(0)
		case os.Kill:
			fmt.Println("SIGKILL Received!")
			stop()
			os.Exit(1)
		}
	}
}

func stop() error {
	return nil
}
