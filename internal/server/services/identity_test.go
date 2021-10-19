package services

import (
	"context"
	"reflect"
	"testing"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/AkashGit21/ms-project/lib/configuration"
	"github.com/AkashGit21/ms-project/lib/persistence/dblayer"
	"github.com/golang/protobuf/ptypes/empty"
)

func TestCreateUser(t *testing.T) {
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	TestIdentitySrv = NewIdentityServer(dbhandler)

	tests := []TestCase{
		{
			name: "invalid_username",
			args: &identitypb.User{
				Username:  "",
				Email:     "test_email@domain.com",
				Password:  "test_pwd",
				Role:      identitypb.Role_NORMAL,
				FirstName: "test_first",
			},
			expected:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = The field `username` is required.",
		},
		{
			name: "invalid_email",
			args: &identitypb.User{
				Username:  "test_username",
				Password:  "test_pwd",
				Role:      identitypb.Role_NORMAL,
				FirstName: "test_first",
			},
			expected:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = The field `email` is required.",
		},
		{
			name: "added_user",
			args: &identitypb.User{
				Username:  "test_username",
				Email:     "test_email@domain.com",
				Password:  "test_pwd",
				Role:      identitypb.Role_NORMAL,
				FirstName: "test_first",
			},
			expected:    "test_username",
			expectedErr: "",
		},
		{
			name: "bad_username",
			args: &identitypb.User{
				Username:  "test_username",
				Email:     "test_email2@domain.com",
				Password:  "test_pwd",
				Role:      identitypb.Role_NORMAL,
				FirstName: "test_first",
			},
			expected:    "",
			expectedErr: "rpc error: code = AlreadyExists desc = A user with username `test_username` already exists!",
		},
		{
			name: "bad_email",
			args: &identitypb.User{
				Username:  "test_username2",
				Email:     "test_email@domain.com",
				Password:  "test_pwd",
				Role:      identitypb.Role_NORMAL,
				FirstName: "test_first",
			},
			expected:    "",
			expectedErr: "rpc error: code = AlreadyExists desc = A user with email `test_email@domain.com` already exists.",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestIdentitySrv.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: tcase.args.(*identitypb.User)})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual.GetUsername(), tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual.GetUsername())
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	TestIdentitySrv = NewIdentityServer(dbhandler)

	userObj := &identitypb.User{
		Username:  "test_username",
		Email:     "test_email@domain.com",
		Password:  "test_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	resp, err := TestIdentitySrv.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: userObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	uname := resp.GetUsername()

	tests := []TestCase{
		{
			name:        "bad_username",
			args:        "test_user",
			expected:    nil,
			expectedErr: "rpc error: code = NotFound desc = A user with username `test_user` not found!",
		},
		{
			name:        "found_user",
			args:        uname,
			expected:    userObj,
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestIdentitySrv.GetUser(context.Background(), &identitypb.GetUserRequest{Username: tcase.args.(string)})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	TestIdentitySrv = NewIdentityServer(dbhandler)

	userObj1 := &identitypb.User{
		Username:  "test_list_user1",
		Email:     "test_list_email1@domain.in",
		Password:  "test_list_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	_, err := TestIdentitySrv.CreateUser(
		context.Background(),
		&identitypb.CreateUserRequest{
			User: userObj1,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	userObj2 := &identitypb.User{
		Username:  "test_list_user2",
		Email:     "test_list_email2@domain.in",
		Password:  "test_list_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	_, err = TestIdentitySrv.CreateUser(
		context.Background(),
		&identitypb.CreateUserRequest{
			User: userObj2,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	tests := []TestCase{
		{
			name:        "found_user1",
			args:        int32(1),
			expected:    []*identitypb.User{userObj1},
			expectedErr: "",
		},
		{
			name:        "found_user2",
			args:        int32(1),
			expected:    []*identitypb.User{userObj2},
			expectedErr: "",
		},
	}

	pageToken := ""
	var req *identitypb.ListUsersRequest
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			if pageToken == "" {
				req = &identitypb.ListUsersRequest{PageSize: tcase.args.(int32)}
			} else {
				req = &identitypb.ListUsersRequest{PageSize: tcase.args.(int32), PageToken: pageToken}
			}
			actual, err := TestIdentitySrv.ListUsers(context.Background(), req)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) &&
				!reflect.DeepEqual(
					actual.GetUsers(), tcase.expected.([]*identitypb.User),
				) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}

			pageToken = actual.GetNextPageToken()
		})
	}
}

// TODO: Tests for UpdateUser is not written
func TestUpdateUser(t *testing.T) {
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	TestIdentitySrv = NewIdentityServer(dbhandler)

	userObj := &identitypb.User{
		Username:  "test_list_user1",
		Email:     "test_list_email1@domain.in",
		Password:  "test_list_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	_, err := TestIdentitySrv.CreateUser(
		context.Background(),
		&identitypb.CreateUserRequest{
			User: userObj,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	tests := []TestCase{
		{
			name:        "demo",
			args:        &identitypb.UpdateUserRequest{},
			expected:    &identitypb.User{},
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestIdentitySrv.UpdateUser(
				context.Background(),
				tcase.args.(*identitypb.UpdateUserRequest))

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	TestIdentitySrv = NewIdentityServer(dbhandler)

	userObj := &identitypb.User{
		Username:  "test_delete_user",
		Email:     "test_delete_email@domain.in",
		Password:  "test_delete_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	resp, err := TestIdentitySrv.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: userObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	uname := resp.GetUsername()

	tests := []TestCase{
		{
			name:        "not_exists",
			args:        "test_user",
			expected:    nil,
			expectedErr: "rpc error: code = NotFound desc = A user with username `test_user` does not exist!",
		},
		{
			name:        "removed_user",
			args:        uname,
			expected:    &empty.Empty{},
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestIdentitySrv.DeleteUser(
				context.Background(),
				&identitypb.DeleteUserRequest{
					Username: tcase.args.(string),
				},
			)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}
