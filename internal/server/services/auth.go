package services

import (
	"context"
	"log"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/internal/server/interceptors"
	"github.com/AkashGit21/ms-project/lib/persistence"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	SecretKey = "secret"
)

// AuthServer is the server for authentication
type authServer struct {
	authpb.UnimplementedAuthServiceServer

	dbhandler     persistence.DatabaseHandler
	identityStore ReadOnlyIdentityServer
	JWT           *server.JWTManager
}

// NewAuthServer returns a new auth server
func NewAuthServer(is *identityServer) *authServer {
	return &authServer{
		identityStore: is,
		dbhandler:     is.dbhandler,
	}
}

// Login is a unary RPC to login user
func (as *authServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {

	log.Println("Beginning of Login! ", req)
	interceptors.CURRENT_USERNAME = req.GetUsername()
	user, err := as.identityStore.GetUser(ctx, &identitypb.GetUserRequest{Username: req.GetUsername()})
	if err != nil {
		log.Println("Error: ", err)
		return nil, status.Errorf(codes.InvalidArgument, "incorrect username/password!")
	}

	if !as.dbhandler.Authenticate(req.GetUsername(), req.GetPassword()) {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect username/password!")
	}

	token, err := as.JWT.GenerateToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token! %v", err)
	}

	log.Println("End of Login request!")
	return &authpb.LoginResponse{AccessToken: token}, nil
}

// TODO: Add Logout functionality
func (as *authServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {

	return &authpb.LogoutResponse{}, nil
}
