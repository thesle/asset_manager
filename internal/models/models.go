package models

import (
	"database/sql"
	"time"
)

// NullTime is a wrapper for sql.NullTime with JSON support
type NullTime struct {
	sql.NullTime
}

// NullString is a wrapper for sql.NullString with JSON support
type NullString struct {
	sql.NullString
}

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        int64     `db:"id" json:"ID"`
	CreatedAt time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt time.Time `db:"updated_at" json:"UpdatedAt"`
	DeletedAt NullTime  `db:"deleted_at" json:"DeletedAt,omitempty"`
}

// User represents an application user for authentication
type User struct {
	BaseModel
	Username     string `db:"username" json:"Username"`
	Email        string `db:"email" json:"Email"`
	PasswordHash string `db:"password_hash" json:"-"`
	IsActive     bool   `db:"is_active" json:"IsActive"`
}

// AssetType represents a category of assets
type AssetType struct {
	BaseModel
	Name        string `db:"name" json:"Name"`
	Description string `db:"description" json:"Description"`
}

// Asset represents a tracked asset
type Asset struct {
	BaseModel
	AssetTypeID   int64    `db:"asset_type_id" json:"AssetTypeID"`
	Name          string   `db:"name" json:"Name"`
	Model         string   `db:"model" json:"Model"`
	SerialNumber  string   `db:"serial_number" json:"SerialNumber"`
	OrderNo       string   `db:"order_no" json:"OrderNo"`
	LicenseNumber string   `db:"license_number" json:"LicenseNumber"`
	Notes         string   `db:"notes" json:"Notes"`
	PurchasedAt   NullTime `db:"purchased_at" json:"PurchasedAt,omitempty"`

	// Joined fields (not stored in assets table)
	AssetTypeName string `db:"asset_type_name" json:"AssetTypeName,omitempty"`
}

// DataType represents the type of data for properties and attributes
type DataType string

const (
	DataTypeString   DataType = "string"
	DataTypeInt      DataType = "int"
	DataTypeDecimal  DataType = "decimal"
	DataTypeBoolean  DataType = "boolean"
	DataTypeDate     DataType = "date"
	DataTypeDatetime DataType = "datetime"
	DataTypeEnum     DataType = "enum"
)

// Property defines a custom property that can be attached to assets
type Property struct {
	BaseModel
	Name        string   `db:"name" json:"Name"`
	DataType    DataType `db:"data_type" json:"DataType"`
	EnumOptions string   `db:"enum_options" json:"EnumOptions,omitempty"` // JSON array for enum options
}

// AssetProperty links an asset to a property value
type AssetProperty struct {
	BaseModel
	AssetID    int64  `db:"asset_id" json:"AssetID"`
	PropertyID int64  `db:"property_id" json:"PropertyID"`
	Value      string `db:"value" json:"Value"`

	// Joined fields
	PropertyName string   `db:"property_name" json:"PropertyName,omitempty"`
	DataType     DataType `db:"data_type" json:"DataType,omitempty"`
}

// Person represents a person who can be assigned assets
type Person struct {
	BaseModel
	Name  string `db:"name" json:"Name"`
	Email string `db:"email" json:"Email"`
	Phone string `db:"phone" json:"Phone"`
}

// Attribute defines a custom attribute that can be attached to persons
type Attribute struct {
	BaseModel
	Name        string   `db:"name" json:"Name"`
	DataType    DataType `db:"data_type" json:"DataType"`
	EnumOptions string   `db:"enum_options" json:"EnumOptions,omitempty"` // JSON array for enum options
}

// PersonAttribute links a person to an attribute value
type PersonAttribute struct {
	BaseModel
	PersonID    int64  `db:"person_id" json:"PersonID"`
	AttributeID int64  `db:"attribute_id" json:"AttributeID"`
	Value       string `db:"value" json:"Value"`

	// Joined fields
	AttributeName string   `db:"attribute_name" json:"AttributeName,omitempty"`
	DataType      DataType `db:"data_type" json:"DataType,omitempty"`
}

// AssetAssignment tracks the assignment of an asset to a person
type AssetAssignment struct {
	BaseModel
	AssetID       int64    `db:"asset_id" json:"AssetID"`
	PersonID      int64    `db:"person_id" json:"PersonID"`
	EffectiveFrom NullTime `db:"effective_from" json:"EffectiveFrom"`
	EffectiveTo   NullTime `db:"effective_to" json:"EffectiveTo,omitempty"`
	Notes         string   `db:"notes" json:"Notes"`

	// Joined fields
	AssetName         string `db:"asset_name" json:"AssetName,omitempty"`
	PersonName        string `db:"person_name" json:"PersonName,omitempty"`
	AssetTypeName     string `db:"asset_type_name" json:"AssetTypeName,omitempty"`
	AssetModel        string `db:"asset_model" json:"AssetModel,omitempty"`
	AssetSerialNumber string `db:"asset_serial_number" json:"AssetSerialNumber,omitempty"`
}

// AssetWithAssignment combines asset info with current assignment
type AssetWithAssignment struct {
	Asset
	AssetTypeName     string   `db:"asset_type_name" json:"AssetTypeName,omitempty"`
	CurrentAssignee   *string  `db:"currentassignee" json:"CurrentAssignee,omitempty"`
	CurrentAssigneeID *int64   `db:"currentassigneeid" json:"CurrentAssigneeID,omitempty"`
	AssignedFrom      NullTime `db:"assignedfrom" json:"AssignedFrom,omitempty"`
}

// LoginRequest represents a login attempt
type LoginRequest struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Remember bool   `json:"Remember"`
}

// LoginResponse represents a successful login
type LoginResponse struct {
	Token     string `json:"Token"`
	ExpiresAt int64  `json:"ExpiresAt"`
	User      User   `json:"User"`
}
