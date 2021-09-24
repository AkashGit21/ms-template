package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/internal/server/services"
	fallback "github.com/googleapis/grpc-fallback-go/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// RuntimeConfig has the run-time settings necessary to run the
// ms-project servers.
type RuntimeConfig struct {
	httpPort     int
	port         string
	fallbackPort string
	// tlsCaCert    string
	// tlsCert      string
	// tlsKey       string
}

// Endpoint defines common operations for any of the various types of
// transport-specific network endpoints the application supports
type Endpoint interface {
	// Serve beings the listen-and-serve loop for this
	// Endpoint. It typically blocks until the server is shut
	// down. The error it returns depends on the underlying
	// implementation.
	Serve() error

	// Shutdown causes the currently running Endpoint to
	// terminate. The error it returns depends on the underlying
	// implementation.
	Shutdown() error
}

// CreateAllEndpoints returns an Endpoint that can serve gRPC and
// HTTP/REST connections (on conf.port) and gRPC-fallback
// connections (on conf.fallbackPort)
func createAllEndpoints(conf RuntimeConfig) Endpoint {
	// Ensure the port is of the right form.
	if !strings.HasPrefix(conf.port, ":") {
		conf.port = ":" + conf.port
	}

	// Start listening.
	lis, err := net.Listen("tcp", conf.port)
	if err != nil {
		log.Fatalf("Server failed to listen on port '%s': %v", conf.port, err)
	}
	stdLog.Printf("Server listening on port: %s", conf.port)

	m := cmux.New(lis)
	httpListener := m.Match(cmux.HTTP1Fast())
	// // cmux.Any() is needed below to get mTLS to work for
	// // gRPC, and that in turn means the order of the matchers matters. See
	// // https://github.com/open-telemetry/opentelemetry-collector/issues/2732
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))

	backend := createBackends()
	gRPCServer := newEndpointGRPC(grpcListener, conf, backend)
	restServer := newEndpointREST(httpListener, conf, backend)
	cmuxServer := newEndpointMux(m, gRPCServer, restServer)
	return cmuxServer
}

// endpointMux is an Endpoint for cmux, the connection multiplexer
// allowing different types of connections on the same port.
//
// We choose not to use grpc.Server.ServeHTTP because it is
// experimental and does not support some gRPC features available
// through grpc.Server.Serve. (cf
// https://godoc.org/google.golang.org/grpc#Server.ServeHTTP)
type endpointMux struct {
	endpoints []Endpoint
	cmux      cmux.CMux
	mux       sync.Mutex
}

func newEndpointMux(cmuxEndpoint cmux.CMux, endpoints ...Endpoint) Endpoint {
	return &endpointMux{
		endpoints: endpoints,
		cmux:      cmuxEndpoint,
	}
}

func (em *endpointMux) String() string {
	return "endpoint multiplexer"
}

func (em *endpointMux) Serve() error {
	g := new(errgroup.Group)
	for idx, endpt := range em.endpoints {
		if endpt != nil {
			stdLog.Printf("Starting endpoint %d: %s", idx, endpt)
			endpoint := endpt
			g.Go(func() error {
				err := endpoint.Serve()
				err2 := em.Shutdown()
				if err != nil {
					return err
				}
				return err2
			})
		}
	}
	if em.cmux != nil {
		stdLog.Printf("Starting %s", em)

		g.Go(func() error {
			err := em.cmux.Serve()
			err2 := em.Shutdown()
			if err != nil {
				return err
			}
			return err2

		})
	}
	return g.Wait()
}

func (em *endpointMux) Shutdown() error {
	em.mux.Lock()
	defer em.mux.Unlock()

	var err error
	if em.cmux != nil {
		// TODO: Wait for https://github.com/soheilhy/cmux/pull/69 (due to
		// https://github.com/soheilhy/cmux/pull/69#issuecomment-712928041.)
		//
		// err = em.mux.Close()
		em.cmux = nil
	}

	for idx, endpt := range em.endpoints {
		if endpt != nil {
			// TODO: Wait for https://github.com/soheilhy/cmux/pull/69
			// newErr := endpt.Shutdown()
			// if err==nil {
			// 	err = newErr
			// }
			em.endpoints[idx] = nil
		}
	}
	return err
}

