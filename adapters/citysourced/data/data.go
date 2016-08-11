package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/open311-gateway/adapters/citysourced/common"
	"github.com/open311-gateway/adapters/citysourced/structs"

	"github.com/davecgh/go-spew/spew"
	log "github.com/jeffizhungry/logrus"
)

var (
	configData ConfigData
)

// ShowConfigData dumps configData using spew.
func ShowConfigData() string {
	return spew.Sdump(configData)
}

// ServicesArea returns a list of all services available for the specified Area.
func ServicesArea(area string) (*structs.NServices, error) {
	larea := strings.ToLower(area)
	log.Debugf("   Services for: %s...\n", larea)
	ccode, ok := configData.isValidCity(larea)
	if !ok {
		msg := fmt.Sprintf("The area: %q is not serviced by this Gateway", area)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	log.Debugf("      data length: %d\n", len(configData.areaServices[ccode]))
	services, ok := configData.areaServices[ccode]
	if !ok {
		msg := fmt.Sprintf("Unable to find requested area: %q", area)
		log.Warning(msg)
		return nil, errors.New(msg)
	}
	return &services, nil
}

// ServicesAll returns a list of ALL services.
func ServicesAll() (*structs.NServices, error) {
	resp := make(structs.NServices, 0)
	for _, v := range configData.areaServices {
		resp = append(resp, v...)
	}
	return &resp, nil
}

// Adapter returns the adapter configuration.
func Adapter() (name, atype, address string) {
	return configData.Adapter.Name, configData.Adapter.Type, configData.Adapter.Address
}

// AdapterName returns the adapter name.
func AdapterName() string {
	return configData.Adapter.Name
}

// MIDProvider returns the Provider data for the specified MidAdpID.
func MIDProvider(MID structs.ServiceID) (Provider, error) {
	log.Debugf("MID: %s", MID.MID())
	return getProvider(MID.AreaID, MID.ProviderID)
}

// RouteProvider returns the Provider data for the specified NRoute.
func RouteProvider(route structs.NRoute) (Provider, error) {
	log.Debugf("Route: %s", route)
	return getProvider(route.AreaID, route.ProviderID)
}

// ServiceFromID returns the NService data for the specified ServiceID.
func ServiceFromID(srvID structs.ServiceID) (nsrv structs.NService, err error) {
	log.Debugf("ServiceID: %s", srvID.MID())
	sm, ok := configData.serviceMID[srvID.MID()]
	if !ok {
		err = fmt.Errorf("invalid ServiceID: %s", srvID.MID())
		return
	}
	nsrv = *sm.service
	log.Debugf("Returning: %s", nsrv.String())
	return
}

// getProvider returns the Provider data for the specified Area and Provider.
func getProvider(AreaID string, ProviderID int) (Provider, error) {
	log.Debugf("AreaID: %v  ProviderID: %v\n", AreaID, ProviderID)
	p, ok := configData.areaProvider[areaProvider{AreaID, ProviderID}]
	// log.Debugf("Provider (%t): %s", ok, *p)
	if !ok {
		return Provider{}, fmt.Errorf("Unable to find Provider for %v-%v", AreaID, ProviderID)
	}
	return *p, nil
}

// GetMonitorAddress returns the Telemetry Address from the config file.
func GetMonitorAddress() string {
	return configData.Monitor.Address
}

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

// Init loads the config files.
func Init(configFile string) error {
	if err := readConfig(configFile); err != nil {
		return err
	}
	return nil
}

func readConfig(filePath string) error {
	if configData.Loaded {
		msg := "Route Data is already loaded"
		fmt.Println(msg)
		return errors.New(msg)
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		msg := fmt.Sprintf("Unable to access the config file - %v.", err)
		log.Error(msg)
		return errors.New(msg)
	}

	return configData.Load(file)
}

// ==============================================================================================================================
//                                      ROUTE DATA
// ==============================================================================================================================

type areaProvider struct {
	areaID     string
	providerID int
}

// ConfigData is a list of all the Service Areas.  It contains an indexed list of all the Service Areas.  The index is the *lowercase* area name.
type ConfigData struct {
	Loaded  bool
	Adapter AdapterData `json:"adapter"`
	Monitor struct {
		Address string `json:"address"`
	} `json:"monitor"`
	Categories []string         `json:"serviceCategories"`
	Areas      map[string]*Area `json:"serviceAreas"`

	serviceMID   map[string]dataIndex      // Service MID -> Service
	serviceID    map[int]*structs.NService // Service ID -> Service
	providerID   map[int]*Provider         // Provider ID -> Provider
	areaProvider map[areaProvider]*Provider
	areaCode     map[string]string // City name to City Code

	areaServices map[string]structs.NServices // City Code (lowercase) -> List of Services
}

// AdapterData contains all of the config data.
type AdapterData struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

type dataIndex struct {
	area     *Area
	provider *Provider
	service  *structs.NService
}

// Load loads the specified byte slice into the ConfigData structures.
func (pd *ConfigData) Load(file []byte) error {
	err := json.Unmarshal(file, pd)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse JSON Route Data.\nError: %v", err)
		fmt.Println(msg)
		return errors.New(msg)
	}
	log.Info("Initializing data...")
	_ = pd.settle()
	_ = pd.index()
	pd.Loaded = true
	log.Debug(ShowConfigData())
	return nil
}

