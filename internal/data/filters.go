package data

import (
	"math"
	"strings"

	"github.com/brGuirra/greenlight/internal/validator"
)

type Filters struct {
	Sort         string
	SortSafeList []string
	Page         int
	PageSize     int
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offsett() int {
	return (f.Page - 1) * f.PageSize
}

// TODO: Test at handler level
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 milion")
	v.Check(f.PageSize > 0, "pageSize", "must be greater than zero")
	v.Check(f.PageSize <= 100, "pageSize", "must be a maximum of 100")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

type Metadata struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
