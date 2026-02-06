package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrPropertyNotFound = errors.New("property not found")

// PropertyRepository handles property data operations
type PropertyRepository struct {
	db *sqlx.DB
}

// NewPropertyRepository creates a new property repository
func NewPropertyRepository(db *sqlx.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

// GetByID retrieves a property by ID
func (r *PropertyRepository) GetByID(ctx context.Context, id int64) (*models.Property, error) {
	var property models.Property
	query := `SELECT id, name, data_type, COALESCE(enum_options, '') as enum_options, created_at, updated_at, deleted_at 
			  FROM properties WHERE id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &property, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPropertyNotFound
	}
	return &property, err
}

// GetAll retrieves all properties
func (r *PropertyRepository) GetAll(ctx context.Context) ([]models.Property, error) {
	var properties []models.Property
	query := `SELECT id, name, data_type, COALESCE(enum_options, '') as enum_options, created_at, updated_at, deleted_at 
			  FROM properties WHERE deleted_at IS NULL ORDER BY name`
	err := r.db.SelectContext(ctx, &properties, query)
	return properties, err
}

// Create creates a new property
func (r *PropertyRepository) Create(ctx context.Context, property *models.Property) error {
	// Handle empty enum_options - MySQL JSON column needs NULL or valid JSON
	var enumOptions interface{}
	if property.EnumOptions == "" {
		enumOptions = nil
	} else {
		enumOptions = property.EnumOptions
	}

	query := `INSERT INTO properties (name, data_type, enum_options) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, property.Name, property.DataType, enumOptions)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	property.ID = id
	return nil
}

// Update updates an existing property
func (r *PropertyRepository) Update(ctx context.Context, property *models.Property) error {
	// Handle empty enum_options - MySQL JSON column needs NULL or valid JSON
	var enumOptions interface{}
	if property.EnumOptions == "" {
		enumOptions = nil
	} else {
		enumOptions = property.EnumOptions
	}

	query := `UPDATE properties SET name = ?, data_type = ?, enum_options = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, property.Name, property.DataType, enumOptions, property.ID)
	return err
}

// Delete soft-deletes a property
func (r *PropertyRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE properties SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
