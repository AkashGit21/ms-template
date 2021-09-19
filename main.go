package main

import (
	"flag"

	"github.com/AkashGit21/ms-project/cmd/server"
	"github.com/AkashGit21/ms-project/lib/configuration"
)

func main() {

	// Setting the flag for configuration file
	confPath := flag.String("config", `.\configuration\config.json`, "flag to set the path for configuration JSON file")
	flag.Parse()

	// Getting the required configuration
	config, _ := configuration.ExtractConfiguration(*confPath)

	server.StartServer(config.RestfulEndpoint)
}
