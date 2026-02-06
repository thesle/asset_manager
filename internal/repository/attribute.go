package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrAttributeNotFound = errors.New("attribute not found")

// AttributeRepository handles attribute data operations
type AttributeRepository struct {
	db *sqlx.DB
}

// NewAttributeRepository creates a new attribute repository
func NewAttributeRepository(db *sqlx.DB) *AttributeRepository {
	return &AttributeRepository{db: db}
}

// GetByID retrieves an attribute by ID
func (r *AttributeRepository) GetByID(ctx context.Context, id int64) (*models.Attribute, error) {
	var attribute models.Attribute
	query := `SELECT id, name, data_type, COALESCE(enum_options, '') as enum_options, created_at, updated_at, deleted_at 
			  FROM attributes WHERE id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &attribute, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAttributeNotFound
	}
	return &attribute, err
}

// GetAll retrieves all attributes
func (r *AttributeRepository) GetAll(ctx context.Context) ([]models.Attribute, error) {
	var attributes []models.Attribute
	query := `SELECT id, name, data_type, COALESCE(enum_options, '') as enum_options, created_at, updated_at, deleted_at 
			  FROM attributes WHERE deleted_at IS NULL ORDER BY name`
	err := r.db.SelectContext(ctx, &attributes, query)
	return attributes, err
}

// Create creates a new attribute
func (r *AttributeRepository) Create(ctx context.Context, attribute *models.Attribute) error {
	// Handle empty enum_options - MySQL JSON column needs NULL or valid JSON
	var enumOptions interface{}
	if attribute.EnumOptions == "" {
		enumOptions = nil
	} else {
		enumOptions = attribute.EnumOptions
	}

	query := `INSERT INTO attributes (name, data_type, enum_options) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, attribute.Name, attribute.DataType, enumOptions)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	attribute.ID = id
	return nil
}

// Update updates an existing attribute
func (r *AttributeRepository) Update(ctx context.Context, attribute *models.Attribute) error {
	// Handle empty enum_options - MySQL JSON column needs NULL or valid JSON
	var enumOptions interface{}
	if attribute.EnumOptions == "" {
		enumOptions = nil
	} else {
		enumOptions = attribute.EnumOptions
	}

	query := `UPDATE attributes SET name = ?, data_type = ?, enum_options = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, attribute.Name, attribute.DataType, enumOptions, attribute.ID)
	return err
}

// Delete soft-deletes an attribute
func (r *AttributeRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE attributes SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
