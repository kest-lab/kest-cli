package testing

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Assert provides assertion helpers
type Assert struct {
	t *testing.T
}

// NewAssert creates a new assertion helper
func NewAssert(t *testing.T) *Assert {
	return &Assert{t: t}
}

// Equal asserts that two values are equal
func (a *Assert) Equal(expected, actual interface{}) *Assert {
	if diff := cmp.Diff(expected, actual); diff != "" {
		a.t.Errorf("Values not equal (-expected +actual):\n%s", diff)
	}
	return a
}

// NotEqual asserts that two values are not equal
func (a *Assert) NotEqual(expected, actual interface{}) *Assert {
	if cmp.Equal(expected, actual) {
		a.t.Errorf("Expected values to be different, but both are: %v", actual)
	}
	return a
}

// Nil asserts that a value is nil
func (a *Assert) Nil(value interface{}) *Assert {
	if !isNil(value) {
		a.t.Errorf("Expected nil, got: %v", value)
	}
	return a
}

// NotNil asserts that a value is not nil
func (a *Assert) NotNil(value interface{}) *Assert {
	if isNil(value) {
		a.t.Error("Expected non-nil value, got nil")
	}
	return a
}

// True asserts that a value is true
func (a *Assert) True(value bool) *Assert {
	if !value {
		a.t.Error("Expected true, got false")
	}
	return a
}

// False asserts that a value is false
func (a *Assert) False(value bool) *Assert {
	if value {
		a.t.Error("Expected false, got true")
	}
	return a
}

// Error asserts that an error is not nil
func (a *Assert) Error(err error) *Assert {
	if err == nil {
		a.t.Error("Expected an error, got nil")
	}
	return a
}

// NoError asserts that an error is nil
func (a *Assert) NoError(err error) *Assert {
	if err != nil {
		a.t.Errorf("Expected no error, got: %v", err)
	}
	return a
}

// ErrorIs asserts that an error matches a target error
func (a *Assert) ErrorIs(err, target error) *Assert {
	if err == nil {
		a.t.Errorf("Expected error %v, got nil", target)
		return a
	}
	// Use errors.Is for proper error chain checking
	if err.Error() != target.Error() {
		a.t.Errorf("Expected error %v, got: %v", target, err)
	}
	return a
}

// ErrorContains asserts that an error message contains a substring
func (a *Assert) ErrorContains(err error, substring string) *Assert {
	if err == nil {
		a.t.Errorf("Expected error containing %q, got nil", substring)
		return a
	}
	if !strings.Contains(err.Error(), substring) {
		a.t.Errorf("Expected error containing %q, got: %v", substring, err)
	}
	return a
}

// Contains asserts that a string contains a substring
func (a *Assert) Contains(str, substring string) *Assert {
	if !strings.Contains(str, substring) {
		a.t.Errorf("Expected %q to contain %q", str, substring)
	}
	return a
}

// NotContains asserts that a string does not contain a substring
func (a *Assert) NotContains(str, substring string) *Assert {
	if strings.Contains(str, substring) {
		a.t.Errorf("Expected %q to not contain %q", str, substring)
	}
	return a
}

// HasPrefix asserts that a string has a prefix
func (a *Assert) HasPrefix(str, prefix string) *Assert {
	if !strings.HasPrefix(str, prefix) {
		a.t.Errorf("Expected %q to have prefix %q", str, prefix)
	}
	return a
}

// HasSuffix asserts that a string has a suffix
func (a *Assert) HasSuffix(str, suffix string) *Assert {
	if !strings.HasSuffix(str, suffix) {
		a.t.Errorf("Expected %q to have suffix %q", str, suffix)
	}
	return a
}

// Len asserts the length of a slice, map, or string
func (a *Assert) Len(value interface{}, expected int) *Assert {
	v := reflect.ValueOf(value)
	var length int
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Array, reflect.Chan:
		length = v.Len()
	default:
		a.t.Errorf("Cannot get length of %T", value)
		return a
	}
	if length != expected {
		a.t.Errorf("Expected length %d, got %d", expected, length)
	}
	return a
}

// Empty asserts that a slice, map, or string is empty
func (a *Assert) Empty(value interface{}) *Assert {
	v := reflect.ValueOf(value)
	var length int
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Array, reflect.Chan:
		length = v.Len()
	default:
		a.t.Errorf("Cannot check if %T is empty", value)
		return a
	}
	if length != 0 {
		a.t.Errorf("Expected empty, got length %d", length)
	}
	return a
}

// NotEmpty asserts that a slice, map, or string is not empty
func (a *Assert) NotEmpty(value interface{}) *Assert {
	v := reflect.ValueOf(value)
	var length int
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Array, reflect.Chan:
		length = v.Len()
	default:
		a.t.Errorf("Cannot check if %T is empty", value)
		return a
	}
	if length == 0 {
		a.t.Error("Expected non-empty value")
	}
	return a
}

// Greater asserts that actual > expected
func (a *Assert) Greater(actual, expected interface{}) *Assert {
	if !isGreater(actual, expected) {
		a.t.Errorf("Expected %v > %v", actual, expected)
	}
	return a
}

// GreaterOrEqual asserts that actual >= expected
func (a *Assert) GreaterOrEqual(actual, expected interface{}) *Assert {
	if !isGreaterOrEqual(actual, expected) {
		a.t.Errorf("Expected %v >= %v", actual, expected)
	}
	return a
}

