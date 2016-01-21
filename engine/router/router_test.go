package router

import (
	"CitySourcedAPI/logs"
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var Debug = true

func TestReadConfig(t *testing.T) {
	fmt.Println("\n\n\n\n============================= [TestReadConfig] =============================")
	logs.Init(true)
	if err := Init("../data/config.json"); err != nil {
		t.Errorf("Init() failed: %s", err)
	}

	fmt.Printf("\n-------------------- ADAPTERS ------------------\n%s", adapters)
	if adapters.loaded != true {
		t.Errorf("Site configuration is not marked as loaded.")
	}

	fmt.Println(spew.Sdump(adapters))

}

func TestServiceDataRefresh(t *testing.T) {
	f := func(run int) {
		fmt.Printf("\n\n\n\n============================= [TestServiceDataRefresh%d] =============================", run)
		servicesData.refresh()
		time.Sleep(2 * time.Second)
		fmt.Println(servicesData)
	}

	for i := 1; i <= 4; i++ {
		f(i)
	}
}
