package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
	"github.com/AkashGit21/ms-project/internal/server"
	"github.com/AkashGit21/ms-project/lib/persistence"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type movieServer struct {
	token   server.TokenGenerator
	authSrv *authServer

	mu        sync.Mutex
	dbhandler persistence.DatabaseHandler
	keys      map[string]int
	Store     []movieEntry

	moviepb.UnimplementedMovieServiceServer
}

type movieEntry struct {
	movie  *moviepb.Movie
	active bool
}

func NewMovieServer(as *authServer) *movieServer {
	return &movieServer{
		token:     server.NewTokenGenerator(),
		authSrv:   as,
		keys:      map[string]int{},
		dbhandler: as.dbhandler,
	}
}

func (ms *movieServer) ListMovies(ctx context.Context,
	req *moviepb.ListMoviesRequest) (*moviepb.ListMoviesResponse, error) {

	start, err := ms.token.GetIndex(req.GetPageToken())
	if err != nil {
		return nil, err
	}

	// Default page size is 12
	var pageSz int32
	if pageSz = req.GetPageSize(); pageSz == 0 || pageSz > 12 {
		pageSz = 12
	}

	numOfRecords := ms.dbhandler.CountMovieRecords()

	movies, err := ms.dbhandler.FindAllMovies(start, pageSz)
	if err != nil {
		return nil, err
	}

	nextToken := ""
	if start+int(pageSz) < numOfRecords {
		nextToken = ms.token.ForIndex(start + int(pageSz))
	}

	return &moviepb.ListMoviesResponse{
		Movies:        movies,
		NextPageToken: nextToken,
	}, nil
}

func (ms *movieServer) GetMovie(ctx context.Context,
	req *moviepb.GetMovieRequest) (*moviepb.Movie, error) {
	log.Println("[DEBUG] Beginning GetMovieRequest: ", req)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	objID := req.GetId()

	// Check if Object exists or not
	// codes.NotFound
	res, err := ms.dbhandler.FindMovieByID(objID)
	if err != nil || !res.Active {

		return nil, status.Errorf(
			codes.NotFound, "Movie Record with ID:%v does not exist!",
			objID)
	}

	mv := &moviepb.Movie{
		Id:      res.Id,
		Name:    res.Name,
		Summary: res.Summary,
		Cast:    res.Cast,
		// TODO: Add the Tags here
		// Tags: res.Tags,
		Director:   res.Director,
		Writers:    res.Writers,
		Active:     res.Active,
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
	}

	return mv, nil
}

func (ms *movieServer) CreateMovie(ctx context.Context,
	req *moviepb.CreateMovieRequest) (*moviepb.CreateMovieResponse, error) {
	log.Println("[DEBUG] Beginning CreateMovieRequest: ", req)

	ms.mu.Lock()
	defer ms.mu.Unlock()
	objID := server.GenerateUUID()

	// Check if Object already exists -
	// codes.AlreadyExists
	_, err := ms.dbhandler.FindMovieByID(objID)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists,
			"Movie Record with ID: %v already exists!", objID)
	} else {
		mvObject := req.GetMovie()
		if mvObject.GetId() != "" {
			return nil, status.Errorf(codes.InvalidArgument,
				"Input is not valid! ID is auto-generated...")
		}
		mvObject.Id = objID

		// Validate format of Input and store the data
		if valid, err := ms.isValidMovie(mvObject); !valid {
			return nil, status.Errorf(codes.InvalidArgument,
				"Input is not valid! %v", err.Error())
		}

		// Assign server generated info.
		now := ptypes.TimestampNow()

		movieObject := persistence.Movie{
			Id:      objID,
			Name:    mvObject.Name,
			Summary: mvObject.Summary,
			Cast:    mvObject.Cast,
			// TODO: Add Tags here
			Director:   mvObject.Director,
			Writers:    mvObject.Writers,
			Active:     true,
			CreateTime: now,
			UpdateTime: now,
		}

		ms.dbhandler.AddMovie(movieObject)
	}

	log.Println("[DEBUG] End CreateMovieRequest!")
	return &moviepb.CreateMovieResponse{
		Id: objID,
	}, nil
}

