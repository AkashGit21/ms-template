package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// Global Refresh duration (time period) for particular interval
	RefreshDuration = 1 * time.Minute
	// Max. number of Queries allowed during the Interval
	QueriesPerInterval = 20
)

// When the DidLimitExceed() return true, the request will be rejected.
// Since the User has passes his/her expected limit.
// Otherwise, the request will pass.
type Limiter interface {
	DidLimitExceed() bool
}

type queryLimiter struct {
	endTime         time.Time
	refreshInterval time.Duration
	requests        int
	queriesAllowed  int
}

func NewRateLimiter() *queryLimiter {
	return &queryLimiter{
		endTime:         time.Now().Add(RefreshDuration),
		refreshInterval: RefreshDuration,
		requests:        0,
		queriesAllowed:  QueriesPerInterval,
	}
}

func (lim *queryLimiter) DidLimitExceed() bool {
	now := time.Now()

	if lim != nil {
		// Check for Limiter endTime is still correct or not
		if !lim.endTime.Before(now) {
			// Check if number of requests are greater than expected
			if lim.requests >= lim.queriesAllowed {
				return true
			}
			lim.requests++
		} else {
			// Refresh the limiter in case new interval starts
			refreshLimiter(lim)
		}

		return false
	} else {
		log.Println("No rate limiter found!")
		return true
	}
}

func refreshLimiter(l *queryLimiter) {
	l.endTime = time.Now().Add(l.refreshInterval)
	l.requests = 1
}

// UnaryRateLimiter returns a new unary server interceptors that manages Rate-limiting of requests.
func (l *queryLimiter) UnaryRateLimiter(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		log.Println("--> rate_limit interceptor: ", info.FullMethod)
		if limiter.DidLimitExceed() {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by the API. Please retry after a while.", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that manages Rate-limiting of requests.
func (lim *queryLimiter) StreamRateLimiter(limiter Limiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		log.Println("--> rate_limit interceptor: ", info.FullMethod)
		if limiter.DidLimitExceed() {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by the API. Please retry after a while.", info.FullMethod)
		}
		return handler(srv, stream)
	}
}
