package filters

import (
	"strings"

	"github.com/hvpaiva/greenlight/pkg/validator"
)

type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type Metadata struct {
	Page  Page `json:"page"`
	Count int  `json:"count"`
}

type Page struct {
	Current int `json:"current"`
	Size    int `json:"size"`
	First   int `json:"first"`
	Last    int `json:"last"`
}

func ZeroValueMetadata() Metadata {
	return Metadata{
		Page: Page{},
	}
}

func (f Filter) Validate(v *validator.Validator) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.Permitted(f.Sort, f.SortSafeList...), "sort", "invalid sort value")

}

func (f Filter) SortColumn() string {
	for _, safe := range f.SortSafeList {
		if f.Sort == safe {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort param: " + f.Sort)
}

func (f Filter) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func CalculateMetadata(total, page, size int) Metadata {
	if total == 0 {
		return ZeroValueMetadata()
	}

	return Metadata{
		Page: Page{
			Current: page,
			Size:    size,
			First:   1,
			Last:    (total + size - 1) / size,
		},
		Count: total,
	}
}
