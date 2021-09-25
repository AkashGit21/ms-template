package server

import (
	"github.com/google/uuid"
)

// Generate a Unique UUID string for the object
func GenerateUUID() string {
	uniqueID := uuid.New()

	return uniqueID.String()
}
