package interceptors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const errMsgFake = "fake error"

func TestUnaryRateLimiter_RateLimitPass(t *testing.T) {

	mockPassLimiter := NewRateLimiter()

	// Testing for valid interval duration
	interceptor := mockPassLimiter.UnaryRateLimiter(mockPassLimiter)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New(errMsgFake)
	}
	info := &grpc.UnaryServerInfo{
		FullMethod: "FakeService",
	}
	resp, err := interceptor(nil, nil, info, handler)
	assert.Nil(t, resp)
	assert.EqualError(t, err, errMsgFake)

	// For new duration with exceeded requests
	mockPassLimiter.endTime = time.Now().Add(-1 * time.Second)
	mockPassLimiter.requests = QueriesPerInterval
	interceptor = mockPassLimiter.UnaryRateLimiter(mockPassLimiter)
	resp, err = interceptor(nil, nil, info, handler)
	assert.Nil(t, resp)
	assert.EqualError(t, err, errMsgFake)

}

func TestUnaryRateLimiter_RateLimitFail(t *testing.T) {

	var mockFailLimiter *queryLimiter

	// Without Rate Limiter
	interceptor := mockFailLimiter.UnaryRateLimiter(mockFailLimiter)
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New(errMsgFake)
	}
	info := &grpc.UnaryServerInfo{
		FullMethod: "FakeService",
	}
	resp, err := interceptor(nil, nil, info, handler)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "rpc error: code = ResourceExhausted desc = FakeService is rejected by the API. Please retry after a while.")

	// With Queries exceeded Rate Limiter
	mockFailLimiter = NewRateLimiter()
	mockFailLimiter.requests = QueriesPerInterval

	interceptor = mockFailLimiter.UnaryRateLimiter(mockFailLimiter)
	handler = func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New(errMsgFake)
	}
	info = &grpc.UnaryServerInfo{
		FullMethod: "FakeService",
	}
	resp, err = interceptor(nil, nil, info, handler)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "rpc error: code = ResourceExhausted desc = FakeService is rejected by the API. Please retry after a while.")
}

func TestStreamRateLimiter_RateLimitPass(t *testing.T) {

	mockPassLimiter := NewRateLimiter()

	interceptor := mockPassLimiter.StreamRateLimiter(mockPassLimiter)
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return errors.New(errMsgFake)
	}
	info := &grpc.StreamServerInfo{
		FullMethod: "FakeService",
	}
	err := interceptor(nil, nil, info, handler)
	assert.EqualError(t, err, errMsgFake)
}

func TestStreamRateLimiter_RateLimitFail(t *testing.T) {
	mockFailLimiter := NewRateLimiter()
	mockFailLimiter.requests = QueriesPerInterval

	interceptor := mockFailLimiter.StreamRateLimiter(mockFailLimiter)
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return errors.New(errMsgFake)
	}
	info := &grpc.StreamServerInfo{
		FullMethod: "FakeService",
	}
	err := interceptor(nil, nil, info, handler)
	assert.EqualError(t, err, "rpc error: code = ResourceExhausted desc = FakeService is rejected by the API. Please retry after a while.")
}
