package support

import "testing"

func TestOptional(t *testing.T) {
	// Test Some
	opt := Some("hello")
	if !opt.IsPresent() {
		t.Error("Some should be present")
	}
	if opt.Value() != "hello" {
		t.Errorf("Value should be 'hello', got %s", opt.Value())
	}

	// Test None
	none := None[string]()
	if none.IsPresent() {
		t.Error("None should not be present")
	}
	if !none.IsEmpty() {
		t.Error("None should be empty")
	}
}

func TestOptionalValueOr(t *testing.T) {
	opt := Some("hello")
	if opt.ValueOr("default") != "hello" {
		t.Error("ValueOr should return value when present")
	}

	none := None[string]()
	if none.ValueOr("default") != "default" {
		t.Error("ValueOr should return default when empty")
	}
}

func TestOptionalValueOrFunc(t *testing.T) {
	called := false
	opt := Some("hello")
	result := opt.ValueOrFunc(func() string {
		called = true
		return "default"
	})
	if result != "hello" || called {
		t.Error("ValueOrFunc should not call func when present")
	}

	none := None[string]()
	result = none.ValueOrFunc(func() string {
		return "default"
	})
	if result != "default" {
		t.Error("ValueOrFunc should call func when empty")
	}
}

func TestOptionalMap(t *testing.T) {
	opt := Some(5)
	mapped := opt.Map(func(n int) int {
		return n * 2
	})
	if mapped.Value() != 10 {
		t.Errorf("Map should transform value, got %d", mapped.Value())
	}

	none := None[int]()
	mapped = none.Map(func(n int) int {
		return n * 2
	})
	if mapped.IsPresent() {
		t.Error("Map on None should return None")
	}
}

func TestOptionalFilter(t *testing.T) {
	opt := Some(10)

	// Passes filter
	filtered := opt.Filter(func(n int) bool {
		return n > 5
	})
	if !filtered.IsPresent() {
		t.Error("Filter should keep value when predicate passes")
	}

	// Fails filter
	filtered = opt.Filter(func(n int) bool {
		return n > 15
	})
	if filtered.IsPresent() {
		t.Error("Filter should remove value when predicate fails")
	}
}

func TestOptionalIfPresent(t *testing.T) {
	called := false
	opt := Some("hello")
	opt.IfPresent(func(s string) {
		called = true
	})
	if !called {
		t.Error("IfPresent should call func when present")
	}

	called = false
	none := None[string]()
	none.IfPresent(func(s string) {
		called = true
	})
	if called {
		t.Error("IfPresent should not call func when empty")
	}
}

func TestOptionalOr(t *testing.T) {
	opt1 := Some("first")
	opt2 := Some("second")
	none := None[string]()

	if opt1.Or(opt2).Value() != "first" {
		t.Error("Or should return first when present")
	}
	if none.Or(opt2).Value() != "second" {
		t.Error("Or should return second when first is empty")
	}
}

func TestOptionalMap_Standalone(t *testing.T) {
	type User struct {
		Name string
	}

	userOpt := Some(User{Name: "John"})
	nameOpt := OptionalMap(userOpt, func(u User) string {
		return u.Name
	})

	if nameOpt.Value() != "John" {
		t.Errorf("OptionalMap should transform type, got %s", nameOpt.Value())
	}
}

func TestOptionalMustValue(t *testing.T) {
	opt := Some("hello")
	if opt.MustValue() != "hello" {
		t.Error("MustValue should return value when present")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustValue should panic when empty")
		}
	}()

	none := None[string]()
	none.MustValue() // should panic
}
