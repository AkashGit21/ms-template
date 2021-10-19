package persistence

import "github.com/golang/protobuf/ptypes/timestamp"

// For Identity service
type User struct {
	// Required. The username of the user. Must be unique and length should be between 6 to 30 characters.
	Username string `bson:"username,omitempty"`
	// Required. The email address of the user. Must be unique
	Email string `bson:"email,omitempty"`
	// Required. The encoded password of the user
	Password string `bson:"password,omitempty"`
	// Role of the user ,i.e. Guest, NORMAL, SUBSCRIBED, ADMIN. Default role is Guest.
	Role Role `bson:"role,omitempty"`
	// Status of the user - Active/Inactive
	Active bool `bson:"Active,omitempty"`
	// The first name of user. For example: 'Harry'
	FirstName string `bson:"first_name,omitempty"`
	// The last name of user. For example: 'Potter'
	LastName *string `bson:"last_name,omitempty"`
	// Output only. The timestamp at which the user was created.
	CreateTime *timestamp.Timestamp `bson:"create_time,omitempty"`
	// Output only. The latest timestamp at which the user was updated.
	UpdateTime *timestamp.Timestamp `bson:"update_time,omitempty"`
	// The age of the user in years.
	Age *int32 `bson:"age,omitempty"`
	// The height of the user in feet.
	HeightInCms *float64 `bson:"height_in_cms,omitempty"`
	// The nickname of the user.
	Nickname *string `bson:"nickname,omitempty"`
	// Enables the receiving of notifications. The default is false if unset.
	EnableNotifications *bool `bson:"enable_notifications,omitempty"`
}

// The roles available for users
type Role int32

const (
	// Every User has this role by default
	Role_GUEST Role = 0
	// Logged in user but with no special fees
	Role_NORMAL Role = 1
	// Logged in user who has subscribed to service
	Role_SUBSCRIBED Role = 2
	// For Developers or Maintainers of the service
	Role_ADMIN Role = 3
)

// Enum value maps for Role.
var (
	Role_name = map[int32]string{
		0: "GUEST",
		1: "NORMAL",
		2: "SUBSCRIBED",
		3: "ADMIN",
	}
	Role_value = map[string]int32{
		"GUEST":      0,
		"NORMAL":     1,
		"SUBSCRIBED": 2,
		"ADMIN":      3,
	}
)

// // For Movies service
// type Movie struct {
// }
