package encryption

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	e := Encryption{Key: []byte("0123456789abcdef")}
	original := "hello world"

	cipher, err := e.Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt returned error: %v", err)
	}

	plain, err := e.Decrypt(cipher)
	if err != nil {
		t.Fatalf("Decrypt returned error: %v", err)
	}

	if plain != original {
		t.Errorf("expected %s, got %s", original, plain)
	}
}
