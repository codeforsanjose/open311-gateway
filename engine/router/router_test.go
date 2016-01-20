package router

import (
	"CitySourcedAPI/logs"
	"fmt"
	"testing"
	"time"

	"Gateway311/engine/structs"

	"github.com/davecgh/go-spew/spew"
)

var Debug = true

func TestReadConfig(t *testing.T) {
	logs.Init(true)
	if err := Init("../data/config.json"); err != nil {
		t.Errorf("Init() failed: %s", err)
	}

	fmt.Printf("\n==================== ADAPTERS ==================\n%s", adapters)
	if adapters.loaded != true {
		t.Errorf("Site configuration is not marked as loaded.")
	}

	fmt.Println(spew.Sdump(adapters))

}

func TestRPC1(t *testing.T) {
	rqst := &structs.NServiceRequest{"all"}
	r, err := newRPCCall("Service.All", "all", rqst, servicesData.merge)
	if err != nil {
		log.Debug(err.Error())
	}
	fmt.Println(r)

	r.run()
	time.Sleep(2 * time.Second)

}
