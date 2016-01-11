package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"Gateway311/gateway/common"
	"Gateway311/integration/citysourced/structs"

	"github.com/davecgh/go-spew/spew"
)

const iFID = "CS"

var (
	providerData ProviderData
)

// ShowProviderData dumps providerData using spew.
func ShowProviderData() string {
	return spew.Sdump(providerData)
}

// ServicesForCity returns a list of all services available for the specified City.
func ServicesForCity(city string) (structs.NServices, error) {
	return providerData.ServicesForCity(city)
}

// // ServiceProviders returns a list of all Service Providers for the specified City.
// func ServiceProviders(city string) ([]*Provider, error) {
// 	return providerData.ServiceProviders(city)
// }

// // ServiceProvider returns a pointer to the Provider for the specified Provider ID.
// func ServiceProvider(sid int) (*Provider, error) {
// 	return providerData.ServiceProvider(sid)
// }

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
	Categories []string         `json:"serviceCategories"`
	Areas      map[string]*Area `json:"serviceAreas"`

	serviceMID   map[string]dataIndex         // Service MID -> Service
	serviceID    map[int]*structs.NService    // Service ID -> Service
	providerID   map[int]*Provider            // Provider ID -> Provider
	cityCode     map[string]string            // City name to City Code
	cityServices map[string]structs.NServices // City Code (lowercase) -> List of Services
}

type dataIndex struct {
	area     *Area
	provider *Provider
	service  *structs.NService
}

// ServicesForCity returns a list of all services available for the specified City.
func (pd *ProviderData) ServicesForCity(city string) (structs.NServices, error) {
	lcity := strings.ToLower(city)
	fmt.Printf("   Services for: %s...\n", lcity)
	ccode, ok := pd.isValidCity(lcity)
	if !ok {
		fmt.Printf("The city: %q is not serviced by this Gateway", city)
		return nil, fmt.Errorf("The city: %q is not serviced by this Gateway", city)
	}
	fmt.Printf("      data length: %d\n", len(pd.cityServices[ccode]))
	services, ok := pd.cityServices[ccode]
	if !ok {
		fmt.Printf("  NO MATCH!\n")
		return nil, fmt.Errorf("Unable to find requested city")
	}
	return services, nil
}

// // ServiceProviders returns a list of all Service Providers for the specified City.
// func (rd *ProviderData) ServiceProviders(city string) ([]*Provider, error) {
// 	lcity := strings.ToLower(city)
// 	fmt.Printf("   ServiceProviders for: %q\n", lcity)
// 	if !pd.isValidCity(lcity) {
// 		return nil, fmt.Errorf("The city: %q is not serviced by this Gateway", city)
// 	}
// 	var p []*Provider
// 	for _, v := range pd.Areas[strings.ToLower(city)] {
// 		p = append(p, v)
// 	}
// 	fmt.Printf("[ServiceProviders] returning %d records.\n", len(p))
// 	return p, nil
// }

// // ServiceProvider returns a pointer to the Provider for the specified Provider ID.
// func (pd *ProviderData) ServiceProvider(id string) (*Provider, error) {
//
// 	p, ok := pd.serviceProvider[id]
// 	if !ok {
// 		return nil, fmt.Errorf("Invalid Service ID")
// 	}
// 	// fmt.Printf("[ServiceProvider] returning %#v\n", *p)
// 	return p, nil
// }

// Load loads the specified byte slice into the ProviderData structures.
func (pd *ProviderData) Load(file []byte) error {
	err := json.Unmarshal(file, pd)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse JSON Route Data.\nError: %v", err)
		fmt.Println(msg)
		return errors.New(msg)
	}

	pd.settle()
	pd.index()
	pd.Loaded = true
	return nil
}

// Index() builds all required map indexes: Services by City,
func (pd *ProviderData) settle() error {
	fmt.Printf("    settling data...\n")
	for areaKey, area := range pd.Areas {
		for _, provider := range area.Providers {
			for _, service := range provider.Services {
				service.IFID = iFID
				service.AreaID = areaKey
				service.ProviderID = provider.ID
			}
		}
	}
	return nil
}

// Index() builds all required map indexes: Services by City,
func (pd *ProviderData) index() error {
	fmt.Printf("[ProviderData] building indexes:\n")
	pd.indexServiceMID()
	pd.indexServiceID()
	pd.indexProviderID()
	pd.indexCityCode()
	// Run indexCityCode() last.
	pd.indexCityServices()
	return nil
}

func (pd *ProviderData) indexServiceMID() error {
	fmt.Printf("    ServiceID...\n")
	pd.serviceMID = make(map[string]dataIndex)
	for _, area := range pd.Areas {
		for _, provider := range area.Providers {
			for _, service := range provider.Services {
				pd.serviceMID[service.MID()] = dataIndex{area, provider, service}
			}
		}
	}
	return nil
}

func (pd *ProviderData) indexServiceID() error {
	fmt.Printf("    ServiceID...\n")
	pd.serviceID = make(map[int]*structs.NService)
	for _, area := range pd.Areas {
		for _, provider := range area.Providers {
			for _, service := range provider.Services {
				pd.serviceID[service.ID] = service
			}
		}
	}
	return nil
}

