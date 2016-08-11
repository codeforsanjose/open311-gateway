package data

import (
	"fmt"
	"testing"

	"github.com/open311-gateway/adapters/email/logs"
	"github.com/open311-gateway/adapters/email/structs"

	"github.com/davecgh/go-spew/spew"
)

var Debug = true

func init() {
	logs.Init(Debug)

	fmt.Println("Reading config...")
	if err := Init("config.json"); err != nil {
		fmt.Printf("Init() failed: %s", err)
	}
}

type testResultS struct {
	input string
	isOK  bool
}

func isOK(e error) bool {
	if e == nil {
		return false
	}
	return true
}

func TestServices(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestServices] =============================\n\n")

	var test1 = [3]testResultS{
		{"Cupertino", true},
		{"San Jose", false},
		{"Sunnyvale", true},
	}

	for _, tt := range test1 {
		svcs, err := ServicesArea(tt.input)

		switch {
		case tt.isOK && err == nil:
			fmt.Printf("svcs for %q:\n%s", tt.input, spew.Sdump(*svcs))
		case tt.isOK && err != nil:
			t.Errorf("ServicesArea() failed for: %q", tt.input)
		case !tt.isOK && err == nil:
			t.Errorf("ServicesArea() should have failed for: %q", tt.input)
		}
	}

	fmt.Printf("----------------------------- [TestServicesAll] -----------------------------\n\n")
	svcs, err := ServicesAll()
	if err != nil {
		t.Errorf("ServicesArea() failed.")
	} else {
		fmt.Printf("svcs:\n%s", spew.Sdump(*svcs))
	}

}

func TestAdapter(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [TestAdapter] =============================\n\n")

	fmt.Printf("----------------------------- [Adapter] -----------------------------\n\n")
	if name, atype, address := Adapter(); name != "EM1" || atype != "Email" {
		t.Errorf("Adapter() failed - name: %q  atype: %q", name, atype)
	} else {
		fmt.Printf("OK - address: %s\n", address)
	}

	fmt.Printf("----------------------------- [AdapterName] -----------------------------\n\n")
	if name := AdapterName(); name != "EM1" {
		t.Errorf("AdapterName() failed - name: %q", name)
	} else {
		fmt.Println("OK!")
	}

	fmt.Printf("\n\n----------------------------- [MIDProvider] -----------------------------\n\n")
	var midTests = []struct {
		n    structs.ServiceID // input
		isOK bool              // expected result
	}{
		{structs.ServiceID{AdpID: "EM1", AreaID: "CU", ProviderID: 1, ID: 1}, true},
		{structs.ServiceID{AdpID: "EM1", AreaID: "CU", ProviderID: 2, ID: 9999}, true},
		{structs.ServiceID{AdpID: "EM1", AreaID: "CU", ProviderID: 3, ID: 999999}, false},
		{structs.ServiceID{AdpID: "EM1", AreaID: "SUN", ProviderID: 1, ID: 1}, true},
		{structs.ServiceID{AdpID: "EM1", AreaID: "SJ", ProviderID: 1, ID: 1}, false},
		{structs.ServiceID{AdpID: "EM1", AreaID: "XXXXXXXXX", ProviderID: 1, ID: 1}, false},
	}

	for _, tt := range midTests {
		prov, err := MIDProvider(tt.n)
		switch {
		case tt.isOK && err == nil:
			fmt.Printf("\nsvcs for %q:\n%s", tt.n.MID(), prov.String())
		case tt.isOK && err != nil:
			t.Errorf("ServicesArea() failed for: %q", tt.n.MID())
		case !tt.isOK && err == nil:
			t.Errorf("ServicesArea() should have failed for: %q", tt.n)
		}
	}

	fmt.Printf("\n\n----------------------------- [RouteProvider] -----------------------------\n\n")
	var routeTests = []struct {
		n    structs.NRoute // input
		isOK bool           // expected result
	}{
		{structs.NRoute{AdpID: "EM1", AreaID: "CU", ProviderID: 1}, true},
		{structs.NRoute{AdpID: "EM1", AreaID: "CU", ProviderID: 2}, true},
		{structs.NRoute{AdpID: "EM1", AreaID: "CU", ProviderID: 3}, false},
		{structs.NRoute{AdpID: "EM1", AreaID: "SUN", ProviderID: 1}, true},
		{structs.NRoute{AdpID: "EM1", AreaID: "SJ", ProviderID: 1}, false},
		{structs.NRoute{AdpID: "EM1", AreaID: "XXXXXXXXX", ProviderID: 1}, false},
	}

	for _, tt := range routeTests {
		prov, err := RouteProvider(tt.n)
		switch {
		case tt.isOK && err == nil:
			fmt.Printf("\nsvcs for %q:\n%s", tt.n.String(), prov.String())
		case tt.isOK && err != nil:
			t.Errorf("ServicesArea() failed for: %q", tt.n.String())
		case !tt.isOK && err == nil:
			t.Errorf("ServicesArea() should have failed for: %q", tt.n)
		}
	}
}

func TestEmail(t *testing.T) {
	fmt.Printf("\n\n\n\n============================= [DisplayEmail] =============================\n\n")

	fmt.Print(configData.Email)

	for areaID, areaData := range configData.Areas {
		fmt.Printf("\n\n---------- Area: %q ----------\n", areaID)
		for _, prov := range areaData.Providers {
			fmt.Printf("\nProvider: %v - %v\n%s", prov.ID, prov.Name, prov.Email.String())
		}
	}
}
