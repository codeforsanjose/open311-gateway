package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"Gateway311/engine/common"
	"Gateway311/engine/structs"
)

var (
	adapters Adapters
	routes   routeData
)

// GetChAreaAdp returns the channel used to update the areaAdapters index.
func GetChAreaAdp() chan map[string][]string {
	return adapters.chUpdate
}

// GetChRouteUpd returns the channel used to update the Route Data (received from Services).
func GetChRouteUpd() chan map[string]structs.NRoutes {
	return routes.chUpdate
}

// GetAreaAdapters returns a list of the Adapters that provide services to the specified
// Area.
func GetAreaAdapters(areaID string) ([]*Adapter, error) {
	return adapters.getAreaAdapters(areaID)
}

// GetAreaID returns the AreaID for a Area (City) name, using the aliases in the config.json file.
func GetAreaID(alias string) (string, error) {
	return adapters.areaID(alias)
}

// GetAdapterID retrieves the AdapterID from a MID.
func GetAdapterID(MID string) (string, error) {
	return adapters.getAdapterID(MID)
}

// GetAdapter retrieves a pointer to the Adapter from a ID.
func GetAdapter(id string) (*Adapter, error) {
	return adapters.getAdapter(id)
}

// GetAreaRoutes returns a list of Routes (structs.NRoutes) for the specified Area.
func GetAreaRoutes(areaID string) (structs.NRoutes, error) {
	return routes.getAreaRoutes(areaID)
}

// ==============================================================================================================================
//                                      ROUTES
// ==============================================================================================================================

// routeData is the list of currently active Routes.
type routeData struct {
	routes    [2]map[string]structs.NRoutes // Index: AreaID
	activeSet int
	chUpdate  chan map[string]structs.NRoutes // Channel to receive updates from Services.
	sync.RWMutex
}

// getArea retrieves the ServiceList for the specified area.
func (r *routeData) getAreaRoutes(areaID string) (structs.NRoutes, error) {
	r.RLock()
	defer r.RUnlock()
	l, ok := r.routes[r.activeSet][areaID]
	if !ok {
		return nil, fmt.Errorf("There are no routes for AreaID: %q", areaID)
	}
	return l, nil
}

// getArea retrieves the ServiceList for the specified area.
func (r *routeData) validateRoute(route structs.NRoute) bool {
	r.RLock()
	defer r.RUnlock()

	return false
}

func (r *routeData) update(upd map[string]structs.NRoutes) {
	r.routes[r.loadSet()] = upd
	r.switchSet()
	log.Debug("Updated routeData!%s", r)
}

func (r *routeData) loadSet() (ls int) {
	if r.activeSet == 0 {
		ls = 1
	}
	return
}

func (r *routeData) clearLoadSet(ds int) {
	r.routes[r.loadSet()] = make(map[string]structs.NRoutes)
}

func (r *routeData) switchSet() {
	r.Lock()
	defer r.Unlock()
	if r.activeSet == 0 {
		log.Debug("Switched from list 0 to 1")
		r.activeSet = 1
		r.clearLoadSet(0)
	} else {
		log.Debug("Switched from list 1 to 0")
		r.activeSet = 0
		r.clearLoadSet(1)
	}
}

// String displays the contents of the Spec_Type custom type.
func (r routeData) String() string {
	ls := new(common.LogString)
	ls.AddF("routeData [%d]\n", r.activeSet)
	for k, v := range r.routes[r.activeSet] {
		ls.AddF("<<<<<Area: %s >>>>>%s", k, v)
	}
	return ls.Box(90)
}

// ==============================================================================================================================
//                                      ADAPTERS
// ==============================================================================================================================

// Adapters is the list of all active Adapters.
type Adapters struct {
	loaded   bool
	loadedAt time.Time
	Adapters map[string]*Adapter `json:"adapters"` // Index: AdpID
	Areas    map[string]*Area    `json:"areas"`    // Index: AreaID
	chUpdate chan map[string][]string

	areaAlias    map[string]*Area      // Index: an alias for an area
	areaAdapters map[string][]*Adapter // Index: AreaID
	sync.RWMutex
}

func (adps *Adapters) areaID(alias string) (string, error) {
	adps.RLock()
	defer adps.RUnlock()
	area, ok := adps.areaAlias[strings.ToLower(alias)]
	if !ok {
		return "", fmt.Errorf("Cannot find area: %q", alias)
	}
	return area.ID, nil
}

// getArea retrieves the ServiceList for the specified area.
func (adps *Adapters) getAreaAdapters(areaID string) ([]*Adapter, error) {
	adps.RLock()
	defer adps.RUnlock()
	l, ok := adps.areaAdapters[areaID]
	if !ok {
		return nil, fmt.Errorf("The requested AreaID: %q is not serviced by this gateway.", areaID)
	}
	return l, nil
}

// getAdapterid retrieves the Adapterid from a Mid.
func (adps *Adapters) getAdapter(id string) (*Adapter, error) {
	adps.RLock()
	defer adps.RUnlock()
	a, ok := adps.Adapters[id]
	// log.Debug("a: %s-%s  ok: %t\n", a.ID, a.Type, ok)
	if !ok {
		return nil, fmt.Errorf("Adapter: %q was not found.", id)
	}
	return a, nil
}

// getAdapterID retrieves the AdapterID from a MID.
func (adps *Adapters) getAdapterID(MID string) (string, error) {
	adps.RLock()
	defer adps.RUnlock()
	AdpID, _, _, _, err := structs.SplitMID(MID)
	if err != nil {
		return "", fmt.Errorf("The requested ServiceID: %q is not serviced by this gateway.", MID)
	}
	// a, ok := adps.Adapters[AdpID]
	// log.Debug("AdpID: %q  a: %s-%s  ok: %t\n", AdpID, a.ID, a.Type, ok)
	_, ok := adps.Adapters[AdpID]
	if !ok {
		return "", fmt.Errorf("The requested ServiceID: %q is not serviced by this gateway.", MID)
	}
	return AdpID, nil
}

