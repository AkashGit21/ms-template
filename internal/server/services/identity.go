package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewIdentityServer returns a new instance of showcase identity server.
func NewIdentityServer() *identityServer {
	return &identityServer{
		token: NewTokenGenerator(),
		keys:  map[string]int{},
	}
}

type userEntry struct {
	user    *identitypb.User
	deleted bool
}

// ReadOnlyIdentityServer provides a read-only interface of an identity server.
type ReadOnlyIdentityServer interface {
	GetUser(context.Context, *identitypb.GetUserRequest) (*identitypb.User, error)
	ListUsers(context.Context, *identitypb.ListUsersRequest) (*identitypb.ListUsersResponse, error)
}

type identityServer struct {
	uid   string // server.UniqID
	token string // services.TokenGenerator

	mu    sync.Mutex
	keys  map[string]int
	users []userEntry
}

// Creates a user.
func (s *identityServer) CreateUser(_ context.Context, in *identitypb.CreateUserRequest) (*identitypb.User, error) {
	log.Println("Beginning CreateUser request: ", in)
	s.mu.Lock()
	defer s.mu.Unlock()

	u := in.GetUser()
	log.Println("User: ", u)

	// Ignore passed in name.
	u.Name = ""

	err := s.validate(u)
	if err != nil {
		return nil, err
	}

	// Assign info.
	id := 3 // seededRand.Intn(len(charset)) // s.uid.Next()
	name := fmt.Sprintf("users/%d", id)
	now := ptypes.TimestampNow()

	u.Name = name
	u.CreateTime = now
	u.UpdateTime = now

	// Insert.
	index := len(s.users)
	s.users = append(s.users, userEntry{user: u})
	s.keys[name] = index

	return u, nil
}

// Retrieves the User with the given uri.
func (s *identityServer) GetUser(_ context.Context, in *identitypb.GetUserRequest) (*identitypb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := in.GetName()
	if i, ok := s.keys[name]; ok {
		entry := s.users[i]
		if !entry.deleted {
			return entry.user, nil
		}
	}

	return nil, status.Errorf(
		codes.NotFound, "A user with name %s not found.",
		name)
}

// Updates a user.
func (s *identityServer) UpdateUser(_ context.Context, in *identitypb.UpdateUserRequest) (*identitypb.User, error) {
	mask := in.GetUpdateMask()
	if mask != nil && len(mask.GetPaths()) > 0 {
		return nil, status.Error(
			codes.Unimplemented,
			"Field masks are currently not supported.")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	u := in.GetUser()
	i, ok := s.keys[u.GetName()]
	if !ok || s.users[i].deleted {
		return nil, status.Errorf(
			codes.NotFound,
			"A user with name %s not found.", u.GetName())
	}

	err := s.validate(u)
	if err != nil {
		return nil, err
	}
	entry := s.users[i]
	// Update store.
	updated := &identitypb.User{
		Name:                u.GetName(),
		DisplayName:         u.GetDisplayName(),
		Email:               u.GetEmail(),
		CreateTime:          entry.user.GetCreateTime(),
		UpdateTime:          ptypes.TimestampNow(),
		Age:                 entry.user.Age,
		EnableNotifications: entry.user.EnableNotifications,
		HeightFeet:          entry.user.HeightFeet,
		Nickname:            entry.user.Nickname,
	}

	// Use direct field access to avoid unwrapping and rewrapping the pointer value.
	//
	// TODO: if field_mask is implemented, do a direct update if included,
	// regardless of if the optional field is nil.
	if u.Age != nil {
		updated.Age = u.Age
	}
	if u.EnableNotifications != nil {
		updated.EnableNotifications = u.EnableNotifications
	}
	if u.HeightFeet != nil {
		updated.HeightFeet = u.HeightFeet
	}
	if u.Nickname != nil {
		updated.Nickname = u.Nickname
	}

	s.users[i] = userEntry{user: updated}
	return updated, nil
}

// Deletes a user, their profile, and all of their authored messages.
func (s *identityServer) DeleteUser(_ context.Context, in *identitypb.DeleteUserRequest) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.keys[in.GetName()]

	if !ok {
		return nil, status.Errorf(
			codes.NotFound,
			"A user with name %s not found.", in.GetName())
	}

	entry := s.users[i]
	s.users[i] = userEntry{user: entry.user, deleted: true}

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

func (s *identityServer) validate(u *identitypb.User) error {
	// Validate Required Fields.
	if u.GetDisplayName() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `display_name` is required.")
	}
	if u.GetEmail() == "" {
		return status.Errorf(
			codes.InvalidArgument,
			"The field `email` is required.")
	}
	// Validate Unique Fields.
	for _, x := range s.users {
		if x.deleted {
			continue
		}
		if (u.GetEmail() == x.user.GetEmail()) &&
			(u.GetName() != x.user.GetName()) {
			return status.Errorf(
				codes.AlreadyExists,
				"A user with email %s already exists.",
				u.GetEmail())
		}
	}
	return nil
}
