package validation

import (
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with custom rules and messages
type Validator struct {
	v        *validator.Validate
	mu       sync.RWMutex
	messages map[string]string
}

var (
	globalValidator *Validator
	once            sync.Once
)

// Global returns the global validator instance
func Global() *Validator {
	once.Do(func() {
		globalValidator = New()
	})
	return globalValidator
}

// New creates a new Validator instance
func New() *Validator {
	v := validator.New()

	// Use JSON tag names for field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		if name == "" {
			return toSnakeCase(fld.Name)
		}
		return name
	})

	val := &Validator{
		v:        v,
		messages: make(map[string]string),
	}

	// Register default custom rules
	val.registerDefaultRules()

	return val
}

// registerDefaultRules registers commonly used custom validation rules
func (v *Validator) registerDefaultRules() {
	// Phone number validation (simple)
	v.RegisterRule("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
		return matched
	}, "must be a valid phone number")

	// International phone
	v.RegisterRule("phone_intl", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
		return matched
	}, "must be a valid international phone number")

	// Username (alphanumeric with underscore)
	v.RegisterRule("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]{2,31}$`, username)
		return matched
	}, "must be 3-32 characters starting with a letter, containing only letters, numbers and underscores")

	// Slug validation
	v.RegisterRule("slug", func(fl validator.FieldLevel) bool {
		slug := fl.Field().String()
		matched, _ := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)*$`, slug)
		return matched
	}, "must be a valid slug (lowercase letters, numbers, and hyphens)")

	// Password strength (at least 8 chars, 1 upper, 1 lower, 1 digit)
	v.RegisterRule("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}
		var hasUpper, hasLower, hasDigit bool
		for _, c := range password {
			switch {
			case unicode.IsUpper(c):
				hasUpper = true
			case unicode.IsLower(c):
				hasLower = true
			case unicode.IsDigit(c):
				hasDigit = true
			}
		}
		return hasUpper && hasLower && hasDigit
	}, "must be at least 8 characters with uppercase, lowercase and digit")

	// Strong password (8+ chars, upper, lower, digit, special)
	v.RegisterRule("password_strong", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}
		var hasUpper, hasLower, hasDigit, hasSpecial bool
		for _, c := range password {
			switch {
			case unicode.IsUpper(c):
				hasUpper = true
			case unicode.IsLower(c):
				hasLower = true
			case unicode.IsDigit(c):
				hasDigit = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSpecial = true
			}
		}
		return hasUpper && hasLower && hasDigit && hasSpecial
	}, "must be at least 8 characters with uppercase, lowercase, digit and special character")

	// Chinese ID card
	v.RegisterRule("id_card", func(fl validator.FieldLevel) bool {
		id := fl.Field().String()
		matched, _ := regexp.MatchString(`^\d{17}[\dXx]$`, id)
		return matched
	}, "must be a valid ID card number")

	// URL without protocol
	v.RegisterRule("url_path", func(fl validator.FieldLevel) bool {
		path := fl.Field().String()
		matched, _ := regexp.MatchString(`^/[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=%]*$`, path)
		return matched
	}, "must be a valid URL path")

	// Domain name
	v.RegisterRule("domain", func(fl validator.FieldLevel) bool {
		domain := fl.Field().String()
		matched, _ := regexp.MatchString(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`, domain)
		return matched
	}, "must be a valid domain name")

	// Safe string (no HTML/script injection)
	v.RegisterRule("safe_string", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		dangerousPatterns := []string{"<script", "<iframe", "javascript:", "onerror=", "onclick="}
		lower := strings.ToLower(s)
		for _, p := range dangerousPatterns {
			if strings.Contains(lower, p) {
				return false
			}
		}
		return true
	}, "contains potentially unsafe content")

	// No whitespace
	v.RegisterRule("no_whitespace", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return !strings.ContainsAny(s, " \t\n\r")
	}, "must not contain whitespace")

	// Trimmed (no leading/trailing whitespace)
	v.RegisterRule("trimmed", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return s == strings.TrimSpace(s)
	}, "must not have leading or trailing whitespace")

	// Alphanumeric with spaces
	v.RegisterRule("alpha_space", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		for _, c := range s {
			if !unicode.IsLetter(c) && !unicode.IsSpace(c) {
				return false
			}
		}
		return true
	}, "must contain only letters and spaces")

	// Positive number
	v.RegisterRule("positive", func(fl validator.FieldLevel) bool {
		switch fl.Field().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fl.Field().Int() > 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fl.Field().Uint() > 0
		case reflect.Float32, reflect.Float64:
			return fl.Field().Float() > 0
		}
		return false
	}, "must be a positive number")

	// Non-negative number
	v.RegisterRule("non_negative", func(fl validator.FieldLevel) bool {
		switch fl.Field().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fl.Field().Int() >= 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		case reflect.Float32, reflect.Float64:
			return fl.Field().Float() >= 0
		}
		return false
	}, "must be a non-negative number")
}

