package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ReportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

type FilterCondition struct {
	Field         string      `json:"Field"`
	Operator      string      `json:"Operator"`
	Value         interface{} `json:"Value"`
	LogicOperator string      `json:"LogicOperator"` // AND or OR
}

type CustomReportRequest struct {
	EntityType string            `json:"EntityType"`
	Filters    []FilterCondition `json:"Filters"`
}

func (r *ReportRepository) ExecuteAssetReport(ctx context.Context, filters []FilterCondition) ([]map[string]interface{}, error) {
	// Separate filters into SQL filters (base fields) and post filters (properties)
	sqlFilters, postFilters := separateFilters(filters, "asset")

	query := `
		SELECT 
			a.id, a.asset_type_id, a.name, a.model, a.serial_number, 
			a.order_no, a.license_number, a.notes, a.purchased_at,
			a.created_at, a.updated_at, a.deleted_at,
			at.name as asset_type_name,
			COALESCE(p.name, 'Unassigned') as current_assignee,
			asgn.person_id as current_assignee_id
		FROM assets a
		LEFT JOIN asset_types at ON a.asset_type_id = at.id
		LEFT JOIN (
			SELECT asset_id, person_id, ROW_NUMBER() OVER (PARTITION BY asset_id ORDER BY effective_from DESC) as rn
			FROM asset_assignments
			WHERE effective_to IS NULL OR effective_to > NOW()
		) asgn ON a.id = asgn.asset_id AND asgn.rn = 1
		LEFT JOIN persons p ON asgn.person_id = p.id
	`

	var args []interface{}
	argCounter := 1

	whereClause, whereArgs := buildWhereClauseWithLogic(sqlFilters, &argCounter)
	if whereClause != "" {
		query += " WHERE " + whereClause
		args = append(args, whereArgs...)
	}

	query += " ORDER BY a.name"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}
		// Convert byte arrays to strings
		convertBytesToStrings(result)
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Load properties for each asset
	for _, result := range results {
		assetID := result["id"]
		properties, err := r.getAssetProperties(ctx, assetID)
		if err != nil {
			continue
		}
		for k, v := range properties {
			result[k] = v
		}
	}

	// Apply post-filters (property filters)
	if len(postFilters) > 0 {
		results = applyPostFilters(results, postFilters)
	}

	return results, nil
}

func (r *ReportRepository) ExecutePersonReport(ctx context.Context, filters []FilterCondition) ([]map[string]interface{}, error) {
	// Separate filters into SQL filters (base fields) and post filters (attributes)
	sqlFilters, postFilters := separateFilters(filters, "person")

	query := `
		SELECT 
			p.id, p.name, p.email, p.phone,
			p.created_at, p.updated_at, p.deleted_at
		FROM persons p
		WHERE p.name != 'Unassigned'
	`

	var args []interface{}
	argCounter := 1

	whereClause, whereArgs := buildWhereClauseWithLogic(sqlFilters, &argCounter)
	if whereClause != "" {
		query += " AND " + whereClause
		args = append(args, whereArgs...)
	}

	query += " ORDER BY p.name"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			return nil, err
		}
		// Convert byte arrays to strings
		convertBytesToStrings(result)
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Load attributes for each person
	for _, result := range results {
		personID := result["id"]
		attributes, err := r.getPersonAttributes(ctx, personID)
		if err != nil {
			continue
		}
		for k, v := range attributes {
			result[k] = v
		}
	}

	// Apply post-filters (attribute filters)
	if len(postFilters) > 0 {
		results = applyPostFilters(results, postFilters)
	}

	return results, nil
}

func buildWhereClauseWithLogic(filters []FilterCondition, argCounter *int) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	var clauses []string
	var args []interface{}

	for i, filter := range filters {
		clause, arg := buildWhereClause(filter, argCounter)
		if clause == "" {
			continue
		}

		// Add logic operator before this clause (except for the first one)
		if len(clauses) > 0 {
			logicOp := filters[i-1].LogicOperator // Use previous filter's operator
			if logicOp != "OR" {
				logicOp = "AND" // Default to AND
			}
			clauses = append(clauses, logicOp)
		}

		// Add the clause
		clauses = append(clauses, clause)
		if arg != nil {
			args = append(args, arg)
		}
	}

	if len(clauses) == 0 {
		return "", nil
	}

	// Join all parts with spaces
	return "(" + strings.Join(clauses, " ") + ")", args
}

func buildWhereClause(filter FilterCondition, argCounter *int) (string, interface{}) {
	field := sanitizeFieldName(filter.Field)
	if field == "" {
		return "", nil
	}

	switch filter.Operator {
	case "=":
		clause := fmt.Sprintf("%s = ?", field)
		*argCounter++
		return clause, filter.Value
	case "!=":
		clause := fmt.Sprintf("%s != ?", field)
		*argCounter++
		return clause, filter.Value
	case ">":
		clause := fmt.Sprintf("%s > ?", field)
		*argCounter++
		return clause, filter.Value
	case "<":
		clause := fmt.Sprintf("%s < ?", field)
		*argCounter++
		return clause, filter.Value
	case ">=":
		clause := fmt.Sprintf("%s >= ?", field)
		*argCounter++
		return clause, filter.Value
	case "<=":
		clause := fmt.Sprintf("%s <= ?", field)
		*argCounter++
		return clause, filter.Value
	case "LIKE":
		clause := fmt.Sprintf("%s LIKE ?", field)
		*argCounter++
		return clause, fmt.Sprintf("%%%v%%", filter.Value)
	case "NOT LIKE":
		clause := fmt.Sprintf("%s NOT LIKE ?", field)
		*argCounter++
		return clause, fmt.Sprintf("%%%v%%", filter.Value)
	case "IS NULL":
		return fmt.Sprintf("%s IS NULL", field), nil
	case "IS NOT NULL":
		return fmt.Sprintf("%s IS NOT NULL", field), nil
	default:
		return "", nil
	}
}

