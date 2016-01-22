package services

import (
	"CitySourcedAPI/logs"
	"fmt"
	"testing"
	"time"

	"Gateway311/engine/router"
)

var Debug = true

func init() {
	logs.Init(Debug)

	fmt.Println("Reading config...")
	if err := router.Init("../data/config.json"); err != nil {
		fmt.Errorf("Init() failed: %s", err)
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
		time.Sleep(2 * time.Second)
		fmt.Println(servicesData)
	}

	time.Sleep(time.Second * 5)
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

	time.Sleep(time.Second * 1)
	for _, c := range cases {
		fmt.Printf(">>>>>>>>>>>> Retrieving data for: %q\n", c.in)
		l, err := GetArea(c.in)
		if isError(err) != c.want {
			t.Errorf("GetArea() failed: %s", err)
		}
		fmt.Println(l)
	}
}

func TestShutdown(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestShutdown] =============================\n\n")
	time.Sleep(time.Second * 1)
	Shutdown()
}
