package auth

import (
	"testing"
	"time"

	"assetManager/internal/models"
)

func TestJWTService_GenerateAndValidate(t *testing.T) {
	secret := "test-secret-key-for-testing"
	service := NewJWTService(secret, 24)

	user := &models.User{
		BaseModel: models.BaseModel{ID: 1},
		Username:  "testuser",
	}

	// Generate token
	token, expiresAt, err := service.GenerateToken(user, false)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}

	if expiresAt <= time.Now().Unix() {
		t.Error("Expected expiry to be in the future")
	}

	// Validate token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %d, got %d", user.ID, claims.UserID)
	}

	if claims.Username != user.Username {
		t.Errorf("Expected Username %s, got %s", user.Username, claims.Username)
	}
}

func TestJWTService_InvalidToken(t *testing.T) {
	service := NewJWTService("test-secret", 24)

	_, err := service.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestJWTService_WrongSecret(t *testing.T) {
	service1 := NewJWTService("secret1", 24)
	service2 := NewJWTService("secret2", 24)

	user := &models.User{
		BaseModel: models.BaseModel{ID: 1},
		Username:  "testuser",
	}

	token, _, err := service1.GenerateToken(user, false)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with different secret
	_, err = service2.ValidateToken(token)
	if err == nil {
		t.Error("Expected error when validating with wrong secret")
	}
}

func TestJWTService_RememberMe(t *testing.T) {
	service := NewJWTService("test-secret", 24)

	user := &models.User{
		BaseModel: models.BaseModel{ID: 1},
		Username:  "testuser",
	}

	// Without remember
	_, expiresAt1, _ := service.GenerateToken(user, false)

	// With remember
	_, expiresAt2, _ := service.GenerateToken(user, true)

	// Remember should have longer expiry (30 days vs 24 hours)
	if expiresAt2 <= expiresAt1 {
		t.Error("Expected remember token to have longer expiry")
	}
}
