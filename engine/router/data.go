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

	log "github.com/jeffizhungry/logrus"
)

var (
	adapters Adapters
	routes   routeData
)

// GetChAreaAdp returns the channel used to update the areaAdapters index.
func GetChAreaAdp() chan map[string][]string {
	return adapters.chUpdate
}

// GetChRouteUpd returns the channel used to update the Route Data indexed by AreaID (received from Services).
func GetChRouteUpd() chan map[structs.NRoute]bool {
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

// GetAllActiveRoutes returns a list of Routes that have been returned by the Service Cache
// queries - hence all "active" routes.
func GetAllActiveRoutes() structs.NRoutes {
	return routes.getAllRoutes()
}

// GetAllRoutes returns a list of all configured routes.
func GetAllRoutes() structs.NRoutes {
	return adapters.getAllRoutes()
}

// GetRouteAdapter gets a pointer to the Adapter for the specified route.
func GetRouteAdapter(route structs.NRoute) (AdpRPCer, error) {
	return adapters.getRouteAdapter(route)
}

// ==============================================================================================================================
//                                      ROUTES
// ==============================================================================================================================

// routeData is the list of currently active Routes.
type routeData struct {
	indArea   [2]map[string]structs.NRoutes // Index: AreaID
	all       [2][]structs.NRoute           // Index: All Routes
	activeSet int
	chUpdate  chan map[structs.NRoute]bool // Channel to receive Area indexed updates from Services.
	sync.RWMutex
}

// getArea retrieves the RouteList for the specified area.
func (r *routeData) getAllRoutes() structs.NRoutes {
	r.RLock()
	defer r.RUnlock()
	log.WithFields(log.Fields{
		"activeSet":  r.activeSet,
		"len(r.all)": len(r.all[r.activeSet]),
	}).Debug("getAllRoutes")
	return r.all[r.activeSet]
}

// getArea retrieves the RouteList for the specified area.
func (r *routeData) getAreaRoutes(areaID string) (structs.NRoutes, error) {
	r.RLock()
	defer r.RUnlock()
	l, ok := r.indArea[r.activeSet][areaID]
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

func (r *routeData) update(upd map[structs.NRoute]bool) {
	r.Lock()
	defer r.Unlock()

	ls := r.loadSet()

	r.clearLoadSet(ls)

	for route := range upd {
		log.Debugf("Route: %v", route)
		r.all[ls] = append(r.all[ls], route)
		if _, ok := r.indArea[ls][route.AreaID]; !ok {
			r.indArea[ls][route.AreaID] = make(structs.NRoutes, 0)
		}
		r.indArea[ls][route.AreaID] = append(r.indArea[ls][route.AreaID], route)
	}

	r.switchSet()
	log.Debug("Updated routeData: " + r.String())
}

func (r *routeData) loadSet() (ls int) {
	if r.activeSet == 0 {
		ls = 1
	}
	return
}

func (r *routeData) clearLoadSet(ds int) {
	r.all[r.loadSet()] = make([]structs.NRoute, 0)
	r.indArea[r.loadSet()] = make(map[string]structs.NRoutes)
}

func (r *routeData) switchSet() {
	if r.activeSet == 0 {
		log.Debug("[RouteData] Switched from list 0 to 1")
		r.activeSet = 1
	} else {
		log.Debug("[RouteData] Switched from list 1 to 0")
		r.activeSet = 0
	}
}

// String displays the contents of the Spec_Type custom type.
func (r routeData) String() string {
	ls := new(common.LogString)
	ls.AddF("routeData [%d]\n", r.activeSet)
	ls.AddF("All Routes: %v\n", r.all[r.activeSet])
	for k, v := range r.indArea[r.activeSet] {
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

func (r *Adapters) getAllRoutes() (routes structs.NRoutes) {
	for _, adp := range r.Adapters {
		routes = append(routes, structs.NRoute{
			AdpID:      adp.ID,
			AreaID:     "all",
			ProviderID: 0,
		})
	}
	return
}

func (r *Adapters) areaID(alias string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	area, ok := r.areaAlias[strings.ToLower(alias)]
	if !ok {
		return "", fmt.Errorf("Cannot find area: %q", alias)
	}
	return area.ID, nil
}

// getRouteAdapter gets a pointer to the Adapter servicing the specifed NRoute.
func (r *Adapters) getRouteAdapter(route structs.NRoute) (*Adapter, error) {
	r.RLock()
	defer r.RUnlock()
	adp, ok := r.Adapters[route.AdpID]
	if !ok {
		return nil, fmt.Errorf("cannot find the Adapter for route: %s", route)
	}
	return adp, nil
}

// getArea retrieves the ServiceList for the specified area.
func (r *Adapters) getAreaAdapters(areaID string) ([]*Adapter, error) {
	r.RLock()
	defer r.RUnlock()
	l, ok := r.areaAdapters[areaID]
	if !ok {
		return nil, fmt.Errorf("The requested AreaID: %q is not serviced by this gateway.", areaID)
	}
	return l, nil
}

// getAdapterid retrieves the Adapterid from a Mid.
func (r *Adapters) getAdapter(id string) (*Adapter, error) {
	r.RLock()
	defer r.RUnlock()
	a, ok := r.Adapters[id]
	if !ok {
		return nil, fmt.Errorf("Adapter: %q was not found.", id)
	}
	return a, nil
}

// getAdapterID retrieves the AdapterID from a MID.
func (r *Adapters) getAdapterID(MID string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	AdpID, _, _, _, err := structs.SplitMID(MID)
	if err != nil {
		return "", fmt.Errorf("The requested ServiceID: %q is not serviced by this gateway.", MID)
	}
	// a, ok := r.Adapters[AdpID]
	_, ok := r.Adapters[AdpID]
	if !ok {
		return "", fmt.Errorf("The requested ServiceID: %q is not serviced by this gateway.", MID)
	}
	return AdpID, nil
}

// Load loads the specified byte slice into the adapters structures.
func (r *Adapters) load(file []byte) error {
	r.Lock()
	defer r.Unlock()
	err := json.Unmarshal(file, r)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal config data file.\nError: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	// Denormalize the Adapters.
	for k, v := range r.Adapters {
		v.ID = k
	}

	// Denormalize the Areas.
	for k, v := range r.Areas {
		v.ID = k
	}

	r.indexAreaAlias()

	r.loaded = true
	r.loadedAt = time.Now()

	// log.Debugf("=================== Adapters ===============\n%s\n\n\n", spew.Sdump(*r))
	// log.Debug("")
	return nil
}

// indexAreaAdapters builds the areaAlias index.
func (r *Adapters) indexAreaAlias() error {
	r.areaAlias = make(map[string]*Area)
	for _, v := range r.Areas {
		for _, alias := range v.Aliases {
			r.areaAlias[alias] = v
		}
	}
	return nil
}

func (r *Adapters) updateAreaAdapters(input map[string][]string) error {
	r.Lock()
	defer r.Unlock()
	r.areaAdapters = make(map[string][]*Adapter)

	for areaID, adpList := range input {
		for _, adpID := range adpList {
			// log.Debugf("AreaID: %q  AdapterID: %q", areaID, adpID)
			if _, ok := r.areaAdapters[areaID]; !ok {
				r.areaAdapters[areaID] = make([]*Adapter, 0)
			}
			r.areaAdapters[areaID] = append(r.areaAdapters[areaID], r.Adapters[adpID])
		}
	}

	log.Debugf("After updateAreaAdapters...\n" + r.String() + "\n")

	return nil
}

// Connect asks each adapter to Dial it's Server.
func (r *Adapters) connect() error {
	for _, v := range r.Adapters {
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
	connected bool
	client    *rpc.Client
}

func (adp *Adapter) connect() error {
	client, err := rpc.DialHTTP("tcp", adp.Address)
	if err != nil {
		log.WithFields(log.Fields{
			"adapter": adp.ID,
			"error":   err.Error(),
		}).Error("Failed to connect to adapter - ")
		return err
	}
	log.WithFields(log.Fields{
		"adapter": adp.ID,
	}).Info("Established connection to adapter - ")
	adp.client = client
	adp.connected = true
	return nil
}

// --------------------------- AdpRPCer Interface ----------------------------------------

// AdpRPCer is an interface to the Adapter RPC status and Client.
type AdpRPCer interface {
	AdpID() string
	Connected() bool
	Call(serviceMethod string, args interface{}, reply interface{}) error
}

// AdpID returns the Adapter ID
func (adp *Adapter) AdpID() string {
	return adp.ID
}

// Connected returns the current connection status of the adapter RPC connection.
func (adp *Adapter) Connected() bool {
	return adp.connected
}

// Call invokes the RPC Client.Call() function (see https://golang.org/pkg/net/rpc/#Client)
func (adp *Adapter) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return adp.client.Call(serviceMethod, args, reply)
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

	routes.chUpdate = make(chan map[structs.NRoute]bool, 1)
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
func (r Adapters) String() string {
	ls := new(common.LogString)
	ls.AddS("Adapters\n")
	for _, v := range r.Adapters {
		ls.AddS(v.String())
	}
	ls.AddS("\n-------Areas---------------\n")
	for _, v := range r.Areas {
		// ls.AddS("   %5s %-7s  %-25s  [%s]\n", k, fmt.Sprintf("(%s)", v.ID), v.Name, fmt.Sprintf("\"%s\"", strings.Join(v.Aliases, "\", \"")))
		ls.AddF("%s\n", v)
	}
	ls.AddS("\n-------AreaAlias-----------\n")
	for k, v := range r.areaAlias {
		ls.AddF("   %-20s  %s\n", k, v.ID)
	}
	ls.AddS("\n-------AreaAdapters--------\n")
	for k, v := range r.areaAdapters {
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
	ls.AddF("%-17s   Type: %s  Address: %s\n", ls.ColorBool(adp.connected, "CONNECTED  ", "UNCONNECTED", "green", "red"), adp.Type, adp.Address)
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
