package urlsigner

import "testing"

func TestGenerateAndVerifyToken(t *testing.T) {
	s := Signer{Secret: []byte("secretkey")}
	token := s.GenerateTokenFromString("http://example.com")
	if token == "" {
		t.Fatal("token should not be empty")
	}
	if !s.VerifyToken(token) {
		t.Fatal("token should be valid")
	}
	if s.Expired(token, 1) {
		t.Fatal("token should not be expired immediately")
	}
}
