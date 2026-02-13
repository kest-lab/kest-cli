package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginator(t *testing.T) {
	items := []string{"a", "b", "c"}
	paginator := NewPaginator(items, 100, 1, 15)

	assert.Equal(t, 3, len(paginator.Items()))
	assert.Equal(t, int64(100), paginator.Total())
	assert.Equal(t, 1, paginator.CurrentPage())
	assert.Equal(t, 15, paginator.PerPage())
	assert.Equal(t, 7, paginator.LastPage())
}

func TestPaginatorBounds(t *testing.T) {
	items := []string{"a"}

	// Test page < 1
	p := NewPaginator(items, 100, 0, 15)
	assert.Equal(t, 1, p.CurrentPage())

	// Test perPage < 1
	p = NewPaginator(items, 100, 1, 0)
	assert.Equal(t, DefaultPerPage, p.PerPage())

	// Test perPage > max
	p = NewPaginator(items, 100, 1, 200)
	assert.Equal(t, MaxPerPage, p.PerPage())
}

func TestPaginatorFromTo(t *testing.T) {
	items := []string{"a", "b", "c"}

	// First page
	p := NewPaginator(items, 100, 1, 15)
	assert.Equal(t, 1, p.From())
	assert.Equal(t, 15, p.To())

	// Second page
	p = NewPaginator(items, 100, 2, 15)
	assert.Equal(t, 16, p.From())
	assert.Equal(t, 30, p.To())

	// Last page (partial)
	p = NewPaginator(items, 23, 2, 15)
	assert.Equal(t, 16, p.From())
	assert.Equal(t, 23, p.To())

	// Empty result
	p = NewPaginator([]string{}, 0, 1, 15)
	assert.Equal(t, 0, p.From())
	assert.Equal(t, 0, p.To())
}

func TestPaginatorNavigation(t *testing.T) {
	items := []string{"a"}

	// First page
	p := NewPaginator(items, 100, 1, 15)
	assert.True(t, p.OnFirstPage())
	assert.False(t, p.OnLastPage())
	assert.True(t, p.HasMorePages())
	assert.True(t, p.HasPages())

	// Last page
	p = NewPaginator(items, 100, 7, 15)
	assert.False(t, p.OnFirstPage())
	assert.True(t, p.OnLastPage())
	assert.False(t, p.HasMorePages())

	// Single page
	p = NewPaginator(items, 10, 1, 15)
	assert.True(t, p.OnFirstPage())
	assert.True(t, p.OnLastPage())
	assert.False(t, p.HasMorePages())
	assert.False(t, p.HasPages())
}

func TestPaginatorURLs(t *testing.T) {
	items := []string{"a"}
	p := NewPaginator(items, 100, 3, 15)
	p.SetPath("/api/users")

	// Basic URLs
	assert.Equal(t, "/api/users?page=1", p.FirstPageURL())
	assert.Equal(t, "/api/users?page=7", p.LastPageURL())
	assert.Equal(t, "/api/users?page=3", p.URL(3))

	// Previous/Next
	prev := p.PreviousPageURL()
	assert.NotNil(t, prev)
	assert.Equal(t, "/api/users?page=2", *prev)

	next := p.NextPageURL()
	assert.NotNil(t, next)
	assert.Equal(t, "/api/users?page=4", *next)

	// First page has no previous
	p = NewPaginator(items, 100, 1, 15)
	p.SetPath("/api/users")
	assert.Nil(t, p.PreviousPageURL())

	// Last page has no next
	p = NewPaginator(items, 100, 7, 15)
	p.SetPath("/api/users")
	assert.Nil(t, p.NextPageURL())
}

func TestPaginatorWithQuery(t *testing.T) {
	items := []string{"a"}
	p := NewPaginator(items, 100, 1, 15)
	p.SetPath("/api/users")
	p.Append("status", "active")
	p.Append("sort", "name")

	url := p.URL(2)
	assert.Contains(t, url, "page=2")
	assert.Contains(t, url, "status=active")
	assert.Contains(t, url, "sort=name")
}

