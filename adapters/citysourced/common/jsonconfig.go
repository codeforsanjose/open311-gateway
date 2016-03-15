package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// JSONConfig provides functionality to read a JSON config file into the parent struct.
type JSONConfig struct {
	parent   interface{}
	filename string
	loaded   bool
}

// Read loads the config file into the parent struct.
func (r *JSONConfig) Read(parent interface{}, filePath string) error {
	r.filename = filePath
	r.parent = parent

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ReadConfig failed - %v.", err)
	}

	err = json.Unmarshal(file, parent)
	if err != nil {
		return fmt.Errorf("ReadConfig failed - %v.", err)
	}

	r.loaded = true

	return nil
}
