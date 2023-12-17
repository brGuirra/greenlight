//go:build integration
// +build integration

package data

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

// createRandomUser it's a helper to populate the database
// with users. It takes in a pointer of `testing.T` and
// a pointer of `Models`, and returns a `User` created
// with fake random data.
func createRandomUser(t *testing.T, m *Models) User {
	user := User{
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Activated: false,
		Password: password{
			plaintext: nil,
			hash:      []byte(gofakeit.Password(true, true, true, true, false, 12)),
		},
	}

	err := m.Users.Insert(&user)

	require.NoError(t, err)

	require.Equal(t, user.Version, int32(1))
	require.False(t, user.Activated)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Name)
	require.NotZero(t, user.Email)
	require.NotZero(t, user.Password.hash)
	require.NotZero(t, user.CreatedAt)

	return user
}

// userModelTestsTeardown it's a helper to truncate the `users`
// table in the database during tests.
func userModelTestsTeardown(t *testing.T) {
	t.Helper()

	query := `TRUNCATE TABLE users RESTART IDENTITY CASCADE`

	_, err := testDB.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserModelInsert(t *testing.T) {
	testModels := NewModels(testDB)

	t.Run("Successfully created an user", func(t *testing.T) {
		createRandomUser(t, &testModels)

		t.Cleanup(func() {
			userModelTestsTeardown(t)
		})
	})

	t.Run("'ErrDuplicateEmail' when given 'Email' is already in use", func(t *testing.T) {
		createdUser := createRandomUser(t, &testModels)

		newUser := User{
			Name:      gofakeit.Name(),
			Email:     createdUser.Email,
			Activated: false,
			Password: password{
				plaintext: nil,
				hash:      []byte(gofakeit.Password(true, true, true, true, false, 12)),
			},
		}

		err := testModels.Users.Insert(&newUser)

		require.ErrorIs(t, err, ErrDuplicateEmail)

		t.Cleanup(func() {
			userModelTestsTeardown(t)
		})
	})
}

func TestUserModelGetByEmail(t *testing.T) {
	testModels := NewModels(testDB)

	t.Run("Successfully returns user data", func(t *testing.T) {
		createdUser := createRandomUser(t, &testModels)
		gotUser, err := testModels.Users.GetByEmail(createdUser.Email)

		require.NoError(t, err)

		require.Equal(t, createdUser.ID, gotUser.ID)
		require.Equal(t, createdUser.Name, gotUser.Name)
		require.Equal(t, createdUser.Email, gotUser.Email)
		require.Equal(t, createdUser.Password.hash, gotUser.Password.hash)
		require.Equal(t, createdUser.Activated, gotUser.Activated)
		require.Equal(t, createdUser.Version, gotUser.Version)
		require.WithinDuration(t, createdUser.CreatedAt, gotUser.CreatedAt, time.Second)

		t.Cleanup(func() {
			userModelTestsTeardown(t)
		})
	})

	t.Run("'ErrRecordNotFound' when given 'ID' does not exist", func(t *testing.T) {
		gotUser, err := testModels.Users.GetByEmail(gofakeit.Email())

		require.ErrorIs(t, err, ErrRecordNotFound)
		require.Zero(t, gotUser)
	})
}

func TestUserModelUdpate(t *testing.T) {
	testModels := NewModels(testDB)

	t.Run("Successfully updates user data", func(t *testing.T) {
		user := createRandomUser(t, &testModels)

		user.Name = gofakeit.Name()
		user.Email = gofakeit.Email()

		err := testModels.Users.Update(&user)
		require.NoError(t, err)

		updatedUser, err := testModels.Users.GetByEmail(user.Email)
		require.NoError(t, err)

		require.Equal(t, user.ID, updatedUser.ID)
		require.Equal(t, user.Name, updatedUser.Name)
		require.Equal(t, user.Email, updatedUser.Email)
		require.Equal(t, user.Password.hash, updatedUser.Password.hash)
		require.Equal(t, user.Activated, updatedUser.Activated)
		require.Equal(t, user.Version, int32(2))
		require.Equal(t, user.Version, updatedUser.Version)
		require.WithinDuration(t, user.CreatedAt, updatedUser.CreatedAt, time.Second)

		t.Cleanup(func() {
			userModelTestsTeardown(t)
		})
	})

	t.Run("'ErrDuplicateEmail' when given 'Email' is already in use", func(t *testing.T) {
		user1 := createRandomUser(t, &testModels)
		user2 := createRandomUser(t, &testModels)

		user1.Name = gofakeit.Name()
		user1.Email = user2.Email

		err := testModels.Users.Update(&user1)
		require.ErrorIs(t, err, ErrDuplicateEmail)

		t.Cleanup(func() {
			userModelTestsTeardown(t)
		})
	})
}

func TestGetForToken(t *testing.T) {
	t.Run("Successfully returns user data", func(t *testing.T) {
		testModels := NewModels(testDB)

		createdUser := createRandomUser(t, &testModels)
		token := createRandomToken(t, &testModels, createdUser.ID)

		gotUser, err := testModels.Users.GetForToken(token.Scope, token.Plaintext)

		require.NoError(t, err)

		require.Equal(t, createdUser.ID, gotUser.ID)
		require.Equal(t, createdUser.Name, gotUser.Name)
		require.Equal(t, createdUser.Email, gotUser.Email)
		require.Equal(t, createdUser.Password.hash, gotUser.Password.hash)
		require.Equal(t, createdUser.Activated, gotUser.Activated)
		require.Equal(t, createdUser.Version, gotUser.Version)
		require.WithinDuration(t, createdUser.CreatedAt, gotUser.CreatedAt, time.Second)

		t.Cleanup(func() {
			userModelTestsTeardown(t)
		})
	})

	t.Run("'ErrRecordNotFound' when given token information does not exist", func(t *testing.T) {
		testModels := NewModels(testDB)

		createRandomUser(t, &testModels)

		gotUser, err := testModels.Users.GetForToken(gofakeit.Noun(), gofakeit.Noun())

		require.ErrorIs(t, err, ErrRecordNotFound)
		require.Nil(t, gotUser)

		t.Cleanup(func() {
			userModelTestsTeardown(t)
		})
	})
}
