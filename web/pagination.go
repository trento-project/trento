package web

import (
	"math"
	"strconv"
)

const (
	defaultPage    int = 1
	defaultPerPage int = 10
)

type Pagination struct {
	Items   int
	Page    int
	PerPage int
	Pages   int
}

type Page struct {
	Index  int
	Active bool
}

func pagesNumber(items, perPage int) int {
	return int((float64(items) + float64(perPage) - 1) / float64(perPage))
}

func NewPaginationWithStrings(items int, page, perPage string) *Pagination {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = defaultPage
	}

	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		perPageInt = defaultPerPage
	}

	return NewPagination(items, pageInt, perPageInt)
}

func NewPagination(items, page, perPage int) *Pagination {
	pNumber := pagesNumber(items, perPage)

	if page < 1 {
		page = 1
	} else if page > pNumber {
		page = pNumber
	}

	return &Pagination{
		Items:   items,
		Page:    page,
		PerPage: perPage,
		Pages:   pNumber,
	}
}

func (p *Pagination) GetPages() []*Page {
	pages := []*Page{}

	for i := 1; i <= p.Pages; i++ {
		newPage := &Page{Index: i, Active: false}
		if i == p.Page {
			newPage.Active = true
			pages = append(pages, newPage)
		} else if i >= p.Page-2 && i <= p.Page+2 {
			pages = append(pages, newPage)
		}

	}

	return pages
}

// Function to get the 1st and last indexes using the current pagination
func (p *Pagination) GetSliceNumbers() (int, int) {
	return (p.Page - 1) * p.PerPage, int(math.Min(float64(p.Items), float64(p.Page*p.PerPage)))
}

func (p *Pagination) GetPerPages() []int {
	return []int{10, 25, 50, 100}
}
