package request

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements error interface
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var msgs []string
	for _, e := range ve {
		msgs = append(msgs, e.Field+": "+e.Message)
	}
	return strings.Join(msgs, "; ")
}

// First returns the first error message
func (ve ValidationErrors) First() string {
	if len(ve) == 0 {
		return ""
	}
	return ve[0].Message
}

// Has checks if a field has validation errors
func (ve ValidationErrors) Has(field string) bool {
	for _, e := range ve {
		if e.Field == field {
			return true
		}
	}
	return false
}

// Get returns all errors for a field
func (ve ValidationErrors) Get(field string) []ValidationError {
	var errs []ValidationError
	for _, e := range ve {
		if e.Field == field {
			errs = append(errs, e)
		}
	}
	return errs
}

// ToMap converts errors to a map keyed by field
func (ve ValidationErrors) ToMap() map[string][]string {
	result := make(map[string][]string)
	for _, e := range ve {
		result[e.Field] = append(result[e.Field], e.Message)
	}
	return result
}

// Validate binds and validates the request body into the given struct
func Validate(c *gin.Context, req interface{}) error {
	// Bind JSON
	if err := c.ShouldBindJSON(req); err != nil {
		return formatBindingError(err)
	}
	return nil
}

// ValidateQuery binds and validates query parameters
func ValidateQuery(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return formatBindingError(err)
	}
	return nil
}

// ValidateURI binds and validates URI parameters
func ValidateURI(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindUri(req); err != nil {
		return formatBindingError(err)
	}
	return nil
}

// ValidateForm binds and validates form data
func ValidateForm(c *gin.Context, req interface{}) error {
	if err := c.ShouldBind(req); err != nil {
		return formatBindingError(err)
	}
	return nil
}

// formatBindingError converts validator errors to ValidationErrors
func formatBindingError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := make(ValidationErrors, 0, len(ve))
		for _, fe := range ve {
			errs = append(errs, ValidationError{
				Field:   toSnakeCase(fe.Field()),
				Message: getErrorMessage(fe),
				Tag:     fe.Tag(),
				Value:   fe.Param(),
			})
		}
		return errs
	}
	// Return generic error wrapped in ValidationErrors
	return ValidationErrors{{
		Field:   "",
		Message: err.Error(),
		Tag:     "binding",
	}}
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(fe validator.FieldError) string {
	field := toSnakeCase(fe.Field())

	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "len":
		return field + " must be exactly " + fe.Param() + " characters"
	case "gte":
		return field + " must be greater than or equal to " + fe.Param()
	case "gt":
		return field + " must be greater than " + fe.Param()
	case "lte":
		return field + " must be less than or equal to " + fe.Param()
	case "lt":
		return field + " must be less than " + fe.Param()
	case "eq":
		return field + " must equal " + fe.Param()
	case "ne":
		return field + " must not equal " + fe.Param()
	case "oneof":
		return field + " must be one of: " + fe.Param()
	case "url":
		return field + " must be a valid URL"
	case "uuid":
		return field + " must be a valid UUID"
	case "alpha":
		return field + " must contain only alphabetic characters"
	case "alphanum":
		return field + " must contain only alphanumeric characters"
	case "numeric":
		return field + " must be a valid number"
	case "eqfield":
		return field + " must match " + toSnakeCase(fe.Param())
	default:
		return field + " failed validation: " + fe.Tag()
	}
}

// toSnakeCase converts PascalCase/camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// Request wraps gin.Context with additional helpers
type Request struct {
	*gin.Context
}

// New creates a new Request wrapper
func New(c *gin.Context) *Request {
	return &Request{Context: c}
}

