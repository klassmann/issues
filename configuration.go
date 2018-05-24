package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

const defaultFileName string = ".issuesrc"

// Query is alias for map type
type Query map[string]string

// Configuration struct
type Configuration struct {
	Queries Query `json:"queries"`
}

func userHome() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func getConfigFilename() string {
	return path.Join(userHome(), defaultFileName)
}

func configurationExists() bool {
	_, err := os.Stat(getConfigFilename())
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func createConfiguration() error {
	var conf = Configuration{}
	conf.Queries = make(Query)
	return saveConfiguration(&conf)
}

func saveConfiguration(conf *Configuration) error {

	if conf == nil {
		return fmt.Errorf("save configuration: conf is nil")
	}

	data, err := json.MarshalIndent(*conf, "", "  ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(getConfigFilename(), data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func loadConfiguration() (*Configuration, error) {
	var conf = Configuration{}

	data, err := ioutil.ReadFile(getConfigFilename())

	if err != nil {
		return nil, fmt.Errorf("loadConfiguration: reading file: %v", err)
	}

	err = json.Unmarshal(data, &conf)

	if err != nil {
		return nil, fmt.Errorf("loadConfiguration: json unmarshiling data: %v", err)
	}

	return &conf, nil
}
