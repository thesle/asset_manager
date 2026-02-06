package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrAssetTypeNotFound = errors.New("asset type not found")

// AssetTypeRepository handles asset type data operations
type AssetTypeRepository struct {
	db *sqlx.DB
}

// NewAssetTypeRepository creates a new asset type repository
func NewAssetTypeRepository(db *sqlx.DB) *AssetTypeRepository {
	return &AssetTypeRepository{db: db}
}

// GetByID retrieves an asset type by ID
func (r *AssetTypeRepository) GetByID(ctx context.Context, id int64) (*models.AssetType, error) {
	var assetType models.AssetType
	query := `SELECT id, name, COALESCE(description, '') as description, created_at, updated_at, deleted_at 
			  FROM asset_types WHERE id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &assetType, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetTypeNotFound
	}
	return &assetType, err
}

// GetAll retrieves all asset types
func (r *AssetTypeRepository) GetAll(ctx context.Context) ([]models.AssetType, error) {
	var assetTypes []models.AssetType
	query := `SELECT id, name, COALESCE(description, '') as description, created_at, updated_at, deleted_at 
			  FROM asset_types WHERE deleted_at IS NULL ORDER BY name`
	err := r.db.SelectContext(ctx, &assetTypes, query)
	return assetTypes, err
}

// Create creates a new asset type
func (r *AssetTypeRepository) Create(ctx context.Context, assetType *models.AssetType) error {
	query := `INSERT INTO asset_types (name, description) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, assetType.Name, assetType.Description)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	assetType.ID = id
	return nil
}

// Update updates an existing asset type
func (r *AssetTypeRepository) Update(ctx context.Context, assetType *models.AssetType) error {
	query := `UPDATE asset_types SET name = ?, description = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, assetType.Name, assetType.Description, assetType.ID)
	return err
}

// Delete soft-deletes an asset type
func (r *AssetTypeRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE asset_types SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
