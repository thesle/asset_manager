package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"assetManager/internal/models"
	"assetManager/internal/repository"
)

// PersonHandler handles person endpoints
type PersonHandler struct {
	repo          *repository.PersonRepository
	attributeRepo *repository.PersonAttributeRepository
}

// NewPersonHandler creates a new person handler
func NewPersonHandler(repo *repository.PersonRepository, attributeRepo *repository.PersonAttributeRepository) *PersonHandler {
	return &PersonHandler{
		repo:          repo,
		attributeRepo: attributeRepo,
	}
}

// GetAll returns all persons
func (h *PersonHandler) GetAll(c *gin.Context) {
	persons, err := h.repo.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch persons"})
		return
	}
	if len(persons) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, persons)
}

// GetByID returns a person by ID
func (h *PersonHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	person, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Person not found"})
		return
	}
	c.JSON(http.StatusOK, person)
}

// Search searches persons
func (h *PersonHandler) Search(c *gin.Context) {
	term := c.Query("q")
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Search term required"})
		return
	}

	persons, err := h.repo.Search(context.Background(), term)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to search persons"})
		return
	}
	if len(persons) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, persons)
}

// Create creates a new person
func (h *PersonHandler) Create(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	if err := h.repo.Create(context.Background(), &person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create person"})
		return
	}
	c.JSON(http.StatusCreated, person)
}

// Update updates a person
func (h *PersonHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	person.ID = id

	if err := h.repo.Update(context.Background(), &person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update person"})
		return
	}
	c.JSON(http.StatusOK, person)
}

// Delete deletes a person
func (h *PersonHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete person"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Person deleted"})
}

// GetAttributes returns attributes for a person
func (h *PersonHandler) GetAttributes(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	attributes, err := h.attributeRepo.GetByPersonID(context.Background(), id)
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

// SetAttribute sets an attribute value for a person
func (h *PersonHandler) SetAttribute(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}

	var pa models.PersonAttribute
	if err := c.ShouldBindJSON(&pa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}
	pa.PersonID = id

	if err := h.attributeRepo.Upsert(context.Background(), &pa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to set attribute"})
		return
	}
	c.JSON(http.StatusOK, pa)
}

// DeleteAttribute deletes an attribute from a person
func (h *PersonHandler) DeleteAttribute(c *gin.Context) {
	attrID, err := strconv.ParseInt(c.Param("attrId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid attribute ID"})
		return
	}

	if err := h.attributeRepo.Delete(context.Background(), attrID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete attribute"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Attribute deleted"})
}
