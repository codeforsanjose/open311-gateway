package jx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// NewJSONConfig reads a JSON file into the specified interface (struct or map), and returns
// a pointer to the JSONConfig struct.
func NewJSONConfig(target interface{}, filename string) (*JSONConfig, error) {
	var jc JSONConfig
	return &jc, jc.Load(target, filename)
}

// JSONConfig provides functionality to read a JSON config file into the target struct.
type JSONConfig struct {
	target   interface{}
	filename string
	loaded   bool
}

// Load reads the JSON file into the target struct.  Load will return an error if
// the struct has already been loaded.
func (r *JSONConfig) Load(target interface{}, filename string) error {
	if r.loaded {
		return fmt.Errorf("JSONConfig already loaded")
	}
	r.filename = filename
	r.target = target
	return r.load()
}

// Reload re-reads the JSON file into the target struct.  Load() must be called
// prior to Reload(), otherwise it will return an error.
func (r *JSONConfig) Reload() error {
	if !r.loaded {
		return fmt.Errorf("JSONConfig.Load() must be called prior to Reload()")
	}
	return r.load()
}

func (r *JSONConfig) load() error {
	file, err := ioutil.ReadFile(r.filename)
	if err != nil {
		return fmt.Errorf("ReadConfig failed - %v.", err)
	}

	err = json.Unmarshal(file, r.target)
	if err != nil {
		return fmt.Errorf("ReadConfig failed - %v.", err)
	}

	r.loaded = true

	return nil
}
