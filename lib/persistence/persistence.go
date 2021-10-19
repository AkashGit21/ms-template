package persistence

import (
	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
)

type DatabaseHandler interface {
	AddUser(User) ([]byte, error)
	FindByUsername(string) (User, error)
	FindAllUsers(int, int32) ([]*identitypb.User, error)
	RemoveByUsername(string) error
	CountUsers() int

	Authenticate(string, string) bool

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
