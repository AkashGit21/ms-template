package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
)

func TestIdentityService(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create some basic configuration for testing
	identitySrv := NewIdentityServer()
	conf := generateIdentityConfig(identitySrv)
	// Apply POST call testing and get ID in case of success
	objectId, err := testPostUser(ctx, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	// Verify that the fields of Response and provided Configuration matches
	err = testGetUserById(ctx, objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Update the configuration for further Testing
	updateIdentityConfig(conf)
	// Apply PUT call testing and get ID in case of success
	// objectId, err = testPutUserById(ctx, objectId, conf)
	// if err != nil {
	// 	t.Errorf(err.Error())
	// 	return
	// }
	// // Verify that the fields of Response and provided Configuration matches
	// err = testGetUserById(ctx, objectId, conf)
	// if err != nil {
	// 	t.Errorf(err.Error())
	// 	return
	// }

	// // TODO: Unit Testing for PATCH call

	// Apply DELETE call testing and get ID in case of success
	err = testDeleteUserById(ctx, objectId, conf)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	err = testGetUserById(ctx, objectId, conf)
	if err == nil {
		t.Errorf("Record is not deleted. DELETE call failed!")
		return
	}
}

func generateIdentityConfig(srv *identityServer) *TestConfig {
	return &TestConfig{
		Server: srv,
		URL:    "/v1/users",
		Body: &identitypb.User{
			Username:  "usrname1",
			Email:     "somemail98@domain.com",
			Role:      identitypb.Role_ADMIN,
			FirstName: "Alice",
			Password:  "abcdef",
		},
	}
}

func updateIdentityConfig(config *TestConfig) {
	config.Body = &identitypb.User{
		Username:  "usrname2",
		Email:     "updatedmail98@domain.com",
		FirstName: "John",
		Role:      identitypb.Role_NORMAL,
		Password:  "abcdef123",
	}
}

func testGetUserById(ctx context.Context, objID string, config *TestConfig) error {

	// Convert srv from interface to server
	srv := config.Server.(*identityServer)

	resp, err := srv.GetUser(ctx, &identitypb.GetUserRequest{Username: objID})
	if err != nil {
		return err
	}

	// Validate output is not empty
	if resp.String() == "" {
		return fmt.Errorf("Empty Response of GetByID call!")
	}

	// Verify the fields in response matches with the config passed
	configUser := config.Body.(*identitypb.User)

	if !reflect.DeepEqual(resp, configUser) {
		return fmt.Errorf("Fields do not match the required configuration!")
	}

	return nil
}

func testPostUser(ctx context.Context, config *TestConfig) (string, error) {

	// Convert srv from interface to server
	srv := config.Server.(*identityServer)
	userObj := config.Body.(*identitypb.User)

	resp, err := srv.CreateUser(ctx, &identitypb.CreateUserRequest{User: userObj})
	if err != nil {
		return "", err
	}

	// Verify Response is not empty
	if resp.String() == "" {
		return "", fmt.Errorf("Empty Response of POST call!")
	}

	id := resp.GetUsername()
	return id, nil
}

// func testPutUserById(ctx context.Context, objectId string, config *TestConfig) (string, error) {

// 	// Convert srv from interface to server
// 	srv := config.Server.(*identityServer)
// 	userObject := config.Body.(*identitypb.User)

// 	resp, err := srv.UpdateUser(ctx, &identitypb.UpdateUserRequest{Id: objectId, User: userObject})
// 	if err != nil {
// 		return "", err
// 	}

// 	// Verify Response is not empty
// 	if resp.String() == "" {
// 		return "", fmt.Errorf("Empty Response of PUT call!")
// 	}
// 	return resp.GetId(), nil
// }

func testDeleteUserById(ctx context.Context, objectId string, config *TestConfig) error {

	// Convert srv from interface to server
	srv := config.Server.(*identityServer)

	resp, err := srv.DeleteUser(ctx, &identitypb.DeleteUserRequest{Username: objectId})
	if err != nil {
		return err
	}

	// Verify Response is empty
	if resp.String() != "" {
		return fmt.Errorf("Expecting Empty Response by DELETE call!")
	}

	return nil
}
