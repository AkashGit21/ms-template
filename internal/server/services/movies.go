package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type movieServer struct {
	token   server.TokenGenerator
	authSrv *authServer

	mu    sync.Mutex
	keys  map[string]int
	Store []movieEntry

	moviepb.UnimplementedMovieServiceServer
}

type movieEntry struct {
	movie  *moviepb.Movie
	active bool
}

func NewMovieServer(as *authServer) *movieServer {
	return &movieServer{
		token:   server.NewTokenGenerator(),
		authSrv: as,
		keys:    map[string]int{},
	}
}

func isAuthorized(service string, role string) bool {
	if strings.HasPrefix(service, "Get") || strings.HasPrefix(service, "List") {
		return true
	} else if (strings.HasPrefix(service, "Create") || strings.HasPrefix(service, "Update") || strings.HasPrefix(service, "Delete")) && (role == "ADMIN" || role == "SUBSCRIBED") {
		return true
	}
	return false
}

func (ms *movieServer) verifyAuthorization(ctx context.Context, serviceName string) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	var role string
	values := md["authorization"]
	if len(values) == 0 {
		role = "GUEST"
	} else {

		accessToken := strings.TrimPrefix(values[0], "Basic ")

		claims, err := ms.authSrv.JWT.GetUserFromToken(accessToken)
		if err != nil {
			return err
		}
		role = claims.Role
	}

	if !isAuthorized(serviceName, role) {
		return fmt.Errorf("not authorized to perform this action!")
	}
	return nil
}

func (ms *movieServer) ListMovies(ctx context.Context, req *moviepb.ListMoviesRequest) (*moviepb.ListMoviesResponse, error) {

	err := ms.verifyAuthorization(ctx, "ListMovies")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized to perform the folloowing action! %v", err)
	}

	start, err := ms.token.GetIndex(req.GetPageToken())
	if err != nil {
		return nil, err
	}

	// Default page size is 12
	var pageSz int32
	if pageSz = req.GetPageSize(); pageSz == 0 || pageSz > 12 {
		pageSz = 12
	}
	offset := 0

	out := []*moviepb.Movie{}

	for _, entry := range ms.Store[start:] {
		offset++
		if !entry.active {
			continue
		}
		out = append(out, entry.movie)
		if len(out) >= int(pageSz) {
			break
		}
	}

	nextToken := ""
	if start+offset < len(ms.Store) {
		nextToken = ms.token.ForIndex(start + offset)
	}

	return &moviepb.ListMoviesResponse{
		Movies:        out,
		NextPageToken: nextToken,
	}, nil
}

func (ms *movieServer) GetMovie(ctx context.Context, req *moviepb.GetMovieRequest) (*moviepb.Movie, error) {

	log.Println("[DEBUG] Beginning GetMovieRequest: ", req)

	err := ms.verifyAuthorization(ctx, "GetMovie")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized to perform the folloowing action! %v", err)
	}

	objID := req.GetId()

	// Check if Object exists or not
	// codes.NotFound
	if obj, ok := ms.keys[objID]; ok && ms.Store[obj].active {
		// Found the required Object, hence return it

		return ms.Store[obj].movie, nil
	} else {
		return nil, status.Errorf(codes.NotFound, "Record with ID:%v does not exist!", objID)
	}
}

func (ms *movieServer) CreateMovie(ctx context.Context, req *moviepb.CreateMovieRequest) (*moviepb.CreateMovieResponse, error) {

	log.Println("[DEBUG] Beginning CreateMovieRequest: ", req)

	err := ms.verifyAuthorization(ctx, "CreateMovie")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized to perform the folloowing action! %v", err)
	}

	objID := server.GenerateUUID()

	// Check if Object already exists -
	// codes.AlreadyExists
	if _, ok := ms.keys[objID]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "Movie Record with ID: %v already exists!", objID)
	} else {
		mvObject := req.GetMovie()
		if mvObject.GetId() != "" {
			return nil, status.Errorf(codes.InvalidArgument, "Input is not valid! ID is auto-generated")
		}
		mvObject.Id = objID

		// Validate format of Input and store the data
		if valid, err := IsValidMovie(mvObject); !valid {
			return nil, status.Errorf(codes.InvalidArgument, "Input is not valid! %v", err.Error())
		}

		ms.mu.Lock()
		defer ms.mu.Unlock()
		index := len(ms.Store)
		ms.keys[objID] = index
		ms.Store = append(ms.Store, movieEntry{movie: mvObject, active: true})
	}

	log.Println("[DEBUG] End CreateMovieRequest!")
	return &moviepb.CreateMovieResponse{
		Id: objID,
	}, nil
}