// endpointGRPC is an Endpoint for gRPC connections to the Showcase
// server.
type endpointGRPC struct {
	server         *grpc.Server
	fallbackServer *fallback.FallbackServer
	listener       net.Listener
	mux            sync.Mutex
}

// createBackends creates services used by both the gRPC and REST servers.
func createBackends() *services.Backend {
	logger := &loggerObserver{}
	observerRegistry := server.ShowcaseObserverRegistry()
	observerRegistry.RegisterUnaryObserver(logger)
	// observerRegistry.RegisterStreamRequestObserver(logger)
	// observerRegistry.RegisterStreamResponseObserver(logger)

	return &services.Backend{
		ObserverRegistry: observerRegistry,
		IdentityServer:   services.NewIdentityServer(),
		MovieServer:      services.NewMovieServer(),
		StdLog:           stdLog,
		ErrLog:           errLog,
	}
}

func newEndpointGRPC(lis net.Listener, config RuntimeConfig, backend *services.Backend) Endpoint {

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(backend.ObserverRegistry.StreamInterceptor),
		grpc.UnaryInterceptor(backend.ObserverRegistry.UnaryInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
	}

	s := grpc.NewServer(opts...)

	// Register Services to the server.
	identitypb.RegisterIdentityServiceServer(s, backend.IdentityServer)
	moviepb.RegisterMovieServiceServer(s, backend.MovieServer)

	fb := fallback.NewServer(config.fallbackPort, "localhost"+config.port)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	return &endpointGRPC{
		server:         s,
		fallbackServer: fb,
		listener:       lis,
	}
}

func (eg *endpointGRPC) String() string {
	return "gRPC endpoint"
}

func (eg *endpointGRPC) Serve() error {
	defer eg.Shutdown()
	if eg.fallbackServer != nil {
		stdLog.Printf("Listening for gRPC-fallback connections")
		eg.fallbackServer.StartBackground()
	}
	if eg.server != nil {
		stdLog.Printf("Listening for gRPC connections")
		return eg.server.Serve(eg.listener)
	}
	return fmt.Errorf("gRPC server not set up")
}

func (eg *endpointGRPC) Shutdown() error {
	eg.mux.Lock()
	defer eg.mux.Unlock()

	if eg.fallbackServer != nil {
		stdLog.Printf("Stopping gRPC-fallback connections")
		eg.fallbackServer.Shutdown()
		eg.fallbackServer = nil
	}

	if eg.server != nil {
		stdLog.Printf("Stopping gRPC connections")
		eg.server.GracefulStop()
		eg.server = nil
	}
	stdLog.Printf("Stopped gRPC")
	return nil
}

// endpointREST is an Endpoint for HTTP/REST connections to the ms-project
// server.
type endpointREST struct {
	server   *http.Server
	listener net.Listener
	mux      sync.Mutex
}

func newEndpointREST(lis net.Listener, config RuntimeConfig, backend *services.Backend) *endpointREST {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()

	err := identitypb.RegisterIdentityServiceHandlerServer(ctx, mux, backend.IdentityServer)
	if err != nil {
		log.Printf("failed to Register Identity server: %v", err)
		return nil
	}

	err = moviepb.RegisterMovieServiceHandlerServer(ctx, mux, backend.MovieServer)
	if err != nil {
		log.Printf("failed to Register Movie server: %v", err)
		return nil
	}

	addr := fmt.Sprintf("localhost:%d", config.httpPort)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &endpointREST{
		server:   httpServer,
		listener: lis,
	}
}

func (er *endpointREST) String() string {
	return "HTTP/REST endpoint"
}

func (er *endpointREST) Serve() error {
	defer er.Shutdown()
	if er.server != nil {
		stdLog.Printf("Listening for REST connections")
		return er.server.Serve(er.listener)
	}
	return nil
}

func (er *endpointREST) Shutdown() error {
	er.mux.Lock()
	defer er.mux.Unlock()
	var err error
	if er.server != nil {
		stdLog.Printf("Stopping REST connections")
		err = er.server.Shutdown(context.Background())
		er.server = nil
	}
	stdLog.Printf("Stopped REST")
	return err
}
