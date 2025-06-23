package api

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordMatches(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), 12)
	match, err := (&Server{}).passwordMatches(string(hash), "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected passwords to match")
	}
	match, _ = (&Server{}).passwordMatches(string(hash), "wrong")
	if match {
		t.Fatal("expected passwords not to match")
	}
}
