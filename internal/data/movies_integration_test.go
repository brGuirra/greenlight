//go:build integration
// +build integration

package data

import (
	"os"
	"regexp"
	"slices"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// sortOrder indicates the sort order, it is represented
// by itegers to help controlling sort direction with `slices.SortFunc`
type sortOrder int

const (
	ASC  sortOrder = 1
	DESC sortOrder = -1
)

// fallbackSort takes in two `Movie` and returns
// and integer indicationg if `a.ID` is smaller,
// bigger than `b.ID`. This helper is used as a
// fallback when other properties of two `Movie`
// elements are equal in other sort functions.
func fallbackSort(a, b Movie) int {
	if a.ID < b.ID {
		return -1
	}

	return 1
}

// sortById it's a helper to sort elements during test cases.
// It takes in a `[]Movie` and a `sortOrder`, and returns a `[]Movie`
// sorted by `ID`.
func sortById(m []Movie, so sortOrder) []Movie {
	cn := slices.Clone(m)

	slices.SortFunc(cn, func(a, b Movie) int {
		if a.ID < b.ID {
			return -1 * int(so)
		}

		return 1 * int(so)
	})

	return cn
}

// sortByTitle it's a helper to sort elements during test cases.
// It takes in a `[]Movie` and a `sortOrder`, and returns a `[]Movie`
// sorted by `Title`. If the title is equal between elements in
// the slice, the `ID` property will determine which element
// should be placed first.
func sortByTitle(m []Movie, so sortOrder) []Movie {
	cn := slices.Clone(m)

	slices.SortFunc(cn, func(a, b Movie) int {
		if a.Title < b.Title {
			return -1 * int(so)
		}

		if a.Title > b.Title {
			return 1 * int(so)
		}

		return fallbackSort(a, b)
	})

	return cn
}

// sortByYear it's a helper to sort elements during test cases.
// It takes in a `[]Movie` and a `sortOrder`, and returns a
// `[]Movie` sorted by `Year`. If the year is equal between
// elements in the slice, the `ID` property will determine
// which element should be placed first.
func sortByYear(m []Movie, so sortOrder) []Movie {
	cn := slices.Clone(m)

	slices.SortFunc(cn, func(a, b Movie) int {
		if a.Year < b.Year {
			return -1 * int(so)
		}

		if a.Year > b.Year {
			return 1 * int(so)
		}

		return fallbackSort(a, b)
	})

	return cn
}

// sortByRuntime it's a helper to sort elements during test cases.
// It takes in a `[]Movie` and a `sortOrder`, and returns a
// `[]Movie` sorted by `Runtime`. If the runtime is equal between
// elements in the slice, the `ID` property will determine
// which element should be placed first.
func sortByRuntime(m []Movie, so sortOrder) []Movie {
	cn := slices.Clone(m)

	slices.SortFunc(cn, func(a, b Movie) int {
		if a.Runtime < b.Runtime {
			return -1 * int(so)
		}

		if a.Runtime > b.Runtime {
			return 1 * int(so)
		}

		return fallbackSort(a, b)
	})

	return cn
}

// filterByTitle it's a helper to filter elemests
// by a title search term during tests.
// It takes in a string with the search and a `[]Movie`,
// and returns a new slice of `Movie` that matches the
// search term.
func filterByTitle(st string, movies []Movie) []Movie {
	results := make([]Movie, 0)

	rx := regexp.MustCompile(`(?i)\b\w*` + st + `\w*\b`)

	for _, m := range movies {
		if rx.MatchString(m.Title) {
			results = append(results, m)
		}
	}

	return results
}

// filterByGenres it's a helper to filter elements
// by a slices of genres during tests.
// It takes in the slice of genres and a `[]Movie`,
// and returns a new slice of `Movie` containing
// only the movie that matched the required genres.
func filterByGenres(genres []string, movies []Movie) []Movie {
	results := make([]Movie, 0)

	for _, m := range movies {
		matchs := 0

		for _, g := range genres {
			if slices.Contains(m.Genres, g) {
				matchs++
			}
		}

		if matchs == len(genres) {
			results = append(results, m)
		}
	}

	return results
}

// createRandomMovie it's a helper to populate the database
// with movies. It takes in a pointer of `testing.T` and
// a pointer of `Models`, and returns a `Movie` created
// with fake random data.
func createRandomMovie(t *testing.T, m *Models) Movie {
	movie := Movie{
		Title:   gofakeit.MovieName(),
		Genres:  []string{gofakeit.MovieGenre()},
		Year:    int32(gofakeit.Year()),
		Runtime: Runtime(gofakeit.Number(90, 180)),
	}

	err := m.Movies.Insert(&movie)

	require.NoError(t, err)

	require.Equal(t, movie.Version, int32(1))
	require.NotZero(t, movie.ID)
	require.NotZero(t, movie.CreatedAt)
	require.NotZero(t, movie.Title)
	require.NotZero(t, movie.Year)
	require.NotZero(t, movie.Runtime)
	require.NotZero(t, movie.Genres)

	return movie
}

// setup it's a helper to populate the database with
// movies. It only takes in a pointer of `*testing.T`
// and returns a `[]Movie`. The data is static and
// came from `./db/seed/movies.sql`.
func setup(t *testing.T) []Movie {
	t.Helper()

	movies := make([]Movie, 0)

	script, err := os.ReadFile("../../db/seed/movies.sql")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := testDB.Query(string(script))
	if err != nil {
		t.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			t.Fatal(err)
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}

	return movies
}

// teardown it's a helper to truncate the `movies`
// table in the database during tests.
func teardown(t *testing.T) {
	t.Helper()

	query := `TRUNCATE TABLE movies RESTART IDENTITY`

	_, err := testDB.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	testModels := NewModels(testDB)

	createRandomMovie(t, &testModels)

	t.Cleanup(func() {
		teardown(t)
	})
}

func TestGetAll(t *testing.T) {
	const titleSearchTerm = "The"
	const noResultsTitleSearchTerm = "Damage"

	testModels := NewModels(testDB)

	movies := setup(t)

	assertMovies := func(t *testing.T, expected []Movie, actual []Movie) {
		require.NotZero(t, actual)
		require.Len(t, actual, len(expected))

		for i := 0; i < len(expected); i++ {
			require.Equal(t, expected[i].Title, actual[i].Title)
			require.Equal(t, expected[i].Genres, actual[i].Genres)
			require.Equal(t, expected[i].Year, actual[i].Year)
			require.Equal(t, expected[i].Runtime, actual[i].Runtime)

			require.Equal(t, actual[i].Version, int32(1))

			require.NotZero(t, actual[i].ID)
			require.NotZero(t, actual[i].CreatedAt)
		}
	}

	assertMetadata := func(t *testing.T, expected Metadata, actual Metadata) {
		require.NotZero(t, actual)

		require.Equal(t, expected.CurrentPage, actual.CurrentPage)
		require.Equal(t, expected.PageSize, actual.PageSize)
		require.Equal(t, expected.FirstPage, actual.FirstPage)
		require.Equal(t, expected.LastPage, actual.LastPage)
		require.Equal(t, expected.TotalRecords, actual.TotalRecords)
	}

	testCases := []struct {
		name    string
		title   string
		genres  []string
		filters Filters
		assert  func(t *testing.T, data []Movie, meta Metadata, err error)
	}{
		{
			name:   "Sort by 'ID' with 'ASC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				require.NoError(t, err)
				assertMovies(t, movies[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'ID' with 'DESC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortById(movies, DESC)

				require.NoError(t, err)
				assertMovies(t, expected[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'Title' with 'ASC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "title",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortByTitle(movies, 1)

				require.NoError(t, err)
				assertMovies(t, expected[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'Title' with 'DESC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-title",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortByTitle(movies, -1)

				require.NoError(t, err)
				assertMovies(t, expected[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'Year' with 'ASC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "year",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortByYear(movies, ASC)

				require.NoError(t, err)
				assertMovies(t, expected[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'Year' with 'DESC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-year",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortByYear(movies, DESC)

				require.NoError(t, err)
				assertMovies(t, expected[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'Runtime' with 'ASC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortByRuntime(movies, ASC)

				require.NoError(t, err)
				assertMovies(t, expected[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'Runtime' with 'DESC'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortByRuntime(movies, DESC)

				require.NoError(t, err)
				assertMovies(t, expected[0:20], data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'ID' with 'ASC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'ID' with 'DESC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = sortById(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'Title' with 'ASC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "title",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = sortByTitle(expected, ASC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'Title' with 'DESC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-title",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = sortByTitle(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'Year' with 'ASC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "year",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = sortByYear(expected, ASC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'Year' with 'DESC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-year",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = sortByYear(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'Runtime' with 'ASC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = sortByRuntime(expected, ASC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and sorted by 'Runtime' with 'DESC'",
			title:  titleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = sortByRuntime(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 14,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'ID' with 'ASC'",
			title:  "",
			genres: []string{"Comedy", "Adventure"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Comedy", "Adventure"}, movies)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 2,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'ID' with 'DESC'",
			title:  "",
			genres: []string{"Sci-Fi"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Sci-Fi"}, movies)
				expected = sortById(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 3,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'Title' with 'ASC'",
			title:  "",
			genres: []string{"Comedy", "Adventure"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "title",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Comedy", "Adventure"}, movies)
				expected = sortByTitle(expected, ASC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 2,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'Title' with 'DESC'",
			title:  "",
			genres: []string{"Sci-Fi"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-title",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Sci-Fi"}, movies)
				expected = sortByTitle(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 3,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'Year' with 'ASC'",
			title:  "",
			genres: []string{"Comedy", "Adventure"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "year",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Comedy", "Adventure"}, movies)
				expected = sortByYear(expected, ASC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 2,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'Year' with 'DESC'",
			title:  "",
			genres: []string{"Sci-Fi"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-year",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Sci-Fi"}, movies)
				expected = sortByYear(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 3,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'Runtime' with 'ASC'",
			title:  "",
			genres: []string{"Comedy", "Adventure"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Comedy", "Adventure"}, movies)
				expected = sortByRuntime(expected, ASC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 2,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Genres' and sorted by 'Runtime' with 'DESC'",
			title:  "",
			genres: []string{"Sci-Fi"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByGenres([]string{"Sci-Fi"}, movies)
				expected = sortByRuntime(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 3,
				}, meta)
			},
		},
		{
			name:   "Sort by 'ID' with 'ASC' 'Page=2'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         2,
				PageSize:     20,
				Sort:         "id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				require.NoError(t, err)
				assertMovies(t, movies[20:40], data)
				assertMetadata(t, Metadata{
					CurrentPage:  2,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Sort by 'ID' with 'DESC' 'Page=2'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         2,
				PageSize:     20,
				Sort:         "-id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := sortById(movies, DESC)

				require.NoError(t, err)
				assertMovies(t, expected[20:40], data)
				assertMetadata(t, Metadata{
					CurrentPage:  2,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     3,
					TotalRecords: 50,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and 'Genres' and sorted by 'Runtime' with 'ASC'",
			title:  titleSearchTerm,
			genres: []string{"Romance"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = filterByGenres([]string{"Romance"}, expected)
				expected = sortByRuntime(expected, ASC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 2,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' and 'Genres' and sorted by 'Runtime' with 'DESC'",
			title:  titleSearchTerm,
			genres: []string{"Sci-Fi", "Adventure"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "-runtime",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				expected := filterByTitle(titleSearchTerm, movies)
				expected = filterByGenres([]string{"Sci-Fi", "Adventure"}, expected)
				expected = sortByRuntime(expected, DESC)

				require.NoError(t, err)
				assertMovies(t, expected, data)
				assertMetadata(t, Metadata{
					CurrentPage:  1,
					PageSize:     20,
					FirstPage:    1,
					LastPage:     1,
					TotalRecords: 1,
				}, meta)
			},
		},
		{
			name:   "Filter by 'Title' with no results",
			title:  noResultsTitleSearchTerm,
			genres: []string{},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				require.NoError(t, err)
				require.Empty(t, data)
				require.Empty(t, meta)
			},
		},
		{
			name:   "Filter by 'Genres' with no results",
			title:  "",
			genres: []string{"Horror"},
			filters: Filters{
				Page:         1,
				PageSize:     20,
				Sort:         "id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				require.NoError(t, err)
				require.Empty(t, data)
				require.Empty(t, meta)
			},
		},
		{
			name:   "Sort by 'ID' with 'ASC' no results in 'PAGE=10'",
			title:  "",
			genres: []string{},
			filters: Filters{
				Page:         10,
				PageSize:     20,
				Sort:         "id",
				SortSafeList: MoviesSortSafeList,
			},
			assert: func(t *testing.T, data []Movie, meta Metadata, err error) {
				require.NoError(t, err)
				require.Empty(t, data)
				require.Empty(t, meta)
			},
		},
	}

	t.Cleanup(func() {
		teardown(t)
	})

	for _, tc := range testCases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				data, meta, err := testModels.Movies.GetAll(tc.title, tc.genres, tc.filters)

				tc.assert(t, data, meta, err)
			},
		)
	}
}

func TestGet(t *testing.T) {
	testModels := NewModels(testDB)

	t.Run("Successfully return movie data", func(t *testing.T) {
		createdMovie := createRandomMovie(t, &testModels)
		gotMovie, err := testModels.Movies.Get(createdMovie.ID)

		require.NoError(t, err)
		require.Equal(t, createdMovie.ID, gotMovie.ID)
		require.Equal(t, createdMovie.Title, gotMovie.Title)
		require.Equal(t, createdMovie.Year, gotMovie.Year)
		require.Equal(t, createdMovie.Runtime, gotMovie.Runtime)
		require.Equal(t, createdMovie.Version, gotMovie.Version)
		require.EqualValues(t, createdMovie.Genres, gotMovie.Genres)
		require.WithinDuration(t, createdMovie.CreatedAt, gotMovie.CreatedAt, time.Second)

		t.Cleanup(func() {
			teardown(t)
		})
	})

	t.Run("'ErrRecordNotFound' when given 'ID' does not exist", func(t *testing.T) {
		gotMovie, err := testModels.Movies.Get(gofakeit.Int64())

		require.ErrorIs(t, err, ErrRecordNotFound)
		require.Zero(t, gotMovie)
	})

	t.Run("'ErrRecordNotFound' when given 'ID' is lower than 1", func(t *testing.T) {
		gotMovie, err := testModels.Movies.Get(0)

		require.ErrorIs(t, err, ErrRecordNotFound)
		require.Zero(t, gotMovie)
	})
}

func TestDelete(t *testing.T) {
	testModels := NewModels(testDB)

	t.Run("Successfully delete a movie", func(t *testing.T) {
		movie := createRandomMovie(t, &testModels)

		err := testModels.Movies.Delete(movie.ID)

		require.NoError(t, err)

		gotMovie, err := testModels.Movies.Get(movie.ID)

		require.Zero(t, gotMovie)
		require.ErrorIs(t, err, ErrRecordNotFound)

		t.Cleanup(func() {
			teardown(t)
		})
	})

	t.Run("'ErrRecordNotFound' when given 'ID' does not exist", func(t *testing.T) {
		err := testModels.Movies.Delete(gofakeit.Int64())

		require.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("'ErrRecordNotFound' when given 'ID' is lower than 1", func(t *testing.T) {
		err := testModels.Movies.Delete(0)

		require.ErrorIs(t, err, ErrRecordNotFound)
	})
}

func TestUpdate(t *testing.T) {
	testModels := NewModels(testDB)

	t.Run("Successfully update a movie", func(t *testing.T) {
		movie := createRandomMovie(t, &testModels)

		movie.Title = "New Title"
		movie.Runtime = movie.Runtime - Runtime(1)
		movie.Year = movie.Year - 1

		err := testModels.Movies.Update(&movie)
		require.NoError(t, err)

		updatedMovie, err := testModels.Movies.Get(movie.ID)
		require.NoError(t, err)

		require.Equal(t, movie.Title, updatedMovie.Title)
		require.Equal(t, movie.Runtime, updatedMovie.Runtime)
		require.Equal(t, movie.Year, updatedMovie.Year)

		t.Cleanup(func() {
			teardown(t)
		})
	})

	t.Run("'ErrEditConflict' when given 'ID' does not exist", func(t *testing.T) {
		movie := Movie{
			ID:        gofakeit.Int64(),
			Title:     gofakeit.MovieName(),
			Genres:    []string{gofakeit.MovieGenre()},
			Year:      int32(gofakeit.Year()),
			Runtime:   Runtime(gofakeit.Number(90, 180)),
			Version:   gofakeit.Int32(),
			CreatedAt: gofakeit.Date(),
		}

		err := testModels.Movies.Update(&movie)
		require.ErrorIs(t, err, ErrEditConflict)
	})
}
