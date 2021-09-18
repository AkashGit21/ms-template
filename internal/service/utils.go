package service

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// isAllowedSummary checks whether the given string is a valid Summary or not, i.e -
// length: [8,1200]
func isAllowedSummary(s string) error {
	if !StringLenBetween(s, 8, 1200) {
		return fmt.Errorf("The summary should be between 8 and 1200 characters.")
	}
	return nil
}

func isValidName(s string) bool {
	return StringLenBetween(s, 1, 120)
}

// StringLenBetween is a fxn which tests if the provided value is of type string and has
// length (excluding whitespace characters at the sides) between min and max (inclusive)
func StringLenBetween(s interface{}, min, max int) bool {
	if v, ok := s.(string); ok {
		v = strings.TrimSpace(v)
		length := len(v)

		return length >= min && length <= max
	}
	return false
}

// Generate a Unique UUID string for the object
func generateUUID() string {
	uniqueID := uuid.New()

	return uniqueID.String()
}
