//go:build integration
// +build integration

package data

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

// permissionsModelTestsTeardown it's a helper to truncate the `users`
// table in the database during tests.
func permissionsModelTestsTeardown(t *testing.T) {
	t.Helper()

	query := `TRUNCATE TABLE users RESTART IDENTITY CASCADE`

	_, err := testDB.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPermissionsModelAddForUser(t *testing.T) {
	testModels := NewModels(testDB)

	user := createRandomUser(t, &testModels)

	err := testModels.Permissions.AddForUser(user.ID, "movies:read")

	require.NoError(t, err)

	t.Cleanup(func() {
		permissionsModelTestsTeardown(t)
	})
}

func TestPermissionsModelGetAllForUser(t *testing.T) {
	testModels := NewModels(testDB)

	t.Run("Succesfully returns user permissions", func(t *testing.T) {
		user := createRandomUser(t, &testModels)
		permission := "movies:read"

		err := testModels.Permissions.AddForUser(user.ID, permission)

		require.NoError(t, err)

		permissions, err := testModels.Permissions.GetAllForUser(user.ID)

		require.NoError(t, err)

		require.Equal(t, 1, len(permissions))

		t.Cleanup(func() {
			permissionsModelTestsTeardown(t)
		})
	})

	t.Run("Returns and empty slice when no permissions are related to an user", func(t *testing.T) {
		permissions, err := testModels.Permissions.GetAllForUser(gofakeit.Int64())

		require.NoError(t, err)

		require.Empty(t, permissions)

		t.Cleanup(func() {
			permissionsModelTestsTeardown(t)
		})
	})
}
