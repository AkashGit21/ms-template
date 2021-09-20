package service

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
)

func TestMoviesService(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create some basic configuration for testing
	conf := generateConfig()
	// Apply POST call testing and get ID in case of success
	objectId, err := testPostMovie(ctx, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	log.Println("POST Call Response is: ", objectId)
	// Verify that the fields of Response and provided Configuration matches
	err = testGetMovieById(ctx, objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Update the configuration for further Testing
	updateConfig(conf)
	// Apply PUT call testing and get ID in case of success
	objectId, err = testPutMovieById(ctx, objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	// Verify that the fields of Response and provided Configuration matches
	err = testGetMovieById(ctx, objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// TODO: Unit Testing for PATCH call

	// Apply DELETE call testing and get ID in case of success
	err = testDeleteMovieById(ctx, objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	err = testGetMovieById(ctx, objectId, conf)
	if err == nil {
		t.Errorf("Record is not deleted. DELETE call failed!")
		return
	}
}

type TestConfig struct {
	Server interface{}
	URL    string
	Body   interface{}
}

func generateConfig() *TestConfig {
	return &TestConfig{
		Server: NewMovieServer(),
		URL:    "/v1/movies",
		Body: &moviepb.Movie{
			Name:    "Movie_test",
			Summary: "Some test summary!",
			Cast:    []string{"Cast_test1", "Cast_test2"},
			Tags:    []moviepb.Tag{moviepb.Tag_Adventure, moviepb.Tag_Fantasy},
		},
	}
}

func updateConfig(config *TestConfig) {

	config.Body = &moviepb.Movie{
		Name:     "Movie_test",
		Summary:  "Random summary: " + String(40),
		Cast:     []string{String(8) + "_test1", String(8) + "_test2"},
		Director: String(10),
		Writers:  []string{String(10), String(8), String(6)},
		Tags:     []moviepb.Tag{moviepb.Tag_Adventure, moviepb.Tag_Fantasy},
	}
}

func testGetMovieById(ctx context.Context, objID string, config *TestConfig) error {

	// Convert srv from interface to server
	srv := config.Server.(*MovieServer)

	resp, err := srv.GetMovie(ctx, &moviepb.GetMovieRequest{Id: objID})
	if err != nil {
		return err
	}

	// Validate output is not empty
	if resp.String() == "" {
		return fmt.Errorf("Empty Response of GetByID call!")
	}

	// Verify the fields in response matches with the config passed
	mvObject := resp.Movie
	configMv := config.Body.(*moviepb.Movie)

	if !reflect.DeepEqual(mvObject, configMv) {
		return fmt.Errorf("Fields do not match the required configuration!")
	}

	return nil
}

func testPostMovie(ctx context.Context, config *TestConfig) (string, error) {

	// Convert srv from interface to server
	srv := config.Server.(*MovieServer)
	mvObj := config.Body.(*moviepb.Movie)

	resp, err := srv.CreateMovie(ctx, &moviepb.CreateMovieRequest{Movie: mvObj})
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

func testPutMovieById(ctx context.Context, objectId string, config *TestConfig) (string, error) {

	// Convert srv from interface to server
	srv := config.Server.(*MovieServer)
	mvObject := config.Body.(*moviepb.Movie)

	resp, err := srv.UpdateMovie(ctx, &moviepb.UpdateMovieRequest{Id: objectId, Movie: mvObject})
	if err != nil {
		return "", err
	}

	// Verify Response is not empty
	if resp.String() == "" {
		return "", fmt.Errorf("Empty Response of PUT call!")
	}
	return resp.GetId(), nil
}

func testDeleteMovieById(ctx context.Context, objectId string, config *TestConfig) error {

	// Convert srv from interface to server
	srv := config.Server.(*MovieServer)

	resp, err := srv.DeleteMovie(ctx, &moviepb.DeleteMovieRequest{Id: objectId})
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
