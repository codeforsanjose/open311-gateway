package router

import (
	"Gateway311/engine/common"
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"
	"time"
)

var (
	adapters Adapters
)

// ==============================================================================================================================
//                                      ADAPTERS
// ==============================================================================================================================

// Adapters is the list of all active Adapters.
type Adapters struct {
	loaded   bool
	loadedAt time.Time
	Adapters []*Adapter `json:"adapters"`

	areaAdapters map[string][]*Adapter
}

// Load loads the specified byte slice into the adapters structures.
func (adps *Adapters) load(file []byte) error {
	// adps.init()
	err := json.Unmarshal(file, adps)
	if err != nil {
		msg := fmt.Sprintf("Unable to unmarshal config data file.\nError: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	// adps.index()
	adps.loaded = true
	adps.loadedAt = time.Now()
	return nil
}

func (adps *Adapters) connect() error {
	for _, v := range adps.Adapters {
		v.connect()
	}
	return nil
}

// String returns a formatted representation of Adapters.
func (adps Adapters) String() string {
	ls := new(common.LogString)
	ls.AddS("Adapters")
	for _, v := range adps.Adapters {
		ls.AddS(v.String())
	}
	return ls.Box(90)
}

// ==============================================================================================================================
//                                      ADAPTER
// ==============================================================================================================================

// Adapter represents an active Adapter.
type Adapter struct {
	Name      string `json:"name"`
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
		log.Errorf("Connection to: %s failed: %s", adp.Name, err)
		return err
	}
	log.Info("Connection to: %q OK!\n", adp.Name)
	adp.Client = client
	adp.Connected = true
	return nil
}

// String returns a formatted representation of Adapter.
func (adp Adapter) String() string {
	ls := new(common.LogString)
	ls.AddF("%s\n", adp.Name)
	ls.AddF("Connected: %t Type: %s  Address: %s\n", adp.Connected, adp.Type, adp.Address)
	ls.AddF("File: %s\n", adp.File)
	ls.AddF("Config: %s\n", adp.Config)
	return ls.Box(80)
}
