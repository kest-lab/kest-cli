package support

// Conditionable provides when/unless conditional execution pattern.
// This is a fluent interface for conditional operations.
type Conditionable[T any] struct {
	value T
}

// NewConditionable creates a new Conditionable wrapper.
func NewConditionable[T any](value T) *Conditionable[T] {
	return &Conditionable[T]{value: value}
}

// When applies the callback if the condition is true.
//
// Example:
//
//	result := NewConditionable(query).
//	    When(hasFilter, func(q *Query) *Query {
//	        return q.Where("status", status)
//	    }).
//	    Value()
func (c *Conditionable[T]) When(condition bool, callback func(T) T) *Conditionable[T] {
	if condition {
		c.value = callback(c.value)
	}
	return c
}

// WhenFunc applies the callback if the condition function returns true.
// The condition is evaluated lazily.
func (c *Conditionable[T]) WhenFunc(condition func() bool, callback func(T) T) *Conditionable[T] {
	if condition() {
		c.value = callback(c.value)
	}
	return c
}

// Unless applies the callback if the condition is false.
//
// Example:
//
//	result := NewConditionable(query).
//	    Unless(isAdmin, func(q *Query) *Query {
//	        return q.Where("user_id", userID)
//	    }).
//	    Value()
func (c *Conditionable[T]) Unless(condition bool, callback func(T) T) *Conditionable[T] {
	return c.When(!condition, callback)
}

// UnlessFunc applies the callback if the condition function returns false.
func (c *Conditionable[T]) UnlessFunc(condition func() bool, callback func(T) T) *Conditionable[T] {
	return c.WhenFunc(func() bool { return !condition() }, callback)
}

// WhenElse applies the callback if condition is true, otherwise applies the default callback.
func (c *Conditionable[T]) WhenElse(condition bool, callback func(T) T, defaultCallback func(T) T) *Conditionable[T] {
	if condition {
		c.value = callback(c.value)
	} else if defaultCallback != nil {
		c.value = defaultCallback(c.value)
	}
	return c
}

// Value returns the wrapped value.
func (c *Conditionable[T]) Value() T {
	return c.value
}

// When is a standalone function for simple conditional execution.
// Returns the result of callback if condition is true, otherwise returns the value unchanged.
//
// Example:
//
//	query := When(hasStatus, query, func(q *Query) *Query {
//	    return q.Where("status", status)
//	})
func When[T any](condition bool, value T, callback func(T) T) T {
	if condition {
		return callback(value)
	}
	return value
}

// Unless is a standalone function for simple conditional execution.
// Returns the result of callback if condition is false, otherwise returns the value unchanged.
func Unless[T any](condition bool, value T, callback func(T) T) T {
	return When(!condition, value, callback)
}

// WhenFilled applies the callback only if the value is filled (not blank).
func WhenFilled[T any](value T, callback func(T) T) T {
	if Filled(value) {
		return callback(value)
	}
	return value
}

// WhenBlank applies the callback only if the value is blank.
func WhenBlank[T any](value T, callback func(T) T) T {
	if Blank(value) {
		return callback(value)
	}
	return value
}
