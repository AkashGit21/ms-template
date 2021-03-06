package services

import (
	"log"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/internal/server/interceptors"
)

// Backend contains the various service backends that will be
// accessible via one or more transport endpoints.
type Backend struct {
	// Application schema
	MovieServer    moviepb.MovieServiceServer
	IdentityServer identitypb.IdentityServiceServer
	AuthServer     authpb.AuthServiceServer

	AuthInterceptor *interceptors.AuthInterceptor

	// Supporting protos

	// Other supporting data structures
	StdLog, ErrLog   *log.Logger
	ObserverRegistry server.GrpcObserverRegistry
}
