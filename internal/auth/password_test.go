package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	if hash == password {
		t.Error("Hash should not equal original password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Correct password
	if !CheckPassword(password, hash) {
		t.Error("Expected password check to succeed")
	}

	// Wrong password
	if CheckPassword("wrongpassword", hash) {
		t.Error("Expected password check to fail for wrong password")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "testpassword123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Bcrypt should generate different hashes for same password
	if hash1 == hash2 {
		t.Error("Expected different hashes for same password (bcrypt salt)")
	}

	// But both should validate
	if !CheckPassword(password, hash1) {
		t.Error("Expected password check to succeed for hash1")
	}
	if !CheckPassword(password, hash2) {
		t.Error("Expected password check to succeed for hash2")
	}
}
