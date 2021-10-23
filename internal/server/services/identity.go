package services

import (
	"context"
	"log"
	"sync"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/internal/server/interceptors"
	"github.com/AkashGit21/ms-project/lib/persistence"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewIdentityServer returns a new instance of application identity server.
func NewIdentityServer(dbHandler persistence.DatabaseHandler) *identityServer {
	return &identityServer{
		token:     server.NewTokenGenerator(),
		keys:      map[string]int{},
		dbhandler: dbHandler,
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

	dbhandler persistence.DatabaseHandler
	identitypb.UnimplementedIdentityServiceServer
}

// ReadOnlyIdentityServer provides a read-only interface of an identity server.
type ReadOnlyIdentityServer interface {
	GetUser(context.Context, *identitypb.GetUserRequest) (*identitypb.User, error)
	ListUsers(context.Context, *identitypb.ListUsersRequest) (*identitypb.ListUsersResponse, error)
}

// Creates a user.
func (is *identityServer) CreateUser(_ context.Context,
	req *identitypb.CreateUserRequest) (*identitypb.CreateUserResponse, error) {
	log.Println("Beginning CreateUser request: ", req)
	is.mu.Lock()
	defer is.mu.Unlock()

	u := req.GetUser()

	// Check if Object already exists -
	// codes.AlreadyExists
	_, err := is.dbhandler.FindByUsername(u.GetUsername())
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists,
			"A user with username `%s` already exists!", u.GetUsername())
	} else {

		// Validate format of Input and store the data
		err := is.validate(u)
		if err != nil {
			return nil, err
		}

		// Assign server generated info.
		now := ptypes.TimestampNow()

		pwd, err := server.HashPassword(u.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		user := persistence.User{
			Username:    u.Username,
			Email:       u.Email,
			Password:    pwd,
			Role:        persistence.Role(u.Role),
			Active:      true,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Age:         u.Age,
			HeightInCms: u.HeightInCms,
			CreateTime:  now,
			UpdateTime:  now,
			Nickname:    u.Nickname,
		}

		is.dbhandler.AddUser(user)
	}
	log.Println("End of CreateUser!")

	return &identitypb.CreateUserResponse{
		Username: u.GetUsername(),
	}, nil
}

// Retrieves the User with the given uri.
func (is *identityServer) GetUser(_ context.Context,
	req *identitypb.GetUserRequest) (*identitypb.User, error) {
	log.Println("Beginning GetUser request: ", req)
	is.mu.Lock()
	defer is.mu.Unlock()

	uname := req.GetUsername()

	// Only ADMIN or the user itself can view his/her information
	if interceptors.CURRENT_ROLE != "ADMIN" && interceptors.CURRENT_USERNAME != uname {
		return nil, status.Error(codes.PermissionDenied,
			"not allowed to perform this operation!")
	}

	res, err := is.dbhandler.FindByUsername(uname)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound, "A user with username `%s` not found!",
			uname)
	}

	user := &identitypb.User{
		Username:    res.Username,
		Email:       res.Email,
		Role:        identitypb.Role(res.Role),
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		CreateTime:  res.CreateTime,
		UpdateTime:  res.UpdateTime,
		Age:         res.Age,
		HeightInCms: res.HeightInCms,
		Nickname:    res.Nickname,
		Active:      res.Active,
	}
	return user, nil
}

// Updates a user.
func (is *identityServer) UpdateUser(_ context.Context,
	req *identitypb.UpdateUserRequest) (*identitypb.User, error) {
	// TODO: Add the working for UpdateUser
	return &identitypb.User{}, nil
}

// Deletes a user, their profile, and all of their authored messages.
func (is *identityServer) DeleteUser(_ context.Context,
	req *identitypb.DeleteUserRequest) (*empty.Empty, error) {
	log.Println("Beginning DeleteUser request: ", req)

	uname := req.GetUsername()

	if interceptors.CURRENT_ROLE != "ADMIN" && interceptors.CURRENT_USERNAME != uname {
		return nil, status.Error(codes.PermissionDenied,
			"not allowed to perform this operation!")
	}

	// Check if object already exists or not
	// codes.NotFound
	if _, err := is.dbhandler.FindByUsername(uname); err != nil {
		return nil, status.Errorf(codes.NotFound, "A user with username `%s` does not exist!", uname)
	}

	if err := is.dbhandler.RemoveByUsername(uname); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "some error while deleting user!")
		// return nil, status.Errorf(codes.Internal, "Internal Server Error!", uname)
	}
	log.Println("[DEBUG] End DeleteUser!")
	return &empty.Empty{}, nil
}

// Lists all users.
func (is *identityServer) ListUsers(_ context.Context,
	in *identitypb.ListUsersRequest) (*identitypb.ListUsersResponse, error) {
	start, err := is.token.GetIndex(in.GetPageToken())
	if err != nil {
		return nil, err
	}

	// Default page size is 12
	var pageSz int32
	if pageSz = in.GetPageSize(); pageSz == 0 || pageSz > 12 {
		pageSz = 12
	}
	// offset := 0

	numOfUsers := is.dbhandler.CountUsers()

	users, err := is.dbhandler.FindAllUsers(start, pageSz)
	if err != nil {
		return nil, err
	}

	nextToken := ""
	if start+int(pageSz) < numOfUsers {
		nextToken = is.token.ForIndex(start + int(pageSz))
	}

	return &identitypb.ListUsersResponse{
		Users:         users,
		NextPageToken: nextToken,
	}, nil
}

// TODO: Add Validation for similar email in DB
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
				"A user with email `%s` already exists.",
				u.GetEmail())
		}
	}
	return nil
}
