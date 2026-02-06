package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"assetManager/internal/models"
)

var (
	ErrAssetAssignmentNotFound = errors.New("asset assignment not found")
	ErrOverlappingAssignment   = errors.New("overlapping assignment exists")
)

// AssetAssignmentRepository handles asset assignment data operations
type AssetAssignmentRepository struct {
	db *sqlx.DB
}

// NewAssetAssignmentRepository creates a new asset assignment repository
func NewAssetAssignmentRepository(db *sqlx.DB) *AssetAssignmentRepository {
	return &AssetAssignmentRepository{db: db}
}

// GetByID retrieves an asset assignment by ID
func (r *AssetAssignmentRepository) GetByID(ctx context.Context, id int64) (*models.AssetAssignment, error) {
	var aa models.AssetAssignment
	query := `SELECT aa.id, aa.asset_id, aa.person_id, aa.effective_from, aa.effective_to, COALESCE(aa.notes, '') as notes,
			  aa.created_at, aa.updated_at, aa.deleted_at,
			  COALESCE(a.name, '') as asset_name, COALESCE(p.name, '') as person_name
			  FROM asset_assignments aa
			  LEFT JOIN assets a ON aa.asset_id = a.id
			  LEFT JOIN persons p ON aa.person_id = p.id
			  WHERE aa.id = ? AND aa.deleted_at IS NULL`
	err := r.db.GetContext(ctx, &aa, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetAssignmentNotFound
	}
	return &aa, err
}

// GetCurrentByAssetID retrieves the current assignment for an asset
func (r *AssetAssignmentRepository) GetCurrentByAssetID(ctx context.Context, assetID int64) (*models.AssetAssignment, error) {
	var aa models.AssetAssignment
	query := `SELECT aa.id, aa.asset_id, aa.person_id, aa.effective_from, aa.effective_to, COALESCE(aa.notes, '') as notes,
			  aa.created_at, aa.updated_at, aa.deleted_at,
			  COALESCE(a.name, '') as asset_name, COALESCE(p.name, '') as person_name
			  FROM asset_assignments aa
			  LEFT JOIN assets a ON aa.asset_id = a.id
			  LEFT JOIN persons p ON aa.person_id = p.id
			  WHERE aa.asset_id = ? AND aa.deleted_at IS NULL
			  AND aa.effective_from <= NOW()
			  AND (aa.effective_to IS NULL OR aa.effective_to > NOW())
			  ORDER BY aa.effective_from DESC LIMIT 1`
	err := r.db.GetContext(ctx, &aa, query, assetID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAssetAssignmentNotFound
	}
	return &aa, err
}

// GetHistoryByAssetID retrieves assignment history for an asset
func (r *AssetAssignmentRepository) GetHistoryByAssetID(ctx context.Context, assetID int64) ([]models.AssetAssignment, error) {
	var aas []models.AssetAssignment
	query := `SELECT aa.id, aa.asset_id, aa.person_id, aa.effective_from, aa.effective_to, COALESCE(aa.notes, '') as notes,
			  aa.created_at, aa.updated_at, aa.deleted_at,
			  COALESCE(a.name, '') as asset_name, COALESCE(p.name, '') as person_name
			  FROM asset_assignments aa
			  LEFT JOIN assets a ON aa.asset_id = a.id
			  LEFT JOIN persons p ON aa.person_id = p.id
			  WHERE aa.asset_id = ? AND aa.deleted_at IS NULL
			  ORDER BY aa.effective_from DESC`
	err := r.db.SelectContext(ctx, &aas, query, assetID)
	return aas, err
}

// GetByPersonID retrieves all assignments for a person
func (r *AssetAssignmentRepository) GetByPersonID(ctx context.Context, personID int64) ([]models.AssetAssignment, error) {
	var aas []models.AssetAssignment
	query := `SELECT aa.id, aa.asset_id, aa.person_id, aa.effective_from, aa.effective_to, COALESCE(aa.notes, '') as notes,
			  aa.created_at, aa.updated_at, aa.deleted_at,
			  COALESCE(a.name, '') as asset_name, COALESCE(p.name, '') as person_name
			  FROM asset_assignments aa
			  LEFT JOIN assets a ON aa.asset_id = a.id
			  LEFT JOIN persons p ON aa.person_id = p.id
			  WHERE aa.person_id = ? AND aa.deleted_at IS NULL
			  ORDER BY aa.effective_from DESC`
	err := r.db.SelectContext(ctx, &aas, query, personID)
	return aas, err
}

