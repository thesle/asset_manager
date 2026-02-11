package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"assetManager/internal/repository"
)

type ReportHandler struct {
	repo *repository.ReportRepository
}

func NewReportHandler(repo *repository.ReportRepository) *ReportHandler {
	return &ReportHandler{
		repo: repo,
	}
}

func (h *ReportHandler) ExecuteCustomReport(c *gin.Context) {
	var req repository.CustomReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	ctx := context.Background()
	var results []map[string]interface{}
	var err error

	switch req.EntityType {
	case "asset":
		results, err = h.repo.ExecuteAssetReport(ctx, req.Filters)
	case "person":
		results, err = h.repo.ExecutePersonReport(ctx, req.Filters)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid entity type. Must be 'asset' or 'person'"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