func (pd *ProviderData) indexProviderID() error {
	fmt.Printf("    ProviderID...\n")
	pd.providerID = make(map[int]*Provider)
	for _, area := range pd.Areas {
		for _, provider := range area.Providers {
			pd.providerID[provider.ID] = provider
		}
	}
	return nil
}

func (pd *ProviderData) indexCityCode() error {
	fmt.Printf("    CityCode...\n")
	pd.cityCode = make(map[string]string)
	for areaKey, area := range pd.Areas {
		pd.cityCode[strings.ToLower(area.Name)] = areaKey
	}
	return nil
}

func (pd *ProviderData) indexCityServices() error {
	fmt.Printf("    CityServices...\n")
	pd.cityServices = make(map[string]structs.NServices)
	for areaKey, area := range pd.Areas {
		pd.cityServices[areaKey] = make(structs.NServices, 0)
		for _, provider := range area.Providers {
			for _, service := range provider.Services {
				pd.cityServices[areaKey] = append(pd.cityServices[areaKey], *service)
			}
		}
	}
	return nil
}

// func (pd *ProviderData) indexCityServices() error {
// 	fmt.Printf("    CityServices...\n")
// 	pd.cityServices = make(map[string]*structs.NServices)
// 	var nSvcs structs.NServices
// 	for areaKey, area := range pd.Areas {
// 		nSvcs = make(structs.NServices, 0)
// 		for _, provider := range area.Providers {
// 			for _, service := range provider.Services {
// 				nSvcs = append(nSvcs, *service)
// 			}
// 		}
// 		fmt.Printf("=============>>>> CityServices <<<<<==============\nArea: %q\n", areaKey)
// 		fmt.Printf(spew.Sdump(&nSvcs))
// 		pd.cityServices[areaKey] = &nSvcs
// 		fmt.Printf(spew.Sdump(pd.cityServices[areaKey]))
// 	}
// 	for areaKey := range pd.Areas {
// 		fmt.Printf(spew.Sdump(pd.cityServices[areaKey]))
// 	}
// 	return nil
// }

// String returns the represeentation of the ProviderData custom type.
func (pd ProviderData) String() string {
	ls := new(common.LogString)
	ls.AddS("ProviderData\n")
	ls.AddF("Loaded: %t\n", pd.Loaded)
	ls.AddS("\n-----------INDEX: serviceID-----------\n")
	for k, v := range pd.serviceID {
		ls.AddF("   %-4d %s\n", k, v.Name)
	}
	ls.AddS("\n-----------INDEX: serviceMID-----------\n")
	ls.AddS("      MID                 Area           Prov  Service\n")
	for k, v := range pd.serviceMID {
		ls.AddF("   %-20s %-15s %4d %4d   %s\n", k, v.area.Name, v.provider.ID, v.service.ID, v.service.Name)
	}
	ls.AddS("\n-----------INDEX: providerID-----------\n")
	for k, v := range pd.providerID {
		ls.AddF("   %-4d  %-40s %s\n", k, v.Name, v.URL)
	}
	ls.AddS("\n-----------INDEX: cityCode-----------\n")
	for k, v := range pd.cityCode {
		ls.AddF("   %-20s %s\n", k, v)
	}
	ls.AddS("\n-----------INDEX: cityServices-----------\n")
	for k, v := range pd.cityServices {
		ls.AddF("   [[%s]]\n", k)
		for _, sv := range v {
			ls.AddF("      %s\n", sv)
		}
	}
	ls.AddS("\n--- CATEGORIES ---\n")
	for i, v := range pd.Categories {
		ls.AddF("   %2d  %s\n", i, v)
	}
	ls.AddS("\n---AREAS ---\n")
	for _, v := range pd.Areas {
		for _, v2 := range v.Providers {
			ls.AddF("%s\n", v2)
		}
	}
	return ls.Box(90)
}

func (pd *ProviderData) isValidCity(city string) (string, bool) {
	code, ok := pd.cityCode[strings.ToLower(city)]
	return code, ok

}

// ------------------------------- Area -------------------------------

// Area is a Service Area.  It contains an index list of all of the Service Providers for this Area.
type Area struct {
	Name      string      `json:"name"`
	Providers []*Provider `json:"providers"`
}

func (a Area) String() string {
	ls := new(common.LogString)
	ls.AddF("%s\n", a.Name)
	for _, v := range a.Providers {
		ls.AddF("%s\n", v)
	}
	return ls.Box(85)
}

// ------------------------------- Provider -------------------------------

// Provider is the data for each Service Provider.  It contains an index list of all of the Services provided by this Provider.
type Provider struct {
	ID         int                 //
	Name       string              `json:"name"`
	URL        string              `json:"url"`
	APIVersion string              `json:"apiVersion"`
	Key        string              `json:"key"`
	Services   []*structs.NService `json:"services"`
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

// SplitMID breaks down an MID, and returns pointers to the Area and Provider.
func SplitMID(mid string) (*Area, *Provider, error) {
	parts := strings.Split(mid, "-")
	fmt.Printf("MID: %+v\n", parts)
	area := providerData.Areas[parts[1]]
	v, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid MID: %s", mid)
	}
	provider := providerData.providerID[int(v)]
	return area, provider, nil
}
