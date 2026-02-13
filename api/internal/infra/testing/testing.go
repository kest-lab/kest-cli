package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestCase provides a fluent API for HTTP testing
type TestCase struct {
	t       *testing.T
	router  *gin.Engine
	method  string
	path    string
	body    interface{}
	headers map[string]string
	cookies []*http.Cookie
}

// NewTestCase creates a new test case
func NewTestCase(t *testing.T, router *gin.Engine) *TestCase {
	gin.SetMode(gin.TestMode)
	return &TestCase{
		t:       t,
		router:  router,
		headers: make(map[string]string),
	}
}

// Get creates a GET request test
func (tc *TestCase) Get(path string) *TestCase {
	tc.method = http.MethodGet
	tc.path = path
	return tc
}

// Post creates a POST request test
func (tc *TestCase) Post(path string) *TestCase {
	tc.method = http.MethodPost
	tc.path = path
	return tc
}

// Put creates a PUT request test
func (tc *TestCase) Put(path string) *TestCase {
	tc.method = http.MethodPut
	tc.path = path
	return tc
}

// Patch creates a PATCH request test
func (tc *TestCase) Patch(path string) *TestCase {
	tc.method = http.MethodPatch
	tc.path = path
	return tc
}

// Delete creates a DELETE request test
func (tc *TestCase) Delete(path string) *TestCase {
	tc.method = http.MethodDelete
	tc.path = path
	return tc
}

// WithBody sets the request body
func (tc *TestCase) WithBody(body interface{}) *TestCase {
	tc.body = body
	return tc
}

// WithJSON sets a JSON body
func (tc *TestCase) WithJSON(data interface{}) *TestCase {
	tc.body = data
	tc.headers["Content-Type"] = "application/json"
	return tc
}

// WithHeader sets a header
func (tc *TestCase) WithHeader(key, value string) *TestCase {
	tc.headers[key] = value
	return tc
}

// WithHeaders sets multiple headers
func (tc *TestCase) WithHeaders(headers map[string]string) *TestCase {
	for k, v := range headers {
		tc.headers[k] = v
	}
	return tc
}

// WithToken sets a Bearer token
func (tc *TestCase) WithToken(token string) *TestCase {
	tc.headers["Authorization"] = "Bearer " + token
	return tc
}

// WithAPIKey sets an API key header
func (tc *TestCase) WithAPIKey(key string) *TestCase {
	tc.headers["X-API-Key"] = key
	return tc
}

// WithCookie adds a cookie
func (tc *TestCase) WithCookie(cookie *http.Cookie) *TestCase {
	tc.cookies = append(tc.cookies, cookie)
	return tc
}

// Call executes the request and returns a TestResponse
func (tc *TestCase) Call() *TestResponse {
	var bodyReader io.Reader
	if tc.body != nil {
		switch v := tc.body.(type) {
		case string:
			bodyReader = bytes.NewBufferString(v)
		case []byte:
			bodyReader = bytes.NewBuffer(v)
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				tc.t.Fatalf("Failed to marshal body: %v", err)
			}
			bodyReader = bytes.NewBuffer(jsonBytes)
		}
	}

	req := httptest.NewRequest(tc.method, tc.path, bodyReader)

	for k, v := range tc.headers {
		req.Header.Set(k, v)
	}

	for _, cookie := range tc.cookies {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	tc.router.ServeHTTP(w, req)

	return &TestResponse{
		t:        tc.t,
		response: w,
	}
}

// TestResponse provides assertions for HTTP responses
type TestResponse struct {
	t        *testing.T
	response *httptest.ResponseRecorder
}

// AssertStatus asserts the response status code
func (r *TestResponse) AssertStatus(status int) *TestResponse {
	if r.response.Code != status {
		r.t.Errorf("Expected status %d, got %d", status, r.response.Code)
	}
	return r
}

// AssertOk asserts a 200 status
func (r *TestResponse) AssertOk() *TestResponse {
	return r.AssertStatus(http.StatusOK)
}

// AssertCreated asserts a 201 status
func (r *TestResponse) AssertCreated() *TestResponse {
	return r.AssertStatus(http.StatusCreated)
}

// AssertNoContent asserts a 204 status
func (r *TestResponse) AssertNoContent() *TestResponse {
	return r.AssertStatus(http.StatusNoContent)
}

// AssertBadRequest asserts a 400 status
func (r *TestResponse) AssertBadRequest() *TestResponse {
	return r.AssertStatus(http.StatusBadRequest)
}