func (ms *movieServer) UpdateMovie(ctx context.Context,
	req *moviepb.UpdateMovieRequest) (*moviepb.UpdateMovieResponse, error) {

	// log.Println("[DEBUG] Beginning UpdateMovieRequest: ", req)

	objID := req.GetId()
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Check if object already exists or not
	// codes.NotFound

	_, err := ms.dbhandler.FindMovieByID(objID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Movie Record with ID:%v does not exist!", objID)
	} else {

		// Validate and update the whole object
		mvObject := req.GetMovie()

		// Verify that the ID is not updated
		newID := mvObject.GetId()
		if newID != "" && !reflect.DeepEqual(newID, objID) {
			return nil, status.Errorf(codes.InvalidArgument, "Cannot update the ID of object!")
		}
		mvObject.Id = objID

		// Validate format of Input and store the data
		if valid, err := ms.isValidMovie(mvObject); !valid {
			return nil, status.Errorf(codes.InvalidArgument, "Input is not valid! %v", err.Error())
		}

		updatedMv := persistence.Movie{
			Name:    mvObject.Name,
			Summary: mvObject.Summary,
			Cast:    mvObject.Cast,
			// Tags: mvObject.Tags,
			Director:   mvObject.Director,
			Writers:    mvObject.Writers,
			UpdateTime: ptypes.TimestampNow(),
		}
		ms.dbhandler.UpdateMovieByID(objID, updatedMv)
	}

	log.Println("[DEBUG] End UpdateMovieRequest!")
	return &moviepb.UpdateMovieResponse{
		Id: objID,
	}, nil
}

// func (ms *movieServer) PartialUpdateMovie(ctx context.Context, req *moviepb.PartialUpdateMovieRequest) (*moviepb.PartialUpdateMovieResponse, error) {

// 	log.Println("[DEBUG] Beginning PartialUpdateMovieRequest: ", req)

// 	objID := req.GetId()

// 	// Check if object already exists or not
// 	// codes.NotFound
// 	if index, ok := ms.keys[objID]; ok && ms.Store[index].active {

// 		// Validate each field individually and Update the fields at a go afterwards
// 		ms.mu.Lock()
// 		defer ms.mu.Unlock()
// 		mvObject := ms.Store[index].movie

// 		// Verify that the ID is not updated
// 		newID := req.GetId()
// 		if newID != "" && newID != objID {
// 			return nil, status.Errorf(codes.InvalidArgument, "Cannot update the ID of object!")
// 		}

// 		if req.GetSummary() != "" {
// 			mvObject.Summary = req.GetSummary()
// 		}

// 		if req.GetDirector() != "" {
// 			mvObject.Director = req.GetDirector()
// 		}

// 		if req.GetCast() != nil {
// 			mvObject.Cast = req.GetCast()
// 		}

// 		if req.GetTags() != nil {
// 			mvObject.Tags = req.GetTags()
// 		}

// 		if req.GetWriters() != nil {
// 			mvObject.Writers = req.GetWriters()
// 		}

// 	} else {
// 		return nil, status.Errorf(codes.NotFound, "Movie Record with ID:%v does not exist!", objID)
// 	}

// 	log.Println("[DEBUG] End PartialUpdateMovieRequest!")
// 	return &moviepb.PartialUpdateMovieResponse{
// 		Id: objID,
// 	}, nil
// }

func (ms *movieServer) DeleteMovie(ctx context.Context,
	req *moviepb.DeleteMovieRequest) (*empty.Empty, error) {
	log.Println("[DEBUG] Beginning DeleteMovieRequest: ", req)

	objID := req.GetId()

	// Check if object already exists or not
	// codes.NotFound
	if _, err := ms.dbhandler.FindMovieByID(objID); err != nil {
		return nil, status.Errorf(
			codes.NotFound, "Movie Record with ID:%v does not exist!",
			objID)
	}

	if err := ms.dbhandler.RemoveMovieByID(objID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "some error while deleting movie!")
	}

	log.Println("[DEBUG] End DeleteMovieRequest!")
	return &empty.Empty{}, nil
}

func (ms *movieServer) isValidMovie(mv *moviepb.Movie) (bool, error) {

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

	// Verify no duplicate entries are made for Movie name
	for _, mvObject := range ms.Store {
		if !mvObject.active {
			continue
		}

		// Check whether Name is Unique or not. Also, ignore if update call is happening
		if mvObject.movie.GetName() == mv.GetName() && mv.GetId() == "" {
			return false, fmt.Errorf("Name should be unique!")
		}
	}

	return true, nil
}