// Load loads the specified byte slice into the adapters structures.
func (adps *Adapters) load(file []byte) error {
	adps.Lock()
	defer adps.Unlock()
	err := json.Unmarshal(file, adps)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal config data file.\nError: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	// Denormalize the Adapters.
	for k, v := range adps.Adapters {
		v.ID = k
	}

	// Denormalize the Areas.
	for k, v := range adps.Areas {
		v.ID = k
	}

	adps.indexAreaAlias()

	adps.loaded = true
	adps.loadedAt = time.Now()

	// log.Debug("=================== Adapters ===============\n%s\n\n\n", spew.Sdump(*adps))
	// log.Debug("")
	return nil
}

// indexAreaAdapters builds the areaAlias index.
func (adps *Adapters) indexAreaAlias() error {
	adps.areaAlias = make(map[string]*Area)
	for _, v := range adps.Areas {
		for _, alias := range v.Aliases {
			adps.areaAlias[alias] = v
		}
	}
	return nil
}

func (adps *Adapters) updateAreaAdapters(input map[string][]string) error {
	adps.Lock()
	defer adps.Unlock()
	adps.areaAdapters = make(map[string][]*Adapter)

	for areaID, adpList := range input {
		for _, adpID := range adpList {
			// log.Debug("AreaID: %q  AdapterID: %q", areaID, adpID)
			if _, ok := adps.areaAdapters[areaID]; !ok {
				adps.areaAdapters[areaID] = make([]*Adapter, 0)
			}
			adps.areaAdapters[areaID] = append(adps.areaAdapters[areaID], adps.Adapters[adpID])
		}
	}

	log.Debug("After updateAreaAdapters...\n%s\n", adps)

	return nil
}

// Connect asks each adapter to Dial it's Server.
func (adps *Adapters) connect() error {
	for _, v := range adps.Adapters {
		v.connect()
	}
	return nil
}

// ==============================================================================================================================
//                                      ADAPTER
// ==============================================================================================================================

// Adapter represents an active Adapter.
type Adapter struct {
	ID        string //
	Type      string `json:"type"`
	Address   string `json:"address"`
	File      string `json:"file"`
	Config    string `json:"config"`
	Connected bool
	Client    *rpc.Client
}

func (adp *Adapter) connect() error {
	client, err := rpc.DialHTTP("tcp", adp.Address)
	if err != nil {
		log.Errorf("Connection to: %s failed: %s", adp.ID, err)
		return err
	}
	log.Info("Connection to: %q OK!\n", adp.ID)
	adp.Client = client
	adp.Connected = true
	return nil
}

// ==============================================================================================================================
//                                      AREA
// ==============================================================================================================================

// Area represents a Service Area.
type Area struct {
	ID      string   //
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

func init() {
	adapters.chUpdate = make(chan map[string][]string, 1)
	go func() {
		for {
			select {
			case aa, ok := <-adapters.chUpdate:
				if !ok {
					break
				}
				adapters.updateAreaAdapters(aa)
			}
		}
	}()

	routes.chUpdate = make(chan map[string]structs.NRoutes, 1)
	go func() {
		for {
			select {
			case upd, ok := <-routes.chUpdate:
				if !ok {
					break
				}
				routes.update(upd)
			}
		}
	}()
}

// ==============================================================================================================================
//                                      STRINGs
// ==============================================================================================================================

// String returns a formatted representation of Adapters.
func (adps Adapters) String() string {
	ls := new(common.LogString)
	ls.AddS("Adapters\n")
	for _, v := range adps.Adapters {
		ls.AddS(v.String())
	}
	ls.AddS("\n-------Areas---------------\n")
	for _, v := range adps.Areas {
		// ls.AddS("   %5s %-7s  %-25s  [%s]\n", k, fmt.Sprintf("(%s)", v.ID), v.Name, fmt.Sprintf("\"%s\"", strings.Join(v.Aliases, "\", \"")))
		ls.AddF("%s\n", v)
	}
	ls.AddS("\n-------AreaAlias-----------\n")
	for k, v := range adps.areaAlias {
		ls.AddF("   %-20s  %s\n", k, v.ID)
	}
	ls.AddS("\n-------AreaAdapters--------\n")
	for k, v := range adps.areaAdapters {
		var s []string
		for _, a := range v {
			s = append(s, a.ID)
		}
		ls.AddF("   %-5s  %s\n", k, strings.Join(s, ", "))
	}
	return ls.Box(90)
}

// String returns a formatted representation of Adapter.
func (adp Adapter) String() string {
	// ls := new(common.LogString)
	ls := common.NewLogString()
	ls.AddF("%s\n", adp.ID)
	ls.AddF("%-17s   Type: %s  Address: %s\n", ls.ColorBool(adp.Connected, "CONNECTED  ", "UNCONNECTED", "green", "red"), adp.Type, adp.Address)
	ls.AddF("File: %s\n", adp.File)
	ls.AddF("Config: %s\n", adp.Config)
	return ls.Box(80)
}

// String returns a formatted representation of Adapter.
func (a Area) String() string {
	ls := new(common.LogString)
	ls.AddF("%s\n", a.ID)
	ls.AddF("Name: %s\n", a.Name)
	ls.AddF("Aliases: \"%s\"\n", strings.Join(a.Aliases, "\", \""))
	return ls.Box(80)
}
