package interceptors

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/AkashGit21/ms-project/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	jwtManager      *server.JWTManager
	accessibleRoles map[string][]string
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(jwtManager *server.JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("--> auth interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> auth interceptor: ", info.FullMethod)

		err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {

	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		return fmt.Errorf("unknown service!")
	}
	for _, role := range accessibleRoles {
		if strings.EqualFold("GUEST", role) {
			// Guest is allowed here
			return nil
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided!")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.GetUserFromToken(accessToken)
	if err != nil {
		log.Println("error: ", err.Error())
		return status.Errorf(codes.Unauthenticated, "bad access token!")
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "not allowed to access this feature!")
}
