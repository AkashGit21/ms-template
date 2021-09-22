package server

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewTokenGenerator provides a new instance of a TokenGenerator.
func NewTokenGenerator() TokenGenerator {
	return &tokenGenerator{salt: strconv.FormatInt(time.Now().Unix(), 10)}
}

// TokenGeneratorWithSalt provieds an instance of a TokenGenerator which
// uses the given salt.
func TokenGeneratorWithSalt(salt string) TokenGenerator {
	return &tokenGenerator{salt}
}

// TokenGenerator generates a page token for a given index.
type TokenGenerator interface {
	ForIndex(int) string
	GetIndex(string) (int, error)
}

// InvalidTokenErr is the error returned if the token provided is not
// parseable by the TokenGenerator.
var InvalidTokenErr = status.Errorf(
	codes.InvalidArgument,
	"The field `page_token` is invalid.")

type tokenGenerator struct {
	salt string
}

func (t *tokenGenerator) ForIndex(i int) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s%d", t.salt, i)))
}

func (t *tokenGenerator) GetIndex(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	bs, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return -1, InvalidTokenErr
	}

	if !strings.HasPrefix(string(bs), t.salt) {
		return -1, InvalidTokenErr
	}

	i, err := strconv.Atoi(strings.TrimPrefix(string(bs), t.salt))
	if err != nil {
		return -1, InvalidTokenErr
	}
	return i, nil
}
