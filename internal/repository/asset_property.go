package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrAssetPropertyNotFound = errors.New("asset property not found")

// AssetPropertyRepository handles asset property data operations
type AssetPropertyRepository struct {
	db *sqlx.DB
}

// NewAssetPropertyRepository creates a new asset property repository
func NewAssetPropertyRepository(db *sqlx.DB) *AssetPropertyRepository {
	return &AssetPropertyRepository{db: db}
}

// GetByID retrieves an asset property by ID
func (r *AssetPropertyRepository) GetByID(ctx context.Context, id int64) (*models.AssetProperty, error) {
	var ap models.AssetProperty
	query := `SELECT ap.id, ap.asset_id, ap.property_id, COALESCE(ap.value, '') as value, ap.created_at, ap.updated_at, ap.deleted_at,
			  COALESCE(p.name, '') as property_name, COALESCE(p.data_type, 'string') as data_type
			  FROM assets_properties ap
			  LEFT JOIN properties p ON ap.property_id = p.id
			  WHERE ap.id = ? AND ap.deleted_at IS NULL`
	err := r.db.GetContext(ctx, &ap, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetPropertyNotFound
	}
	return &ap, err
}

// GetByAssetID retrieves all properties for an asset
func (r *AssetPropertyRepository) GetByAssetID(ctx context.Context, assetID int64) ([]models.AssetProperty, error) {
	var aps []models.AssetProperty
	query := `SELECT ap.id, ap.asset_id, ap.property_id, COALESCE(ap.value, '') as value, ap.created_at, ap.updated_at, ap.deleted_at,
			  COALESCE(p.name, '') as property_name, COALESCE(p.data_type, 'string') as data_type
			  FROM assets_properties ap
			  LEFT JOIN properties p ON ap.property_id = p.id
			  WHERE ap.asset_id = ? AND ap.deleted_at IS NULL
			  ORDER BY p.name`
	err := r.db.SelectContext(ctx, &aps, query, assetID)
	return aps, err
}

// Create creates a new asset property
func (r *AssetPropertyRepository) Create(ctx context.Context, ap *models.AssetProperty) error {
	query := `INSERT INTO assets_properties (asset_id, property_id, value) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, ap.AssetID, ap.PropertyID, ap.Value)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ap.ID = id
	return nil
}

// Update updates an existing asset property
func (r *AssetPropertyRepository) Update(ctx context.Context, ap *models.AssetProperty) error {
	query := `UPDATE assets_properties SET value = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, ap.Value, ap.ID)
	return err
}

// Delete soft-deletes an asset property
func (r *AssetPropertyRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE assets_properties SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteByAssetID soft-deletes all properties for an asset
func (r *AssetPropertyRepository) DeleteByAssetID(ctx context.Context, assetID int64) error {
	query := `UPDATE assets_properties SET deleted_at = NOW() WHERE asset_id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, assetID)
	return err
}

// Upsert creates or updates an asset property
func (r *AssetPropertyRepository) Upsert(ctx context.Context, ap *models.AssetProperty) error {
	// Check if exists
	var existing models.AssetProperty
	query := `SELECT id FROM assets_properties WHERE asset_id = ? AND property_id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &existing, query, ap.AssetID, ap.PropertyID)
	if errors.Is(err, sql.ErrNoRows) {
		return r.Create(ctx, ap)
	}
	if err != nil {
		return err
	}
	ap.ID = existing.ID
	return r.Update(ctx, ap)
}
