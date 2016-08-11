package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/open311-gateway/engine/router"

	"github.com/davecgh/go-spew/spew"
	log "github.com/jeffizhungry/logrus"
)

var Debug = true

func init() {
	log.Setup(false, log.DebugLevel)
	fmt.Println("Reading config...")
	if err := router.Init("../data/config.json"); err != nil {
		log.Errorf("Init() failed: %s", err)
	}
}

func isError(e error) bool {
	if e == nil {
		return false
	}
	return true
}

func TestServiceDataRefresh(t *testing.T) {
	f := func(run int) {
		fmt.Printf("\n\n\n\n============================= [TestServiceDataRefresh%d] =============================\n\n", run)
		Refresh()
		time.Sleep(300 * time.Millisecond)
		fmt.Println(servicesData)
	}

	time.Sleep(time.Second * 2)
	for i := 1; i <= 4; i++ {
		f(i)
	}
}

func TestRetrieve(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestRetrieve] =============================\n\n")

	cases := []struct {
		in   string
		want bool
	}{
		{"SJ", false},
		{"SC", false},
		{"XX", true},
	}

	for _, c := range cases {
		fmt.Printf(">>>>>>>>>>>> Retrieving data for: %q\n", c.in)
		l, err := GetArea(c.in)
		if isError(err) != c.want {
			t.Errorf("GetArea() failed: %s", err)
		}
		fmt.Println(l)
	}
}

func TestAreaAdapters(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestAreaAdapters] =============================\n\n")
	fmt.Printf("---------- San Jose -----------\n%s\n", spew.Sdump(router.GetAreaAdapters("SJ")))
	fmt.Printf("---------- Santa Clara --------\n%s\n", spew.Sdump(router.GetAreaAdapters("SC")))
	fmt.Printf("---------- XX -----------------\n%s\n", spew.Sdump(router.GetAreaAdapters("SJ")))
	fmt.Println(spew.Sdump(router.GetAreaAdapters("SJ")))
	fmt.Println(spew.Sdump(router.GetAreaAdapters("XX")))
}

func TestShutdown(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestShutdown] =============================\n\n")
	time.Sleep(time.Second * 1)
	Shutdown()
}
