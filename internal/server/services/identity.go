package services

import (
	"context"
	"log"
	"sync"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewIdentityServer returns a new instance of application identity server.
func NewIdentityServer() *identityServer {
	return &identityServer{
		token:       server.NewTokenGenerator(),
		keys:        map[string]int{},
		userEntries: []userEntry{},
	}
}

type userEntry struct {
	user   *identitypb.User
	active bool
}

type identityServer struct {
	token server.TokenGenerator

	mu          sync.Mutex
	keys        map[string]int
	userEntries []userEntry

	identitypb.UnimplementedIdentityServiceServer
}

// ReadOnlyIdentityServer provides a read-only interface of an identity server.
type ReadOnlyIdentityServer interface {
	GetUser(context.Context, *identitypb.GetUserRequest) (*identitypb.User, error)
	ListUsers(context.Context, *identitypb.ListUsersRequest) (*identitypb.ListUsersResponse, error)
}

// Creates a user.
func (is *identityServer) CreateUser(_ context.Context, req *identitypb.CreateUserRequest) (*identitypb.CreateUserResponse, error) {
	log.Println("Beginning CreateUser request: ", req)
	is.mu.Lock()
	defer is.mu.Unlock()

	user := req.GetUser()
	var uname string

	// Check if Object already exists -
	// codes.AlreadyExists

	if _, ok := is.keys[user.GetUsername()]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "A user with username `%s` already exists!", user.GetUsername())
	} else {

		uname = user.GetUsername()
		if uname == "" {
			return nil, status.Errorf(codes.InvalidArgument, "Input is not valid! Username is required.")
		}

		// Validate format of Input and store the data
		err := is.validate(user)
		if err != nil {
			return nil, err
		}

		// Assign server generated info.
		now := ptypes.TimestampNow()

		pwd, err := server.HashPassword(user.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		user.Password = pwd
		user.CreateTime = now
		user.UpdateTime = now

		// Insert.
		index := len(is.userEntries)
		is.userEntries = append(is.userEntries, userEntry{user: user, active: true})
		is.keys[user.GetUsername()] = index
	}
	log.Println("End of CreateUser!")

	return &identitypb.CreateUserResponse{Username: uname}, nil
}

// Retrieves the User with the given uri.
func (is *identityServer) GetUser(_ context.Context, req *identitypb.GetUserRequest) (*identitypb.User, error) {
	log.Println("Beginning GetUser request: ", req)
	is.mu.Lock()
	defer is.mu.Unlock()

	uname := req.GetUsername()
	// Check if Object exists or not
	// codes.NotFound
	if obj, ok := is.keys[uname]; ok {
		entry := is.userEntries[obj]
		if entry.active {
			return entry.user, nil
		}
	}

	return nil, status.Errorf(
		codes.NotFound, "A user with username `%s` not found!",
		uname)
}

// Updates a user.
func (is *identityServer) UpdateUser(_ context.Context, req *identitypb.UpdateUserRequest) (*identitypb.User, error) {
	// TODO: Add the working for UpdateUser
	return &identitypb.User{}, nil
}

// Deletes a user, their profile, and all of their authored messages.
func (is *identityServer) DeleteUser(_ context.Context, req *identitypb.DeleteUserRequest) (*empty.Empty, error) {
	log.Println("Beginning DeleteUser request: ", req)
	is.mu.Lock()
	defer is.mu.Unlock()

	index, ok := is.keys[req.GetUsername()]

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"A user with username `%s` not found.", req.GetUsername())
	}

	entry := is.userEntries[index]
	is.userEntries[index] = userEntry{user: entry.user, active: false}

	return &empty.Empty{}, nil
}

// Lists all users.
func (is *identityServer) ListUsers(_ context.Context, in *identitypb.ListUsersRequest) (*identitypb.ListUsersResponse, error) {
	start, err := is.token.GetIndex(in.GetPageToken())
	if err != nil {
		return nil, err
	}

	// Default page size is 12
	var pageSz int32
	if pageSz = in.GetPageSize(); pageSz == 0 || pageSz > 12 {
		pageSz = 12
	}
	offset := 0

	users := []*identitypb.User{}
	for _, entry := range is.userEntries[start:] {
		offset++
		if !entry.active {
			continue
		}
		users = append(users, entry.user)
		if len(users) >= int(pageSz) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(is.userEntries) {
		nextToken = is.token.ForIndex(start + offset)
	}

	return &identitypb.ListUsersResponse{
		Users:         users,
		NextPageToken: nextToken,
	}, nil
}

func (is *identityServer) validate(u *identitypb.User) error {
	// Validate Required Fields.
	if u.GetUsername() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `username` is required.")
	}
	if u.GetEmail() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `email` is required.")
	}
	// Validate Unique Fields.
	for _, x := range is.userEntries {
		if !x.active {
			continue
		}
		if (u.GetEmail() == x.user.GetEmail()) &&
			(u.GetUsername() != x.user.GetUsername()) {
			return status.Errorf(
				codes.AlreadyExists,
				"A user with email %s already exists.",
				u.GetEmail())
		}
	}
	return nil
}
