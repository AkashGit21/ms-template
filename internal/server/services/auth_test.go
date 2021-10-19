package services

import (
	"context"
	"testing"
	"time"

	authpb "github.com/AkashGit21/ms-project/internal/grpc/auth"
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/lib/configuration"
	"github.com/AkashGit21/ms-project/lib/persistence/dblayer"
)

func TestLogin(t *testing.T) {

	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	TestIdentitySrv = NewIdentityServer(dbhandler)
	TestAuthSrv = NewAuthServer(TestIdentitySrv)
	TestAuthSrv.JWT = server.NewJWTManager(SecretKey, 2*time.Minute)

	userObj := &identitypb.User{
		Username:  "test_login_username",
		Email:     "test_login_email@domain.com",
		Password:  "test_login_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	_, err := TestIdentitySrv.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: userObj})
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
			expectedErr: "rpc error: code = NotFound desc = incorrect username/password!",
		},
		{
			name: "bad_combo",
			args: &authpb.LoginRequest{
				Username: "test_login_username",
				Password: "test_login_password",
			},
			expected:    5,
			expectedErr: "rpc error: code = NotFound desc = incorrect username/password!",
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
			actual, err := TestAuthSrv.Login(context.Background(), tcase.args.(*authpb.LoginRequest))

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && (len(actual.String()) != tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, len(actual.String()))
			}
		})
	}
}
