package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// PropertyHandler handles property endpoints
type PropertyHandler struct {
	repo *repository.PropertyRepository
}

// NewPropertyHandler creates a new property handler
func NewPropertyHandler(repo *repository.PropertyRepository) *PropertyHandler {
	return &PropertyHandler{repo: repo}
}

// GetAll returns all properties
func (h *PropertyHandler) GetAll(c *gin.Context) {
	properties, err := h.repo.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch properties"})
		return
	}
	if len(properties) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, properties)
}

// GetByID returns a property by ID
func (h *PropertyHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	property, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Property not found"})
		return
	}
	c.JSON(http.StatusOK, property)
}

// Create creates a new property
func (h *PropertyHandler) Create(c *gin.Context) {
	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.repo.Create(context.Background(), &property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create property: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, property)
}

// Update updates a property
func (h *PropertyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	property.ID = id

	if err := h.repo.Update(context.Background(), &property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update property"})
		return
	}
	c.JSON(http.StatusOK, property)
}

// Delete deletes a property
func (h *PropertyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete property"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Property deleted"})
}
