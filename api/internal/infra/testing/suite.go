package testing

import (
	"fmt"
	"testing"
)

// Suite provides BDD-style test organization with hooks
type Suite struct {
	t          *testing.T
	name       string
	beforeAll  []func()
	afterAll   []func()
	beforeEach []func()
	afterEach  []func()
	tests      []testEntry
	focused    bool
	skipped    bool
}

type testEntry struct {
	name    string
	fn      func(*testing.T)
	focused bool
	skipped bool
}

// NewSuite creates a new test suite
func NewSuite(t *testing.T, name string) *Suite {
	return &Suite{
		t:    t,
		name: name,
	}
}

// Describe creates a new test suite (alias for NewSuite)
func Describe(t *testing.T, name string) *Suite {
	return NewSuite(t, name)
}

// BeforeAll adds a hook that runs once before all tests
func (s *Suite) BeforeAll(fn func()) *Suite {
	s.beforeAll = append(s.beforeAll, fn)
	return s
}

// AfterAll adds a hook that runs once after all tests
func (s *Suite) AfterAll(fn func()) *Suite {
	s.afterAll = append(s.afterAll, fn)
	return s
}

// BeforeEach adds a hook that runs before each test
func (s *Suite) BeforeEach(fn func()) *Suite {
	s.beforeEach = append(s.beforeEach, fn)
	return s
}

// AfterEach adds a hook that runs after each test
func (s *Suite) AfterEach(fn func()) *Suite {
	s.afterEach = append(s.afterEach, fn)
	return s
}

// It adds a test case
func (s *Suite) It(name string, fn func(*testing.T)) *Suite {
	s.tests = append(s.tests, testEntry{name: name, fn: fn})
	return s
}

// FIt adds a focused test case (only focused tests will run)
func (s *Suite) FIt(name string, fn func(*testing.T)) *Suite {
	s.tests = append(s.tests, testEntry{name: name, fn: fn, focused: true})
	s.focused = true
	return s
}

// XIt adds a skipped test case
func (s *Suite) XIt(name string, fn func(*testing.T)) *Suite {
	s.tests = append(s.tests, testEntry{name: name, fn: fn, skipped: true})
	return s
}

// Skip marks the entire suite as skipped
func (s *Suite) Skip() *Suite {
	s.skipped = true
	return s
}

// Focus marks the entire suite as focused
func (s *Suite) Focus() *Suite {
	s.focused = true
	return s
}

// Run executes all tests in the suite
func (s *Suite) Run() {
	s.t.Run(s.name, func(t *testing.T) {
		if s.skipped {
			t.Skip("Suite skipped")
			return
		}

		// Run beforeAll hooks
		for _, fn := range s.beforeAll {
			fn()
		}

		// Ensure afterAll hooks run
		defer func() {
			for i := len(s.afterAll) - 1; i >= 0; i-- {
				s.afterAll[i]()
			}
		}()

		// Check if any tests are focused
		hasFocused := false
		for _, test := range s.tests {
			if test.focused {
				hasFocused = true
				break
			}
		}

		// Run tests
		for _, test := range s.tests {
			test := test // capture range variable

			// Skip non-focused tests if there are focused tests
			if hasFocused && !test.focused {
				continue
			}

			t.Run(test.name, func(t *testing.T) {
				if test.skipped {
					t.Skip("Test skipped")
					return
				}

				// Run beforeEach hooks
				for _, fn := range s.beforeEach {
					fn()
				}

				// Ensure afterEach hooks run
				defer func() {
					for i := len(s.afterEach) - 1; i >= 0; i-- {
						s.afterEach[i]()
					}
				}()

				test.fn(t)
			})
		}
	})
}

// TableTests provides table-driven test support
type TableTests[T any] struct {
	t     *testing.T
	cases []TableTestCase[T]
}

// TableTestCase represents a single test case in table-driven tests
type TableTestCase[T any] struct {
	Name     string
	Input    T
	Expected interface{}
	Skip     bool
}

// NewTableTests creates a new table-driven test builder
func NewTableTests[T any](t *testing.T) *TableTests[T] {
	return &TableTests[T]{t: t}
}

