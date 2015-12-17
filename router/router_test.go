package router

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	if err := readConfig("config.json"); err == nil {
		t.Errorf("There should have been an error here.")
	}

	fmt.Printf("\n==================== ROUTE DATA ==================\n%s", routeData)
	if routeData.Loaded != true {
		t.Errorf("Site configuration is not marked as loaded.")
	}

	city := "San Jose"
	p := GetServiceProviders(city)
	fmt.Printf("\n==================== SERVICE PROVIDERS ==================\nFor: %s\n%#v\n", city, p)
	for i, v := range p {
		fmt.Printf("%2d  %s\n", i, v.Name)
	}

	s := GetServices(city)
	fmt.Printf("\n==================== SERVICES ===========================\nFor: %s\n", city)
	for i, v := range s {
		fmt.Printf("%2d:%v\n", i+1, v)
	}

}