func convertBytesToStrings(m map[string]interface{}) {
	for k, v := range m {
		if b, ok := v.([]byte); ok {
			m[k] = string(b)
		}
	}
}

func separateFilters(filters []FilterCondition, entityType string) (sqlFilters, postFilters []FilterCondition) {
	for _, filter := range filters {
		// Check if this is a property or attribute filter
		if strings.HasPrefix(filter.Field, "prop_") || strings.HasPrefix(filter.Field, "attr_") {
			postFilters = append(postFilters, filter)
		} else {
			sqlFilters = append(sqlFilters, filter)
		}
	}
	return sqlFilters, postFilters
}

func applyPostFilters(results []map[string]interface{}, filters []FilterCondition) []map[string]interface{} {
	var filtered []map[string]interface{}

	for _, result := range results {
		if matchesPostFilters(result, filters) {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

func matchesPostFilters(record map[string]interface{}, filters []FilterCondition) bool {
	// Track the current match state and logic operators
	if len(filters) == 0 {
		return true
	}

	// Evaluate filters with AND/OR logic
	currentMatch := true
	nextOperator := "AND"

	for i, filter := range filters {
		fieldMatch := matchesFilter(record, filter)

		// Apply the logic operator
		if i == 0 {
			currentMatch = fieldMatch
		} else {
			if nextOperator == "OR" {
				currentMatch = currentMatch || fieldMatch
			} else { // AND
				currentMatch = currentMatch && fieldMatch
			}
		}

		// Set next operator (from current filter's LogicOperator)
		nextOperator = filter.LogicOperator
		if nextOperator != "OR" {
			nextOperator = "AND"
		}
	}

	return currentMatch
}

func matchesFilter(record map[string]interface{}, filter FilterCondition) bool {
	fieldValue, exists := record[filter.Field]

	switch filter.Operator {
	case "IS NULL":
		return !exists || fieldValue == nil || fieldValue == ""
	case "IS NOT NULL":
		return exists && fieldValue != nil && fieldValue != ""
	}

	if !exists {
		return false
	}

	// Convert to string for comparison
	fieldStr := fmt.Sprintf("%v", fieldValue)
	filterStr := fmt.Sprintf("%v", filter.Value)

	switch filter.Operator {
	case "=":
		return fieldStr == filterStr
	case "!=":
		return fieldStr != filterStr
	case "LIKE":
		return strings.Contains(strings.ToLower(fieldStr), strings.ToLower(filterStr))
	case "NOT LIKE":
		return !strings.Contains(strings.ToLower(fieldStr), strings.ToLower(filterStr))
	case ">":
		return fieldStr > filterStr
	case "<":
		return fieldStr < filterStr
	case ">=":
		return fieldStr >= filterStr
	case "<=":
		return fieldStr <= filterStr
	default:
		return false
	}
}

func sanitizeFieldName(field string) string {
	// Map frontend field names to database column names
	fieldMap := map[string]string{
		"ID":              "a.id",
		"Name":            "a.name",
		"AssetTypeName":   "at.name",
		"Model":           "a.model",
		"SerialNumber":    "a.serial_number",
		"OrderNo":         "a.order_no",
		"LicenseNumber":   "a.license_number",
		"Notes":           "a.notes",
		"PurchasedAt":     "a.purchased_at",
		"CurrentAssignee": "p.name",
		"Email":           "p.email",
		"Phone":           "p.phone",
		"PersonName":      "p.name",
		"PersonEmail":     "p.email",
		"PersonPhone":     "p.phone",
	}

	if mapped, ok := fieldMap[field]; ok {
		return mapped
	}

	// Property and attribute fields are prefixed with prop_ or attr_
	// These cannot be filtered in the main query and should be ignored
	// (they need to be filtered at the application level after loading)
	return ""
}

func (r *ReportRepository) getAssetProperties(ctx context.Context, assetID interface{}) (map[string]interface{}, error) {
	query := `
		SELECT prop.name, ap.value
		FROM assets_properties ap
		JOIN properties prop ON ap.property_id = prop.id
		WHERE ap.asset_id = ?
	`

	rows, err := r.db.QueryxContext(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	properties := make(map[string]interface{})
	for rows.Next() {
		var name string
		var value []byte
		if err := rows.Scan(&name, &value); err != nil {
			continue
		}
		// Convert byte array to string for TEXT columns
		properties["prop_"+name] = string(value)
	}

	return properties, nil
}

func (r *ReportRepository) getPersonAttributes(ctx context.Context, personID interface{}) (map[string]interface{}, error) {
	query := `
		SELECT attr.name, pa.value
		FROM persons_attributes pa
		JOIN attributes attr ON pa.attribute_id = attr.id
		WHERE pa.person_id = ?
	`

	rows, err := r.db.QueryxContext(ctx, query, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attributes := make(map[string]interface{})
	for rows.Next() {
		var name string
		var value []byte
		if err := rows.Scan(&name, &value); err != nil {
			continue
		}
		// Convert byte array to string for TEXT columns
		attributes["attr_"+name] = string(value)
	}

	return attributes, nil
}
