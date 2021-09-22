package server

import (
	"encoding/base64"
	"testing"
)

func Test_tokenGenerator_ForIndex(t *testing.T) {
	salt := "salt"
	index := 1
	want := base64.StdEncoding.EncodeToString(
		[]byte("salt1"))
	tok := TokenGeneratorWithSalt(salt)
	if got := tok.ForIndex(index); got != want {
		t.Errorf("tokenGenerator.ForIndex() = %v, want %v", got, want)
	}
}

func Test_tokenGenerator_GetIndex_notParseable(t *testing.T) {
	tok := NewTokenGenerator()
	_, err := tok.GetIndex("invalid")
	if err == nil {
		t.Error("GetIndex: want error for invalid token.")
	}
}
func Test_tokenGenerator_GetIndex_noSalt(t *testing.T) {
	tok := NewTokenGenerator()
	_, err := tok.GetIndex(base64.StdEncoding.EncodeToString([]byte("invalid")))
	if err == nil {
		t.Error("GetIndex: want error for invalid token.")
	}
}

func Test_tokenGenerator_GetIndex_invalidIndex(t *testing.T) {
	tok := TokenGeneratorWithSalt("salt")
	_, err := tok.GetIndex(base64.StdEncoding.EncodeToString([]byte("saltinvalid")))
	if err == nil {
		t.Error("GetIndex: want error for invalid token.")
	}
}

func Test_tokenGenerator_GetIndex(t *testing.T) {
	tok := TokenGeneratorWithSalt("salt")
	i, err := tok.GetIndex(base64.StdEncoding.EncodeToString([]byte("salt1")))
	if err != nil {
		t.Error("GetIndex: unexpected err")
	}
	if i != 1 {
		t.Errorf("GetIndex: want 1, got %d", i)
	}
}
