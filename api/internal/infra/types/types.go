// Package types provides common type aliases for convenience.
// Inspired by GoFrame's type alias pattern.
package types

import "context"

// Common type aliases for cleaner code
type (
	// Map is an alias for map[string]any, commonly used for JSON-like data.
	Map = map[string]any

	// List is an alias for []Map, commonly used for collections of records.
	List = []Map

	// Slice is an alias for []any, commonly used for mixed-type arrays.
	Slice = []any

	// Ctx is an alias for context.Context.
	Ctx = context.Context

	// Strings is an alias for []string.
	Strings = []string

	// Ints is an alias for []int.
	Ints = []int

	// Int64s is an alias for []int64.
	Int64s = []int64
)

// H is a shorthand for Map, similar to gin.H.
// Useful for quick JSON responses.
//
// Example:
//
//	c.JSON(200, types.H{"status": "ok", "data": user})
type H = Map

// NewMap creates a new Map with optional initial values.
func NewMap(pairs ...any) Map {
	m := make(Map)
	for i := 0; i < len(pairs)-1; i += 2 {
		if key, ok := pairs[i].(string); ok {
			m[key] = pairs[i+1]
		}
	}
	return m
}

// NewList creates a new List with optional initial maps.
func NewList(maps ...Map) List {
	return maps
}

// Result represents a generic result with data and error.
type Result[T any] struct {
	Data  T
	Error error
}

// Ok creates a successful Result.
func Ok[T any](data T) Result[T] {
	return Result[T]{Data: data}
}

// Err creates a failed Result.
func Err[T any](err error) Result[T] {
	return Result[T]{Error: err}
}

// IsOk returns true if the result is successful.
func (r Result[T]) IsOk() bool {
	return r.Error == nil
}

// IsErr returns true if the result is an error.
func (r Result[T]) IsErr() bool {
	return r.Error != nil
}

// Unwrap returns the data or panics if there's an error.
func (r Result[T]) Unwrap() T {
	if r.Error != nil {
		panic(r.Error)
	}
	return r.Data
}

// UnwrapOr returns the data or the default value if there's an error.
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.Error != nil {
		return defaultValue
	}
	return r.Data
}

// Optional represents an optional value that may or may not be present.
type Optional[T any] struct {
	value   T
	present bool
}

// Some creates an Optional with a value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// None creates an empty Optional.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// IsPresent returns true if the value is present.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty returns true if the value is not present.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the value and a boolean indicating if it's present.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// OrElse returns the value if present, otherwise returns the default.
func (o Optional[T]) OrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// OrElseGet returns the value if present, otherwise calls the supplier.
func (o Optional[T]) OrElseGet(supplier func() T) T {
	if o.present {
		return o.value
	}
	return supplier()
}

// Map transforms the value if present.
func (o Optional[T]) Map(mapper func(T) T) Optional[T] {
	if o.present {
		return Some(mapper(o.value))
	}
	return o
}

// Filter returns the Optional if the predicate is true, otherwise returns None.
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if o.present && predicate(o.value) {
		return o
	}
	return None[T]()
}
