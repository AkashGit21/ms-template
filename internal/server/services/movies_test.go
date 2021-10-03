package services

import (
	"context"
	"reflect"
	"testing"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/golang/protobuf/ptypes/empty"
)

func getMovieServer() *movieServer {
	authSrv := NewAuthServer(NewIdentityServer())
	return NewMovieServer(authSrv)
}

func TestCreateMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

	type testCase struct {
		name        string
		args        *moviepb.Movie
		expected    string
		expectedErr string
	}

	tests := []testCase{
		{
			name: "invalid_ID_error",
			args: &moviepb.Movie{
				Id:      "test_id",
				Name:    "test_name",
				Summary: "test_summary",
				Cast:    []string{"test_cast1"},
				Tags:    []moviepb.Tag{moviepb.Tag(0)},
				Writers: []string{"test_writer1", "test_writer2"},
			},
			expected:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = Input is not valid! ID is auto-generated...",
		},
		{
			name: "invalid_data_error",
			args: &moviepb.Movie{
				Name:    " ",
				Summary: "test_summary",
				Cast:    []string{"test_cast1"},
				Tags:    []moviepb.Tag{moviepb.Tag(0)},
				Writers: []string{"test_writer1", "test_writer2"},
			},
			expected:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = Input is not valid! The name should be between 1 and 120 characters.",
		},
		{
			name: "created_movie",
			args: &moviepb.Movie{
				Name:    "test_create_movie",
				Summary: "test_summary",
				Cast:    []string{"test_cast1"},
				Tags:    []moviepb.Tag{moviepb.Tag(0)},
				Writers: []string{"test_writer1", "test_writer2"},
			},
			expected:    "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: tcase.args})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if len(actual.GetId()) != len(tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", len(tcase.expected), len(actual.GetId()))
			}
		})
	}
}

func TestGetMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

	type testCase struct {
		name        string
		args        string
		expected    *moviepb.Movie
		expectedErr string
	}

	mvObj := &moviepb.Movie{
		Name:    "test_get_movie",
		Summary: "test_get_movie_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	resp, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	movieID := resp.GetId()

	tests := []testCase{
		{
			name:        "not_exists",
			args:        "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
			expected:    nil,
			expectedErr: "rpc error: code = NotFound desc = Record with ID:9e6f9248-e147-4cbe-9c4f-e3d06c79e361 does not exist!",
		},
		{
			name:        "found_movie",
			args:        movieID,
			expected:    mvObj,
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestMovieSrv.GetMovie(context.Background(), &moviepb.GetMovieRequest{Id: tcase.args})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

func TestListMovies(t *testing.T) {
	TestMovieSrv = getMovieServer()

	type testCase struct {
		name        string
		args        string
		expected    *moviepb.ListMoviesResponse
		expectedErr string
	}

	mvObj := &moviepb.Movie{
		Name:    "test_list_movies",
		Summary: "test_list_movies_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	_, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	// movieID := resp.GetId()

	tests := []testCase{
		{
			name:        "found_movies",
			args:        "",
			expected:    &moviepb.ListMoviesResponse{Movies: []*moviepb.Movie{mvObj}},
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestMovieSrv.ListMovies(context.Background(), &moviepb.ListMoviesRequest{})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

func TestUpdateMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

	type testCase struct {
		name        string
		args        *moviepb.UpdateMovieRequest
		expected    string
		expectedErr string
	}

	mvObj := &moviepb.Movie{
		Name:    "test_update_movie",
		Summary: "test_update_movie_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	resp, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	movieID := resp.GetId()

	tests := []testCase{
		{
			name: "not_exists",
			args: &moviepb.UpdateMovieRequest{
				Id:    "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
				Movie: &moviepb.Movie{}},
			expected:    "",
			expectedErr: "rpc error: code = NotFound desc = Movie Record with ID:9e6f9248-e147-4cbe-9c4f-e3d06c79e361 does not exist!",
		},
		{
			name: "bad_id_update",
			args: &moviepb.UpdateMovieRequest{
				Id: movieID,
				Movie: &moviepb.Movie{
					Id:      "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
					Name:    "test_update_movie",
					Summary: "test_update_movie_summary",
					Cast:    []string{"test_cast1", "test_cast2"},
				},
			},
			expected:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = Cannot update the ID of object!",
		},
		{
			name: "bad_input",
			args: &moviepb.UpdateMovieRequest{
				Id: movieID,
				Movie: &moviepb.Movie{
					Name:    "test_update_movie",
					Summary: "summary",
					Cast:    []string{"test_cast1", "test_cast2"},
				},
			},
			expected:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = Input is not valid! The summary should be between 8 and 1200 characters.",
		},
		{
			name: "updated_movie",
			args: &moviepb.UpdateMovieRequest{
				Id: movieID,
				Movie: &moviepb.Movie{
					Name:    "test_update_movie",
					Summary: "test_update_movie_summary_updated",
					Cast:    []string{"test_cast1", "test_cast2"},
				},
			},
			expected:    movieID,
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestMovieSrv.UpdateMovie(context.Background(), tcase.args)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if !reflect.DeepEqual(actual.GetId(), tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

func TestDeleteMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

	type testCase struct {
		name        string
		args        string
		expected    *empty.Empty
		expectedErr string
	}

	mvObj := &moviepb.Movie{
		Name:    "test_delete_movie",
		Summary: "test_delete_movie_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	resp, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	movieID := resp.GetId()

	tests := []testCase{
		{
			name:        "not_exists",
			args:        "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
			expected:    nil,
			expectedErr: "rpc error: code = NotFound desc = Movie Record with ID:9e6f9248-e147-4cbe-9c4f-e3d06c79e361 does not exist!",
		},
		{
			name:        "removed_movie",
			args:        movieID,
			expected:    &empty.Empty{},
			expectedErr: "",
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := TestMovieSrv.DeleteMovie(
				context.Background(),
				&moviepb.DeleteMovieRequest{
					Id: tcase.args,
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
