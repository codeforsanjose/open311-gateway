package router

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestReadConfig(t *testing.T) {
	if err := Init(); err != nil {
		t.Errorf("Init() failed: %s", err)
	}

	fmt.Printf("\n==================== ADAPTERS ==================\n%s", adapters)
	if adapters.loaded != true {
		t.Errorf("Site configuration is not marked as loaded.")
	}

	fmt.Println(spew.Sdump(adapters))

}
