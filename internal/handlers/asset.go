package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// AssetHandler handles asset endpoints
type AssetHandler struct {
	repo         *repository.AssetRepository
	propertyRepo *repository.AssetPropertyRepository
}

// NewAssetHandler creates a new asset handler
func NewAssetHandler(repo *repository.AssetRepository, propertyRepo *repository.AssetPropertyRepository) *AssetHandler {
	return &AssetHandler{
		repo:         repo,
		propertyRepo: propertyRepo,
	}
}

// GetAll returns all assets
func (h *AssetHandler) GetAll(c *gin.Context) {
	includeDeleted := c.Query("include_deleted") == "true"
	assets, err := h.repo.GetAll(context.Background(), includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch assets"})
		return
	}
	if len(assets) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assets)
}

// GetWithAssignments returns all assets with current assignment info
func (h *AssetHandler) GetWithAssignments(c *gin.Context) {
	includeDeleted := c.Query("include_deleted") == "true"
	assets, err := h.repo.GetWithCurrentAssignment(context.Background(), includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch assets"})
		return
	}
	if len(assets) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assets)
}

// GetByID returns an asset by ID
func (h *AssetHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	asset, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Asset not found"})
		return
	}
	c.JSON(http.StatusOK, asset)
}

// GetByAssetType returns assets by asset type
func (h *AssetHandler) GetByAssetType(c *gin.Context) {
	typeID, err := strconv.ParseInt(c.Param("typeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid type ID"})
		return
	}

	assets, err := h.repo.GetByAssetType(context.Background(), typeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch assets"})
		return
	}
	if len(assets) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assets)
}

// Search searches assets
func (h *AssetHandler) Search(c *gin.Context) {
	term := c.Query("q")
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Search term required"})
		return
	}

	assets, err := h.repo.Search(context.Background(), term)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to search assets"})
		return
	}
	if len(assets) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assets)
}

// Create creates a new asset
func (h *AssetHandler) Create(c *gin.Context) {
	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	if err := h.repo.Create(context.Background(), &asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create asset"})
		return
	}
	c.JSON(http.StatusCreated, asset)
}

// Update updates an asset
func (h *AssetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	asset.ID = id

	if err := h.repo.Update(context.Background(), &asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update asset"})
		return
	}
	c.JSON(http.StatusOK, asset)
}

// Delete deletes an asset
func (h *AssetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete asset"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Asset deleted"})
}

// GetProperties returns properties for an asset
func (h *AssetHandler) GetProperties(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	properties, err := h.propertyRepo.GetByAssetID(context.Background(), id)
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

// SetProperty sets a property value for an asset
func (h *AssetHandler) SetProperty(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var ap models.AssetProperty
	if err := c.ShouldBindJSON(&ap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	ap.AssetID = id

	if err := h.propertyRepo.Upsert(context.Background(), &ap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to set property"})
		return
	}
	c.JSON(http.StatusOK, ap)
}

// DeleteProperty deletes a property from an asset
func (h *AssetHandler) DeleteProperty(c *gin.Context) {
	propID, err := strconv.ParseInt(c.Param("propId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid property ID"})
		return
	}

	if err := h.propertyRepo.Delete(context.Background(), propID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete property"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Property deleted"})
}
