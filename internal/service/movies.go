package service

import (
	"context"
	"fmt"
	"log"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MovieServer struct {
	moviepb.UnimplementedMovieServiceServer
	Store map[string]*moviepb.Movie
}

func NewMovieServer() *MovieServer {
	return &MovieServer{
		Store: make(map[string]*moviepb.Movie),
	}
}
func (ms *MovieServer) ListMovies(context.Context, *moviepb.ListMoviesRequest) (*moviepb.ListMoviesResponse, error) {
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

func (ms *MovieServer) GetMovie(ctx context.Context, req *moviepb.GetMovieRequest) (*moviepb.GetMovieResponse, error) {

	log.Println("[DEBUG] Beginning GetMovieRequest: ", req)

	objID := req.GetId()
	log.Println("Object ID is: ", objID)

	log.Println("Store: ", ms.Store)

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

func (ms *MovieServer) CreateMovie(ctx context.Context, req *moviepb.CreateMovieRequest) (*moviepb.CreateMovieResponse, error) {

	log.Println("[DEBUG] Beginning CreateMovieRequest: ", req)

	objID := req.GetMovie().GetId()
	log.Println("Object ID is: ", objID)

	// Check if Object already exists -
	// codes.AlreadyExists
	if _, ok := ms.Store[objID]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "Movie Record with ID: %v already exists!", objID)
	} else {
		mvObject := req.GetMovie()

		// Validate format of Input and store the data
		IsValidMovie(mvObject)
		ms.Store[objID] = mvObject
	}

	log.Println("[DEBUG] End CreateMovieRequest!")
	return &moviepb.CreateMovieResponse{
		Id: objID,
	}, nil
}

func (ms *MovieServer) UpdateMovie(ctx context.Context, req *moviepb.UpdateMovieRequest) (*moviepb.UpdateMovieResponse, error) {

	log.Println("[DEBUG] Beginning UpdateMovieRequest: ", req)

	objID := req.GetId()
	log.Println("Object ID is: ", objID)

	// Check if object already exists or not
	// codes.NotFound
	if _, ok := ms.Store[objID]; ok {
		// Validate and update the whole object

	} else {
		return nil, status.Errorf(codes.NotFound, "Movie Record with ID:%v does not exist!", objID)
	}

	log.Println("[DEBUG] End UpdateMovieRequest!")
	return &moviepb.UpdateMovieResponse{
		Id: "23",
	}, nil
}

func (ms *MovieServer) PartialUpdateMovie(ctx context.Context, req *moviepb.PartialUpdateMovieRequest) (*moviepb.PartialUpdateMovieResponse, error) {

	log.Println("[DEBUG] Beginning PartialUpdateMovieRequest: ", req)

	objID := req.GetId()
	log.Println("Object ID is: ", objID)

	// Check if object already exists or not
	// codes.NotFound
	if _, ok := ms.Store[objID]; ok {
		// Validate each field individually and Update the fields at a go afterwards

	} else {
		return nil, status.Errorf(codes.NotFound, "Movie Record with ID:%v does not exist!", objID)
	}

	log.Println("[DEBUG] End PartialUpdateMovieRequest!")
	return &moviepb.PartialUpdateMovieResponse{
		Id: "123",
	}, nil
}

func (ms *MovieServer) DeleteMovie(ctx context.Context, req *moviepb.DeleteMovieRequest) (*empty.Empty, error) {

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
	return nil, nil
}

func IsValidMovie(mv *moviepb.Movie) (bool, error) {
	if !StringLenBetween(mv.Name, 1, 120) {
		return false, fmt.Errorf("The name should be between 1 and 120 characters.")
	}

	if err := isAllowedSummary(mv.Summary); err != nil {
		return false, err
	}

	return true, nil
}
