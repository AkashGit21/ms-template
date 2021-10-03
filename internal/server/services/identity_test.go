package services

import (
	"context"
	"reflect"
	"testing"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/golang/protobuf/ptypes/empty"
)

func TestCreateUser(t *testing.T) {
	TestIdentitySrv = NewIdentityServer()

	type testCase struct {
		name        string
		args        *identitypb.User
		expected    string
		expectedErr string
	}

	tests := []testCase{
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
			actual, err := TestIdentitySrv.CreateUser(context.Background(), &identitypb.CreateUserRequest{User: tcase.args})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if actual.GetUsername() != tcase.expected {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual.GetUsername())
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	TestIdentitySrv = NewIdentityServer()

	type testCase struct {
		name        string
		args        string
		expected    *identitypb.User
		expectedErr string
	}

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

	tests := []testCase{
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
			actual, err := TestIdentitySrv.GetUser(context.Background(), &identitypb.GetUserRequest{Username: tcase.args})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	TestIdentitySrv = NewIdentityServer()

	type testCase struct {
		name        string
		args        string
		expected    *identitypb.ListUsersResponse
		expectedErr string
	}

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
	// uname1 := resp.GetUsername()

	tests := []testCase{
		{
			name: "found_users",
			args: "",
			expected: &identitypb.ListUsersResponse{
				Users: []*identitypb.User{userObj},
			},
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestIdentitySrv.ListUsers(context.Background(), &identitypb.ListUsersRequest{})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

// TODO: Tests for UpdateUser is not written
// func TestUpdateUser(t *testing.T) {
// 	TestIdentitySrv = NewIdentityServer()

// 	type testCase struct {
// 		name        string
// 		args        *moviepb.UpdateMovieRequest
// 		expected    string
// 		expectedErr string
// 	}

// 	mvObj := &moviepb.Movie{
// 		Name:    "test_update_movie",
// 		Summary: "test_update_movie_summary",
// 		Cast:    []string{"test_cast1", "test_cast2"},
// 	}
// 	resp, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: mvObj})
// 	if err != nil {
// 		t.Fatalf("Failed to create pre-requisite object!")
// 	}
// 	movieID := resp.GetId()

// 	tests := []testCase{
// 		{
// 			name: "not_exists",
// 			args: &moviepb.UpdateMovieRequest{
// 				Id:    "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
// 				Movie: &moviepb.Movie{}},
// 			expected:    "",
// 			expectedErr: "rpc error: code = NotFound desc = Movie Record with ID:9e6f9248-e147-4cbe-9c4f-e3d06c79e361 does not exist!",
// 		},
// 		{
// 			name: "bad_id_update",
// 			args: &moviepb.UpdateMovieRequest{
// 				Id: movieID,
// 				Movie: &moviepb.Movie{
// 					Id:      "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
// 					Name:    "test_update_movie",
// 					Summary: "test_update_movie_summary",
// 					Cast:    []string{"test_cast1", "test_cast2"},
// 				},
// 			},
// 			expected:    "",
// 			expectedErr: "rpc error: code = InvalidArgument desc = Cannot update the ID of object!",
// 		},
// 		{
// 			name: "bad_input",
// 			args: &moviepb.UpdateMovieRequest{
// 				Id: movieID,
// 				Movie: &moviepb.Movie{
// 					Name:    "test_update_movie",
// 					Summary: "summary",
// 					Cast:    []string{"test_cast1", "test_cast2"},
// 				},
// 			},
// 			expected:    "",
// 			expectedErr: "rpc error: code = InvalidArgument desc = Input is not valid! The summary should be between 8 and 1200 characters.",
// 		},
// 		{
// 			name: "updated_movie",
// 			args: &moviepb.UpdateMovieRequest{
// 				Id: movieID,
// 				Movie: &moviepb.Movie{
// 					Name:    "test_update_movie",
// 					Summary: "test_update_movie_summary_updated",
// 					Cast:    []string{"test_cast1", "test_cast2"},
// 				},
// 			},
// 			expected:    movieID,
// 			expectedErr: "",
// 		},
// 	}

// 	for _, tcase := range tests {
// 		t.Run(tcase.name, func(t *testing.T) {
// 			actual, err := TestMovieSrv.UpdateMovie(context.Background(), tcase.args)

// 			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
// 				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
// 			}
// 			if !reflect.DeepEqual(actual.GetId(), tcase.expected) {
// 				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
// 			}
// 		})
// 	}
// }

func TestDeleteUser(t *testing.T) {
	TestIdentitySrv = NewIdentityServer()

	type testCase struct {
		name        string
		args        string
		expected    *empty.Empty
		expectedErr string
	}

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

	tests := []testCase{
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
					Username: tcase.args,
				},
			)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}
