package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/internal/server/services"
	fallback "github.com/googleapis/grpc-fallback-go/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// RuntimeConfig has the run-time settings necessary to run the
// ms-project servers.
type RuntimeConfig struct {
	port         int
	httpPort     int
	fallbackPort int
	tlsCaCert    string
	tlsCert      string
	tlsKey       string
}

func accessRoles() map[string][]string {
	const movieServicePath = "/movie.MovieService/"
	const identityServicePath = "/identity.IdentityService/"
	const authServicePath = "/auth.AuthService/"

	return map[string][]string{

		// Roles for IdentityService
		identityServicePath + "ListUsers": {"ADMIN"},
		// identityServicePath + "GetUser":    {"ADMIN", "GUEST", "NORMAL", "SUBSCRIBED"},
		// identityServicePath + "CreateUser": {"ADMIN", "GUEST", "NORMAL", "SUBSCRIBED"},
		identityServicePath + "UpdateUser": {"ADMIN", "NORMAL", "SUBSCRIBED"},
		identityServicePath + "DeleteUser": {"ADMIN", "NORMAL", "SUBSCRIBED"},

		// Roles for AuthService
		// authServicePath + "Login": {"GUEST"},

		// Roles for MovieService
		// movieServicePath + "ListMovies":  {"ADMIN", "GUEST", "NORMAL", "SUBSCRIBED"},
		// movieServicePath + "GetMovie":    {"ADMIN", "GUEST", "NORMAL", "SUBSCRIBED"},
		movieServicePath + "CreateMovie": {"ADMIN", "SUBSCRIBED"},
		movieServicePath + "UpdateMovie": {"ADMIN", "SUBSCRIBED"},
		movieServicePath + "DeleteMovie": {"ADMIN", "SUBSCRIBED"},
	}
}

// createBackends creates services used by both the gRPC and REST servers.
func createBackends() *services.Backend {

	identitySrv := services.NewIdentityServer()
	authSrv := services.NewAuthServer(identitySrv)
	movieSrv := services.NewMovieServer(authSrv)

	jm := server.NewJWTManager(services.SecretKey, 5*time.Minute)
	authSrv.JWT = jm
	authI := server.NewAuthInterceptor(jm, accessRoles())

	logger := &loggerObserver{}
	observerRegistry := server.ShowcaseObserverRegistry()
	observerRegistry.RegisterUnaryObserver(logger)
	observerRegistry.RegisterStreamRequestObserver(logger)
	observerRegistry.RegisterStreamResponseObserver(logger)

	return &services.Backend{
		IdentityServer: identitySrv,
		AuthServer:     authSrv,
		MovieServer:    movieSrv,
		Interceptor:    authI,

		StdLog: stdLog,
		ErrLog: errLog,

		ObserverRegistry: observerRegistry,
	}
}

type Servers struct {
	Backend        *services.Backend
	gRPCServer     *grpc.Server
	httpServer     *http.Server
	fallbackServer *fallback.FallbackServer

	httpListener net.Listener
	gRPCListener net.Listener
}

func (s *Servers) initiateServers(endpoint string, config RuntimeConfig, backend *services.Backend) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.initiateGRPCServer(endpoint, config)
	})

	g.Go(func() error {
		return s.initiateHTTPServer(endpoint, config)
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	log.Println("received shutdown signal")
	s.closeServers()
}

// Clean-up process for both HTTP and gRPC servers
func (s *Servers) closeServers() {
	// Shutdown the REST server
	if s.httpServer != nil {
		log.Println("Shutting httpServer!")
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			stdLog.Printf("Error while stopping REST server! %v", err)
		}
		stdLog.Printf("Stopped the REST server!")
		s.httpServer = nil
	}

	// Shutdown the gRPC server
	if s.gRPCServer != nil {
		log.Println("Shutting gRPCServer!")
		s.gRPCServer.GracefulStop()
		stdLog.Printf("Stopped the gRPC server!")
		s.gRPCServer = nil
	}
}

// Starts the gRPC Server
func (s *Servers) initiateGRPCServer(endpoint string, config RuntimeConfig) error {
	addr := fmt.Sprintf("%s:%d", endpoint, config.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("gRPC server: failed to listen: %v", err)
		os.Exit(2)
	}
	s.gRPCListener = ln

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			s.Backend.Interceptor.Unary(),
			s.Backend.ObserverRegistry.UnaryInterceptor,
		),
		// MaxConnectionAge is just to avoid long connection, to facilitate load balancing
		// MaxConnectionAgeGrace will torn them, default to infinity
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
	}

	// load mutual TLS cert/key and root CA cert
	if config.tlsCaCert != "" && config.tlsCert != "" && config.tlsKey != "" {
		keyPair, err := tls.LoadX509KeyPair(config.tlsCert, config.tlsKey)
		if err != nil {
			log.Fatalf("Failed to load server TLS cert/key with error:%v", err)
		}

		cert, err := ioutil.ReadFile(config.tlsCaCert)
		if err != nil {
			log.Fatalf("Failed to load root CA cert file with error:%v", err)
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(cert)

		ta := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{keyPair},
			ClientCAs:    pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		})

		opts = append(opts, grpc.Creds(ta))
	}

	s.gRPCServer = grpc.NewServer(opts...)

	s.registerGRPCService()
	log.Printf("gRPC server serving at %s", addr)

	fb := fallback.NewServer(fmt.Sprintf("%s:%d", endpoint, config.fallbackPort), fmt.Sprintf("%s:%d", endpoint, config.port))
	s.fallbackServer = fb

	// Register reflection service on gRPC server.
	reflection.Register(s.gRPCServer)

	return s.gRPCServer.Serve(ln)
}

// Starts the REST/HTTP Server
func (s *Servers) initiateHTTPServer(endpoint string, config RuntimeConfig) error {
	addr := fmt.Sprintf("%s:%d", endpoint, config.httpPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("REST server: failed to listen: %v", err)
	}
	s.httpListener = lis

	mux := runtime.NewServeMux()
	dialAddr := fmt.Sprintf(":%d", config.port)
	s.registerHTTPService(dialAddr, mux)

	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	s.httpServer = httpSrv
	log.Printf("REST server serving at %s", addr)

	if err = s.httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Register all the services required for gRPC server
func (s *Servers) registerGRPCService() {

	identitypb.RegisterIdentityServiceServer(s.gRPCServer, s.Backend.IdentityServer)
	authpb.RegisterAuthServiceServer(s.gRPCServer, s.Backend.AuthServer)
	moviepb.RegisterMovieServiceServer(s.gRPCServer, s.Backend.MovieServer)
}

// Register all the services required for HTTP/REST server
func (s *Servers) registerHTTPService(endpoint string, mux *runtime.ServeMux) error {

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithTimeout(2 * time.Second),
	}

	err := identitypb.RegisterIdentityServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Printf("failed to Register Identity HTTP server: %v", err)
		return err
	}
	err = authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Printf("failed to Register Auth HTTP server: %v", err)
		return err
	}
	err = moviepb.RegisterMovieServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Printf("failed to Register Movie HTTP server: %v", err)
		return err
	}
	return nil
}
