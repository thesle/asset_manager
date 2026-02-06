package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"assetManager/internal/auth"
	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *repository.UserRepository, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	user, err := h.userRepo.GetByUsername(context.Background(), req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "User account is disabled"})
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid credentials"})
		return
	}

	token, expiresAt, err := h.jwtService.GenerateToken(user, req.Remember)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt64("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		return
	}

	var req struct {
		CurrentPassword string `json:"CurrentPassword"`
		NewPassword     string `json:"NewPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	user, err := h.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
		return
	}

	if !auth.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Current password is incorrect"})
		return
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to hash password"})
		return
	}

	if err := h.userRepo.UpdatePassword(context.Background(), userID, newHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Password updated successfully"})
}

// Me returns the current user info
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt64("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		return
	}

	user, err := h.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
