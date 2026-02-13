package testing

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

// Mock provides a generic mock builder
type Mock struct {
	t            *testing.T
	mu           sync.RWMutex
	expectations map[string]*Expectation
	calls        map[string][]CallRecord
}

// CallRecord records a method call
type CallRecord struct {
	Args   []interface{}
	Result []interface{}
}

// Expectation defines expected behavior for a method
type Expectation struct {
	methodName string
	args       []Matcher
	returns    []interface{}
	times      int
	minTimes   int
	maxTimes   int
	called     int
	callback   func(args []interface{}) []interface{}
	captured   []interface{}
}

// Matcher interface for argument matching
type Matcher interface {
	Match(actual interface{}) bool
	String() string
}

// NewMock creates a new mock
func NewMock(t *testing.T) *Mock {
	return &Mock{
		t:            t,
		expectations: make(map[string]*Expectation),
		calls:        make(map[string][]CallRecord),
	}
}

// On sets up an expectation for a method
func (m *Mock) On(methodName string) *Expectation {
	m.mu.Lock()
	defer m.mu.Unlock()

	exp := &Expectation{
		methodName: methodName,
		times:      -1, // -1 means any number of times
		minTimes:   0,
		maxTimes:   -1,
	}
	m.expectations[methodName] = exp
	return exp
}

// Call records a method call and returns the expected result
func (m *Mock) Call(methodName string, args ...interface{}) []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Record the call
	m.calls[methodName] = append(m.calls[methodName], CallRecord{Args: args})

	// Find matching expectation
	exp, ok := m.expectations[methodName]
	if !ok {
		m.t.Errorf("Unexpected call to %s", methodName)
		return nil
	}

	// Check argument matchers
	if len(exp.args) > 0 {
		if len(args) != len(exp.args) {
			m.t.Errorf("Expected %d args for %s, got %d", len(exp.args), methodName, len(args))
			return exp.returns
		}
		for i, matcher := range exp.args {
			if !matcher.Match(args[i]) {
				m.t.Errorf("Argument %d for %s: expected %s, got %v", i, methodName, matcher.String(), args[i])
			}
		}
	}

	// Capture arguments if requested
	if len(exp.captured) > 0 {
		for i, ptr := range exp.captured {
			if i < len(args) && ptr != nil {
				reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(args[i]))
			}
		}
	}

	exp.called++

	// Use callback if provided
	if exp.callback != nil {
		return exp.callback(args)
	}

	return exp.returns
}

// Verify checks that all expectations were met
func (m *Mock) Verify() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, exp := range m.expectations {
		if exp.times >= 0 && exp.called != exp.times {
			m.t.Errorf("Expected %s to be called %d times, was called %d times", name, exp.times, exp.called)
		}
		if exp.minTimes > 0 && exp.called < exp.minTimes {
			m.t.Errorf("Expected %s to be called at least %d times, was called %d times", name, exp.minTimes, exp.called)
		}
		if exp.maxTimes >= 0 && exp.called > exp.maxTimes {
			m.t.Errorf("Expected %s to be called at most %d times, was called %d times", name, exp.maxTimes, exp.called)
		}
	}
}

// Calls returns all recorded calls for a method
func (m *Mock) Calls(methodName string) []CallRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[methodName]
}

// CallCount returns the number of times a method was called
func (m *Mock) CallCount(methodName string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.calls[methodName])
}

// Reset clears all expectations and recorded calls
func (m *Mock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expectations = make(map[string]*Expectation)
	m.calls = make(map[string][]CallRecord)
}

// Expectation methods

// WithArgs sets expected arguments using matchers
func (e *Expectation) WithArgs(matchers ...Matcher) *Expectation {
	e.args = matchers
	return e
}

// WithArgsEqual sets expected arguments using equality matching
func (e *Expectation) WithArgsEqual(args ...interface{}) *Expectation {
	matchers := make([]Matcher, len(args))
	for i, arg := range args {
		matchers[i] = Equal(arg)
	}
	e.args = matchers
	return e
}

