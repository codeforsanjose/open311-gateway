package router

import (
	"CitySourcedAPI/logs"
	"fmt"
	"testing"

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
