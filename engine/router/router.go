package router

import (
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/jeffizhungry/logrus"
)

// ==============================================================================================================================
//                                      INIT
// ==============================================================================================================================

// Init loads the config files.
func Init(configFile string) error {
	if err := readConfig(configFile); err != nil {
		return err
	}
	log.Debug("Adapters: " + adapters.String())
	err := adapters.connect()
	if err != nil {
		return err
	}
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

	return nil
}