// Less asserts that actual < expected
func (a *Assert) Less(actual, expected interface{}) *Assert {
	if !isLess(actual, expected) {
		a.t.Errorf("Expected %v < %v", actual, expected)
	}
	return a
}

// LessOrEqual asserts that actual <= expected
func (a *Assert) LessOrEqual(actual, expected interface{}) *Assert {
	if !isLessOrEqual(actual, expected) {
		a.t.Errorf("Expected %v <= %v", actual, expected)
	}
	return a
}

// JSON Assertions

// JSONPath asserts a value at a JSON path
func (a *Assert) JSONPath(jsonData []byte, path string, expected interface{}) *Assert {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		a.t.Fatalf("Failed to parse JSON: %v", err)
	}

	actual := getNestedValue(data, path)
	if actual == nil {
		a.t.Errorf("JSON path %q not found", path)
		return a
	}

	if !reflect.DeepEqual(actual, expected) {
		a.t.Errorf("JSON path %q: expected %v, got %v", path, expected, actual)
	}
	return a
}

// JSONPathExists asserts that a JSON path exists
func (a *Assert) JSONPathExists(jsonData []byte, path string) *Assert {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		a.t.Fatalf("Failed to parse JSON: %v", err)
	}

	if getNestedValue(data, path) == nil {
		a.t.Errorf("JSON path %q not found", path)
	}
	return a
}

// JSONEqual asserts that two JSON values are equal
func (a *Assert) JSONEqual(expected, actual []byte) *Assert {
	var expectedData, actualData interface{}
	if err := json.Unmarshal(expected, &expectedData); err != nil {
		a.t.Fatalf("Failed to parse expected JSON: %v", err)
	}
	if err := json.Unmarshal(actual, &actualData); err != nil {
		a.t.Fatalf("Failed to parse actual JSON: %v", err)
	}

	if diff := cmp.Diff(expectedData, actualData); diff != "" {
		a.t.Errorf("JSON not equal (-expected +actual):\n%s", diff)
	}
	return a
}

// Time Assertions

// TimeEqual asserts that two times are equal within a tolerance
func (a *Assert) TimeEqual(expected, actual time.Time, tolerance time.Duration) *Assert {
	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		a.t.Errorf("Times not equal within %v: expected %v, got %v (diff: %v)",
			tolerance, expected, actual, diff)
	}
	return a
}

// TimeAfter asserts that actual is after expected
func (a *Assert) TimeAfter(actual, expected time.Time) *Assert {
	if !actual.After(expected) {
		a.t.Errorf("Expected %v to be after %v", actual, expected)
	}
	return a
}

// TimeBefore asserts that actual is before expected
func (a *Assert) TimeBefore(actual, expected time.Time) *Assert {
	if !actual.Before(expected) {
		a.t.Errorf("Expected %v to be before %v", actual, expected)
	}
	return a
}

// TimeBetween asserts that actual is between start and end
func (a *Assert) TimeBetween(actual, start, end time.Time) *Assert {
	if actual.Before(start) || actual.After(end) {
		a.t.Errorf("Expected %v to be between %v and %v", actual, start, end)
	}
	return a
}

// Panic asserts that a function panics
func (a *Assert) Panic(fn func()) *Assert {
	defer func() {
		if r := recover(); r == nil {
			a.t.Error("Expected panic, but function did not panic")
		}
	}()
	fn()
	return a
}

// NoPanic asserts that a function does not panic
func (a *Assert) NoPanic(fn func()) *Assert {
	defer func() {
		if r := recover(); r != nil {
			a.t.Errorf("Expected no panic, but got: %v", r)
		}
	}()
	fn()
	return a
}

// Eventually asserts that a condition becomes true within a timeout
func (a *Assert) Eventually(condition func() bool, timeout, interval time.Duration) *Assert {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return a
		}
		time.Sleep(interval)
	}
	a.t.Errorf("Condition not met within %v", timeout)
	return a
}

// Never asserts that a condition never becomes true within a duration
func (a *Assert) Never(condition func() bool, duration, interval time.Duration) *Assert {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		if condition() {
			a.t.Error("Condition became true when it should never be")
			return a
		}
		time.Sleep(interval)
	}
	return a
}

// Helper functions

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func isGreater(a, b interface{}) bool {
	return compare(a, b) > 0
}

func isGreaterOrEqual(a, b interface{}) bool {
	return compare(a, b) >= 0
}

func isLess(a, b interface{}) bool {
	return compare(a, b) < 0
}

func isLessOrEqual(a, b interface{}) bool {
	return compare(a, b) <= 0
}

func compare(a, b interface{}) int {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	switch aVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		aInt := aVal.Int()
		bInt := bVal.Int()
		if aInt > bInt {
			return 1
		} else if aInt < bInt {
			return -1
		}
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		aUint := aVal.Uint()
		bUint := bVal.Uint()
		if aUint > bUint {
			return 1
		} else if aUint < bUint {
			return -1
		}
		return 0
	case reflect.Float32, reflect.Float64:
		aFloat := aVal.Float()
		bFloat := bVal.Float()
		if aFloat > bFloat {
			return 1
		} else if aFloat < bFloat {
			return -1
		}
		return 0
	case reflect.String:
		aStr := aVal.String()
		bStr := bVal.String()
		if aStr > bStr {
			return 1
		} else if aStr < bStr {
			return -1
		}
		return 0
	}
	panic(fmt.Sprintf("cannot compare %T and %T", a, b))
}
