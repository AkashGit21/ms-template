package server

import (
	"fmt"
	"time"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string
	Role     string
}

func NewJWTManager(sk string, td time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     sk,
		tokenDuration: td,
	}
}

type jwtTokenGenerator struct {
	token string
}

func (jm *JWTManager) GenerateToken(user *identitypb.User) (string, error) {

	// Create the Claims
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jm.tokenDuration).Unix(),
		},
		Username: user.Username,
		Role:     user.Role.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jm.secretKey)
}

func (jm *JWTManager) GetUserFromToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, fmt.Errorf("Unexpected Token signing method!")
			}
			return []byte(jm.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Invalid token! %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims!")
	}
	return claims, nil

}

// Verify verifies the access token string and return a user claim if the token is valid
func (jm *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(jm.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// TODO
// To hash the password for security reasons
func HashPassword(s string) string {
	return ""
}
