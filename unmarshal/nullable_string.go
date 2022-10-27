package unmarshal

import (
	"encoding/json"
)

// NullableString is a string with <nil> default value
type NullableString struct {
	Value *string
}

// String describes the value
// Returns <nil> when string is undefined
func (s NullableString) String() (v string) {
	if s.Value != nil {
		return *s.Value
	}
	return `<nil>`
}

// UnmarshalJSON parses the JSON
// null value in JSON are mapped into empty string
// when value doesnt specified specify the values is stated as nil
func (s *NullableString) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)

	s.Value = &str

	return err
}

// String converts native string into NullableString
func String(v string) NullableString {
	s := NullableString{Value: &v}
	return s
}
