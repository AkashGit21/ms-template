package services

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"testing"
// 	"time"

// 	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
// 	"github.com/AkashGit21/ms-project/internal/server"
// )

// func TestAuthService(t *testing.T) {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	identitySrv := NewIdentityServer()
// 	authSrv := NewAuthServer(identitySrv)
// 	_, err := validateUser(ctx, identitySrv, authSrv)
// 	if err != nil {
// 		t.Errorf(err.Error())
// 		return
// 	}
// }

// func generateAuthConfig(srv *authServer) *TestConfig {
// 	return &TestConfig{
// 		Server: srv,
// 		URL:    "/v1/auth",
// 		Body: &authpb.LoginRequest{
// 			Username: "usrname1",
// 			Password: "abcdef",
// 		},
// 	}
// }

// func testLogin(ctx context.Context, config *TestConfig) (string, error) {

// 	// Convert srv from interface to server
// 	srv := config.Server.(*authServer)
// 	creds := config.Body.(*authpb.LoginRequest)

// 	resp, err := srv.Login(ctx, &authpb.LoginRequest{
// 		Username: creds.Username,
// 		Password: creds.Password,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	// Verify Response is not empty
// 	if resp.AccessToken == "" {
// 		return "", fmt.Errorf("Empty Response of POST call!")
// 	}

// 	return resp.AccessToken, nil
// }

// func validateUser(ctx context.Context, iSrv *identityServer, srv *authServer) (string, error) {

// 	userConf := generateIdentityConfig(iSrv)
// 	_, err := testPostUser(ctx, userConf)
// 	if err != nil {
// 		return "", err
// 	}

// 	srv.JWT = server.NewJWTManager(SecretKey, 5*time.Minute)
// 	authConf := generateAuthConfig(srv)
// 	token, err := testLogin(ctx, authConf)
// 	if err != nil {
// 		return "", err
// 	}

// 	if _, err = srv.JWT.GetUserFromToken(token); err != nil {
// 		log.Printf("Unable to get user")
// 		return "", err
// 	}
// 	return token, nil
// }
