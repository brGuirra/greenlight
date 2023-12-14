package validator

import (
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	v := New()

	require.NotZero(t, v)
	require.Empty(t, v.Errors)
}

func TestValid(t *testing.T) {
	t.Run("'Valid' for valid Validator", func(t *testing.T) {
		fakeKey := gofakeit.Noun()
		fakeMessage := gofakeit.Phrase()

		v := New()

		v.Check(true, fakeKey, fakeMessage)

		isValid := v.Valid()

		require.True(t, isValid)
	})

	t.Run("'Valid' for brand new Validator", func(t *testing.T) {
		v := New()

		isValid := v.Valid()

		require.True(t, isValid)
	})

	t.Run("'Valid' for invalid Validator", func(t *testing.T) {
		fakeKey := gofakeit.Noun()
		fakeMessage := gofakeit.Phrase()

		v := New()

		v.Check(false, fakeKey, fakeMessage)

		isValid := v.Valid()

		require.False(t, isValid)
	})
}

func TestAddError(t *testing.T) {
	t.Run("'AddError' for brand new Validator", func(t *testing.T) {
		fakeKey := gofakeit.Noun()
		fakeMessage := gofakeit.Phrase()

		v := New()

		v.AddError(fakeKey, fakeMessage)

		require.NotEmpty(t, v.Errors)
		require.Equal(t, v.Errors[fakeKey], fakeMessage)
		require.Equal(t, len(v.Errors), 1)
	})

	t.Run("'AddError' non existing key", func(t *testing.T) {
		fakeKey := gofakeit.Noun()
		fakeMessage := gofakeit.Phrase()

		fakeNewKey := gofakeit.Noun()
		fakeNewMessage := gofakeit.Phrase()

		v := New()

		v.Check(false, fakeKey, fakeMessage)

		v.AddError(fakeNewKey, fakeNewMessage)

		require.NotEmpty(t, v.Errors)
		require.Equal(t, v.Errors[fakeKey], fakeMessage)
		require.Equal(t, len(v.Errors), 2)
	})

	t.Run("'AddError' existing key", func(t *testing.T) {
		fakeKey := gofakeit.Noun()
		fakeMessage := gofakeit.Phrase()

		fakeNewMessage := gofakeit.Phrase()

		v := New()

		v.Check(false, fakeKey, fakeMessage)

		v.AddError(fakeKey, fakeNewMessage)

		require.NotEmpty(t, v.Errors)
		require.Equal(t, v.Errors[fakeKey], fakeMessage)
		require.Equal(t, len(v.Errors), 1)
	})
}

func TestCheck(t *testing.T) {
	fakeKey := gofakeit.Noun()
	fakeMessage := gofakeit.Phrase()

	v := New()

	v.Check(false, fakeKey, fakeMessage)

	require.NotEmpty(t, v.Errors)
	require.Equal(t, v.Errors[fakeKey], fakeMessage)
	require.Equal(t, len(v.Errors), 1)
}

func TestPermittedValue(t *testing.T) {
	t.Run("Permitted value", func(t *testing.T) {
		permittedValues := []string{gofakeit.Noun(), gofakeit.Noun()}
		value := permittedValues[0]

		isPermitted := PermittedValue(value, permittedValues...)

		require.True(t, isPermitted)
	})

	t.Run("Not permitted value", func(t *testing.T) {
		permittedValues := []string{gofakeit.Noun(), gofakeit.Noun()}
		value := gofakeit.Verb()

		isPermitted := PermittedValue(value, permittedValues...)

		require.False(t, isPermitted)
	})

	t.Run("With zero value and empty list of permitted values", func(t *testing.T) {
		permittedValues := []string{}
		value := ""

		isPermitted := PermittedValue(value, permittedValues...)

		require.False(t, isPermitted)
	})
}

func TestUnique(t *testing.T) {
	t.Run("Unique values", func(t *testing.T) {
		values := []string{gofakeit.Noun(), gofakeit.Noun()}

		isUnique := Unique(values)

		require.True(t, isUnique)
	})

	t.Run("Non unique values", func(t *testing.T) {
		repeated := gofakeit.Noun()
		values := []string{repeated, gofakeit.Noun(), repeated}

		isUnique := Unique(values)

		require.False(t, isUnique)
	})

	t.Run("With empty list of values", func(t *testing.T) {
		values := []string{}

		isUnique := Unique(values)

		require.True(t, isUnique)
	})
}

func TestMatches(t *testing.T) {
	t.Run("With match", func(t *testing.T) {
		value := "any_value"

		hasMatched := Matches(value, regexp.MustCompile(`any`))

		require.True(t, hasMatched)
	})

	t.Run("With no match", func(t *testing.T) {
		value := ""

		hasMatched := Matches(value, regexp.MustCompile(`any`))

		require.False(t, hasMatched)
	})
}