// RegisterRule registers a custom validation rule
func (v *Validator) RegisterRule(tag string, fn validator.Func, message string) error {
	if err := v.v.RegisterValidation(tag, fn); err != nil {
		return err
	}
	v.mu.Lock()
	v.messages[tag] = message
	v.mu.Unlock()
	return nil
}

// RegisterRuleWithParam registers a custom validation rule that accepts parameters
func (v *Validator) RegisterRuleWithParam(tag string, fn validator.Func, message string) error {
	return v.RegisterRule(tag, fn, message)
}

// SetMessage sets a custom message for a validation tag
func (v *Validator) SetMessage(tag, message string) {
	v.mu.Lock()
	v.messages[tag] = message
	v.mu.Unlock()
}

// SetMessages sets multiple custom messages
func (v *Validator) SetMessages(messages map[string]string) {
	v.mu.Lock()
	for tag, msg := range messages {
		v.messages[tag] = msg
	}
	v.mu.Unlock()
}

// Validate validates a struct and returns ValidationErrors
func (v *Validator) Validate(s interface{}) ValidationErrors {
	err := v.v.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors, ok := err.(validator.ValidationErrors); ok {
		validationErrors = errors
	} else {
		return ValidationErrors{{
			Field:   "",
			Message: err.Error(),
			Tag:     "validation",
		}}
	}

	errs := make(ValidationErrors, 0, len(validationErrors))
	for _, fe := range validationErrors {
		errs = append(errs, ValidationError{
			Field:   fe.Field(),
			Message: v.getMessage(fe),
			Tag:     fe.Tag(),
			Param:   fe.Param(),
			Value:   fe.Value(),
		})
	}

	return errs
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.v.Var(field, tag)
}

// getMessage returns the error message for a validation error
func (v *Validator) getMessage(fe validator.FieldError) string {
	v.mu.RLock()
	msg, ok := v.messages[fe.Tag()]
	v.mu.RUnlock()

	field := fe.Field()

	if ok {
		// Replace placeholders
		msg = strings.ReplaceAll(msg, ":field", field)
		msg = strings.ReplaceAll(msg, ":param", fe.Param())
		return msg
	}

	// Default messages
	return getDefaultMessage(fe)
}

