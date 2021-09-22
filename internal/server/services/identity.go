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
		token: server.NewJWTTokenGenerator(),
		keys:  map[string]int{},
	}
}

type userEntry struct {
	user   *identitypb.User
	active bool
}

type identityServer struct {
	token server.JWTTokenGenerator

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

	objID := server.GenerateUUID()

	// Check if Object already exists -
	// codes.AlreadyExists

	if _, ok := is.keys[objID]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "User Record with ID: %v already exists!", objID)
	} else {

		u := req.GetUser()
		if u.GetId() != "" {
			return nil, status.Errorf(codes.InvalidArgument, "Input is not valid! ID is auto-generated")
		}

		// Validate format of Input and store the data
		err := is.validate(u)
		if err != nil {
			return nil, err
		}

		// Assign server generated info.
		now := ptypes.TimestampNow()

		u.Id = objID
		u.CreateTime = now
		u.UpdateTime = now

		// Insert.
		index := len(is.userEntries)
		is.userEntries = append(is.userEntries, userEntry{user: u, active: true})
		is.keys[objID] = index

	}

	return &identitypb.CreateUserResponse{Id: objID}, nil
}

// Retrieves the User with the given uri.
func (is *identityServer) GetUser(_ context.Context, req *identitypb.GetUserRequest) (*identitypb.User, error) {
	log.Println("Beginning GetUser request: ", req)
	is.mu.Lock()
	defer is.mu.Unlock()

	objID := req.GetId()

	log.Println("Keys: ", is.keys)
	log.Println("Entries: ", is.userEntries)
	// Check if Object exists or not
	// codes.NotFound
	if obj, ok := is.keys[objID]; ok {
		entry := is.userEntries[obj]
		if entry.active {
			return entry.user, nil
		}
	}

	return nil, status.Errorf(
		codes.NotFound, "A user with id %s not found!",
		objID)
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

	i, ok := is.keys[req.GetId()]

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"A user with id %s not found.", req.GetId())
	}

	entry := is.userEntries[i]
	is.userEntries[i] = userEntry{user: entry.user, active: false}

	return &empty.Empty{}, nil
}

// Lists all users.
func (s *identityServer) ListUsers(_ context.Context, in *identitypb.ListUsersRequest) (*identitypb.ListUsersResponse, error) {
	// start, err := s.token.GetIndex(in.GetPageToken())
	// if err != nil {
	// 	return nil, err
	// }

	// offset := 0
	// users := []*identitypb.User{}
	// for _, entry := range s.users[start:] {
	// 	offset++
	// 	if entry.deleted {
	// 		continue
	// 	}
	// 	users = append(users, entry.user)
	// 	if len(users) >= int(in.GetPageSize()) {
	// 		break
	// 	}
	// }

	// nextToken := ""
	// if start+offset < len(s.users) {
	// 	nextToken = s.token.ForIndex(start + offset)
	// }

	// return &identitypb.ListUsersResponse{Users: users, NextPageToken: nextToken}, nil
	return &identitypb.ListUsersResponse{}, nil
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
			(u.GetId() != x.user.GetId()) {
			return status.Errorf(
				codes.AlreadyExists,
				"A user with email %s already exists.",
				u.GetEmail())
		}
	}
	return nil
}
