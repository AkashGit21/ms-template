package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/AkashGit21/ms-project/internal/server/interceptors"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func TestCreateUser(t *testing.T) {

	// Tests to be checked
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
		// {
		// 	name: "bad_email",
		// 	args: &identitypb.User{
		// 		Username:  "test_username2",
		// 		Email:     "test_email@domain.com",
		// 		Password:  "test_pwd",
		// 		Role:      identitypb.Role_NORMAL,
		// 		FirstName: "test_first",
		// 	},
		// 	expected:    "",
		// 	expectedErr: "rpc error: code = AlreadyExists desc = A user with email `test_email@domain.com` already exists.",
		// },
	}

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := identitypb.NewIdentityServiceClient(conn)

	// Start checking tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			actual, err := client.CreateUser(ctx, &identitypb.CreateUserRequest{User: tcase.args.(*identitypb.User)})
			// actual, err := TestIdentitySrv.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: tcase.args.(*identitypb.User)})

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

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := identitypb.NewIdentityServiceClient(conn)

	userObj := &identitypb.User{
		Username:  "test_get_username",
		Email:     "test_email@domain.com",
		Password:  "test_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	resp, err := client.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: userObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	uname := resp.GetUsername()

	// Tests to be checked
	tests := []TestCase{
		{
			name:        "bad_username",
			args:        "test_get_user",
			expected:    nil,
			expectedErr: "rpc error: code = NotFound desc = A user with username `test_get_user` not found!",
		},
		{
			name:        "found_user",
			args:        uname,
			expected:    userObj,
			expectedErr: "",
		},
	}

	// To allow the deletion, assigned ADMIN role temporarily
	interceptors.CURRENT_ROLE = "ADMIN"

	// Start running tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := client.GetUser(context.Background(), &identitypb.GetUserRequest{Username: tcase.args.(string)})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}

			if (tcase.expected != nil && (actual == nil)) || (tcase.expected == nil && actual != nil) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)

			} else if tcase.expected != nil && actual != nil {

				gotUser := tcase.expected.(*identitypb.User)
				if !reflect.DeepEqual(actual.Username, gotUser.Username) {
					t.Errorf("\n\texpected: %v \n\tactual: %v", gotUser.Username, actual.Username)
				}
				if !reflect.DeepEqual(actual.Username, gotUser.Username) {
					t.Errorf("\n\texpected: %v \n\tactual: %v", gotUser.Email, actual.Email)
				}
				if !reflect.DeepEqual(actual.Username, gotUser.Username) {
					t.Errorf("\n\texpected: %v \n\tactual: %v", gotUser.Role, actual.Role)
				}
				if !reflect.DeepEqual(actual.Username, gotUser.Username) {
					t.Errorf("\n\texpected: %v \n\tactual: %v", gotUser.FirstName, actual.FirstName)
				}
			}
		})
	}
}

func TestListUsers(t *testing.T) {

	// Mock server using client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := identitypb.NewIdentityServiceClient(conn)

	// Creating users to be viewed
	userObj1 := &identitypb.User{
		Username:  "test_list_user1",
		Email:     "test_list_email1@domain.in",
		Password:  "test_list_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	_, err = client.CreateUser(
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
	_, err = client.CreateUser(
		context.Background(),
		&identitypb.CreateUserRequest{
			User: userObj2,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	// Tests to be checked
	tests := []TestCase{
		{
			name:        "list_users_size1",
			args:        int32(2),
			expected:    2,
			expectedErr: "",
		},
		{
			name:        "list_users_size2",
			args:        int32(1),
			expected:    1,
			expectedErr: "",
		},
	}

	// To allow the deletion, assigned ADMIN role temporarily
	interceptors.CURRENT_ROLE = "ADMIN"

	// Start running tests
	var req *identitypb.ListUsersRequest
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			req = &identitypb.ListUsersRequest{PageSize: tcase.args.(int32)}
			actual, err := client.ListUsers(context.Background(), req)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) &&
				!reflect.DeepEqual(
					len(actual.GetUsers()), tcase.expected.(int),
				) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected.(int), len(actual.GetUsers()))
			}

		})
	}
}

// TODO: Tests for UpdateUser is not written
// func TestUpdateUser(t *testing.T) {
// 	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

// 	TestIdentitySrv = NewIdentityServer(dbhandler)

// 	userObj := &identitypb.User{
// 		Username:  "test_list_user1",
// 		Email:     "test_list_email1@domain.in",
// 		Password:  "test_list_pwd",
// 		Role:      identitypb.Role_NORMAL,
// 		FirstName: "test_first",
// 	}
// 	_, err := TestIdentitySrv.CreateUser(
// 		context.Background(),
// 		&identitypb.CreateUserRequest{
// 			User: userObj,
// 		},
// 	)
// 	if err != nil {
// 		t.Fatalf("Failed to create pre-requisite object!")
// 	}

// 	tests := []TestCase{
// 		{
// 			name:        "demo",
// 			args:        &identitypb.UpdateUserRequest{},
// 			expected:    &identitypb.User{},
// 			expectedErr: "",
// 		},
// 	}

// 	for _, tcase := range tests {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			actual, err := TestIdentitySrv.UpdateUser(
// 				context.Background(),
// 				tcase.args.(*identitypb.UpdateUserRequest))

// 			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
// 				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
// 			}
// 			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual, tcase.expected) {
// 				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
// 			}
// 		})
// 	}
// }

func TestDeleteUser(t *testing.T) {

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := identitypb.NewIdentityServiceClient(conn)

	// User ccreation for deletion test
	userObj := &identitypb.User{
		Username:  "test_delete_username",
		Email:     "test_delete_email@domain.in",
		Password:  "test_delete_pwd",
		Role:      identitypb.Role_NORMAL,
		FirstName: "test_first",
	}
	resp, err := client.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: userObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	uname := resp.GetUsername()

	// Tests to be checked
	tests := []TestCase{
		{
			name:        "not_exists",
			args:        "test_delete_user",
			expected:    nil,
			expectedErr: "rpc error: code = NotFound desc = A user with username `test_delete_user` does not exist!",
		},
		{
			name:        "removed_user",
			args:        uname,
			expected:    &empty.Empty{},
			expectedErr: "",
		},
	}

	// To allow the deletion, assigned ADMIN role temporarily
	interceptors.CURRENT_ROLE = "ADMIN"

	// start running the tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := client.DeleteUser(
				context.Background(),
				&identitypb.DeleteUserRequest{
					Username: tcase.args.(string),
				},
			)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}

			fmt.Println("Actual:", actual, " Type:", reflect.TypeOf(actual))
			fmt.Println("Expected:", tcase.expected, " Type:", reflect.TypeOf(tcase.expected))

			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual, tcase.expected) {
				// fmt.Println("Equal: ", reflect.DeepEqual(tcase.expected, actual))
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}
