package configuration

import (
	"encoding/json"
	"os"
)

var (
	RestfulEPDefault = "localhost"
)

type ServiceConfig struct {
	RestfulEndpoint string `json:"restfulapi_endpoint"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	config := ServiceConfig{
		RestfulEPDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}

	err = json.NewDecoder(file).Decode(&config)
	return config, err
}
