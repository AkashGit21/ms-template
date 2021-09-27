package main

import (
	"testing"
	"time"
)

func TestServers(t *testing.T) {

	testServer := &Servers{}
	conf := RuntimeConfig{
		port:         8082,
		httpPort:     8081,
		fallbackPort: 8084,
	}

	testServer.Backend = createBackends()
	go func() {
		testServer.initiateServers("", conf, testServer.Backend)
	}()

	time.Sleep(1 * time.Second)
	testServer.closeServers()
}
