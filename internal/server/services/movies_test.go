package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/lib/configuration"
	"github.com/AkashGit21/ms-project/lib/persistence/dblayer"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func getMovieServer() *movieServer {
	dbhandler, _ := dblayer.NewPersistenceLayer(configuration.DBTypeDefault, configuration.DBConnectionDefault)

	authSrv := NewAuthServer(NewIdentityServer(dbhandler))
	return NewMovieServer(authSrv)
}

func TestCreateMovie(t *testing.T) {

	// Tests to be checked
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

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(),
		grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	movieClient := moviepb.NewMovieServiceClient(conn)

	// Start checking tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := movieClient.CreateMovie(context.Background(),
				&moviepb.CreateMovieRequest{Movie: tcase.args.(*moviepb.Movie)})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if len(actual.GetId()) != len(tcase.expected.(string)) {
				t.Errorf("\n\texpected: %v \n\tactual: %v",
					len(tcase.expected.(string)), len(actual.GetId()))
			}
		})
	}
}

func TestGetMovie(t *testing.T) {

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(),
		grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	movieClient := moviepb.NewMovieServiceClient(conn)

	// Object to be created for tests
	mvObj := &moviepb.Movie{
		Name:    "test_get_movie",
		Summary: "test_get_movie_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	resp, err := movieClient.CreateMovie(context.Background(),
		&moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	movieID := resp.GetId()

	// Tests to be checked
	tests := []TestCase{
		{
			name:        "not_exists",
			args:        "9e6f9248-e147-4cbe-9c4f-e3d06c79e361",
			expected:    nil,
			expectedErr: "rpc error: code = NotFound desc = Movie Record with ID:9e6f9248-e147-4cbe-9c4f-e3d06c79e361 does not exist!",
		},
		{
			name:        "found_movie",
			args:        movieID,
			expected:    mvObj,
			expectedErr: "",
		},
	}

	// Start checking tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := movieClient.GetMovie(context.Background(),
				&moviepb.GetMovieRequest{Id: tcase.args.(string)})

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}

			if (tcase.expected != nil && (actual == nil)) || (tcase.expected == nil && actual != nil) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual)

			} else if tcase.expected != nil && actual != nil {

				gotMovie := tcase.expected.(*moviepb.Movie)
				if !reflect.DeepEqual(actual.Name, gotMovie.Name) {
					t.Errorf("\n\texpected: %v \n\tactual: %v", gotMovie.Name, actual.Name)
				}
				if !reflect.DeepEqual(actual.Summary, gotMovie.Summary) {
					t.Errorf("\n\texpected: %v \n\tactual: %v", gotMovie.Summary, actual.Summary)
				}
				if !reflect.DeepEqual(actual.Cast, gotMovie.Cast) {
					t.Errorf("\n\texpected: %v \n\tactual: %v", gotMovie.Cast, actual.Cast)
				}
			}
		})
	}
}

func TestListMovies(t *testing.T) {

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(),
		grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	movieClient := moviepb.NewMovieServiceClient(conn)

	// Objects to be created for testing ListMovies
	mvObj1 := &moviepb.Movie{
		Name:    "test_list_movie1",
		Summary: "test_list_movies_summary1",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	_, err = movieClient.CreateMovie(
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
	_, err = movieClient.CreateMovie(
		context.Background(),
		&moviepb.CreateMovieRequest{
			Movie: mvObj2,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}

	// Tests to be checked
	tests := []TestCase{
		{
			name:        "found_movie1",
			args:        int32(2),
			expected:    2,
			expectedErr: "",
		},
		{
			name:        "found_movie2",
			args:        int32(1),
			expected:    1,
			expectedErr: "",
		},
	}

	var req *moviepb.ListMoviesRequest
	// Start checking tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {

			req = &moviepb.ListMoviesRequest{PageSize: tcase.args.(int32)}
			actual, err := TestMovieSrv.ListMovies(context.Background(), req)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) &&
				!reflect.DeepEqual(
					len(actual.GetMovies()), tcase.expected.(int)) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected.(int), len(actual.GetMovies()))
			}

		})
	}
}

func TestUpdateMovie(t *testing.T) {

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(),
		grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	movieClient := moviepb.NewMovieServiceClient(conn)

	// Object to be created for testing UpdateMovie
	mvObj := &moviepb.Movie{
		Name:    "test_update_movie",
		Summary: "test_update_movie_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	resp, err := movieClient.CreateMovie(context.Background(),
		&moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	movieID := resp.GetId()

	// Tests to be checked
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

	// Start checking tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := movieClient.UpdateMovie(context.Background(),
				tcase.args.(*moviepb.UpdateMovieRequest))

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}
			if (tcase.expected != nil || (actual != nil)) &&
				!reflect.DeepEqual(actual.GetId(), tcase.expected) {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expected, actual.GetId())
			}
		})
	}
}

func TestDeleteMovie(t *testing.T) {

	// Mock server using Client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(),
		grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	movieClient := moviepb.NewMovieServiceClient(conn)

	// Object to be created for testing DeleteMovie
	mvObj := &moviepb.Movie{
		Name:    "test_delete_movie",
		Summary: "test_delete_movie_summary",
		Cast:    []string{"test_cast1", "test_cast2"},
	}
	resp, err := movieClient.CreateMovie(context.Background(),
		&moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		t.Fatalf("Failed to create pre-requisite object!")
	}
	movieID := resp.GetId()

	// Tests to be checked
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

	// Start checking tests
	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			actual, err := movieClient.DeleteMovie(
				context.Background(),
				&moviepb.DeleteMovieRequest{
					Id: tcase.args.(string),
				},
			)

			if (err == nil || (err.Error() != tcase.expectedErr)) && tcase.expectedErr != "" {
				t.Errorf("\n\texpected: %v \n\tactual: %v", tcase.expectedErr, err)
			}

			fmt.Println("Actual:", actual, " Type:", reflect.TypeOf(actual))
			fmt.Println("Expected:", tcase.expected, " Type:", reflect.TypeOf(tcase.expected))

			if (tcase.expected != nil || (actual != nil)) &&
				!reflect.DeepEqual(actual, tcase.expected) {
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
	resp, err := TestMovieSrv.CreateMovie(context.Background(),
		&moviepb.CreateMovieRequest{Movie: mvObj})
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
		resp, err := TestMovieSrv.CreateMovie(context.Background(),
			&moviepb.CreateMovieRequest{Movie: mvObj})
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
