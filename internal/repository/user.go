package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepository handles user data operations
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at, deleted_at 
			  FROM users WHERE id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at, deleted_at 
			  FROM users WHERE username = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &user, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at, deleted_at 
			  FROM users WHERE deleted_at IS NULL ORDER BY username`
	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, is_active) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.IsActive)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username = ?, email = ?, is_active = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	return err
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, passwordHash, id)
	return err
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
