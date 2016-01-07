package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"Gateway311/gateway/common"

	"github.com/davecgh/go-spew/spew"
)

const beIface = "CS"

var (
	providerData ProviderData
)

// ShowProviderData dumps providerData using spew.
func ShowProviderData() string {
	return spew.Sdump(providerData)
}

// Services returns a list of all services available for the specified City.
func Services(city string) ([]*Service, error) {
	return providerData.Services(city)
}

// func Services(city string) ([]*request.NService, error) {
// 	return providerData.Services(city)
// }

// ServiceProviders returns a list of all Service Providers for the specified City.
func ServiceProviders(city string) ([]*Provider, error) {
	return providerData.ServiceProviders(city)
}

// ServiceProvider returns a pointer to the Provider for the specified Provider ID.
func ServiceProvider(sid int) (*Provider, error) {
	return providerData.ServiceProvider(sid)
}

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

func init() {
	if err := readConfig("/Users/james/Dropbox/Development/go/src/Gateway311/integration/citysourced/data/config.json"); err != nil {
		fmt.Printf("Error %v occurred when reading the config - ReadConfig()", err)
	}
}

func readConfig(filePath string) error {
	if providerData.Loaded {
		msg := "Route Data is already loaded"
		fmt.Println(msg)
		return errors.New(msg)
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		msg := fmt.Sprintf("Unable to access the providerData Config file - specified at: %q.\nError: %v", filePath, err)
		fmt.Println(msg)
		return errors.New(msg)
	}

	return providerData.Load(file)
}

// ==============================================================================================================================
//                                      ROUTE DATA
// ==============================================================================================================================

// ProviderData is a list of all the Service Areas.  It contains an indexed list of all the Service Areas.  The index is the *lowercase* city name.
type ProviderData struct {
	Loaded     bool
	Categories []string                        `json:"serviceCategories"`
	Areas      map[string]map[string]*Provider `json:"serviceAreas"`

	providerID   map[int]*Provider
	serviceID    map[int]*Service
	cityServices map[string][]*Service
	// serviceID       map[int]*request.NService
	// cityServices    map[string][]*request.NService
	serviceProvider map[int]*Provider
}

// Services returns a list of all services available for the specified City.
func (rd *ProviderData) Services(city string) ([]*Service, error) {
	lcity := strings.ToLower(city)
	fmt.Printf("   Services for: %s...\n", lcity)
	if !rd.isValidCity(lcity) {
		return nil, fmt.Errorf("The city: %q is not serviced by this Gateway", city)
	}
	fmt.Printf("      data length: %d\n", len(rd.cityServices[lcity]))
	return rd.cityServices[lcity], nil
}

// func (rd *ProviderData) Services(city string) ([]*request.NService, error) {
// 	lcity := strings.ToLower(city)
// 	fmt.Printf("   Services for: %s...\n", lcity)
// 	if !rd.isValidCity(lcity) {
// 		return nil, fmt.Errorf("The city: %q is not serviced by this Gateway", city)
// 	}
// 	fmt.Printf("      data length: %d\n", len(rd.cityServices[lcity]))
// 	return rd.cityServices[lcity], nil
// }

// ServiceProviders returns a list of all Service Providers for the specified City.
func (rd *ProviderData) ServiceProviders(city string) ([]*Provider, error) {
	lcity := strings.ToLower(city)
	fmt.Printf("   ServiceProviders for: %q\n", lcity)
	if !rd.isValidCity(lcity) {
		return nil, fmt.Errorf("The city: %q is not serviced by this Gateway", city)
	}
	var p []*Provider
	for _, v := range rd.Areas[strings.ToLower(city)] {
		p = append(p, v)
	}
	fmt.Printf("[ServiceProviders] returning %d records.\n", len(p))
	return p, nil
}

// ServiceProvider returns a pointer to the Provider for the specified Provider ID.
func (rd *ProviderData) ServiceProvider(id int) (*Provider, error) {
	p, ok := rd.serviceProvider[id]
	if !ok {
		return nil, fmt.Errorf("Invalid Service ID")
	}
	// fmt.Printf("[ServiceProvider] returning %#v\n", *p)
	return p, nil
}

