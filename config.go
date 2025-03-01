package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	HttpServerPort int    `json:"httpServerPort"`
	FmaUsername    string `json:"fmaUsername"`
	FmaPassword    string `json:"fmaPassword"`
	FmaURL         string `json:"fmaUrl"`
	DB             string `json:"db"`
}

func (c *config) getHttpServerAddr() string {
	return fmt.Sprintf(":%d", c.HttpServerPort)
}

const configFilename = "config.json"

func loadConfigFile() (config, error) {
	file, err := os.Open(configFilename)
	if err != nil {
		return config{}, fmt.Errorf("error while loading config file: %w", err)
	}

	var result config
	decodeErr := json.NewDecoder(file).Decode(&result)

	if decodeErr != nil {
		return config{}, fmt.Errorf("error while loading config file: config file json parse error: %w", decodeErr)
	}

	return result, nil
}
