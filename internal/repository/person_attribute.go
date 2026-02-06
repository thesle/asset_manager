package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrPersonAttributeNotFound = errors.New("person attribute not found")

// PersonAttributeRepository handles person attribute data operations
type PersonAttributeRepository struct {
	db *sqlx.DB
}

// NewPersonAttributeRepository creates a new person attribute repository
func NewPersonAttributeRepository(db *sqlx.DB) *PersonAttributeRepository {
	return &PersonAttributeRepository{db: db}
}

// GetByID retrieves a person attribute by ID
func (r *PersonAttributeRepository) GetByID(ctx context.Context, id int64) (*models.PersonAttribute, error) {
	var pa models.PersonAttribute
	query := `SELECT pa.id, pa.person_id, pa.attribute_id, COALESCE(pa.value, '') as value, pa.created_at, pa.updated_at, pa.deleted_at,
			  COALESCE(a.name, '') as attribute_name, COALESCE(a.data_type, 'string') as data_type
			  FROM persons_attributes pa
			  LEFT JOIN attributes a ON pa.attribute_id = a.id
			  WHERE pa.id = ? AND pa.deleted_at IS NULL`
	err := r.db.GetContext(ctx, &pa, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPersonAttributeNotFound
	}
	return &pa, err
}

// GetByPersonID retrieves all attributes for a person
func (r *PersonAttributeRepository) GetByPersonID(ctx context.Context, personID int64) ([]models.PersonAttribute, error) {
	var pas []models.PersonAttribute
	query := `SELECT pa.id, pa.person_id, pa.attribute_id, COALESCE(pa.value, '') as value, pa.created_at, pa.updated_at, pa.deleted_at,
			  COALESCE(a.name, '') as attribute_name, COALESCE(a.data_type, 'string') as data_type
			  FROM persons_attributes pa
			  LEFT JOIN attributes a ON pa.attribute_id = a.id
			  WHERE pa.person_id = ? AND pa.deleted_at IS NULL
			  ORDER BY a.name`
	err := r.db.SelectContext(ctx, &pas, query, personID)
	return pas, err
}

// Create creates a new person attribute
func (r *PersonAttributeRepository) Create(ctx context.Context, pa *models.PersonAttribute) error {
	query := `INSERT INTO persons_attributes (person_id, attribute_id, value) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, pa.PersonID, pa.AttributeID, pa.Value)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	pa.ID = id
	return nil
}

// Update updates an existing person attribute
func (r *PersonAttributeRepository) Update(ctx context.Context, pa *models.PersonAttribute) error {
	query := `UPDATE persons_attributes SET value = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, pa.Value, pa.ID)
	return err
}

// Delete soft-deletes a person attribute
func (r *PersonAttributeRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE persons_attributes SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteByPersonID soft-deletes all attributes for a person
func (r *PersonAttributeRepository) DeleteByPersonID(ctx context.Context, personID int64) error {
	query := `UPDATE persons_attributes SET deleted_at = NOW() WHERE person_id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, personID)
	return err
}

// Upsert creates or updates a person attribute
func (r *PersonAttributeRepository) Upsert(ctx context.Context, pa *models.PersonAttribute) error {
	var existing models.PersonAttribute
	query := `SELECT id FROM persons_attributes WHERE person_id = ? AND attribute_id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &existing, query, pa.PersonID, pa.AttributeID)
	if errors.Is(err, sql.ErrNoRows) {
		return r.Create(ctx, pa)
	}
	if err != nil {
		return err
	}
	pa.ID = existing.ID
	return r.Update(ctx, pa)
}
