package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func runGRPCServer(mSrv moviepb.MovieServiceServer, lis net.Listener) error {

	gRPCServer := grpc.NewServer()
	moviepb.RegisterMovieServiceServer(gRPCServer, mSrv)

	err := gRPCServer.Serve(lis)
	if err != nil {
		log.Fatal("Failed to Serve\n ")
	}

	return err
}

func runRESTServer(mSrv moviepb.MovieServiceServer, lis net.Listener, endpoint string) error {

	mux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := moviepb.RegisterMovieServiceHandlerServer(ctx, mux, mSrv)
	if err != nil {
		log.Println("Failed to Register Movie Handler Server")
		return err
	}
	log.Printf("Started REST server at %s", lis.Addr().String())

	return http.Serve(lis, mux)
}

func StartServer() {
	port := flag.Int("port", 8080, "Server Port to be used")
	endpoint := flag.String("endpoint", "localhost", "gRPC endpoint for the REST API")
	serverType := flag.String("type", "rest", "Type of Server - grpc/rest")
	flag.Parse()

	log.Println("Initiating Server...")

	mServer := service.NewMovieServer()

	address := fmt.Sprintf("%s:%d", *endpoint, *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if *serverType == "rest" {
		runRESTServer(mServer, listener, *endpoint)
	} else if *serverType == "grpc" {
		runGRPCServer(mServer, listener)
	}
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
