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

func (l *loggerObserver) ObserveStreamRequest(
	ctx context.Context,
	req interface{},
	info *grpc.StreamServerInfo,
	_ error) {
	stdLog.Printf("%s Stream for Method: %s\n", streamType(info), info.FullMethod)
	if Verbose {
		dumpIncomingHeaders(ctx)
	}
	stdLog.Printf("    Receiving Message:  %v\n", req)
	stdLog.Println("")
}

func (l *loggerObserver) ObserveStreamResponse(
	_ context.Context,
	resp interface{},
	info *grpc.StreamServerInfo,
	_ error) {
	stdLog.Printf("%s Stream for Method: %s\n", streamType(info), info.FullMethod)
	stdLog.Printf("    Sending Message:  %+v\n", resp)
	stdLog.Println("")
}

func streamType(info *grpc.StreamServerInfo) string {
	if info.IsClientStream && info.IsServerStream {
		return "Bi-directional"
	} else if info.IsClientStream {
		return "Client"
	}
	return "Server"
}
