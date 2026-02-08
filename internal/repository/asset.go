package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var ErrAssetNotFound = errors.New("asset not found")

// AssetRepository handles asset data operations
type AssetRepository struct {
	db *sqlx.DB
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(db *sqlx.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

// GetByID retrieves an asset by ID
func (r *AssetRepository) GetByID(ctx context.Context, id int64) (*models.Asset, error) {
	var asset models.Asset
	query := `SELECT a.id, a.asset_type_id, a.name, 
			  COALESCE(a.model, '') as model, 
			  COALESCE(a.serial_number, '') as serial_number, 
			  COALESCE(a.order_no, '') as order_no, 
			  COALESCE(a.license_number, '') as license_number, 
			  COALESCE(a.notes, '') as notes, 
			  a.created_at, a.updated_at, a.deleted_at,
			  COALESCE(at.name, '') as asset_type_name
			  FROM assets a
			  LEFT JOIN asset_types at ON a.asset_type_id = at.id
			  WHERE a.id = ? AND a.deleted_at IS NULL`
	err := r.db.GetContext(ctx, &asset, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetNotFound
	}
	return &asset, err
}

// GetAll retrieves all assets. If includeDeleted is true, returns only soft-deleted records.
func (r *AssetRepository) GetAll(ctx context.Context, includeDeleted bool) ([]models.Asset, error) {
	var assets []models.Asset
	deletedFilter := "a.deleted_at IS NULL"
	if includeDeleted {
		deletedFilter = "a.deleted_at IS NOT NULL"
	}
	query := `SELECT a.id, a.asset_type_id, a.name, 
			  COALESCE(a.model, '') as model, 
			  COALESCE(a.serial_number, '') as serial_number, 
			  COALESCE(a.order_no, '') as order_no, 
			  COALESCE(a.license_number, '') as license_number, 
			  COALESCE(a.notes, '') as notes, 
			  a.created_at, a.updated_at, a.deleted_at,
			  COALESCE(at.name, '') as asset_type_name
			  FROM assets a
			  LEFT JOIN asset_types at ON a.asset_type_id = at.id
			  WHERE ` + deletedFilter + ` ORDER BY a.name`
	err := r.db.SelectContext(ctx, &assets, query)
	return assets, err
}

// GetByAssetType retrieves all assets of a specific type
func (r *AssetRepository) GetByAssetType(ctx context.Context, assetTypeID int64) ([]models.Asset, error) {
	var assets []models.Asset
	query := `SELECT a.id, a.asset_type_id, a.name, 
			  COALESCE(a.model, '') as model, 
			  COALESCE(a.serial_number, '') as serial_number, 
			  COALESCE(a.order_no, '') as order_no, 
			  COALESCE(a.license_number, '') as license_number, 
			  COALESCE(a.notes, '') as notes, 
			  a.created_at, a.updated_at, a.deleted_at,
			  COALESCE(at.name, '') as asset_type_name
			  FROM assets a
			  LEFT JOIN asset_types at ON a.asset_type_id = at.id
			  WHERE a.asset_type_id = ? AND a.deleted_at IS NULL ORDER BY a.name`
	err := r.db.SelectContext(ctx, &assets, query, assetTypeID)
	return assets, err
}

// GetWithCurrentAssignment retrieves all assets with their current assignment. If includeDeleted is true, returns only soft-deleted records.
func (r *AssetRepository) GetWithCurrentAssignment(ctx context.Context, includeDeleted bool) ([]models.AssetWithAssignment, error) {
	var assets []models.AssetWithAssignment
	deletedFilter := "a.deleted_at IS NULL"
	if includeDeleted {
		deletedFilter = "a.deleted_at IS NOT NULL"
	}
	query := `SELECT a.id, a.asset_type_id, a.name, 
			  COALESCE(a.model, '') as model, 
			  COALESCE(a.serial_number, '') as serial_number, 
			  COALESCE(a.order_no, '') as order_no, 
			  COALESCE(a.license_number, '') as license_number, 
			  COALESCE(a.notes, '') as notes, 
			  a.created_at, a.updated_at, a.deleted_at,
			  COALESCE(at.name, '') as asset_type_name,
			  p.name as currentassignee, p.id as currentassigneeid, aa.effective_from as assignedfrom
			  FROM assets a
			  LEFT JOIN asset_types at ON a.asset_type_id = at.id
			  LEFT JOIN asset_assignments aa ON a.id = aa.asset_id 
			      AND aa.deleted_at IS NULL 
			      AND aa.effective_from <= NOW() 
			      AND (aa.effective_to IS NULL OR aa.effective_to > NOW())
			  LEFT JOIN persons p ON aa.person_id = p.id
			  WHERE ` + deletedFilter + ` 
			  ORDER BY a.name`
	err := r.db.SelectContext(ctx, &assets, query)
	return assets, err
}

// Create creates a new asset
func (r *AssetRepository) Create(ctx context.Context, asset *models.Asset) error {
	query := `INSERT INTO assets (asset_type_id, name, model, serial_number, order_no, license_number, notes) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, asset.AssetTypeID, asset.Name, asset.Model,
		asset.SerialNumber, asset.OrderNo, asset.LicenseNumber, asset.Notes)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	asset.ID = id
	return nil
}

// Update updates an existing asset
func (r *AssetRepository) Update(ctx context.Context, asset *models.Asset) error {
	query := `UPDATE assets SET asset_type_id = ?, name = ?, model = ?, serial_number = ?, 
			  order_no = ?, license_number = ?, notes = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, asset.AssetTypeID, asset.Name, asset.Model,
		asset.SerialNumber, asset.OrderNo, asset.LicenseNumber, asset.Notes, asset.ID)
	return err
}

// Delete soft-deletes an asset
func (r *AssetRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE assets SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Search searches assets by name, serial number, or model
func (r *AssetRepository) Search(ctx context.Context, term string) ([]models.Asset, error) {
	var assets []models.Asset
	searchTerm := "%" + term + "%"
	query := `SELECT a.id, a.asset_type_id, a.name, 
			  COALESCE(a.model, '') as model, 
			  COALESCE(a.serial_number, '') as serial_number, 
			  COALESCE(a.order_no, '') as order_no, 
			  COALESCE(a.license_number, '') as license_number, 
			  COALESCE(a.notes, '') as notes, 
			  a.created_at, a.updated_at, a.deleted_at,
			  COALESCE(at.name, '') as asset_type_name
			  FROM assets a
			  LEFT JOIN asset_types at ON a.asset_type_id = at.id
			  WHERE a.deleted_at IS NULL 
			  AND (a.name LIKE ? OR a.serial_number LIKE ? OR a.model LIKE ?)
			  ORDER BY a.name`
	err := r.db.SelectContext(ctx, &assets, query, searchTerm, searchTerm, searchTerm)
	return assets, err
}
