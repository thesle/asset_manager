package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// MarshalJSON for NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		nt.Valid = false
		return nil
	}

	// Try standard RFC3339 format first
	var t time.Time
	if err := json.Unmarshal(data, &t); err == nil {
		nt.Time = t
		nt.Valid = true
		return nil
	}

	// Try date-only format (YYYY-MM-DD)
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	nt.Time = t
	nt.Valid = true
	return nil
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ns.String = s
	ns.Valid = true
	return nil
}

// NewNullTime creates a NullTime from a time.Time
func NewNullTime(t time.Time) NullTime {
	return NullTime{sql.NullTime{Time: t, Valid: true}}
}

// NewNullString creates a NullString from a string
func NewNullString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: s != ""}}
}
