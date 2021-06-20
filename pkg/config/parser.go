package config

import (
	"encoding/json"
	"io/ioutil"
)

func ReadConfig(configPath string) (*Config, error) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(configFile, config); err != nil {
		return nil, err
	}
	return config, nil
}
