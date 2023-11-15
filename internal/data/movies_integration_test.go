//go:build integration
// +build integration

package data

import (
	"testing"

	"github.com/brGuirra/greenlight/internal/assert"
)

func TestGetAll(t *testing.T) {
	_, _, err := testModels.Movies.GetAll("", []string{}, Filters{
		Sort:         "id",
		SortSafeList: []string{"id"},
		PageSize:     20,
		Page:         1,
	})
	assert.NilError(t, err)
}
