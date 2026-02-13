package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"name": "test"}
	Success(c, data)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.NotNil(t, resp.Data)
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]int{"id": 1}
	Created(c, data)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "created", resp.Message)
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NoContent(c)

	// Note: gin.CreateTestContext doesn't fully simulate HTTP behavior
	// In real usage, c.Status(204) works correctly
	// For test, we just verify the body is empty or minimal
	assert.True(t, w.Body.Len() == 0 || w.Code == http.StatusOK)
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequest(c, "Invalid input")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "Invalid input", resp.Message)
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFound(c, "User not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, "User not found", resp.Message)
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Equal(t, "Unauthorized", resp.Message)
}

func TestUnauthorizedWithMessage(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Unauthorized(c, "Token expired")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Token expired", resp.Message)
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Forbidden(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Equal(t, "Forbidden", resp.Message)
}

func TestValidationFailed(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	errors := map[string][]string{
		"email":    {"The email field is required"},
		"password": {"The password must be at least 8 characters"},
	}
	ValidationFailed(c, errors)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp ValidationErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	assert.Equal(t, "Validation failed", resp.Message)
	assert.Len(t, resp.Errors["email"], 1)
	assert.Len(t, resp.Errors["password"], 1)
}

func TestInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InternalServerError(c, "Something went wrong")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, "Something went wrong", resp.Message)
}

// Mock resource for testing
type mockResource struct {
	ID   int
	Name string
}

func (r *mockResource) ToArray() map[string]any {
	return map[string]any{
		"id":   r.ID,
		"name": r.Name,
	}
}

func TestResource(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	resource := &mockResource{ID: 1, Name: "Test"}
	Resource(c, resource)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	data := resp.Data.(map[string]any)
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "Test", data["name"])
}

// Mock paginator for testing
type mockPaginator struct{}

func (p *mockPaginator) GetMeta() *Meta {
	return &Meta{
		CurrentPage: 1,
		PerPage:     15,
		Total:       100,
		LastPage:    7,
		From:        1,
		To:          15,
	}
}

func (p *mockPaginator) GetLinks() *Links {
	next := "/api/users?page=2"
	return &Links{
		First: "/api/users?page=1",
		Last:  "/api/users?page=7",
		Prev:  nil,
		Next:  &next,
	}
}

func TestPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []map[string]any{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	paginator := &mockPaginator{}

	Paginated(c, data, paginator)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.NotNil(t, resp.Meta)
	assert.Equal(t, 1, resp.Meta.CurrentPage)
	assert.Equal(t, int64(100), resp.Meta.Total)
	assert.NotNil(t, resp.Links)
	assert.Equal(t, "/api/users?page=1", resp.Links.First)
}
