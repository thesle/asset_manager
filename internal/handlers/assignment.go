package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// AssignmentHandler handles asset assignment endpoints
type AssignmentHandler struct {
	repo       *repository.AssetAssignmentRepository
	personRepo *repository.PersonRepository
}

// NewAssignmentHandler creates a new assignment handler
func NewAssignmentHandler(repo *repository.AssetAssignmentRepository, personRepo *repository.PersonRepository) *AssignmentHandler {
	return &AssignmentHandler{
		repo:       repo,
		personRepo: personRepo,
	}
}

// GetByAssetID returns assignment history for an asset
func (h *AssignmentHandler) GetByAssetID(c *gin.Context) {
	assetID, err := strconv.ParseInt(c.Param("assetId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid asset ID"})
		return
	}

	assignments, err := h.repo.GetHistoryByAssetID(context.Background(), assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch assignments"})
		return
	}
	if len(assignments) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assignments)
}

// GetCurrentByAssetID returns the current assignment for an asset
func (h *AssignmentHandler) GetCurrentByAssetID(c *gin.Context) {
	assetID, err := strconv.ParseInt(c.Param("assetId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid asset ID"})
		return
	}

	assignment, err := h.repo.GetCurrentByAssetID(context.Background(), assetID)
	if err != nil {
		if err == repository.ErrAssetAssignmentNotFound {
			c.JSON(http.StatusOK, nil)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch assignment"})
		return
	}
	c.JSON(http.StatusOK, assignment)
}

// GetByPersonID returns assignments for a person
func (h *AssignmentHandler) GetByPersonID(c *gin.Context) {
	personID, err := strconv.ParseInt(c.Param("personId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid person ID"})
		return
	}

	assignments, err := h.repo.GetByPersonID(context.Background(), personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch assignments"})
		return
	}
	if len(assignments) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assignments)
}

// GetCurrentByPersonID returns current assignments for a person
func (h *AssignmentHandler) GetCurrentByPersonID(c *gin.Context) {
	personID, err := strconv.ParseInt(c.Param("personId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid person ID"})
		return
	}

	assignments, err := h.repo.GetCurrentByPersonID(context.Background(), personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch assignments"})
		return
	}
	if len(assignments) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, assignments)
}

// AssignAsset assigns an asset to a person
func (h *AssignmentHandler) AssignAsset(c *gin.Context) {
	var req struct {
		AssetID       int64      `json:"AssetID"`
		PersonID      int64      `json:"PersonID"`
		Notes         string     `json:"Notes"`
		EffectiveDate *time.Time `json:"EffectiveDate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	// Default to now if no date provided
	effectiveDate := time.Now()
	if req.EffectiveDate != nil {
		effectiveDate = *req.EffectiveDate
	}

	if err := h.repo.AssignAsset(context.Background(), req.AssetID, req.PersonID, req.Notes, effectiveDate); err != nil {
		if err == repository.ErrOverlappingAssignment {
			c.JSON(http.StatusConflict, gin.H{"Error": "Overlapping assignment exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to assign asset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Asset assigned successfully"})
}

// UnassignAsset unassigns an asset (assigns to 'Unassigned' person)
func (h *AssignmentHandler) UnassignAsset(c *gin.Context) {
	assetID, err := strconv.ParseInt(c.Param("assetId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid asset ID"})
		return
	}

	var req struct {
		EffectiveDate string `json:"EffectiveDate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	// Parse the effective date
	effectiveDate := time.Now()
	if req.EffectiveDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EffectiveDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid date format"})
			return
		}
		effectiveDate = parsed
	}

	// Get the 'Unassigned' person
	unassigned, err := h.personRepo.GetUnassigned(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to find Unassigned person"})
		return
	}

	if err := h.repo.AssignAsset(context.Background(), assetID, unassigned.ID, "Unassigned", effectiveDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to unassign asset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Asset unassigned successfully"})
}

// EndAssignment ends an assignment
func (h *AssignmentHandler) EndAssignment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var req struct {
		EndDate *time.Time `json:"EndDate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	endDate := time.Now()
	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	if err := h.repo.EndAssignment(context.Background(), id, endDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to end assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Assignment ended"})
}

// Create creates a new assignment with custom dates
func (h *AssignmentHandler) Create(c *gin.Context) {
	var aa models.AssetAssignment
	if err := c.ShouldBindJSON(&aa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	if err := h.repo.Create(context.Background(), &aa); err != nil {
		if err == repository.ErrOverlappingAssignment {
			c.JSON(http.StatusConflict, gin.H{"Error": "Overlapping assignment exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusCreated, aa)
}

// Update updates an assignment
func (h *AssignmentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var aa models.AssetAssignment
	if err := c.ShouldBindJSON(&aa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	aa.ID = id

	if err := h.repo.Update(context.Background(), &aa); err != nil {
		if err == repository.ErrOverlappingAssignment {
			c.JSON(http.StatusConflict, gin.H{"Error": "Overlapping assignment exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update assignment"})
		return
	}

	c.JSON(http.StatusOK, aa)
}

// Delete deletes an assignment
func (h *AssignmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Assignment deleted"})
}