// Add adds a test case
func (tc *TableTests[T]) Add(name string, input T, expected interface{}) *TableTests[T] {
	tc.cases = append(tc.cases, TableTestCase[T]{
		Name:     name,
		Input:    input,
		Expected: expected,
	})
	return tc
}

// AddSkipped adds a skipped test case
func (tc *TableTests[T]) AddSkipped(name string, input T, expected interface{}) *TableTests[T] {
	tc.cases = append(tc.cases, TableTestCase[T]{
		Name:     name,
		Input:    input,
		Expected: expected,
		Skip:     true,
	})
	return tc
}

// Run executes all test cases with the provided test function
func (tc *TableTests[T]) Run(fn func(t *testing.T, input T, expected interface{})) {
	for _, c := range tc.cases {
		c := c // capture range variable
		tc.t.Run(c.Name, func(t *testing.T) {
			if c.Skip {
				t.Skip("Test case skipped")
				return
			}
			fn(t, c.Input, c.Expected)
		})
	}
}

// Cases returns all test cases for manual iteration
func (tc *TableTests[T]) Cases() []TableTestCase[T] {
	return tc.cases
}

// Context provides nested test context (like Ginkgo's Context)
type Context struct {
	t          *testing.T
	name       string
	parent     *Context
	beforeEach []func()
	afterEach  []func()
}

// NewContext creates a new test context
func NewContext(t *testing.T, name string) *Context {
	return &Context{t: t, name: name}
}

// BeforeEach adds a hook that runs before each test in this context
func (c *Context) BeforeEach(fn func()) *Context {
	c.beforeEach = append(c.beforeEach, fn)
	return c
}

// AfterEach adds a hook that runs after each test in this context
func (c *Context) AfterEach(fn func()) *Context {
	c.afterEach = append(c.afterEach, fn)
	return c
}

// Context creates a nested context
func (c *Context) Context(name string, fn func(*Context)) {
	child := &Context{
		t:      c.t,
		name:   name,
		parent: c,
	}
	c.t.Run(name, func(t *testing.T) {
		fn(child)
	})
}

// It runs a test within this context
func (c *Context) It(name string, fn func(*testing.T)) {
	c.t.Run(name, func(t *testing.T) {
		// Run all beforeEach hooks from parent to child
		c.runBeforeEach()
		defer c.runAfterEach()
		fn(t)
	})
}

func (c *Context) runBeforeEach() {
	if c.parent != nil {
		c.parent.runBeforeEach()
	}
	for _, fn := range c.beforeEach {
		fn()
	}
}

func (c *Context) runAfterEach() {
	// Run in reverse order, child to parent
	for i := len(c.afterEach) - 1; i >= 0; i-- {
		c.afterEach[i]()
	}
	if c.parent != nil {
		c.parent.runAfterEach()
	}
}

// Spec is a shorthand for creating and running a suite
func Spec(t *testing.T, name string, fn func(*Suite)) {
	suite := NewSuite(t, name)
	fn(suite)
	suite.Run()
}

// Given/When/Then style helpers for BDD
type GWT struct {
	t      *testing.T
	given  string
	when   string
	then   string
	setup  func()
	action func() interface{}
}

// Given starts a Given/When/Then test
func Given(t *testing.T, description string) *GWT {
	return &GWT{t: t, given: description}
}

// Setup adds setup logic for the Given clause
func (g *GWT) Setup(fn func()) *GWT {
	g.setup = fn
	return g
}

// When adds the When clause
func (g *GWT) When(description string, action func() interface{}) *GWT {
	g.when = description
	g.action = action
	return g
}

// Then adds the Then clause and runs the test
func (g *GWT) Then(description string, assertion func(result interface{})) {
	testName := fmt.Sprintf("Given %s, When %s, Then %s", g.given, g.when, description)
	g.t.Run(testName, func(t *testing.T) {
		if g.setup != nil {
			g.setup()
		}
		var result interface{}
		if g.action != nil {
			result = g.action()
		}
		assertion(result)
	})
}
