package services

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/lib/configuration"
	"github.com/AkashGit21/ms-project/lib/persistence/dblayer"
	"github.com/golang/protobuf/ptypes/empty"
)

func getMovieServer() *movieServer {
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	authSrv := NewAuthServer(NewIdentityServer(dbhandler))
	return NewMovieServer(authSrv)
}

func TestCreateMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

	tests := []TestCase{
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
			actual, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: tcase.args.(*moviepb.Movie)})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if len(actual.GetId()) != len(tcase.expected.(string)) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", len(tcase.expected.(string)), len(actual.GetId()))
			}
		})
	}
}

func TestGetMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

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

	tests := []TestCase{
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
			actual, err := TestMovieSrv.GetMovie(context.Background(), &moviepb.GetMovieRequest{Id: tcase.args.(string)})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual, tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}
		})
	}
}

func TestListMovies(t *testing.T) {
	TestMovieSrv = getMovieServer()

	mvObj1 := &moviepb.Movie{
		Name:    "test_list_movie1",
		Summary: "test_list_movies_summary1",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	_, err := TestMovieSrv.CreateMovie(
		context.Background(),
		&moviepb.CreateMovieRequest{
			Movie: mvObj1,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	mvObj2 := &moviepb.Movie{
		Name:    "test_list_movie2",
		Summary: "test_list_movies_summary2",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	_, err = TestMovieSrv.CreateMovie(
		context.Background(),
		&moviepb.CreateMovieRequest{
			Movie: mvObj2,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	tests := []TestCase{
		{
			name:        "found_movie1",
			args:        int32(1),
			expected:    []*moviepb.Movie{mvObj1},
			expectedErr: "",
		},
		{
			name:        "found_movie2",
			args:        int32(1),
			expected:    []*moviepb.Movie{mvObj2},
			expectedErr: "",
		},
	}

	pageToken := ""
	var req *moviepb.ListMoviesRequest
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			if pageToken == "" {
				req = &moviepb.ListMoviesRequest{PageSize: tcase.args.(int32)}
			} else {
				req = &moviepb.ListMoviesRequest{PageSize: tcase.args.(int32), PageToken: pageToken}
			}
			actual, err := TestMovieSrv.ListMovies(context.Background(), req)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) &&
				!reflect.DeepEqual(actual.GetMovies(), tcase.expected.([]*moviepb.Movie)) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)
			}

			pageToken = actual.GetNextPageToken()
		})
	}
}

func TestUpdateMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

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

	tests := []TestCase{
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
			actual, err := TestMovieSrv.UpdateMovie(context.Background(), tcase.args.(*moviepb.UpdateMovieRequest))

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) && !reflect.DeepEqual(actual.GetId(), tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual.GetId())
			}
		})
	}
}

func TestDeleteMovie(t *testing.T) {
	TestMovieSrv = getMovieServer()

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

	tests := []TestCase{
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
					Id: tcase.args.(string),
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

func BenchmarkCreateMovie(b *testing.B) {

	TestMovieSrv = getMovieServer()
	for i := 0; i < b.N; i++ {
		mvObj := &moviepb.Movie{
			Name:    "test_create_movie" + strconv.FormatInt(int64(i), 10),
			Summary: "test_create_movie_summary" + strconv.FormatInt(int64(i), 10),
			Cast:    []string{"test_cast1", "test_cast2"},
		}
		if _, err := TestMovieSrv.CreateMovie(context.Background(),
			&moviepb.CreateMovieRequest{Movie: mvObj}); err != nil {
			b.Errorf("Failed to create pre-requisite object!")
		}
	}
}

func BenchmarkGetMovie(b *testing.B) {

	TestMovieSrv = getMovieServer()
	mvObj := &moviepb.Movie{
		Name:    "test_get_movie",
		Summary: "test_get_movie_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	resp, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		b.Errorf("Failed to create pre-requisite object!")
	}
	movieID := resp.GetId()

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if _, err := TestMovieSrv.GetMovie(context.Background(),
			&moviepb.GetMovieRequest{Id: movieID}); err != nil {
			b.Errorf("Failed to get pre-requisite object!")
		}
	}
}

func BenchmarkDeleteMovie(b *testing.B) {

	TestMovieSrv = getMovieServer()
	for i := 0; i < b.N; i++ {
		mvObj := &moviepb.Movie{
			Name:    "test_delete_movie",
			Summary: "test_delete_movie_summary",
			Cast:    []string{"test_cast1", "test_cast2"},
		}
		resp, err := TestMovieSrv.CreateMovie(context.Background(), &moviepb.CreateMovieRequest{Movie: mvObj})
		if err != nil {
			b.Errorf("Failed to create pre-requisite object!")
		}
		movieID := resp.GetId()

		b.StopTimer()
		b.StartTimer()
		if _, err := TestMovieSrv.DeleteMovie(context.Background(),
			&moviepb.DeleteMovieRequest{Id: movieID}); err != nil {
			b.Errorf("Failed to delete pre-requisite object!")
		}
	}
}
