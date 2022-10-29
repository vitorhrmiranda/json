//go:build unit
// +build unit

package unmarshal_test

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/vitorhrmiranda/json/unmarshal"
)

type parseable struct {
	Name unmarshal.NullableString `json:"name"`
}

func parse(t *testing.T, b []byte) (p parseable) {
	t.Helper()

	err := json.Unmarshal(b, &p)
	assert.NoError(t, err)

	return p
}

func String(v string) *string {
	return &v
}

func TestParse(t *testing.T) {
	tcases := []struct {
		when, json string
		expect     parseable
	}{
		{
			when: "when send with value",
			json: `{"name":"test"}`,
			expect: parseable{
				Name: unmarshal.String("test"),
			}},
		{
			when: "when send with empty value",
			json: `{"name":""}`,
			expect: parseable{
				Name: unmarshal.String(""),
			}},
		{
			when: "when doesnt send any field",
			json: `{}`,
			expect: parseable{
				Name: unmarshal.NullableString{},
			}},
		{
			when: "when send with empty value as null",
			json: `{"name":null}`,
			expect: parseable{
				Name: unmarshal.NullableString{Valid: true},
			}},
	}

	for _, tt := range tcases {
		t.Run(tt.when, func(t *testing.T) {
			got := parse(t, []byte(tt.json))
			assert.Equal(t, tt.expect, got)
		})
	}
}
