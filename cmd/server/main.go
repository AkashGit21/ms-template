package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	version = "no version"

	httpPort   = flag.Int("httpPort", 8888, "http port")
	grpcPort   = flag.Int("grpcPort", 9200, "grpc port")
	healthPort = flag.Int("healthPort", 6666, "grpc health port")

	logLevel     = flag.String("logLevel", "INFO", "Log level, INFO|WARNING|DEBUG|ERROR")
	gcpProjectID = flag.String("gcpProjectID", "none", "GCP Project ID")

	dev = flag.Bool("dev", false, "Enable development tooling")

	grpcServer       *grpc.Server
	grpcHealthServer *grpc.Server
	httpServer       *http.Server
)

func StartServer() {
	// port := flag.Int("port", 8080, "Server Port to be used")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		addr := fmt.Sprintf(":%d", *grpcPort)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("gRPC server: failed to listen: %v", err)
			os.Exit(2)
		}

		server := service.NewMovieServer()
		grpcServer = grpc.NewServer(
			// MaxConnectionAge is just to avoid long connection, to facilitate load balancing
			// MaxConnectionAgeGrace will torn them, default to infinity
			grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
		)
		moviepb.RegisterMovieServiceServer(grpcServer, server)

		log.Printf("gRPC server serving at %s", addr)

		return grpcServer.Serve(ln)
	})

	g.Go(func() error {
		addr := fmt.Sprintf(":%d", *httpPort)

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("REST server: failed to listen: %v", err)
		}

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithTimeout(2 * time.Second),
		}

		err = moviepb.RegisterMovieServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", *grpcPort), opts)
		if err != nil {
			log.Printf("failed to Register Movie server: %v", err)
			return err
		}

		httpServer = &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}
		log.Printf("REST server serving at %s", addr)

		if err = httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Println("received shutdown signal")

	cancel()

	if httpServer != nil {
		log.Println("Shutting httpServer!")
		_ = httpServer.Shutdown(ctx)
	}
	if grpcServer != nil {
		log.Println("Shutting gRPCServer!")
		grpcServer.GracefulStop()
	}

	err := g.Wait()
	if err != nil {
		log.Println("server returning an error: ", err)
		os.Exit(2)
	}

}