// Input gets an input value from query, form, or JSON body
func (r *Request) Input(key string, defaultValue ...string) string {
	// Try query first
	if val := r.Query(key); val != "" {
		return val
	}

	// Try form
	if val := r.PostForm(key); val != "" {
		return val
	}

	// Try JSON body
	var body map[string]interface{}
	if err := r.ShouldBindJSON(&body); err == nil {
		if val, ok := body[key]; ok {
			if s, ok := val.(string); ok {
				return s
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// InputInt gets an integer input value
func (r *Request) InputInt(key string, defaultValue ...int) int {
	val := r.Input(key)
	if val == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return i
}

// InputBool gets a boolean input value
func (r *Request) InputBool(key string, defaultValue ...bool) bool {
	val := r.Input(key)
	if val == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return b
}

// Has checks if an input key exists
func (r *Request) Has(keys ...string) bool {
	for _, key := range keys {
		if r.Input(key) == "" {
			return false
		}
	}
	return true
}

// Filled checks if an input key exists and is not empty
func (r *Request) Filled(key string) bool {
	return strings.TrimSpace(r.Input(key)) != ""
}

// Missing checks if an input key is missing
func (r *Request) Missing(key string) bool {
	return !r.Has(key)
}

// Only returns only the specified keys from input
func (r *Request) Only(keys ...string) map[string]string {
	result := make(map[string]string)
	for _, key := range keys {
		result[key] = r.Input(key)
	}
	return result
}

// Except returns all input except the specified keys
func (r *Request) Except(keys ...string) map[string]string {
	excluded := make(map[string]bool)
	for _, k := range keys {
		excluded[k] = true
	}

	result := make(map[string]string)

	// Collect from query
	for k, v := range r.Request.URL.Query() {
		if !excluded[k] && len(v) > 0 {
			result[k] = v[0]
		}
	}

	// Collect from form
	if err := r.Request.ParseForm(); err == nil {
		for k, v := range r.Request.PostForm {
			if !excluded[k] && len(v) > 0 {
				result[k] = v[0]
			}
		}
	}

	return result
}

// Merge merges additional data with input
func (r *Request) Merge(data map[string]string) map[string]string {
	result := make(map[string]string)

	// Get all current input
	for k, v := range r.Request.URL.Query() {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}

	// Merge provided data
	for k, v := range data {
		result[k] = v
	}

	return result
}

// Boolean returns boolean value from input with multiple true representations
func (r *Request) Boolean(key string) bool {
	val := strings.ToLower(r.Input(key))
	return val == "1" || val == "true" || val == "on" || val == "yes"
}

// IP returns the client IP address
func (r *Request) IP() string {
	return r.ClientIP()
}

// UserAgent returns the User-Agent header
func (r *Request) UserAgent() string {
	return r.GetHeader("User-Agent")
}

// IsAjax checks if the request is an AJAX request
func (r *Request) IsAjax() bool {
	return r.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

// IsJSON checks if the request expects JSON response
func (r *Request) IsJSON() bool {
	accept := r.GetHeader("Accept")
	return strings.Contains(accept, "application/json")
}

// WantsJSON checks if the request wants a JSON response
func (r *Request) WantsJSON() bool {
	return r.IsJSON() || r.IsAjax()
}

// IsMethod checks the request method
func (r *Request) IsMethod(method string) bool {
	return r.Request.Method == strings.ToUpper(method)
}

// IsGet checks if request method is GET
func (r *Request) IsGet() bool {
	return r.IsMethod(http.MethodGet)
}

// IsPost checks if request method is POST
func (r *Request) IsPost() bool {
	return r.IsMethod(http.MethodPost)
}

// IsPut checks if request method is PUT
func (r *Request) IsPut() bool {
	return r.IsMethod(http.MethodPut)
}

// IsDelete checks if request method is DELETE
func (r *Request) IsDelete() bool {
	return r.IsMethod(http.MethodDelete)
}

// Fingerprint returns a unique request fingerprint
func (r *Request) Fingerprint() string {
	return r.IP() + "|" + r.UserAgent() + "|" + r.Request.URL.Path
}

// BearerToken extracts the bearer token from the Authorization header
func (r *Request) BearerToken() string {
	auth := r.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:]
	}
	return ""
}

// ValidatedData stores validated request data in context
func ValidatedData(c *gin.Context, key string, data interface{}) {
	c.Set("validated:"+key, data)
}

// GetValidated retrieves validated data from context
func GetValidated[T any](c *gin.Context, key string) (T, bool) {
	val, exists := c.Get("validated:" + key)
	if !exists {
		var zero T
		return zero, false
	}
	if typed, ok := val.(T); ok {
		return typed, true
	}
	var zero T
	return zero, false
}

// Struct2Map converts a struct to map[string]interface{}
func Struct2Map(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Get json tag name or use field name
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = toSnakeCase(field.Name)
		} else {
			// Handle json tag with options like `json:"name,omitempty"`
			if idx := strings.Index(name, ","); idx != -1 {
				name = name[:idx]
			}
		}

		if value.CanInterface() {
			result[name] = value.Interface()
		}
	}

	return result
}
