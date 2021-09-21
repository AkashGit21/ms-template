package services

import (
	"context"
	"fmt"
	"log"
	"reflect"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type movieServer struct {
	Store map[string]*moviepb.Movie
}

func NewMovieServer() *movieServer {
	return &movieServer{
		Store: make(map[string]*moviepb.Movie),
	}
}

func (ms *movieServer) ListMovies(context.Context, *moviepb.ListMoviesRequest) (*moviepb.ListMoviesResponse, error) {
	pageSize := 4
	var out []*moviepb.Movie

	objectsIn := 0
	for _, object := range ms.Store {
		out = append(out, object)
		if objectsIn++; objectsIn > pageSize {
			break
		}
	}

	return &moviepb.ListMoviesResponse{
		Movies: out,
	}, nil
}

func (ms *movieServer) GetMovie(ctx context.Context, req *moviepb.GetMovieRequest) (*moviepb.GetMovieResponse, error) {

	log.Println("[DEBUG] Beginning GetMovieRequest: ", req)

	objID := req.GetId()
	log.Println("Object ID is: ", objID)

	// Check if Object exists or not
	// codes.NotFound
	if obj, ok := ms.Store[objID]; ok {
		// Found the required Object, hence return it
		return &moviepb.GetMovieResponse{
			Movie: obj,
		}, nil
	} else {
		return nil, status.Errorf(codes.NotFound, "Record with ID:%v does not exist!", objID)
	}
}

func (ms *movieServer) CreateMovie(ctx context.Context, req *moviepb.CreateMovieRequest) (*moviepb.CreateMovieResponse, error) {

	log.Println("[DEBUG] Beginning CreateMovieRequest: ", req)
	objID := generateUUID()

	// Check if Object already exists -
	// codes.AlreadyExists
	if _, ok := ms.Store[objID]; ok {
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

		ms.Store[objID] = mvObject
	}

	log.Println("[DEBUG] End CreateMovieRequest!")
	return &moviepb.CreateMovieResponse{
		Id: objID,
	}, nil
}

func (ms *movieServer) UpdateMovie(ctx context.Context, req *moviepb.UpdateMovieRequest) (*moviepb.UpdateMovieResponse, error) {

	log.Println("[DEBUG] Beginning UpdateMovieRequest: ", req)

	objID := req.GetId()
	log.Println("Object ID is: ", objID)

	// Check if object already exists or not
	// codes.NotFound
	if _, ok := ms.Store[objID]; ok {
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
		ms.Store[objID] = mvObject

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
	log.Println("Object ID is: ", objID)

	// Check if object already exists or not
	// codes.NotFound
	if mvObject, ok := ms.Store[objID]; ok {
		// Validate each field individually and Update the fields at a go afterwards

		// Verify that the ID is not updated
		newID := mvObject.GetId()
		if newID != "" && newID != objID {
			return nil, status.Errorf(codes.InvalidArgument, "Cannot update the ID of object!")
		}

		if smry := req.GetSummary(); smry != "" && smry != mvObject.GetSummary() {
			mvObject.Summary = smry
		}

		if dir := req.GetDirector(); dir != "" && dir != mvObject.GetDirector() {
			mvObject.Director = dir
		}

		if cast := req.GetCast(); cast != nil && reflect.DeepEqual(cast, mvObject.GetCast()) {
			mvObject.Cast = cast
		}

		if tags := req.GetTags(); tags != nil && reflect.DeepEqual(tags, mvObject.GetTags()) {
			mvObject.Tags = tags
		}

		if writers := req.GetWriters(); writers != nil && reflect.DeepEqual(writers, mvObject.GetCast()) {
			mvObject.Writers = writers
		}

		ms.Store[objID] = mvObject
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

	objID := req.GetId()
	log.Println("Object ID is: ", objID)

	// Check if object already exists or not
	// codes.NotFound
	if _, ok := ms.Store[objID]; ok {
		delete(ms.Store, objID)
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
