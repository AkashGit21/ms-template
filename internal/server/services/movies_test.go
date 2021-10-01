package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/stretchr/testify/assert"
)

func TestMoviesService(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create some basic configuration for testing
	conf, err := generateMoviesConfig(ctx)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Apply POST call testing and get ID in case of success
	objectId, err := testPostMovie(conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	log.Println("POST Call Response is: ", objectId)
	// Verify that the fields of Response and provided Configuration matches
	err = testGetMovieById(objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Update the configuration for further Testing
	updateMoviesConfig(conf)
	// Apply PUT call testing and get ID in case of success
	objectId, err = testPutMovieById(objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	// Verify that the fields of Response and provided Configuration matches
	err = testGetMovieById(objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// TODO: Unit Testing for PATCH call

	// Apply DELETE call testing and get ID in case of success
	err = testDeleteMovieById(objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	err = testGetMovieById(objectId, conf)
	if err == nil {
		t.Errorf("Record is not deleted. DELETE call failed!")
		return
	}
}

func TestCreateMovie_Created(t *testing.T) {

	authSrv := NewAuthServer(NewIdentityServer())
	TestMovieSrv = NewMovieServer(authSrv)

	obj := &moviepb.CreateMovieRequest{
		Movie: &moviepb.Movie{
			Name:     "movie_test",
			Summary:  "movie summary to be added",
			Cast:     []string{"cast_test1", "cast_test2"},
			Tags:     []moviepb.Tag{moviepb.Tag_Action, moviepb.Tag_Adventure},
			Director: "director_test1",
			Writers:  []string{"director_test1", "director_test2"},
		},
	}

	resp, err := TestMovieSrv.CreateMovie(context.Background(), obj)
	assert.Nil(t, err)
	log.Println("Response without error!")
	assert.NotEmpty(t, resp)
}

func generateMoviesConfig(ctx context.Context) (*TestConfig, error) {

	iSrv := NewIdentityServer()
	authSrv := NewAuthServer(iSrv)
	// token, err := validateUser(ctx, iSrv, authSrv)
	// if err != nil {
	// 	return nil, err
	// }

	// ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
	// // if _, ok := metadata.FromIncomingContext(ctx); !ok {
	// // 	return nil, status.Errorf(codes.Unauthenticated, "metadata is not added!")
	// // }

	return &TestConfig{
		Server: NewMovieServer(authSrv),
		URL:    "/v1/movies",
		Body: &moviepb.Movie{
			Name:    "Movie_test",
			Summary: "Some test summary!",
			Cast:    []string{"Cast_test1", "Cast_test2"},
			Tags:    []moviepb.Tag{moviepb.Tag_Adventure, moviepb.Tag_Fantasy},
		},
	}, nil
}

func updateMoviesConfig(config *TestConfig) {

	config.Body = &moviepb.Movie{
		Name:     "Movie_test",
		Summary:  "Random summary: " + String(40),
		Cast:     []string{String(8) + "_test1", String(8) + "_test2"},
		Director: String(10),
		Writers:  []string{String(10), String(8), String(6)},
		Tags:     []moviepb.Tag{moviepb.Tag_Adventure, moviepb.Tag_Fantasy},
	}
}

func testGetMovieById(objID string, config *TestConfig) error {

	// Convert srv from interface to server
	srv := config.Server.(*movieServer)

	resp, err := srv.GetMovie(config.Context, &moviepb.GetMovieRequest{Id: objID})
	if err != nil {
		return err
	}

	// Validate output is not empty
	if resp.String() == "" {
		return fmt.Errorf("Empty Response of GetByID call!")
	}

	// Verify the fields in response matches with the config passed
	configMv := config.Body.(*moviepb.Movie)

	if !reflect.DeepEqual(resp, configMv) {
		return fmt.Errorf("Fields do not match the required configuration!")
	}

	return nil
}

func testPostMovie(config *TestConfig) (string, error) {

	// Convert srv from interface to server
	srv := config.Server.(*movieServer)
	mvObj := config.Body.(*moviepb.Movie)

	resp, err := srv.CreateMovie(config.Context, &moviepb.CreateMovieRequest{Movie: mvObj})
	if err != nil {
		return "", err
	}

	// Verify Response is not empty
	if resp.String() == "" {
		return "", fmt.Errorf("Empty Response of POST call!")
	}

	id := resp.GetId()
	return id, nil
}

func testPutMovieById(objectId string, config *TestConfig) (string, error) {

	// Convert srv from interface to server
	srv := config.Server.(*movieServer)
	mvObject := config.Body.(*moviepb.Movie)

	resp, err := srv.UpdateMovie(config.Context, &moviepb.UpdateMovieRequest{Id: objectId, Movie: mvObject})
	if err != nil {
		return "", err
	}

	// Verify Response is not empty
	if resp.String() == "" {
		return "", fmt.Errorf("Empty Response of PUT call!")
	}
	return resp.GetId(), nil
}

func testDeleteMovieById(objectId string, config *TestConfig) error {

	// Convert srv from interface to server
	srv := config.Server.(*movieServer)

	resp, err := srv.DeleteMovie(config.Context, &moviepb.DeleteMovieRequest{Id: objectId})
	if err != nil {
		return err
	}

	// Verify Response is empty
	if resp.String() != "" {
		return fmt.Errorf("Expecting Empty Response by DELETE call!")
	}

	return nil
}

// Steps in brief
// 1. Initiate with a POST call
// 2. Check whether the POST call was successful by GET call and asserting fields match
// 3. Update the data using PUT call
// 4. Repeat Step 2
// 5. Update the data using PATCH call
// 6. Repeat Step 2
// 7. Delete the object using DELETE call
// 8. Repeat Step 2 to verify deletion
