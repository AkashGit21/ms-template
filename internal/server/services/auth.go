package services

import (
	"context"
	"log"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/AkashGit21/ms-project/internal/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	SecretKey = "secret"
)

// AuthServer is the server for authentication
type authServer struct {
	authpb.UnimplementedAuthServiceServer

	identityStore ReadOnlyIdentityServer
	JWT           *server.JWTManager
}

// NewAuthServer returns a new auth server
func NewAuthServer(is *identityServer) *authServer {
	return &authServer{
		identityStore: is,
	}
}

// Login is a unary RPC to login user
func (as *authServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {

	log.Println("Login start!")
	user, err := as.identityStore.GetUser(ctx, &identitypb.GetUserRequest{Username: req.GetUsername()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	log.Println("Trying to match Password!")

	if user == nil || !server.DoesPasswordMatch(req.GetPassword(), user.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	log.Println("Generating token!")

	token, err := as.JWT.GenerateToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token! %v", err)
	}

	log.Println("End of Login request!")
	return &authpb.LoginResponse{AccessToken: token}, nil
}
