package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrPersonNotFound = errors.New("person not found")

// PersonRepository handles person data operations
type PersonRepository struct {
	db *sqlx.DB
}

// NewPersonRepository creates a new person repository
func NewPersonRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

// GetByID retrieves a person by ID
func (r *PersonRepository) GetByID(ctx context.Context, id int64) (*models.Person, error) {
	var person models.Person
	query := `SELECT id, name, COALESCE(email, '') as email, COALESCE(phone, '') as phone, created_at, updated_at, deleted_at 
			  FROM persons WHERE id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &person, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPersonNotFound
	}
	return &person, err
}

// GetAll retrieves all persons
func (r *PersonRepository) GetAll(ctx context.Context) ([]models.Person, error) {
	var persons []models.Person
	query := `SELECT id, name, COALESCE(email, '') as email, COALESCE(phone, '') as phone, created_at, updated_at, deleted_at 
			  FROM persons WHERE deleted_at IS NULL ORDER BY name`
	err := r.db.SelectContext(ctx, &persons, query)
	return persons, err
}

// GetUnassigned retrieves the special 'Unassigned' person
func (r *PersonRepository) GetUnassigned(ctx context.Context) (*models.Person, error) {
	var person models.Person
	query := `SELECT id, name, COALESCE(email, '') as email, COALESCE(phone, '') as phone, created_at, updated_at, deleted_at 
			  FROM persons WHERE name = 'Unassigned' AND deleted_at IS NULL LIMIT 1`
	err := r.db.GetContext(ctx, &person, query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPersonNotFound
	}
	return &person, err
}

// Create creates a new person
func (r *PersonRepository) Create(ctx context.Context, person *models.Person) error {
	query := `INSERT INTO persons (name, email, phone) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, person.Name, person.Email, person.Phone)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	person.ID = id
	return nil
}

// Update updates an existing person
func (r *PersonRepository) Update(ctx context.Context, person *models.Person) error {
	query := `UPDATE persons SET name = ?, email = ?, phone = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, person.Name, person.Email, person.Phone, person.ID)
	return err
}

// Delete soft-deletes a person
func (r *PersonRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE persons SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Search searches persons by name or email
func (r *PersonRepository) Search(ctx context.Context, term string) ([]models.Person, error) {
	var persons []models.Person
	searchTerm := "%" + term + "%"
	query := `SELECT id, name, COALESCE(email, '') as email, COALESCE(phone, '') as phone, created_at, updated_at, deleted_at 
			  FROM persons WHERE deleted_at IS NULL 
			  AND (name LIKE ? OR email LIKE ?)
			  ORDER BY name`
	err := r.db.SelectContext(ctx, &persons, query, searchTerm, searchTerm)
	return persons, err
}
