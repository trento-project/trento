package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPagination(t *testing.T) {
	p := NewPagination(10, 1, 10)
	assert.Equal(t, p.Items, 10)
	assert.Equal(t, p.Page, 1)
	assert.Equal(t, p.PerPage, 10)
	assert.Equal(t, p.Pages, 1)

	// Check Pages number
	p = NewPagination(10, 1, 5)
	assert.Equal(t, p.Pages, 2)

	p = NewPagination(21, 1, 10)
	assert.Equal(t, p.Pages, 3)

	p = NewPagination(25, 1, 5)
	assert.Equal(t, p.Pages, 5)

	// Check page sanity check
	p = NewPagination(25, -1, 5)
	assert.Equal(t, p.Page, 1)

	p = NewPagination(0, 1, 10)
	assert.Equal(t, p.Pages, 0)
	assert.Equal(t, p.Page, 1)

	p = NewPagination(25, 7, 5)
	assert.Equal(t, p.Page, 5)
}

func TestNewPaginationWithStrings(t *testing.T) {
	p := NewPaginationWithStrings(10, "1", "10")
	assert.Equal(t, p.Items, 10)
	assert.Equal(t, p.Page, 1)
	assert.Equal(t, p.PerPage, 10)
	assert.Equal(t, p.Pages, 1)

	p = NewPaginationWithStrings(10, "a", "b")
	assert.Equal(t, p.Items, 10)
	assert.Equal(t, p.Page, 1)
	assert.Equal(t, p.PerPage, 10)
	assert.Equal(t, p.Pages, 1)
}

func TestGetPages(t *testing.T) {
	p := NewPagination(111, 1, 10)
	pages := p.GetPages()

	expectedPages := []*Page{
		&Page{Index: 1, Active: true},
		&Page{Index: 2, Active: false},
		&Page{Index: 3, Active: false},
	}

	assert.ElementsMatch(t, expectedPages, pages)

	p = NewPagination(111, 2, 10)
	pages = p.GetPages()

	expectedPages = []*Page{
		&Page{Index: 1, Active: false},
		&Page{Index: 2, Active: true},
		&Page{Index: 3, Active: false},
		&Page{Index: 4, Active: false},
	}

	assert.ElementsMatch(t, expectedPages, pages)

	p = NewPagination(111, 3, 10)
	pages = p.GetPages()

	expectedPages = []*Page{
		&Page{Index: 1, Active: false},
		&Page{Index: 2, Active: false},
		&Page{Index: 3, Active: true},
		&Page{Index: 4, Active: false},
		&Page{Index: 5, Active: false},
	}

	assert.ElementsMatch(t, expectedPages, pages)

	p = NewPagination(111, 4, 10)
	pages = p.GetPages()

	expectedPages = []*Page{
		&Page{Index: 2, Active: false},
		&Page{Index: 3, Active: false},
		&Page{Index: 4, Active: true},
		&Page{Index: 5, Active: false},
		&Page{Index: 6, Active: false},
	}

	assert.ElementsMatch(t, expectedPages, pages)

	p = NewPagination(111, 10, 10)
	pages = p.GetPages()

	expectedPages = []*Page{
		&Page{Index: 8, Active: false},
		&Page{Index: 9, Active: false},
		&Page{Index: 10, Active: true},
		&Page{Index: 11, Active: false},
		&Page{Index: 12, Active: false},
	}

	assert.ElementsMatch(t, expectedPages, pages)

	p = NewPagination(111, 11, 10)
	pages = p.GetPages()

	expectedPages = []*Page{
		&Page{Index: 9, Active: false},
		&Page{Index: 10, Active: false},
		&Page{Index: 11, Active: true},
		&Page{Index: 12, Active: false},
	}

	assert.ElementsMatch(t, expectedPages, pages)

	p = NewPagination(111, 12, 10)
	pages = p.GetPages()

	expectedPages = []*Page{
		&Page{Index: 10, Active: false},
		&Page{Index: 11, Active: false},
		&Page{Index: 12, Active: true},
	}

	assert.ElementsMatch(t, expectedPages, pages)
}

func TestGetSliceNumbers(t *testing.T) {
	p := NewPagination(111, 4, 10)
	first, last := p.GetSliceNumbers()

	assert.Equal(t, 30, first)
	assert.Equal(t, 40, last)

	p = NewPagination(151, 4, 25)
	first, last = p.GetSliceNumbers()

	assert.Equal(t, 75, first)
	assert.Equal(t, 100, last)
}

func TestGetPerPages(t *testing.T) {
	p := &Pagination{}
	perPages := p.GetPerPages()
	assert.ElementsMatch(t, []int{10, 25, 50, 100}, perPages)
}
