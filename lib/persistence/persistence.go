package persistence

import (
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	moviepb "github.com/AkashGit21/ms-project/internal/grpc/movie"
)

type DatabaseHandler interface {
	AddUser(User) ([]byte, error)
	FindByUsername(string) (User, error)
	FindAllUsers(int, int32) ([]*identitypb.User, error)
	RemoveByUsername(string) error
	CountUsers() int

	Authenticate(string, string) bool

	AddMovie(Movie) ([]byte, error)
	FindMovieByID(string) (Movie, error)
	FindAllMovies(int, int32) ([]*moviepb.Movie, error)
	UpdateMovieByID(string, Movie) ([]byte, error)
	RemoveMovieByID(string) error
	CountMovieRecords() int

	// AddEvent(Event) ([]byte, error)
	// AddBookingForUser([]byte, Booking) error
	// AddLocation(Location) (Location, error)
	// FindBookingsForUser([]byte) ([]Booking, error)
	// FindEvent([]byte) (Event, error)
	// FindEventByName(string) (Event, error)
	// FindAllAvailableEvents() ([]Event, error)
	// FindLocation(string) (Location, error)
	// FindAllLocations() ([]Location, error)
}