// GetCurrentByPersonID retrieves current assignments for a person
func (r *AssetAssignmentRepository) GetCurrentByPersonID(ctx context.Context, personID int64) ([]models.AssetAssignment, error) {
	var aas []models.AssetAssignment
	query := `SELECT aa.id, aa.asset_id, aa.person_id, aa.effective_from, aa.effective_to, COALESCE(aa.notes, '') as notes,
			  aa.created_at, aa.updated_at, aa.deleted_at,
			  COALESCE(a.name, '') as asset_name, COALESCE(p.name, '') as person_name
			  FROM asset_assignments aa
			  LEFT JOIN assets a ON aa.asset_id = a.id
			  LEFT JOIN persons p ON aa.person_id = p.id
			  WHERE aa.person_id = ? AND aa.deleted_at IS NULL
			  AND aa.effective_from <= NOW()
			  AND (aa.effective_to IS NULL OR aa.effective_to > NOW())
			  ORDER BY aa.effective_from DESC`
	err := r.db.SelectContext(ctx, &aas, query, personID)
	return aas, err
}

// CheckOverlap checks if there's an overlapping assignment for an asset
func (r *AssetAssignmentRepository) CheckOverlap(ctx context.Context, assetID int64, from, to time.Time, excludeID int64) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM asset_assignments 
			  WHERE asset_id = ? AND deleted_at IS NULL AND id != ?
			  AND effective_from < ? 
			  AND (effective_to IS NULL OR effective_to > ?)`
	err := r.db.GetContext(ctx, &count, query, assetID, excludeID, to, from)
	return count > 0, err
}

// Create creates a new asset assignment
func (r *AssetAssignmentRepository) Create(ctx context.Context, aa *models.AssetAssignment) error {
	// Check for overlapping assignments
	toTime := time.Now().AddDate(100, 0, 0) // Far future if no end date
	if aa.EffectiveTo.Valid {
		toTime = aa.EffectiveTo.Time
	}
	overlap, err := r.CheckOverlap(ctx, aa.AssetID, aa.EffectiveFrom.Time, toTime, 0)
	if err != nil {
		return err
	}
	if overlap {
		return ErrOverlappingAssignment
	}

	query := `INSERT INTO asset_assignments (asset_id, person_id, effective_from, effective_to, notes) 
			  VALUES (?, ?, ?, ?, ?)`
	var effectiveTo interface{}
	if aa.EffectiveTo.Valid {
		effectiveTo = aa.EffectiveTo.Time
	}
	result, err := r.db.ExecContext(ctx, query, aa.AssetID, aa.PersonID, aa.EffectiveFrom.Time, effectiveTo, aa.Notes)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	aa.ID = id
	return nil
}

// Update updates an existing asset assignment
func (r *AssetAssignmentRepository) Update(ctx context.Context, aa *models.AssetAssignment) error {
	// Check for overlapping assignments (excluding this one)
	toTime := time.Now().AddDate(100, 0, 0)
	if aa.EffectiveTo.Valid {
		toTime = aa.EffectiveTo.Time
	}
	overlap, err := r.CheckOverlap(ctx, aa.AssetID, aa.EffectiveFrom.Time, toTime, aa.ID)
	if err != nil {
		return err
	}
	if overlap {
		return ErrOverlappingAssignment
	}

	query := `UPDATE asset_assignments SET person_id = ?, effective_from = ?, effective_to = ?, 
			  notes = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	var effectiveTo interface{}
	if aa.EffectiveTo.Valid {
		effectiveTo = aa.EffectiveTo.Time
	}
	_, err = r.db.ExecContext(ctx, query, aa.PersonID, aa.EffectiveFrom.Time, effectiveTo, aa.Notes, aa.ID)
	return err
}

// EndAssignment ends an assignment by setting the effective_to date
func (r *AssetAssignmentRepository) EndAssignment(ctx context.Context, id int64, endDate time.Time) error {
	query := `UPDATE asset_assignments SET effective_to = ?, updated_at = NOW() 
			  WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, endDate, id)
	return err
}

// Delete soft-deletes an asset assignment
func (r *AssetAssignmentRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE asset_assignments SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// AssignAsset assigns an asset to a person, ending any current assignment
func (r *AssetAssignmentRepository) AssignAsset(ctx context.Context, assetID, personID int64, notes string, effectiveDate time.Time) error {
	// End current assignment if exists
	current, err := r.GetCurrentByAssetID(ctx, assetID)
	if err == nil && current != nil {
		if err := r.EndAssignment(ctx, current.ID, effectiveDate); err != nil {
			return err
		}
	}

	// Create new assignment
	aa := &models.AssetAssignment{
		AssetID:       assetID,
		PersonID:      personID,
		EffectiveFrom: models.NewNullTime(effectiveDate),
		Notes:         notes,
	}
	return r.Create(ctx, aa)
}
