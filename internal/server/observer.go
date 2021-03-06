package server

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

// UnaryObserver provides an interface for observing unary requests and responses.
type UnaryObserver interface {
	GetName() string
	ObserveUnary(
		ctx context.Context,
		req interface{},
		resp interface{},
		info *grpc.UnaryServerInfo,
		err error)
}

// StreamRequestObserver provides an interface for observing streaming requests.
type StreamRequestObserver interface {
	GetName() string
	ObserveStreamRequest(
		ctx context.Context,
		req interface{},
		info *grpc.StreamServerInfo,
		err error)
}

// StreamResponseObserver provides an interface for observing streaming responses.
type StreamResponseObserver interface {
	GetName() string
	ObserveStreamResponse(
		ctx context.Context,
		resp interface{},
		info *grpc.StreamServerInfo,
		err error)
}

// GrpcObserverRegistry is a registry of observers. These observers are hooked into the
// grpc interceptors that are provided by this interface.
type GrpcObserverRegistry interface {
	// UnaryInterceptor implements the grpc.UnaryServerInterceptor type to allow the
	// registry to hook into unary grpc methods.
	UnaryInterceptor(
		context.Context,
		interface{},
		*grpc.UnaryServerInfo,
		grpc.UnaryHandler) (interface{}, error)
	// StreamInterceptor implements the grpc.StreamServerInterceptor type to allow the
	// registry to hook into streaming grpc methods.
	StreamInterceptor(
		interface{},
		grpc.ServerStream,
		*grpc.StreamServerInfo,
		grpc.StreamHandler) error
	RegisterUnaryObserver(UnaryObserver)
	DeleteUnaryObserver(name string)
	RegisterStreamRequestObserver(StreamRequestObserver)
	DeleteStreamRequestObserver(name string)
	RegisterStreamResponseObserver(StreamResponseObserver)
	DeleteStreamResponseObserver(name string)
}

// ShowcaseObserverRegistry returns the showcase specific observer registry.
func ShowcaseObserverRegistry() GrpcObserverRegistry {
	return &showcaseObserverRegistry{
		uObservers:     map[string]UnaryObserver{},
		sReqObservers:  map[string]StreamRequestObserver{},
		sRespObservers: map[string]StreamResponseObserver{},
	}
}

// showcaseObserverRegistry is an implementation of the ObserverRegistry. This registry
// automatically handles DeleteTest requests and deletes the appropriate observers
// for that request.
type showcaseObserverRegistry struct {
	mu             sync.Mutex
	uObservers     map[string]UnaryObserver
	sReqObservers  map[string]StreamRequestObserver
	sRespObservers map[string]StreamResponseObserver
}

func (r *showcaseObserverRegistry) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	resp, err := handler(ctx, req)

	for _, obs := range r.uObservers {
		obs.ObserveUnary(ctx, req, resp, info, err)
	}

	return resp, err
}

type showcaseStream struct {
	info     *grpc.StreamServerInfo
	registry *showcaseObserverRegistry

	grpc.ServerStream
}

func (s *showcaseStream) SendMsg(m interface{}) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	err := s.ServerStream.SendMsg(m)
	for _, obs := range s.registry.sRespObservers {
		obs.ObserveStreamResponse(s.ServerStream.Context(), m, s.info, err)
	}
	return err
}

func (s *showcaseStream) RecvMsg(m interface{}) error {
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()

	err := s.ServerStream.RecvMsg(m)
	for _, obs := range s.registry.sReqObservers {
		obs.ObserveStreamRequest(s.ServerStream.Context(), m, s.info, err)
	}
	return err
}

func (r *showcaseObserverRegistry) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return handler(srv, &showcaseStream{info, r, ss})
}

// RegisterUnaryObserver registers a unary observer. If an observer of the same name
// has already been registered, the new observer will override it.
func (r *showcaseObserverRegistry) RegisterUnaryObserver(obs UnaryObserver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.uObservers[obs.GetName()] = obs
}

func (r *showcaseObserverRegistry) DeleteUnaryObserver(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.uObservers, name)
}

// RegisterStreamRequestObserver registers a stream observer. If an observer of the same name
// has already been registered, the new observer will override it.
func (r *showcaseObserverRegistry) RegisterStreamRequestObserver(obs StreamRequestObserver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sReqObservers[obs.GetName()] = obs
}

func (r *showcaseObserverRegistry) DeleteStreamRequestObserver(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sReqObservers, name)
}

// RegisterStreamResponseObserver registers a stream observer. If an observer of the same name
// has already been registered, the new observer will override it.
func (r *showcaseObserverRegistry) RegisterStreamResponseObserver(obs StreamResponseObserver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sRespObservers[obs.GetName()] = obs
}

func (r *showcaseObserverRegistry) DeleteStreamResponseObserver(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sRespObservers, name)
}
