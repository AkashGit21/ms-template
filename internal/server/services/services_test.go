package services

import (
	"log"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/server"
)

type TestBackend struct {
	// Application schema
	MovieServer    moviepb.MovieServiceServer
	IdentityServer identitypb.IdentityServiceServer
	AuthServer     authpb.AuthServiceServer

	Interceptor *server.AuthInterceptor

	// Supporting protos

	// Other supporting data structures
	StdLog, ErrLog   *log.Logger
	ObserverRegistry server.GrpcObserverRegistry
}
