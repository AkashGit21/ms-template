package interceptors

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	testingpb "github.com/AkashGit21/ms-project/internal/grpc/testing"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	secretKey       = "secret"
	mockAccessRoles = map[string][]string{
		"/testing.TestService/PingEmpty":  {"ADMIN", "SUBSCRIBED"},
		"/testing.TestService/PingList":   {"ADMIN", "SUBSCRIBED"},
		"/testing.TestService/PingStream": {"ADMIN", "SUBSCRIBED"},
	}

	authToken   = "some_bad_token"
	headerParam = "access_token"

	pingReq = &testingpb.PingRequest{Value: "random ping", SleepTimeMs: 9876, ErrorCodeReturned: 16}
)

type TestInterceptorSuite struct {
	TestService testingpb.TestServiceServer
	T           *testing.T

	ServerOpts     []grpc.ServerOption
	ClientOpts     []grpc.DialOption
	serverAddr     string
	ServerListener net.Listener
	Server         *grpc.Server
	Client         testingpb.TestServiceClient

	restartServerWithDelayedStart chan time.Duration
	serverRunning                 chan bool
}

func (ms *TestInterceptorSuite) SetupSuite() {
	ms.restartServerWithDelayedStart = make(chan time.Duration)
	ms.serverRunning = make(chan bool)

	go func() {
		for {
			var err error
			ms.ServerListener, err = net.Listen("tcp", ms.serverAddr)
			if err != nil {
				ms.T.Fatalf("unable to generate test certificate/key: " + err.Error())
			}

			ms.serverAddr = ms.ServerListener.Addr().String()
			require.NoError(ms.T, err, "must be able to allocate a port for serverListener")

			ms.Server = grpc.NewServer(ms.ServerOpts...)

			// Create a service of the instantiator hasn't provided one.
			if ms.TestService == nil {
				ms.TestService = &server.TestPingService{T: ms.T}
			}
			testingpb.RegisterTestServiceServer(ms.Server, ms.TestService)

			go func() {
				ms.Server.Serve(ms.ServerListener)
			}()
			if ms.Client == nil {
				ms.Client = ms.NewClient(ms.ClientOpts...)
			}

			ms.serverRunning <- true

			d := <-ms.restartServerWithDelayedStart
			ms.Server.Stop()
			time.Sleep(d)
		}
	}()

	select {
	case <-ms.serverRunning:
	case <-time.After(1 * time.Second):
		ms.T.Fatal("server failed to start before deadline")
	}
}

func (ms *TestInterceptorSuite) NewClient(dialOpts ...grpc.DialOption) testingpb.TestServiceClient {
	newDialOpts := append(dialOpts, grpc.WithBlock())
	newDialOpts = append(newDialOpts, grpc.WithInsecure())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	clientConn, err := grpc.DialContext(ctx, ms.serverAddr, newDialOpts...)
	require.NoError(ms.T, err, "must not error on client Dial")
	return testingpb.NewTestServiceClient(clientConn)
}

func (ms *TestInterceptorSuite) pingReqWithAuth(choice int, accessLevel int) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to gRPC Request.
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", addAuthToken(accessLevel))

	switch choice {
	case 1:
		return ms.Client.PingEmpty(ctx, &emptypb.Empty{})
	case 2:
		return ms.Client.PingList(ctx, pingReq)
	default:
		return nil, fmt.Errorf("check your choice!")
	}
}

func addAuthToken(accessLevel int) string {
	var token string

	// If accessLevel is allowed, generate token else pass default invalid token
	if accessLevel < 0 && accessLevel < 4 {
		token = authToken
	} else {
		mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
		token, _ = mockJWTmngr.GenerateToken(&identitypb.User{
			Username: "usrname1",
			Role:     identitypb.Role(accessLevel),
		})
	}

	return token
}

func TestUnary_BadService(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthPassInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(mockAuthPassInterceptor.Unary()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8089",
	}
	testSuite.SetupSuite()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := testSuite.Client.Ping(ctx, pingReq)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = unknown service!")
}

func TestUnary_NoAuth(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthFailInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(mockAuthFailInterceptor.Unary()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8081",
	}
	testSuite.SetupSuite()

	resp, err := testSuite.Client.PingEmpty(context.Background(), &empty.Empty{})
	assert.Nil(t, resp)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = authorization token is not provided!")
}

func TestUnary_BadAuth(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthFailInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(mockAuthFailInterceptor.Unary()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8082",
	}
	testSuite.SetupSuite()

	resp, err := testSuite.pingReqWithAuth(1, -1)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = bad access token!")
}

func TestUnary_BadAuthPermissionDenied(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthFailInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(mockAuthFailInterceptor.Unary()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8083",
	}
	testSuite.SetupSuite()

	resp, err := testSuite.pingReqWithAuth(1, 1)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = not allowed to perform this operation!")
}

func TestUnary_AuthPasses(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthPassInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(mockAuthPassInterceptor.Unary()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8084",
	}
	testSuite.SetupSuite()

	resp, err := testSuite.pingReqWithAuth(1, 2)
	assert.Equal(t, "default_response_value", resp.(*testingpb.PingResponse).GetValue())
	assert.NoError(t, err)
}

func TestStream_NoAuth(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthFailInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(mockAuthFailInterceptor.Stream()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8085",
	}
	testSuite.SetupSuite()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := testSuite.Client.PingList(ctx, pingReq)
	_, err = stream.Recv()
	assert.Error(t, err, "there must be an error")
	assert.Equal(t, codes.Unauthenticated, status.Code(err), "must error with unauthenticated")
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = authorization token is not provided!")
}

// TODO: TestStream_BadAuth not working as expected
func TestStream_BadAuth(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthFailInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(mockAuthFailInterceptor.Stream()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8086",
	}
	testSuite.SetupSuite()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to gRPC Request.
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authToken)

	stream, err := testSuite.Client.PingList(ctx, pingReq)
	_, err = stream.Recv()
	assert.Error(t, err, "there must be an error")
	assert.Equal(t, codes.Unauthenticated, status.Code(err), "must error with unauthenticated")
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = bad access token!")
}

func TestStream_BadAuthPermissionDenied(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthFailInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(mockAuthFailInterceptor.Stream()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8087",
	}
	testSuite.SetupSuite()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to gRPC Request.
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", addAuthToken(1))

	stream, err := testSuite.Client.PingList(ctx, pingReq)
	_, err = stream.Recv()
	assert.Error(t, err, "there must be an error")
	assert.Equal(t, codes.PermissionDenied, status.Code(err), "must error with unauthenticated")
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = not allowed to perform this operation!")
}

func TestStream_AuthPasses(t *testing.T) {

	mockJWTmngr := server.NewJWTManager(secretKey, 2*time.Minute)
	mockAuthFailInterceptor := NewAuthInterceptor(mockJWTmngr, mockAccessRoles)

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(mockAuthFailInterceptor.Stream()),
	}

	testSuite := TestInterceptorSuite{
		T:          t,
		ServerOpts: opts,
		serverAddr: "127.0.0.1:8088",
	}
	testSuite.SetupSuite()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add token to gRPC Request.
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", addAuthToken(2))

	streamReq := &testingpb.PingRequest{
		ErrorCodeReturned: 0,
		Value:             "stream_request",
		SleepTimeMs:       10000,
	}
	stream, err := testSuite.Client.PingList(ctx, streamReq)
	assert.NoError(t, err)
	for ind := 0; ind < 100; ind++ {
		resp, err := stream.Recv()
		assert.Equal(t, ind, int(resp.Counter))
		assert.NoError(t, err)
	}
}
