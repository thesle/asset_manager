package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// AttributeHandler handles attribute endpoints
type AttributeHandler struct {
	repo *repository.AttributeRepository
}

// NewAttributeHandler creates a new attribute handler
func NewAttributeHandler(repo *repository.AttributeRepository) *AttributeHandler {
	return &AttributeHandler{repo: repo}
}

// GetAll returns all attributes
func (h *AttributeHandler) GetAll(c *gin.Context) {
	attributes, err := h.repo.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch attributes"})
		return
	}
	if len(attributes) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, attributes)
}

// GetByID returns an attribute by ID
func (h *AttributeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	attribute, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Attribute not found"})
		return
	}
	c.JSON(http.StatusOK, attribute)
}

// Create creates a new attribute
func (h *AttributeHandler) Create(c *gin.Context) {
	var attribute models.Attribute
	if err := c.ShouldBindJSON(&attribute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	if err := h.repo.Create(context.Background(), &attribute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create attribute"})
		return
	}
	c.JSON(http.StatusCreated, attribute)
}

// Update updates an attribute
func (h *AttributeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var attribute models.Attribute
	if err := c.ShouldBindJSON(&attribute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	attribute.ID = id

	if err := h.repo.Update(context.Background(), &attribute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update attribute"})
		return
	}
	c.JSON(http.StatusOK, attribute)
}

// Delete deletes an attribute
func (h *AttributeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete attribute"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Attribute deleted"})
}
