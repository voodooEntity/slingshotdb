package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"slingshot/types"
)

var Conf = types.Config{}

func init() {
	err := getConfig()
	if nil != err {
		fmt.Println("> Database exiting")
		os.Exit(1)
	}
}

func getConfig() error {
	// first we read the json data
	data, err := ioutil.ReadFile("config.json")
	if nil != err {
		fmt.Println("> Config file could not be found or is not readable")
		os.Exit(1)
	}
	// now we parse the config contents
	// lets see if the body json is valid tho
	err = json.Unmarshal(data, &Conf)
	if nil != err {
		return errors.New("> Config file could not be found or is not readable")
	}
	return nil
}
