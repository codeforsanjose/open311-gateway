package router

import (
	"Gateway311/common"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	routeData RouteData
)

// GetServices returns a list of all services available for the specified City.
func GetServices(city string) (int, []*Service, error) {
	return routeData.GetServices(city)
}

// GetServiceProviders returns a list of all Service Providers for the specified City.
func GetServiceProviders(city string) ([]*Provider, error) {
	return routeData.GetServiceProviders(city)
}

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

func init() {
	if err := readConfig("/Users/james/Dropbox/Development/go/src/Gateway311/router/config.json"); err != nil {
		fmt.Printf("Error %v occurred when reading the config - ReadConfig()", err)
	}
}

func readConfig(filePath string) error {
	if routeData.Loaded {
		msg := "Route Data is already loaded"
		fmt.Println(msg)
		return errors.New(msg)
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		msg := fmt.Sprintf("Unable to access the routeData Config file - specified at: %q.\nError: %v", filePath, err)
		fmt.Println(msg)
		return errors.New(msg)
	}

	return routeData.Load(file)
}

// ==============================================================================================================================
//                                      ROUTE DATA
// ==============================================================================================================================

// RouteData is a list of all the Service Areas.  It contains an indexed list of all the Service Areas.  The index is the *lowercase* city name.
type RouteData struct {
	Loaded          bool
	Categories      []string         `json:"serviceCategories"`
	Areas           map[string]*Area `json:"serviceAreas"`
	areaID          map[int]*Area
	providerID      map[int]*Provider
	serviceID       map[int]*Service
	cityServices    map[string][]*Service
	providerService map[int]*Provider
}

// GetServices returns a list of all services available for the specified City.
func (rd *RouteData) GetServices(city string) (int, []*Service, error) {
	lcity := strings.ToLower(city)
	fmt.Printf("   GetServices for: %s...\n", lcity)
	if !rd.isValidCity(lcity) {
		return 0, nil, fmt.Errorf("The city: %q is not serviced by this Gateway", city)
	}
	fmt.Printf("      data length: %d\n", len(rd.cityServices[lcity]))
	return rd.Areas[lcity].ID, rd.cityServices[lcity], nil
}

// GetServiceProviders returns a list of all Service Providers for the specified City.
func (rd *RouteData) GetServiceProviders(city string) ([]*Provider, error) {
	lcity := strings.ToLower(city)
	fmt.Printf("   GetServiceProviders for: %q\n", lcity)
	if !rd.isValidCity(lcity) {
		return nil, fmt.Errorf("The city: %q is not serviced by this Gateway", city)
	}
	var p []*Provider
	for _, v := range rd.Areas[strings.ToLower(city)].Providers {
		p = append(p, v)
	}
	fmt.Printf("[getServiceProviders] returning %d records.\n", len(p))
	return p, nil
}

// Load loads the specified byte slice into the RouteData structures.
func (rd *RouteData) Load(file []byte) error {
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

func (rd *RouteData) init() {
	rd.areaID = map[int]*Area{}
	rd.providerID = map[int]*Provider{}
	rd.serviceID = map[int]*Service{}
	rd.cityServices = map[string][]*Service{}
	rd.providerService = map[int]*Provider{}
}

func (rd *RouteData) isValidCity(city string) bool {
	_, ok := rd.Areas[strings.ToLower(city)]
	return ok
}

// String returns the represeentation of the RouteData custom type.
func (rd RouteData) String() string {
	ls := new(common.LogString)
	ls.AddS("RouteData\n")
	ls.AddF("Loaded: %t\n", rd.Loaded)
	ls.AddS("INDEX: Area\n")
	for k, v := range rd.areaID {
		ls.AddF("   %-4d  %s\n", k, v.Name)
	}
	ls.AddS("INDEX: Provider\n")
	for k, v := range rd.providerID {
		ls.AddF("   %-4d  %-40s %s\n", k, v.Name, v.URL)
	}
	ls.AddS("INDEX: Service\n")
	for k, v := range rd.serviceID {
		ls.AddF("   %-4d (%t)  %s\n", k, k == v.ID, v.Name)
	}
	ls.AddS("INDEX: City Services\n")
	for k, v := range rd.cityServices {
		ls.AddF("   [[%s]]\n", k)
		for _, sv := range v {
			ls.AddF("      %s\n", sv)
		}
	}
	ls.AddS("INDEX: Provider Service\n")
	for k, v := range rd.providerService {
		ls.AddF("   %-4d %-50.50s  %-40s\n", rd.serviceID[k].ID, rd.serviceID[k].Name, v.Name)
	}
	ls.AddS("Categories\n")
	for i, v := range rd.Categories {
		ls.AddF("   %2d  %s\n", i, v)
	}
	ls.AddS("---AREAS:\n")
	for k, v := range rd.Areas {
		ls.AddF("[%s]\n%s\n", k, v)
	}
	return ls.Box(90)
}

// Index() builds all required map indexes: Services by City,
func (rd *RouteData) index() error {
	fmt.Printf("[RouteData] building all indexes:\n")
	rd.indexServicesByLoc()
	rd.indexServiceProviderByID()
	rd.indexID()
	return nil
}

func (rd *RouteData) indexServicesByLoc() error {
	fmt.Printf("    building ServicesByLoc index...\n")

	return nil
}

func (rd *RouteData) indexServiceProviderByID() error {
	fmt.Printf("    building ServiceProviderByID index...\n")
	return nil
}

func (rd *RouteData) indexID() error {
	fmt.Printf("    building areaID index...\n")
	for areaKey, area := range rd.Areas {
		rd.areaID[area.ID] = area
		rd.cityServices[areaKey] = make([]*Service, 0)
		for _, provider := range area.Providers {
			rd.providerID[provider.ID] = provider
			// fmt.Println("*** BUILDING SERVICES")
			for _, service := range provider.Services {
				service.ProviderID = provider.ID
				rd.cityServices[areaKey] = append(rd.cityServices[areaKey], service)
				// fmt.Printf("   %s ===> %+v\n", serviceName, service)
				rd.serviceID[service.ID] = service
				rd.providerService[service.ID] = provider
			}
		}
	}
	return nil
}

// ------------------------------- Area -------------------------------

// Area is a Service Area.  It contains an index list of all of the Service Providers for this Area.
type Area struct {
	ID        int                  `json:"id"`
	Name      string               `json:"name"`
	Providers map[string]*Provider `json:"serviceProviders"`
}

func (a Area) String() string {
	ls := new(common.LogString)
	ls.AddF("%s (ID: %d)\n", a.Name, a.ID)
	for k, v := range a.Providers {
		ls.AddF("[%s]\n%s\n", k, v)
	}
	return ls.Box(85)
}

// ------------------------------- Provider -------------------------------

// Provider is the data for each Service Provider.  It contains an index list of all of the Services provided by this Provider.
type Provider struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	InterfaceType string     `json:"interfaceType"`
	URL           string     `json:"url"`
	Key           string     `json:"key"`
	Services      []*Service `json:"services"`
}

func (p Provider) String() string {
	ls := new(common.LogString)
	ls.AddF("%s (ID: %d)\n", p.Name, p.ID)
	ls.AddF("InterfaceType: %s  URL: %s  Key: %s\n", p.InterfaceType, p.URL, p.Key)
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