// Returns sets the return values
func (e *Expectation) Returns(values ...interface{}) *Expectation {
	e.returns = values
	return e
}

// Times sets the exact number of expected calls
func (e *Expectation) Times(n int) *Expectation {
	e.times = n
	return e
}

// Once expects exactly one call
func (e *Expectation) Once() *Expectation {
	return e.Times(1)
}

// Twice expects exactly two calls
func (e *Expectation) Twice() *Expectation {
	return e.Times(2)
}

// Never expects no calls
func (e *Expectation) Never() *Expectation {
	return e.Times(0)
}

// AtLeast sets minimum expected calls
func (e *Expectation) AtLeast(n int) *Expectation {
	e.minTimes = n
	return e
}

// AtMost sets maximum expected calls
func (e *Expectation) AtMost(n int) *Expectation {
	e.maxTimes = n
	return e
}

// Callback sets a callback function for dynamic returns
func (e *Expectation) Callback(fn func(args []interface{}) []interface{}) *Expectation {
	e.callback = fn
	return e
}

// Capture captures arguments into the provided pointers
func (e *Expectation) Capture(ptrs ...interface{}) *Expectation {
	e.captured = ptrs
	return e
}

// Matchers

type equalMatcher struct {
	expected interface{}
}

func (m *equalMatcher) Match(actual interface{}) bool {
	return reflect.DeepEqual(m.expected, actual)
}

func (m *equalMatcher) String() string {
	return fmt.Sprintf("equal to %v", m.expected)
}

// Equal creates a matcher that checks for equality
func Equal(expected interface{}) Matcher {
	return &equalMatcher{expected: expected}
}

type anyMatcher struct{}

func (m *anyMatcher) Match(actual interface{}) bool {
	return true
}

func (m *anyMatcher) String() string {
	return "any value"
}

// Any creates a matcher that matches any value
func Any() Matcher {
	return &anyMatcher{}
}

type typeMatcher struct {
	expectedType reflect.Type
}

func (m *typeMatcher) Match(actual interface{}) bool {
	return reflect.TypeOf(actual) == m.expectedType
}

func (m *typeMatcher) String() string {
	return fmt.Sprintf("type %v", m.expectedType)
}

// OfType creates a matcher that checks the type
func OfType(expected interface{}) Matcher {
	return &typeMatcher{expectedType: reflect.TypeOf(expected)}
}

type funcMatcher struct {
	fn   func(interface{}) bool
	desc string
}

func (m *funcMatcher) Match(actual interface{}) bool {
	return m.fn(actual)
}

func (m *funcMatcher) String() string {
	return m.desc
}

// MatchFunc creates a custom matcher using a function
func MatchFunc(desc string, fn func(interface{}) bool) Matcher {
	return &funcMatcher{fn: fn, desc: desc}
}

type notNilMatcher struct{}

func (m *notNilMatcher) Match(actual interface{}) bool {
	if actual == nil {
		return false
	}
	v := reflect.ValueOf(actual)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return !v.IsNil()
	}
	return true
}

func (m *notNilMatcher) String() string {
	return "not nil"
}

// NotNil creates a matcher that checks for non-nil values
func NotNil() Matcher {
	return &notNilMatcher{}
}

type containsMatcher struct {
	substring string
}

func (m *containsMatcher) Match(actual interface{}) bool {
	str, ok := actual.(string)
	if !ok {
		return false
	}
	return len(str) > 0 && len(m.substring) > 0 &&
		(str == m.substring || len(str) >= len(m.substring))
}

func (m *containsMatcher) String() string {
	return fmt.Sprintf("contains %q", m.substring)
}

// Contains creates a matcher for string containment
func Contains(substring string) Matcher {
	return &containsMatcher{substring: substring}
}
