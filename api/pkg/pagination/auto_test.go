package pagination

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func createTestContext(path string, query string) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		URL: &url.URL{Path: path, RawQuery: query},
	}
	return c
}

func TestResult_NilItems(t *testing.T) {
	// Test that nil items don't cause panic
	var nilItems []*string
	result := NewResult(nilItems, 0, 1, 15)

	assert.NotNil(t, result)
	assert.Nil(t, result.Items)
	assert.Equal(t, int64(0), result.Total)

	// ToPaginator should not panic
	c := createTestContext("/api/users", "")
	paginator := result.ToPaginator(c)

	assert.NotNil(t, paginator)
	assert.Equal(t, 0, len(paginator.Items()))
}

func TestResult_EmptyItems(t *testing.T) {
	items := make([]*string, 0)
	result := NewResult(items, 0, 1, 15)

	c := createTestContext("/api/users", "")
	paginator := result.ToPaginator(c)

	assert.NotNil(t, paginator)
	assert.Equal(t, 0, len(paginator.Items()))
	assert.Equal(t, int64(0), paginator.Total())
}

func TestResult_ToPaginator_PreservesPath(t *testing.T) {
	items := []string{"a", "b", "c"}
	result := NewResult(items, 100, 2, 15)

	c := createTestContext("/api/v1/users", "keyword=test&status=active")
	paginator := result.ToPaginator(c)

	// Check path is set
	assert.Contains(t, paginator.FirstPageURL(), "/api/v1/users")

	// Check query params are preserved (except page)
	assert.Contains(t, paginator.FirstPageURL(), "keyword=test")
	assert.Contains(t, paginator.FirstPageURL(), "status=active")
}

func TestResult_ToPaginator_PageParams(t *testing.T) {
	items := []string{"a", "b", "c"}
	result := NewResult(items, 100, 2, 15)

	c := createTestContext("/api/users", "page=2&per_page=15")
	paginator := result.ToPaginator(c)

	// page and per_page should NOT be duplicated in query
	meta := paginator.GetMeta()
	assert.Equal(t, 2, meta.CurrentPage)
	assert.Equal(t, 15, meta.PerPage)
}

func TestResult_InvalidPage(t *testing.T) {
	items := []string{"a", "b", "c"}

	// Page 0 should be corrected to 1
	result := NewResult(items, 100, 0, 15)
	c := createTestContext("/api/users", "")
	paginator := result.ToPaginator(c)

	assert.Equal(t, 1, paginator.CurrentPage())
}

func TestResult_InvalidPerPage(t *testing.T) {
	items := []string{"a", "b", "c"}

	// PerPage 0 should be corrected to default
	result := NewResult(items, 100, 1, 0)
	c := createTestContext("/api/users", "")
	paginator := result.ToPaginator(c)

	assert.Equal(t, DefaultPerPage, paginator.PerPage())
}

func TestResult_ExceedsMaxPerPage(t *testing.T) {
	items := []string{"a", "b", "c"}

	// PerPage > MaxPerPage should be capped
	result := NewResult(items, 100, 1, 1000)
	c := createTestContext("/api/users", "")
	paginator := result.ToPaginator(c)

	assert.Equal(t, MaxPerPage, paginator.PerPage())
}

func TestPaginator_GetItems_ReturnsCorrectType(t *testing.T) {
	items := []string{"a", "b", "c"}
	paginator := NewPaginator(items, 3, 1, 15)

	// GetItems should return the same slice
	gotItems := paginator.GetItems()
	assert.Equal(t, items, gotItems)
}

func TestPaginator_EmptyTotal_FromTo(t *testing.T) {
	items := []string{}
	paginator := NewPaginator(items, 0, 1, 15)

	// From and To should be 0 when total is 0
	assert.Equal(t, 0, paginator.From())
	assert.Equal(t, 0, paginator.To())
}

func TestPaginator_Links_NilPrevOnFirstPage(t *testing.T) {
	items := []string{"a", "b", "c"}
	paginator := NewPaginator(items, 100, 1, 15)
	paginator.SetPath("/api/users")

	links := paginator.GetLinks()
	assert.Nil(t, links.Prev)
	assert.NotNil(t, links.Next)
}

func TestPaginator_Links_NilNextOnLastPage(t *testing.T) {
	items := []string{"a", "b", "c"}
	paginator := NewPaginator(items, 30, 2, 15) // 30 items, 15 per page = 2 pages
	paginator.SetPath("/api/users")

	links := paginator.GetLinks()
	assert.NotNil(t, links.Prev)
	assert.Nil(t, links.Next)
}
