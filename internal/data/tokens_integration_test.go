//go:build integration
// +build integration

package data

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

// createRandomToken it's a helper to populate the database
// with users. It takes in a pointer of `testing.T` and
// a pointer of `Models`, and returns a `Token` created
// with fake random data.
func createRandomToken(t *testing.T, m *Models, userID int64) *Token {
	ttl := time.Hour * 24
	scope := gofakeit.Noun()

	token, err := m.Tokens.New(userID, ttl, scope)

	require.NoError(t, err)

	require.Equal(t, token.UserID, userID)
	require.Equal(t, token.Scope, scope)
	require.WithinDuration(t, token.Expiry, time.Now().Add(ttl), time.Second)
	require.NotZero(t, token.Plaintext)
	require.NotZero(t, token.Hash)

	return token
}

// tokenModelTestsTeardown it's a helper to truncate the `users`
// table in the database during tests.
func tokenModelTestsTeardown(t *testing.T) {
	t.Helper()

	query := `TRUNCATE TABLE users RESTART IDENTITY CASCADE`

	_, err := testDB.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTokenModelNew(t *testing.T) {
	testModels := NewModels(testDB)

	user := createRandomUser(t, &testModels)
	createRandomToken(t, &testModels, user.ID)

	t.Cleanup(func() {
		tokenModelTestsTeardown(t)
	})
}

func TestTokenModelDeleteAllForUser(t *testing.T) {
	testModels := NewModels(testDB)

	user := createRandomUser(t, &testModels)
	token := createRandomToken(t, &testModels, user.ID)

	err := testModels.Tokens.DeleteAllForUser(token.Scope, token.UserID)

	require.NoError(t, err)
}
