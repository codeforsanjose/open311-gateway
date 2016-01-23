package router

import (
	"_sketches/spew"
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
	"time"

	"Gateway311/engine/common"
)

var (
	adapters Adapters
)

// GetChAreaAdp returns the channel used to update the areaAdapters index.
func GetChAreaAdp() chan map[string][]string {
	return adapters.chAreaAdp
}

// GetAreaAdapters returns a list of the Adapters that provide services to the specified
// Area.
func GetAreaAdapters(areaID string) ([]*Adapter, error) {
	return adapters.getAreaAdapters(areaID)
}

// GetAreaID returns the AreaID for a city name, using the aliases in the config.json file.
func GetAreaID(alias string) (string, error) {
	return adapters.areaID(alias)
}

// ==============================================================================================================================
//                                      ADAPTERS
// ==============================================================================================================================

// Adapters is the list of all active Adapters.
type Adapters struct {
	loaded    bool
	loadedAt  time.Time
	Adapters  map[string]*Adapter `json:"adapters"`
	Areas     map[string]*Area    `json:"areas"`
	chAreaAdp chan map[string][]string

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

	fmt.Printf("=================== Adapters ===============\n%s\n\n\n", spew.Sdump(*adps))
	fmt.Printf("")
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
			log.Debug("AreaID: %q  AdapterID: %q", areaID, adpID)
			if _, ok := adps.areaAdapters[areaID]; !ok {
				adps.areaAdapters[areaID] = make([]*Adapter, 0)
			}
			adps.areaAdapters[areaID] = append(adps.areaAdapters[areaID], adps.Adapters[adpID])
		}
	}

	log.Debug("%s\n", adps)

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
	adapters.chAreaAdp = make(chan map[string][]string, 1)

	go func() {
		for {
			select {
			case aa, ok := <-adapters.chAreaAdp:
				if !ok {
					break
				}
				adapters.updateAreaAdapters(aa)
			}
		}
	}()
}

// ==============================================================================================================================
//                                      string
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
	ls := new(common.LogString)
	ls.AddF("%s\n", adp.ID)
	ls.AddF("Connected: %t Type: %s  Address: %s\n", adp.Connected, adp.Type, adp.Address)
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
