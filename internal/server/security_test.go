package server

import (
	"fmt"
	"log"
	"testing"
	"time"

	identitypb "github.com/AkashGit21/ms-project/internal/grpc/identity"
	"github.com/golang-jwt/jwt"
)

var (
	SecretKey = "secret"
)

// TODO: RSA token generation and parsing
func TestGetUserFromToken_BadSigning(t *testing.T) {

	jm := NewJWTManager(SecretKey, 2*time.Minute)
	u := &identitypb.User{Username: "usrname1", Role: identitypb.Role_NORMAL}

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jm.tokenDuration).Unix(),
		},
		Username: u.Username,
		Role:     u.Role.String(),
	}

	out := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err := out.SignedString([]byte(jm.secretKey))
	if err != nil {
		t.Error("token is not generating!")
	}
	token = "Basic " + token

	user, err := jm.GetUserFromToken(token)
	if err != fmt.Errorf("Unexpected Token signing method!") {
		t.Error("expecting invalid signing method here")
		log.Print("Error: ", err)
	}
	if user != nil {
		t.Error("UserClaims should be nil!")
	}
}

func TestGetUserFromToken_BadToken(t *testing.T) {

}

func TestGetUserFromToken_ValidToken(t *testing.T) {

	jm := NewJWTManager(SecretKey, 2*time.Minute)
	u := &identitypb.User{Username: "usrname1", Role: identitypb.Role_NORMAL}

	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jm.tokenDuration).Unix(),
		},
		Username: u.Username,
		Role:     u.Role.String(),
	}

	out := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := out.SignedString([]byte(jm.secretKey))
	if err != nil {
		t.Error("token is not generating!")
	}
	token = "Basic " + token

	user, err := jm.GetUserFromToken(token)
	if err != nil {
		t.Error("expecting no error")
		log.Print("Error: ", err)
	}
	if user == nil {
		t.Error("UserClaims should return something!")
	}
}
