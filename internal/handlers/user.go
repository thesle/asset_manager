package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"assetManager/internal/auth"
	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// UserHandler handles user management endpoints
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetAll returns all users
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.repo.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch users"})
		return
	}
	if len(users) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetByID returns a user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	user, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Create creates a new user
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"Username"`
		Email    string `json:"Email"`
		Password string `json:"Password"`
		IsActive bool   `json:"IsActive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to hash password"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		IsActive:     req.IsActive,
	}

	if err := h.repo.Create(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Update updates a user
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	user.ID = id

	if err := h.repo.Update(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ResetPassword resets a user's password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var req struct {
		Password string `json:"Password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to hash password"})
		return
	}

	if err := h.repo.UpdatePassword(context.Background(), id, passwordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to reset password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Password reset successfully"})
}

// Delete deletes a user
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "User deleted"})
}
