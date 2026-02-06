package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// AssetTypeHandler handles asset type endpoints
type AssetTypeHandler struct {
	repo *repository.AssetTypeRepository
}

// NewAssetTypeHandler creates a new asset type handler
func NewAssetTypeHandler(repo *repository.AssetTypeRepository) *AssetTypeHandler {
	return &AssetTypeHandler{repo: repo}
}

// GetAll returns all asset types
func (h *AssetTypeHandler) GetAll(c *gin.Context) {
	assetTypes, err := h.repo.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch asset types"})
		return
	}
	if len(assetTypes) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assetTypes)
}

// GetByID returns an asset type by ID
func (h *AssetTypeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	assetType, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Asset type not found"})
		return
	}
	c.JSON(http.StatusOK, assetType)
}

// Create creates a new asset type
func (h *AssetTypeHandler) Create(c *gin.Context) {
	var assetType models.AssetType
	if err := c.ShouldBindJSON(&assetType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	if err := h.repo.Create(context.Background(), &assetType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create asset type"})
		return
	}
	c.JSON(http.StatusCreated, assetType)
}

// Update updates an asset type
func (h *AssetTypeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var assetType models.AssetType
	if err := c.ShouldBindJSON(&assetType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	assetType.ID = id

	if err := h.repo.Update(context.Background(), &assetType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update asset type"})
		return
	}
	c.JSON(http.StatusOK, assetType)
}

// Delete deletes an asset type
func (h *AssetTypeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete asset type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Asset type deleted"})
}
