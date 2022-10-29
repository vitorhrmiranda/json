//go:build integration
// +build integration

package test_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorhrmiranda/json/test/with"
	"github.com/vitorhrmiranda/json/unmarshal"
)

// user is a model to test
type user struct {
	ID   int `gorm:"primaryKey"`
	Name unmarshal.NullableString
}

func String(v string) *string {
	return &v
}

func TestWithPostgres(t *testing.T) {
	const DefaultColumnValue = "default"
	ctx := context.TODO()

	db := with.Postgres(t, ctx).Connect(t, ctx)
	_ = db.AutoMigrate(&user{})

	tcases := []struct {
		json     string
		expected *string
	}{
		{json: `{"name":"test"}`, expected: String("test")},
		{json: `{"name":""}`, expected: String("")},
		{json: `{}`, expected: String(DefaultColumnValue)},
		{json: `{"name":null}`, expected: nil},
	}

	for _, tcase := range tcases {
		t.Run(tcase.json, func(t *testing.T) {
			// create with default value
			u := &user{Name: unmarshal.String(DefaultColumnValue)}
			_ = db.Create(&u)

			// parses JSON and update register
			u = &user{ID: u.ID}
			err := json.Unmarshal([]byte(tcase.json), u)
			assert.NoError(t, err)

			_ = db.Updates(&u)

			// get register and assert it
			u = &user{ID: u.ID}
			db.Find(u)
			assert.Equal(t, tcase.expected, u.Name.Val)

			t.Cleanup(func() { db.Delete(u) })
		})
	}
}
