package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var stdLog, errLog *log.Logger

func init() {
	stdLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errLog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

type loggerObserver struct{}

func (l *loggerObserver) GetName() string { return "loggerObserver" }

func (l *loggerObserver) ObserveUnary(
	ctx context.Context,
	req interface{},
	resp interface{},
	info *grpc.UnaryServerInfo,
	err error) {
	stdLog.Printf("Received Unary Request for Method: %s\n", info.FullMethod)
	if Verbose {
		dumpIncomingHeaders(ctx)
	}
	stdLog.Printf("    Request:  %+v\n", req)
	if err == nil {
		stdLog.Printf("    Returning Response: %+v\n", resp)
	} else {
		stdLog.Printf("    Returning Error: %+v\n", err)
	}
	stdLog.Println("")
}

func dumpIncomingHeaders(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		stdLog.Printf("Cannot get metadata from the context.")
		return
	}

	stdLog.Printf("    Request headers:")
	for key, values := range md {
		for _, value := range values {
			stdLog.Printf("      %s: %s\n", key, value)
		}
	}
}
