package services

import (
	"log"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
)

// Backend contains the various service backends that will be
// accessible via one or more transport endpoints.
type Backend struct {
	// Application schema
	MovieServer    moviepb.MovieServiceServer
	IdentityServer identitypb.IdentityServiceServer

	// Supporting protos

	// Other supporting data structures
	StdLog, ErrLog *log.Logger
}
