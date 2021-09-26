package main

import (
	"github.com/spf13/cobra"
)

func init() {
	config := RuntimeConfig{
		port:         8082,
		httpPort:     8081,
		fallbackPort: 8084,
	}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Starts running the application server",
		Run: func(cmd *cobra.Command, args []string) {

			srvs := &Servers{}
			backend := createBackends()
			srvs.Backend = backend
			srvs.InitiateServers("", config, backend)
		},
	}

	rootCmd.AddCommand(runCmd)
}

func message(err error) string {
	if err == nil {
		return "ok"
	}
	return err.Error()
}
