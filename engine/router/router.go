package router

import (
	"errors"
	"fmt"
	"io/ioutil"
)

// map: city -> list of Services (MID)

// func: MID -> Adapter

// loaded from config file:
// map: Gateway Code -> Adapter API Call
// map: city name -> City Code (MID)

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

// Init loads the config file, and connects to all Adapters.
func Init() error {
	err := readConfig("/Users/james/Dropbox/Development/go/src/Gateway311/engine/router/config.json")
	if err != nil {
		return fmt.Errorf("Error %v occurred when reading the config - ReadConfig()", err)
	}
	return nil
}

func readConfig(filePath string) error {

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		msg := fmt.Sprintf("Unable to access the adapters Config file - specified at: %q.\nError: %v", filePath, err)
		fmt.Println(msg)
		return errors.New(msg)
	}

	err = adapters.load(file)
	if err != nil {
		return err
	}

	err = adapters.connect()
	if err != nil {
		return err
	}
	return nil
}
