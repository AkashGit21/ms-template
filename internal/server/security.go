package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

	out, err := token.SignedString([]byte(jm.secretKey))
	return fmt.Sprintf("Basic %v", out), err
}

func (jm *JWTManager) GetUserFromToken(accessToken string) (*UserClaims, error) {
	accessToken = strings.TrimPrefix(accessToken, "Basic ")
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
	accessToken = strings.TrimPrefix(accessToken, "Basic ")
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

// To hash the password for security reasons
func HashPassword(s string) (string, error) {

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to hash password! %v", err)
	}
	return string(hashedPwd), nil
}

func DoesPasswordMatch(given string, need string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(need), []byte(given))
	if err != nil {
		log.Println("Password doesn't match! ", err)
		return false
	}
	return true
}