// Index() builds all required map indexes: Services by City,
func (pd *ConfigData) settle() error {
	log.Info("   Denormalizing service keys...\n")
	for areaKey, area := range pd.Areas {
		area.ID = areaKey
		for _, provider := range area.Providers {
			for _, service := range provider.Services {
				service.AdpID = pd.Adapter.Name
				service.AreaID = areaKey
				service.ProviderID = provider.ID
				if service.ResponseType == "" {
					service.ResponseType = provider.ResponseType
				}
			}
		}
	}
	return nil
}

// Index() builds all required map indexes: Services by City,
func (pd *ConfigData) index() error {
	log.Info("   Building indexes:\n")
	_ = pd.indexServiceMID()
	_ = pd.indexServiceID()
	_ = pd.indexProviderID()
	_ = pd.indexCityCode()
	_ = pd.indexAreaProvider()
	// Run indexCityCode() last.
	_ = pd.indexCityServices()
	return nil
}

func (pd *ConfigData) indexServiceMID() error {
	log.Info("       Indexing ServiceID...\n")
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

func (pd *ConfigData) indexServiceID() error {
	log.Info("       Indexing ServiceID...\n")
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

func (pd *ConfigData) indexProviderID() error {
	log.Info("       Indexing ProviderID...\n")
	pd.providerID = make(map[int]*Provider)
	for _, area := range pd.Areas {
		for _, provider := range area.Providers {
			pd.providerID[provider.ID] = provider
		}
	}
	return nil
}

func (pd *ConfigData) indexCityCode() error {
	log.Info("       Indexing CityCode...\n")
	pd.areaCode = make(map[string]string)
	for areaKey, area := range pd.Areas {
		pd.areaCode[strings.ToLower(area.Name)] = areaKey
	}
	return nil
}

func (pd *ConfigData) indexAreaProvider() error {
	log.Info("       Indexing AreaProvider...\n")
	pd.areaProvider = make(map[areaProvider]*Provider)
	for _, area := range pd.Areas {
		for _, provider := range area.Providers {
			pd.areaProvider[areaProvider{area.ID, provider.ID}] = provider
		}
	}
	return nil
}

func (pd *ConfigData) indexCityServices() error {
	log.Info("       Indexing CityServices...\n")
	pd.areaServices = make(map[string]structs.NServices)
	for areaKey, area := range pd.Areas {
		pd.areaServices[areaKey] = make(structs.NServices, 0)
		for _, provider := range area.Providers {
			for _, service := range provider.Services {
				pd.areaServices[areaKey] = append(pd.areaServices[areaKey], *service)
			}
		}
	}
	return nil
}

// String returns the represeentation of the ConfigData custom type.
func (pd ConfigData) String() string {
	ls := new(common.LogString)
	ls.AddF("[%s] ConfigData\n", pd.Adapter.Name)
	ls.AddF("Loaded: %t\n", pd.Loaded)
	ls.AddF("Adapter: %s   Type: %s   Address: %s\n", pd.Adapter.Name, pd.Adapter.Type, pd.Adapter.Address)
	ls.AddF("Monitor - address: %s\n", pd.Monitor.Address)
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
	ls.AddS("\n-----------INDEX: areaCode-----------\n")
	for k, v := range pd.areaCode {
		ls.AddF("   %-20s %s\n", k, v)
	}
	ls.AddS("\n-----------INDEX: areaProvider-----------\n")
	for k, v := range pd.areaProvider {
		ls.AddF("   %-20s  %-40s %s\n", fmt.Sprintf("%s-%d", k.areaID, k.providerID), v.Name, v.URL)
	}
	ls.AddS("\n-----------INDEX: areaServices-----------\n")
	for k, v := range pd.areaServices {
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
		ls.AddF("%s\n", v)
	}
	return ls.Box(90)
}

func (pd *ConfigData) isValidCity(area string) (string, bool) {
	code, ok := pd.areaCode[strings.ToLower(area)]
	return code, ok

}

// ------------------------------- Area -------------------------------

// Area is a Service Area.  It contains an index list of all of the Service Providers for this Area.
type Area struct {
	ID        string
	Name      string      `json:"name"`
	Providers []*Provider `json:"providers"`
}

func (a Area) String() string {
	ls := new(common.LogString)
	ls.AddF("%s (%s)\n", a.Name, a.ID)
	for _, v := range a.Providers {
		ls.AddF("%s\n", v)
	}
	return ls.Box(85)
}

// ------------------------------- Provider -------------------------------

// Provider is the data for each Service Provider.  It contains an index list of all of the Services provided by this Provider.
type Provider struct {
	ID           int                 //
	Name         string              `json:"name"`
	URL          string              `json:"url"`
	APIVersion   string              `json:"apiVersion"`
	Key          string              `json:"key"`
	ResponseType string              `json:"responseType"`
	Services     []*structs.NService `json:"services"`
}

func (p Provider) String() string {
	ls := new(common.LogString)
	ls.AddF("%s (ID: %d)\n", p.Name, p.ID)
	ls.AddF("URL: %s  ver: %s  Key: %s\n", p.URL, p.APIVersion, p.Key)
	ls.AddS("---SERVICES:\n")
	for _, v := range p.Services {
		ls.AddF("   %s\n", v)
	}
	return ls.Box(80)
}

// SplitMID breaks down an MID, and returns pointers to the Area and Provider.
func SplitMID(mid string) (*Area, *Provider, error) {
	parts := strings.Split(mid, "-")
	log.Debugf("MID: %+v\n", parts)
	area := configData.Areas[parts[1]]
	v, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("Invalid MID: %s", mid)
	}
	provider := configData.providerID[int(v)]
	return area, provider, nil
}
