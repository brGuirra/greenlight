//go:build unit
// +build unit

package data

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestUserModelIsAnonymous(t *testing.T) {
	testCases := []struct {
		name   string
		user   *User
		assert func(t *testing.T, isAnonymous bool)
	}{
		{
			name: "AnonymousUser",
			user: AnonymousUser,
			assert: func(t *testing.T, isAnonymous bool) {
				require.True(t, isAnonymous)
			},
		},
		{
			name: "Zero value User",
			user: &User{},
			assert: func(t *testing.T, isAnonymous bool) {
				require.False(t, isAnonymous)
			},
		},
		{
			name: "User with values",
			user: &User{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
			},
			assert: func(t *testing.T, isAnonymous bool) {
				require.False(t, isAnonymous)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				tc.assert(t, tc.user.IsAnonymous())
			},
		)
	}
}
