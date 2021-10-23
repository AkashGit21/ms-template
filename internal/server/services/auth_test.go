package services

import (
	"context"
	"log"
	"testing"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"google.golang.org/grpc"
)

func TestLogin(t *testing.T) {

	// Mock Client for testing
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	identityClient := identitypb.NewIdentityServiceClient(conn)
	authClient := authpb.NewAuthServiceClient(conn)

	userObj := &identitypb.User{
		Username:  "test_login_username",
		Email:     "test_login_email@domain.com",
		Password:  "test_login_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	_, err = identityClient.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: userObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	tests := []TestCase{
		{
			name: "not_found",
			args: &authpb.LoginRequest{
				Username: "test_login_user",
				Password: "test_login_pwd",
			},
			expected:    5,
			expectedErr: "rpc error: code = InvalidArgument desc = incorrect username/password!",
		},
		{
			name: "bad_combo",
			args: &authpb.LoginRequest{
				Username: "test_login_username",
				Password: "test_login_password",
			},
			expected:    5,
			expectedErr: "rpc error: code = InvalidArgument desc = incorrect username/password!",
		},
		{
			name: "logged_in",
			args: &authpb.LoginRequest{
				Username: "test_login_username",
				Password: "test_login_pwd",
			},
			expected:    192,
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := authClient.Login(context.Background(), tcase.args.(*authpb.LoginRequest))

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && (len(actual.String()) != tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, len(actual.String()))
			}
		})
	}
}
