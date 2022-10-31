package unmarshal

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullableString is a string with <nil> default value
type NullableString struct {
	Val   *string
	Valid bool
}

// String describes the value
// Returns <nil> when string is undefined
func (ns NullableString) String() (v string) {
	if ns.Val != nil {
		return *ns.Val
	}
	return `<nil>`
}

// UnmarshalJSON parses the JSON
// null value in JSON are mapped into nil with Valid flag as True
// when value doesnt specified specify the values is stated as nil
func (ns *NullableString) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)

	ns.Val = &str
	ns.Valid = true

	if string(b) == "null" {
		ns.Val = nil
	}

	return err
}

// Scan extends a Scanner from sql.NullString
func (ns *NullableString) Scan(value any) error {
	s := ns.nullString()
	err := s.Scan(value)

	ns.Val = &s.String
	ns.Valid = s.Valid

	if !s.Valid {
		ns.Val = nil
	}

	return err
}

// Values extends a Valuer from sql.NullString
func (ns NullableString) Value() (driver.Value, error) {
	return ns.nullString().Value()
}

// String converts native string into NullableString
func String(v string) NullableString {
	s := NullableString{Val: &v, Valid: true}
	return s
}

func (ns NullableString) nullString() *sql.NullString {
	if ns.Val == nil {
		return &sql.NullString{}
	}

	return &sql.NullString{String: *ns.Val, Valid: ns.Valid}
}
