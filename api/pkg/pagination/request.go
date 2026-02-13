package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Default pagination settings
const (
	DefaultPage    = 1
	DefaultPerPage = 15
	MaxPerPage     = 100
)

// Request represents pagination parameters from the client.
// Supports both query string and JSON body binding.
//
// Query Parameters:
//   - page: Current page number (default: 1)
//   - per_page: Items per page (default: 15, max: 100)
//   - keyword: Optional search keyword
//   - sort: Sort field (e.g., "created_at")
//   - order: Sort order ("asc" or "desc")
//
// Example URL: /api/users?page=2&per_page=20&keyword=john&sort=created_at&order=desc
type Request struct {
	Page    int    `form:"page" json:"page"`
	PerPage int    `form:"per_page" json:"per_page"`
	Keyword string `form:"keyword" json:"keyword"`
	Sort    string `form:"sort" json:"sort"`
	Order   string `form:"order" json:"order"`
}

// GetPage returns the current page, minimum 1.
func (r *Request) GetPage() int {
	if r.Page < 1 {
		return DefaultPage
	}
	return r.Page
}

// GetPerPage returns items per page, with min/max bounds.
func (r *Request) GetPerPage() int {
	if r.PerPage < 1 {
		return DefaultPerPage
	}
	if r.PerPage > MaxPerPage {
		return MaxPerPage
	}
	return r.PerPage
}

// GetPageSize is an alias for GetPerPage for backward compatibility.
func (r *Request) GetPageSize() int {
	return r.GetPerPage()
}

// GetOffset calculates the SQL offset for the current page.
func (r *Request) GetOffset() int {
	return (r.GetPage() - 1) * r.GetPerPage()
}

// GetSort returns the sort field, or empty string if not set.
func (r *Request) GetSort() string {
	return r.Sort
}

// GetOrder returns the sort order, defaults to "desc".
func (r *Request) GetOrder() string {
	if r.Order == "asc" {
		return "asc"
	}
	return "desc"
}

// GetOrderBy returns the full ORDER BY clause.
// Returns empty string if no sort field is specified.
//
// Example:
//
//	orderBy := req.GetOrderBy() // "created_at desc"
//	db.Order(orderBy)
func (r *Request) GetOrderBy() string {
	if r.Sort == "" {
		return ""
	}
	return r.Sort + " " + r.GetOrder()
}

// HasKeyword returns true if a search keyword is provided.
func (r *Request) HasKeyword() bool {
	return r.Keyword != ""
}

// FromContext extracts pagination request from Gin context.
// Reads from query string parameters.
//
// Example:
//
//	func (h *Handler) List(c *gin.Context) {
//	    req := pagination.FromContext(c)
//	    // req.Page, req.PerPage, req.Keyword, etc.
//	}
func FromContext(c *gin.Context) *Request {
	req := &Request{}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		req.Page = page
	}

	if perPage, err := strconv.Atoi(c.DefaultQuery("per_page", strconv.Itoa(DefaultPerPage))); err == nil {
		req.PerPage = perPage
	}

	req.Keyword = c.Query("keyword")
	req.Sort = c.Query("sort")
	req.Order = c.DefaultQuery("order", "desc")

	return req
}

// FromQuery creates pagination request from a query map.
// Useful for testing or non-Gin contexts.
func FromQuery(query map[string][]string) *Request {
	req := &Request{}

	if pages, ok := query["page"]; ok && len(pages) > 0 {
		if page, err := strconv.Atoi(pages[0]); err == nil {
			req.Page = page
		}
	}

	if perPages, ok := query["per_page"]; ok && len(perPages) > 0 {
		if perPage, err := strconv.Atoi(perPages[0]); err == nil {
			req.PerPage = perPage
		}
	}

	if keywords, ok := query["keyword"]; ok && len(keywords) > 0 {
		req.Keyword = keywords[0]
	}

	if sorts, ok := query["sort"]; ok && len(sorts) > 0 {
		req.Sort = sorts[0]
	}

	if orders, ok := query["order"]; ok && len(orders) > 0 {
		req.Order = orders[0]
	}

	return req
}

// NewRequest creates a new pagination request with specified values.
func NewRequest(page, perPage int) *Request {
	return &Request{
		Page:    page,
		PerPage: perPage,
	}
}

// WithKeyword sets the search keyword.
func (r *Request) WithKeyword(keyword string) *Request {
	r.Keyword = keyword
	return r
}

// WithSort sets the sort field and order.
func (r *Request) WithSort(field, order string) *Request {
	r.Sort = field
	r.Order = order
	return r
}
