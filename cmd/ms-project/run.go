package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func init() {
	config := RuntimeConfig{
		httpPort:     8081,
		port:         "8080",
		fallbackPort: "8082",
	}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Starts running the application server",
		Run: func(cmd *cobra.Command, args []string) {
			appMuxServer := createAllEndpoints(config)
			done := make(chan os.Signal, 2)
			signal.Notify(done, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
			go func() {
				sig := <-done
				stdLog.Printf("Got signal %q", sig)
				stdLog.Printf("Shutting down server: %s", message(appMuxServer.Shutdown()))

				os.Exit(1)
			}()

			stdLog.Printf("Server finished: %s", message(appMuxServer.Serve()))
		},
	}

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(
		&config.port,
		"port",
		"p",
		":7469",
		"The port that application will be served on.")
	runCmd.Flags().IntVarP(
		&config.httpPort,
		"httpPort",
		"r",
		8081,
		"The port that REST will be served on.")
	runCmd.Flags().StringVarP(
		&config.fallbackPort,
		"fallback-port",
		"f",
		":1337",
		"The port that the fallback-proxy will be served on.")
}

func message(err error) string {
	if err == nil {
		return "ok"
	}
	return err.Error()
}