func TestPaginatorMeta(t *testing.T) {
	items := []string{"a", "b", "c"}
	p := NewPaginator(items, 100, 2, 15)

	meta := p.GetMeta()
	assert.Equal(t, 2, meta.CurrentPage)
	assert.Equal(t, 15, meta.PerPage)
	assert.Equal(t, int64(100), meta.Total)
	assert.Equal(t, 7, meta.LastPage)
	assert.Equal(t, 16, meta.From)
	assert.Equal(t, 30, meta.To)
}

func TestPaginatorLinks(t *testing.T) {
	items := []string{"a"}
	p := NewPaginator(items, 100, 3, 15)
	p.SetPath("/api/users")

	links := p.GetLinks()
	assert.Equal(t, "/api/users?page=1", links.First)
	assert.Equal(t, "/api/users?page=7", links.Last)
	assert.NotNil(t, links.Prev)
	assert.Equal(t, "/api/users?page=2", *links.Prev)
	assert.NotNil(t, links.Next)
	assert.Equal(t, "/api/users?page=4", *links.Next)
}

func TestPaginatorPageLinks(t *testing.T) {
	items := []string{"a"}

	// Small number of pages - show all
	p := NewPaginator(items, 50, 3, 10)
	p.SetPath("/api/users")
	links := p.Links()

	// Should have: Previous, 1, 2, 3, 4, 5, Next
	assert.Equal(t, 7, len(links))
	assert.Equal(t, "Previous", links[0].Label)
	assert.Equal(t, "1", links[1].Label)
	assert.Equal(t, "3", links[3].Label)
	assert.True(t, links[3].Active) // Page 3 is active
	assert.Equal(t, "Next", links[6].Label)

	// Large number of pages - show window
	p = NewPaginator(items, 1000, 50, 10)
	p.SetPath("/api/users")
	links = p.Links()

	// Should have ellipsis
	hasEllipsis := false
	for _, link := range links {
		if link.Label == "..." {
			hasEllipsis = true
			break
		}
	}
	assert.True(t, hasEllipsis)
}

func TestRequest(t *testing.T) {
	req := &Request{Page: 2, PerPage: 20}

	assert.Equal(t, 2, req.GetPage())
	assert.Equal(t, 20, req.GetPerPage())
	assert.Equal(t, 20, req.GetOffset())

	// Defaults
	req = &Request{}
	assert.Equal(t, DefaultPage, req.GetPage())
	assert.Equal(t, DefaultPerPage, req.GetPerPage())

	// Bounds
	req = &Request{Page: -1, PerPage: 200}
	assert.Equal(t, DefaultPage, req.GetPage())
	assert.Equal(t, MaxPerPage, req.GetPerPage())
}

func TestRequestSort(t *testing.T) {
	req := &Request{Sort: "created_at", Order: "asc"}
	assert.Equal(t, "created_at asc", req.GetOrderBy())

	req = &Request{Sort: "created_at", Order: "desc"}
	assert.Equal(t, "created_at desc", req.GetOrderBy())

	req = &Request{Sort: "created_at"} // No order specified
	assert.Equal(t, "created_at desc", req.GetOrderBy())

	req = &Request{} // No sort
	assert.Equal(t, "", req.GetOrderBy())
}

func TestFromQuery(t *testing.T) {
	query := map[string][]string{
		"page":     {"3"},
		"per_page": {"25"},
		"keyword":  {"test"},
		"sort":     {"name"},
		"order":    {"asc"},
	}

	req := FromQuery(query)
	assert.Equal(t, 3, req.Page)
	assert.Equal(t, 25, req.PerPage)
	assert.Equal(t, "test", req.Keyword)
	assert.Equal(t, "name", req.Sort)
	assert.Equal(t, "asc", req.Order)
}

func TestNewRequest(t *testing.T) {
	req := NewRequest(2, 20).
		WithKeyword("search").
		WithSort("created_at", "desc")

	assert.Equal(t, 2, req.Page)
	assert.Equal(t, 20, req.PerPage)
	assert.Equal(t, "search", req.Keyword)
	assert.Equal(t, "created_at", req.Sort)
	assert.Equal(t, "desc", req.Order)
}
