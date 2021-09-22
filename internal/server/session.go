package server

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type JWTTokenGenerator interface {
	generateToken()
	validateToken()
	getUserFromToken()
}

func NewJWTTokenGenerator() JWTTokenGenerator {
	return &jwtTokenGenerator{
		token: strconv.FormatInt(time.Now().UnixNano(), 16),
	}
}

type jwtTokenGenerator struct {
	token string
}

func (tg *jwtTokenGenerator) generateToken() {

}

func (tg *jwtTokenGenerator) validateToken() {

}

func (tg *jwtTokenGenerator) getUserFromToken() {

}

// Generate a Unique UUID string for the object
func GenerateUUID() string {
	uniqueID := uuid.New()

	return uniqueID.String()
}
