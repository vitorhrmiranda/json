package unmarshal_test

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/vitorhrmiranda/json/unmarshal"
)

type parseable struct {
	Name unmarshal.NullableString `json:"name"`
	Last *string                  `json:"last"`
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
			json: `{"name":"test", "last":"test"}`,
			expect: parseable{
				Name: unmarshal.String("test"),
				Last: String("test")}},
		{
			when: "when send with empty value",
			json: `{"name":"", "last":""}`,
			expect: parseable{
				Name: unmarshal.String(""),
				Last: String("")}},
		{
			when: "when doesnt send any field",
			json: `{}`,
			expect: parseable{
				Name: unmarshal.NullableString{},
				Last: nil}},
		{
			when: "when send with empty value as null",
			json: `{"name":null, "last":null}`,
			expect: parseable{
				Name: unmarshal.String(""),
				Last: nil}},
	}

	for _, tt := range tcases {
		t.Run(tt.when, func(t *testing.T) {
			got := parse(t, []byte(tt.json))
			assert.Equal(t, tt.expect, got)
		})
	}
}
