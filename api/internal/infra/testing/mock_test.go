package testing

import (
	"testing"
)

func TestMock_BasicUsage(t *testing.T) {
	mock := NewMock(t)
	mock.On("GetUser").Returns("John", nil)

	result := mock.Call("GetUser")
	if result[0] != "John" {
		t.Errorf("Expected 'John', got %v", result[0])
	}
	if result[1] != nil {
		t.Errorf("Expected nil, got %v", result[1])
	}
}

func TestMock_WithArgs(t *testing.T) {
	mock := NewMock(t)
	mock.On("Add").WithArgs(Equal(1), Equal(2)).Returns(3)

	result := mock.Call("Add", 1, 2)
	if result[0] != 3 {
		t.Errorf("Expected 3, got %v", result[0])
	}
}

func TestMock_Times(t *testing.T) {
	mock := NewMock(t)
	mock.On("DoSomething").Returns(nil).Times(2)

	mock.Call("DoSomething")
	mock.Call("DoSomething")

	mock.Verify()
}

func TestMock_Once(t *testing.T) {
	mock := NewMock(t)
	mock.On("DoOnce").Returns(nil).Once()

	mock.Call("DoOnce")

	mock.Verify()
}

func TestMock_Never(t *testing.T) {
	mock := NewMock(t)
	mock.On("NeverCalled").Returns(nil).Never()

	// Don't call it
	mock.Verify()
}

func TestMock_AtLeast(t *testing.T) {
	mock := NewMock(t)
	mock.On("CallMany").Returns(nil).AtLeast(2)

	mock.Call("CallMany")
	mock.Call("CallMany")
	mock.Call("CallMany")

	mock.Verify()
}

func TestMock_AtMost(t *testing.T) {
	mock := NewMock(t)
	mock.On("CallFew").Returns(nil).AtMost(3)

	mock.Call("CallFew")
	mock.Call("CallFew")

	mock.Verify()
}

func TestMock_Callback(t *testing.T) {
	mock := NewMock(t)
	mock.On("Calculate").Callback(func(args []interface{}) []interface{} {
		a := args[0].(int)
		b := args[1].(int)
		return []interface{}{a + b}
	})

	result := mock.Call("Calculate", 5, 3)
	if result[0] != 8 {
		t.Errorf("Expected 8, got %v", result[0])
	}
}

func TestMock_Capture(t *testing.T) {
	mock := NewMock(t)
	var captured string
	mock.On("Save").Capture(&captured).Returns(nil)

	mock.Call("Save", "test-value")

	if captured != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", captured)
	}
}

func TestMock_CallCount(t *testing.T) {
	mock := NewMock(t)
	mock.On("Count").Returns(nil)

	mock.Call("Count")
	mock.Call("Count")
	mock.Call("Count")

	if mock.CallCount("Count") != 3 {
		t.Errorf("Expected 3 calls, got %d", mock.CallCount("Count"))
	}
}

func TestMock_Calls(t *testing.T) {
	mock := NewMock(t)
	mock.On("Track").Returns(nil)

	mock.Call("Track", "arg1")
	mock.Call("Track", "arg2")

	calls := mock.Calls("Track")
	if len(calls) != 2 {
		t.Errorf("Expected 2 calls, got %d", len(calls))
	}
	if calls[0].Args[0] != "arg1" {
		t.Errorf("Expected 'arg1', got %v", calls[0].Args[0])
	}
	if calls[1].Args[0] != "arg2" {
		t.Errorf("Expected 'arg2', got %v", calls[1].Args[0])
	}
}

func TestMock_Reset(t *testing.T) {
	mock := NewMock(t)
	mock.On("Method").Returns(nil)
	mock.Call("Method")

	mock.Reset()

	if mock.CallCount("Method") != 0 {
		t.Error("Expected 0 calls after reset")
	}
}

func TestMatcher_Any(t *testing.T) {
	mock := NewMock(t)
	mock.On("Accept").WithArgs(Any(), Any()).Returns(true)

	result := mock.Call("Accept", "anything", 12345)
	if result[0] != true {
		t.Error("Expected true")
	}
}

func TestMatcher_OfType(t *testing.T) {
	mock := NewMock(t)
	mock.On("TypeCheck").WithArgs(OfType("")).Returns(true)

	result := mock.Call("TypeCheck", "a string")
	if result[0] != true {
		t.Error("Expected true")
	}
}

func TestMatcher_NotNil(t *testing.T) {
	mock := NewMock(t)
	mock.On("NotNilCheck").WithArgs(NotNil()).Returns(true)

	result := mock.Call("NotNilCheck", "not nil")
	if result[0] != true {
		t.Error("Expected true")
	}
}

func TestMatcher_MatchFunc(t *testing.T) {
	mock := NewMock(t)
	mock.On("Custom").WithArgs(MatchFunc("positive number", func(v interface{}) bool {
		n, ok := v.(int)
		return ok && n > 0
	})).Returns(true)

	result := mock.Call("Custom", 42)
	if result[0] != true {
		t.Error("Expected true")
	}
}
