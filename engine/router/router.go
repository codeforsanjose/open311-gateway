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
	// log.Debug(adapters.String())
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

	// log.Debug(adapters.String())
	return nil
}
