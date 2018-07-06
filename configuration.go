package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

const defaultFileName string = ".issuesrc"

// Query is alias for map type
type Query map[string]string

// Configuration struct
type Configuration struct {
	Queries Query  `json:"queries"` // Queries is a dictionary of queries and query commands
	Cache   string `json:"cache"`   // Cache is the duration of the cache, use nocache to no use this feature
}

// CacheDuration parses Cache configuration
// and returns the duration as time.Duration
func (c *Configuration) CacheDuration() time.Duration {
	if c.Cache == "nocache" {
		return time.Duration(0)
	}

	d, err := time.ParseDuration(c.Cache)

	if err != nil {
		return time.Duration(0)
	}
	return d
}

func getConfigFilename() string {
	return path.Join(userHome(), defaultFileName)
}

func configurationExists() bool {
	return fileExists(getConfigFilename())
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