// getDefaultMessage returns a default message for common validation tags
func getDefaultMessage(fe validator.FieldError) string {
	field := fe.Field()
	param := fe.Param()

	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "required_if":
		return field + " is required"
	case "required_unless":
		return field + " is required"
	case "required_with":
		return field + " is required when " + param + " is present"
	case "required_without":
		return field + " is required when " + param + " is not present"
	case "email":
		return field + " must be a valid email address"
	case "url":
		return field + " must be a valid URL"
	case "uri":
		return field + " must be a valid URI"
	case "min":
		return field + " must be at least " + param
	case "max":
		return field + " must be at most " + param
	case "len":
		return field + " must be exactly " + param + " characters"
	case "eq":
		return field + " must equal " + param
	case "ne":
		return field + " must not equal " + param
	case "gt":
		return field + " must be greater than " + param
	case "gte":
		return field + " must be greater than or equal to " + param
	case "lt":
		return field + " must be less than " + param
	case "lte":
		return field + " must be less than or equal to " + param
	case "oneof":
		return field + " must be one of: " + param
	case "contains":
		return field + " must contain " + param
	case "containsany":
		return field + " must contain at least one of: " + param
	case "excludes":
		return field + " must not contain " + param
	case "startswith":
		return field + " must start with " + param
	case "endswith":
		return field + " must end with " + param
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "numeric":
		return field + " must be numeric"
	case "number":
		return field + " must be a number"
	case "hexadecimal":
		return field + " must be hexadecimal"
	case "hexcolor":
		return field + " must be a valid hex color"
	case "rgb":
		return field + " must be a valid RGB color"
	case "rgba":
		return field + " must be a valid RGBA color"
	case "hsl":
		return field + " must be a valid HSL color"
	case "hsla":
		return field + " must be a valid HSLA color"
	case "e164":
		return field + " must be a valid E.164 phone number"
	case "json":
		return field + " must be valid JSON"
	case "jwt":
		return field + " must be a valid JWT"
	case "uuid":
		return field + " must be a valid UUID"
	case "uuid3":
		return field + " must be a valid UUID v3"
	case "uuid4":
		return field + " must be a valid UUID v4"
	case "uuid5":
		return field + " must be a valid UUID v5"
	case "ulid":
		return field + " must be a valid ULID"
	case "ascii":
		return field + " must contain only ASCII characters"
	case "printascii":
		return field + " must contain only printable ASCII characters"
	case "base64":
		return field + " must be valid base64"
	case "ip":
		return field + " must be a valid IP address"
	case "ipv4":
		return field + " must be a valid IPv4 address"
	case "ipv6":
		return field + " must be a valid IPv6 address"
	case "cidr":
		return field + " must be a valid CIDR notation"
	case "mac":
		return field + " must be a valid MAC address"
	case "latitude":
		return field + " must be a valid latitude"
	case "longitude":
		return field + " must be a valid longitude"
	case "ssn":
		return field + " must be a valid SSN"
	case "isbn":
		return field + " must be a valid ISBN"
	case "isbn10":
		return field + " must be a valid ISBN-10"
	case "isbn13":
		return field + " must be a valid ISBN-13"
	case "creditcard":
		return field + " must be a valid credit card number"
	case "datetime":
		return field + " must be a valid datetime"
	case "timezone":
		return field + " must be a valid timezone"
	case "eqfield":
		return field + " must match " + param
	case "nefield":
		return field + " must not match " + param
	case "gtfield":
		return field + " must be greater than " + param
	case "gtefield":
		return field + " must be greater than or equal to " + param
	case "ltfield":
		return field + " must be less than " + param
	case "ltefield":
		return field + " must be less than or equal to " + param
	case "unique":
		return field + " must have unique values"
	case "dive":
		return field + " contains invalid items"
	default:
		return field + " failed validation: " + fe.Tag()
	}
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Tag     string      `json:"tag"`
	Param   string      `json:"param,omitempty"`
	Value   interface{} `json:"-"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
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

// FirstField returns the first error's field name
func (ve ValidationErrors) FirstField() string {
	if len(ve) == 0 {
		return ""
	}
	return ve[0].Field
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

// ToSimpleMap returns a map with single error per field
func (ve ValidationErrors) ToSimpleMap() map[string]string {
	result := make(map[string]string)
	for _, e := range ve {
		if _, exists := result[e.Field]; !exists {
			result[e.Field] = e.Message
		}
	}
	return result
}

// Fields returns all field names with errors
func (ve ValidationErrors) Fields() []string {
	seen := make(map[string]bool)
	var fields []string
	for _, e := range ve {
		if !seen[e.Field] {
			seen[e.Field] = true
			fields = append(fields, e.Field)
		}
	}
	return fields
}

// IsEmpty returns true if there are no validation errors
func (ve ValidationErrors) IsEmpty() bool {
	return len(ve) == 0
}

// --- Convenience functions using global validator ---

// Validate validates a struct using the global validator
func Validate(s interface{}) ValidationErrors {
	return Global().Validate(s)
}

// ValidateVar validates a single variable using the global validator
func ValidateVar(field interface{}, tag string) error {
	return Global().ValidateVar(field, tag)
}

// RegisterRule registers a custom rule on the global validator
func RegisterRule(tag string, fn validator.Func, message string) error {
	return Global().RegisterRule(tag, fn, message)
}

// SetMessage sets a custom message on the global validator
func SetMessage(tag, message string) {
	Global().SetMessage(tag, message)
}

// SetMessages sets multiple custom messages on the global validator
func SetMessages(messages map[string]string) {
	Global().SetMessages(messages)
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
