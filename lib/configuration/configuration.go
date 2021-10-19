package configuration

import (
	"encoding/json"
	"os"

	"github.com/AkashGit21/ms-project/lib/persistence/dblayer"
)

var (
	DBTypeDefault       = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault    = "localhost"
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
