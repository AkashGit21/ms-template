package server

import (
	"context"
	"io"
	"testing"

	testingpb "github.com/AkashGit21/ms-project/internal/grpc/testing"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	// DefaultResponseValue is the default value used for Response.
	DefaultResponseValue = "default_response_value"
	// ListResponseCount is the expected number of responses to PingList
	ListResponseCount = 100
)

type TestPingService struct {
	T *testing.T

	testingpb.UnimplementedTestServiceServer
}

func NewTestPingService() *TestPingService {
	return &TestPingService{}
}

func (ts *TestPingService) PingEmpty(_ context.Context, _ *empty.Empty) (*testingpb.PingResponse, error) {

	return &testingpb.PingResponse{Value: DefaultResponseValue, Counter: 42}, nil
}

func (ts *TestPingService) Ping(_ context.Context, p *testingpb.PingRequest) (*testingpb.PingResponse, error) {

	return &testingpb.PingResponse{Value: p.Value, Counter: 42}, nil
}

func (ts *TestPingService) PingError(_ context.Context, p *testingpb.PingRequest) (*empty.Empty, error) {
	code := codes.Code(p.ErrorCodeReturned)
	return &emptypb.Empty{}, status.Errorf(code, "Userspace error.")
}

func (ts *TestPingService) PingList(p *testingpb.PingRequest, stream testingpb.TestService_PingListServer) error {
	if p.ErrorCodeReturned != 0 {
		return status.Errorf(codes.Code(p.ErrorCodeReturned), "khatra")
	}
	// Send user trailers and headers.
	for i := 0; i < ListResponseCount; i++ {
		stream.Send(&testingpb.PingResponse{Value: p.Value, Counter: int32(i)})
	}
	return nil
}

func (ts *TestPingService) PingStream(stream testingpb.TestService_PingStreamServer) error {
	count := 0
	for true {
		ping, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		stream.Send(&testingpb.PingResponse{Value: ping.Value, Counter: int32(count)})
		count += 1
	}
	return nil
}
