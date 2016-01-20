package router

import (
	"errors"
	"fmt"
	"io/ioutil"

	"Gateway311/engine/logs"
)

var (
	log = logs.Log
)

// map: city -> list of Services (MID)

// func: MID -> Adapter

// loaded from config file:
// map: Gateway Code -> Adapter API Call
// map: city name -> City Code (MID)

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

// Init loads the config files.
func Init(configFile string) error {
	if err := readConfig(configFile); err != nil {
		return err
	}
	err := adapters.connect()
	if err != nil {
		return err
	}
	log.Debug(adapters.String())

	// servicesData.Refresh()
	return nil
}

func readConfig(filePath string) error {

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		msg := fmt.Sprintf("Unable to access the config file - %v.", err)
		log.Error(msg)
		return errors.New(msg)
	}

	err = adapters.load(file)
	if err != nil {
		return err
	}

	log.Debug(adapters.String())
	return nil
}
