package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/open311-gateway/engine/request"
	"github.com/open311-gateway/engine/router"
	"github.com/open311-gateway/engine/services"
	"github.com/open311-gateway/engine/telemetry"

	"github.com/ant0ine/go-json-rest/rest"
	log "github.com/jeffizhungry/logrus"
)

var (
	configFile string

	// Debug switches on some debugging statements.
	Debug = false
)

func main() {
	log.Debug("Version 0.1.1")
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			// return origin == "http://my.other.host"
			return true
		},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	restrouter, err := rest.MakeRouter(
		rest.Get("/v1/services.json", request.Services),
		rest.Post("/v1/requests.json", request.Create),
		rest.Get("/v1/requests.json", request.Search),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(restrouter)

	addr, prot, cert, key := router.GetNetworkConfig()
	switch prot {
	case "http":
		log.Fatal(http.ListenAndServe(addr, api.MakeHandler()))
	case "https":
		log.Fatal(http.ListenAndServeTLS(addr, cert, key, api.MakeHandler()))
	default:
		log.Fatalf("Invalid network protocol: %s specified in config.", prot)
	}
}

func init() {
	log.Setup(false, log.DebugLevel)

	flag.BoolVar(&Debug, "debug", false, "Activates debug logging.")
	flag.StringVar(&configFile, "config", "config.json", "Config file. This is a full or relative path.")
	flag.Parse()

	fmt.Printf("Debug: %v  Config: %v\n", Debug, configFile)

	if err := router.Init(configFile); err != nil {
		log.Fatal("Unable to start - data initilization failed.\n")
	}

	if err := request.Init(); err != nil {
		log.Fatal("Unable to start - initialization of Request package failed.\n")
	}

	telemetry.Init(router.GetMonitorAddress())

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