func (ms *movieServer) UpdateMovie(ctx context.Context, req *moviepb.UpdateMovieRequest) (*moviepb.UpdateMovieResponse, error) {

	log.Println("[DEBUG] Beginning UpdateMovieRequest: ", req)

	err := ms.verifyAuthorization(ctx, "UpdateMovie")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized to perform the folloowing action! %v", err)
	}

	objID := req.GetId()

	// Check if object already exists or not
	// codes.NotFound
	if index, ok := ms.keys[objID]; ok && ms.Store[index].active {
		// Validate and update the whole object
		mvObject := req.GetMovie()

		// Verify that the ID is not updated
		newID := mvObject.GetId()
		if newID != "" && !reflect.DeepEqual(newID, objID) {
			return nil, status.Errorf(codes.InvalidArgument, "Cannot update the ID of object!")
		}

		// Validate format of Input and store the data
		if valid, err := IsValidMovie(mvObject); !valid {
			return nil, status.Errorf(codes.InvalidArgument, "Input is not valid! %v", err.Error())
		}

		mvObject.Id = objID
		ms.Store[ms.keys[objID]] = movieEntry{movie: mvObject, active: true}

	} else {
		return nil, status.Errorf(codes.NotFound, "Movie Record with ID:%v does not exist!", objID)
	}

	log.Println("[DEBUG] End UpdateMovieRequest!")
	return &moviepb.UpdateMovieResponse{
		Id: objID,
	}, nil
}

func (ms *movieServer) PartialUpdateMovie(ctx context.Context, req *moviepb.PartialUpdateMovieRequest) (*moviepb.PartialUpdateMovieResponse, error) {

	log.Println("[DEBUG] Beginning PartialUpdateMovieRequest: ", req)

	objID := req.GetId()

	// Check if object already exists or not
	// codes.NotFound
	if index, ok := ms.keys[objID]; ok && ms.Store[index].active {

		// Validate each field individually and Update the fields at a go afterwards
		ms.mu.Lock()
		defer ms.mu.Unlock()
		mvObject := ms.Store[index].movie

		// Verify that the ID is not updated
		newID := req.GetId()
		if newID != "" && newID != objID {
			return nil, status.Errorf(codes.InvalidArgument, "Cannot update the ID of object!")
		}

		if req.GetSummary() != "" {
			mvObject.Summary = req.GetSummary()
		}

		if req.GetDirector() != "" {
			mvObject.Director = req.GetDirector()
		}

		if req.GetCast() != nil {
			mvObject.Cast = req.GetCast()
		}

		if req.GetTags() != nil {
			mvObject.Tags = req.GetTags()
		}

		if req.GetWriters() != nil {
			mvObject.Writers = req.GetWriters()
		}

	} else {
		return nil, status.Errorf(codes.NotFound, "Movie Record with ID:%v does not exist!", objID)
	}

	log.Println("[DEBUG] End PartialUpdateMovieRequest!")
	return &moviepb.PartialUpdateMovieResponse{
		Id: objID,
	}, nil
}

func (ms *movieServer) DeleteMovie(ctx context.Context, req *moviepb.DeleteMovieRequest) (*empty.Empty, error) {

	log.Println("[DEBUG] Beginning DeleteMovieRequest: ", req)

	err := ms.verifyAuthorization(ctx, "DeleteMovie")
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized to perform the folloowing action! %v", err)
	}

	objID := req.GetId()

	// Check if object already exists or not
	// codes.NotFound
	if index, ok := ms.keys[objID]; ok && ms.Store[index].active {
		ms.mu.Lock()
		defer ms.mu.Unlock()

		ms.Store[index].active = false

	} else {
		return nil, status.Errorf(codes.NotFound, "Movie Record with ID:%v does not exist!", objID)
	}

	log.Println("[DEBUG] End DeleteMovieRequest!")
	return &empty.Empty{}, nil
}

func IsValidMovie(mv *moviepb.Movie) (bool, error) {

	if !isValidName(mv.Name) {
		return false, fmt.Errorf("The name should be between 1 and 120 characters.")
	}

	if err := isAllowedSummary(mv.Summary); err != nil {
		return false, err
	}

	if len(mv.Cast) > 0 {
		for _, crew := range mv.Cast {
			if !isValidName(crew) {
				return false, fmt.Errorf("Name of Cast members should be between 1 and 120 characters.")
			}
		}
	}

	if mv.Director != "" && !isValidName(mv.Director) {
		return false, fmt.Errorf("Director field should be between 1 and 120 characters.")
	}

	if len(mv.Cast) > 0 {
		for _, wr := range mv.Cast {
			if !isValidName(wr) {
				return false, fmt.Errorf("Writers' name should be between 1 and 120 characters.")
			}
		}
	}

	return true, nil
}