// Load loads the specified byte slice into the ProviderData structures.
func (rd *ProviderData) Load(file []byte) error {
	rd.init()
	err := json.Unmarshal(file, rd)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse JSON Route Data.\nError: %v", err)
		fmt.Println(msg)
		return errors.New(msg)
	}

	rd.index()
	rd.Loaded = true
	return nil
}

// Index() builds all required map indexes: Services by City,
func (rd *ProviderData) index() error {
	fmt.Printf("[ProviderData] building all indexes:\n")
	rd.indexID()
	return nil
}

func (rd *ProviderData) indexID() error {
	fmt.Printf("    building indexes...\n")
	for areaKey, area := range rd.Areas {
		rd.cityServices[areaKey] = make([]*Service, 0)
		// rd.cityServices[areaKey] = make([]*request.NService, 0)
		for _, provider := range area {
			rd.providerID[provider.ID] = provider
			// fmt.Println("*** BUILDING SERVICES")
			for _, service := range provider.Services {
				service.ProviderID = provider.ID
				rd.cityServices[areaKey] = append(rd.cityServices[areaKey], service)
				// fmt.Printf("   %s ===> %+v\n", serviceName, service)
				rd.serviceID[service.ID] = service
				rd.serviceProvider[service.ID] = provider
			}
		}
	}
	return nil
}

func (rd *ProviderData) init() {
	rd.providerID = map[int]*Provider{}
	rd.serviceID = map[int]*Service{}
	rd.cityServices = map[string][]*Service{}
	// rd.serviceID = map[int]*request.NService{}
	// rd.cityServices = map[string][]*request.NService{}
	rd.serviceProvider = map[int]*Provider{}
}

// String returns the represeentation of the ProviderData custom type.
func (rd ProviderData) String() string {
	ls := new(common.LogString)
	ls.AddS("ProviderData\n")
	ls.AddF("Loaded: %t\n", rd.Loaded)
	ls.AddS("INDEX: providerID\n")
	for k, v := range rd.providerID {
		ls.AddF("   %-4d  %-40s %s\n", k, v.Name, v.URL)
	}
	ls.AddS("INDEX: serviceID\n")
	for k, v := range rd.serviceID {
		ls.AddF("   %-4d (%t)  %s\n", k, k == v.ID, v.Name)
	}
	ls.AddS("INDEX: cityServices\n")
	for k, v := range rd.cityServices {
		ls.AddF("   [[%s]]\n", k)
		for _, sv := range v {
			ls.AddF("      %s\n", sv)
		}
	}
	ls.AddS("INDEX: serviceProvider\n")
	for k, v := range rd.serviceProvider {
		ls.AddF("   %-4d %-50.50s  %-40s\n", rd.serviceID[k].ID, rd.serviceID[k].Name, v.Name)
	}
	ls.AddS("\n--- CATEGORIES ---\n")
	for i, v := range rd.Categories {
		ls.AddF("   %2d  %s\n", i, v)
	}
	ls.AddS("\n---AREAS ---\n")
	for _, v := range rd.Areas {
		for k2, v2 := range v {
			ls.AddF("[%s]\n%s\n", k2, v2)
		}
	}
	return ls.Box(90)
}

func (rd *ProviderData) isValidCity(city string) bool {
	_, ok := rd.Areas[strings.ToLower(city)]
	return ok
}

// ------------------------------- Provider -------------------------------

// Provider is the data for each Service Provider.  It contains an index list of all of the Services provided by this Provider.
type Provider struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	APIVersion string     `json:"apiVersion"`
	Key        string     `json:"key"`
	Services   []*Service `json:"services"`
	// Services   []*request.NService `json:"services"`
}

func (p Provider) String() string {
	ls := new(common.LogString)
	ls.AddF("%s (ID: %d)\n", p.Name, p.ID)
	ls.AddF("URL: %s  ver: %s  Key: %s\n", p.APIVersion, p.URL, p.Key)
	ls.AddS("---SERVICES:\n")
	for _, v := range p.Services {
		ls.AddF("   %s\n", v)
	}
	return ls.Box(80)
}

// ------------------------------- Service -------------------------------

// Service is a map of the
type Service struct {
	ID         int      `json:"id"`
	ProviderID int      `json:"providerId"`
	Name       string   `json:"name"`
	Categories []string `json:"catg"`
}

func (s Service) String() string {
	r := fmt.Sprintf("   %4d %4d  %-40s  %v", s.ID, s.ProviderID, s.Name, s.Categories)
	return r
}