// AssertUnauthorized asserts a 401 status
func (r *TestResponse) AssertUnauthorized() *TestResponse {
	return r.AssertStatus(http.StatusUnauthorized)
}

// AssertForbidden asserts a 403 status
func (r *TestResponse) AssertForbidden() *TestResponse {
	return r.AssertStatus(http.StatusForbidden)
}

// AssertNotFound asserts a 404 status
func (r *TestResponse) AssertNotFound() *TestResponse {
	return r.AssertStatus(http.StatusNotFound)
}

// AssertUnprocessable asserts a 422 status
func (r *TestResponse) AssertUnprocessable() *TestResponse {
	return r.AssertStatus(http.StatusUnprocessableEntity)
}

// AssertServerError asserts a 500 status
func (r *TestResponse) AssertServerError() *TestResponse {
	return r.AssertStatus(http.StatusInternalServerError)
}

// AssertJSON asserts the response is JSON
func (r *TestResponse) AssertJSON() *TestResponse {
	contentType := r.response.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" && contentType != "application/json" {
		r.t.Errorf("Expected JSON content type, got %s", contentType)
	}
	return r
}

// AssertJSONPath asserts a JSON path has the expected value.
// Supports nested paths using dot notation (e.g., "data.user.name")
func (r *TestResponse) AssertJSONPath(path string, expected interface{}) *TestResponse {
	var data map[string]interface{}
	if err := json.Unmarshal(r.response.Body.Bytes(), &data); err != nil {
		r.t.Fatalf("Failed to parse JSON: %v", err)
	}

	actual := getNestedValue(data, path)
	if actual == nil {
		r.t.Errorf("JSON path '%s' not found", path)
		return r
	}

	if actual != expected {
		r.t.Errorf("Expected %v at path '%s', got %v", expected, path, actual)
	}
	return r
}

// AssertJSONStructure asserts the response contains the expected JSON keys.
// Supports nested paths using dot notation (e.g., "data.access_token")
func (r *TestResponse) AssertJSONStructure(keys []string) *TestResponse {
	var data map[string]interface{}
	if err := json.Unmarshal(r.response.Body.Bytes(), &data); err != nil {
		r.t.Fatalf("Failed to parse JSON: %v", err)
	}

	for _, key := range keys {
		if getNestedValue(data, key) == nil {
			r.t.Errorf("Expected JSON key '%s' not found", key)
		}
	}
	return r
}

// getNestedValue retrieves a value from a nested map using dot notation
func getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
			if current == nil {
				return nil
			}
		} else {
			return nil
		}
	}
	return current
}

// AssertHeader asserts a header value
func (r *TestResponse) AssertHeader(key, value string) *TestResponse {
	actual := r.response.Header().Get(key)
	if actual != value {
		r.t.Errorf("Expected header %s=%s, got %s", key, value, actual)
	}
	return r
}

// AssertHeaderExists asserts a header exists
func (r *TestResponse) AssertHeaderExists(key string) *TestResponse {
	if r.response.Header().Get(key) == "" {
		r.t.Errorf("Expected header %s to exist", key)
	}
	return r
}

// AssertCookie asserts a cookie exists
func (r *TestResponse) AssertCookie(name string) *TestResponse {
	cookies := r.response.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return r
		}
	}
	r.t.Errorf("Expected cookie '%s' not found", name)
	return r
}

// AssertBodyContains asserts the body contains a string
func (r *TestResponse) AssertBodyContains(substring string) *TestResponse {
	body := r.response.Body.String()
	if !bytes.Contains([]byte(body), []byte(substring)) {
		r.t.Errorf("Expected body to contain '%s'", substring)
	}
	return r
}

// JSON returns the response body as a map
func (r *TestResponse) JSON() map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal(r.response.Body.Bytes(), &data); err != nil {
		r.t.Fatalf("Failed to parse JSON: %v", err)
	}
	return data
}

// Body returns the response body as string
func (r *TestResponse) Body() string {
	return r.response.Body.String()
}

// Status returns the response status code
func (r *TestResponse) Status() int {
	return r.response.Code
}

// Dump prints the response for debugging
func (r *TestResponse) Dump() *TestResponse {
	r.t.Logf("Status: %d", r.response.Code)
	r.t.Logf("Headers: %v", r.response.Header())
	r.t.Logf("Body: %s", r.response.Body.String())
	return r
}
