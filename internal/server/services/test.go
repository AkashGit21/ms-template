package services

import (
	"context"
	"log"
	"net"
	"time"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/lib/configuration"
	"github.com/AkashGit21/ms-project/lib/persistence/dblayer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	TestIdentitySrv *identityServer
	TestAuthSrv     *authServer
	TestMovieSrv    *movieServer
)

type TestCase struct {
	name        string
	args        interface{}
	expected    interface{}
	expectedErr string
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	grpcServer := grpc.NewServer()
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	TestIdentitySrv = NewIdentityServer(dbhandler)
	TestAuthSrv = NewAuthServer(TestIdentitySrv)
	TestAuthSrv.JWT = server.NewJWTManager(SecretKey, 2*time.Minute)
	TestMovieSrv = NewMovieServer(TestAuthSrv)

	identitypb.RegisterIdentityServiceServer(grpcServer, TestIdentitySrv)
	authpb.RegisterAuthServiceServer(grpcServer, TestAuthSrv)
	moviepb.RegisterMovieServiceServer(grpcServer, TestMovieSrv)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func removeTestData() {

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	identityClient := identitypb.NewIdentityServiceClient(conn)
	// movieClient := moviepb.NewMovieServiceClient(conn)

	identityClient.DeleteUser(ctx, &identitypb.DeleteUserRequest{Username: "test_username"})
	identityClient.DeleteUser(ctx, &identitypb.DeleteUserRequest{Username: "test_get_username"})
	identityClient.DeleteUser(ctx, &identitypb.DeleteUserRequest{Username: "test_list_user1"})
	identityClient.DeleteUser(ctx, &identitypb.DeleteUserRequest{Username: "test_list_user2"})

	// movieClient.DeleteMovie(ctx, &moviepb.DeleteMovieRequest{Id: })

}
