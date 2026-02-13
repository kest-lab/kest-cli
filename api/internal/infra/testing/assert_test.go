package testing

import (
	"errors"
	"testing"
	"time"
)

func TestAssert_Equal(t *testing.T) {
	a := NewAssert(t)
	a.Equal(1, 1)
	a.Equal("hello", "hello")
	a.Equal([]int{1, 2, 3}, []int{1, 2, 3})
}

func TestAssert_NotEqual(t *testing.T) {
	a := NewAssert(t)
	a.NotEqual(1, 2)
	a.NotEqual("hello", "world")
}

func TestAssert_Nil(t *testing.T) {
	a := NewAssert(t)
	var ptr *int
	a.Nil(nil)
	a.Nil(ptr)
}

func TestAssert_NotNil(t *testing.T) {
	a := NewAssert(t)
	val := 42
	a.NotNil(&val)
	a.NotNil("hello")
}

func TestAssert_True(t *testing.T) {
	a := NewAssert(t)
	a.True(true)
	a.True(1 == 1)
}

func TestAssert_False(t *testing.T) {
	a := NewAssert(t)
	a.False(false)
	a.False(1 == 2)
}

func TestAssert_Error(t *testing.T) {
	a := NewAssert(t)
	a.Error(errors.New("some error"))
}

func TestAssert_NoError(t *testing.T) {
	a := NewAssert(t)
	a.NoError(nil)
}

func TestAssert_ErrorContains(t *testing.T) {
	a := NewAssert(t)
	a.ErrorContains(errors.New("connection failed: timeout"), "timeout")
}

func TestAssert_Contains(t *testing.T) {
	a := NewAssert(t)
	a.Contains("hello world", "world")
}

func TestAssert_NotContains(t *testing.T) {
	a := NewAssert(t)
	a.NotContains("hello world", "foo")
}

func TestAssert_HasPrefix(t *testing.T) {
	a := NewAssert(t)
	a.HasPrefix("hello world", "hello")
}

func TestAssert_HasSuffix(t *testing.T) {
	a := NewAssert(t)
	a.HasSuffix("hello world", "world")
}

func TestAssert_Len(t *testing.T) {
	a := NewAssert(t)
	a.Len([]int{1, 2, 3}, 3)
	a.Len("hello", 5)
	a.Len(map[string]int{"a": 1, "b": 2}, 2)
}

func TestAssert_Empty(t *testing.T) {
	a := NewAssert(t)
	a.Empty([]int{})
	a.Empty("")
	a.Empty(map[string]int{})
}

func TestAssert_NotEmpty(t *testing.T) {
	a := NewAssert(t)
	a.NotEmpty([]int{1})
	a.NotEmpty("hello")
	a.NotEmpty(map[string]int{"a": 1})
}

func TestAssert_Greater(t *testing.T) {
	a := NewAssert(t)
	a.Greater(5, 3)
	a.Greater(5.5, 3.3)
	a.Greater("b", "a")
}

func TestAssert_GreaterOrEqual(t *testing.T) {
	a := NewAssert(t)
	a.GreaterOrEqual(5, 3)
	a.GreaterOrEqual(5, 5)
}

func TestAssert_Less(t *testing.T) {
	a := NewAssert(t)
	a.Less(3, 5)
	a.Less(3.3, 5.5)
}

func TestAssert_LessOrEqual(t *testing.T) {
	a := NewAssert(t)
	a.LessOrEqual(3, 5)
	a.LessOrEqual(5, 5)
}

func TestAssert_JSONPath(t *testing.T) {
	a := NewAssert(t)
	json := []byte(`{"user": {"name": "John", "age": 30}}`)
	a.JSONPath(json, "user.name", "John")
	a.JSONPath(json, "user.age", float64(30)) // JSON numbers are float64
}

func TestAssert_JSONPathExists(t *testing.T) {
	a := NewAssert(t)
	json := []byte(`{"user": {"name": "John"}}`)
	a.JSONPathExists(json, "user.name")
	a.JSONPathExists(json, "user")
}

func TestAssert_JSONEqual(t *testing.T) {
	a := NewAssert(t)
	json1 := []byte(`{"a": 1, "b": 2}`)
	json2 := []byte(`{"b": 2, "a": 1}`)
	a.JSONEqual(json1, json2)
}

func TestAssert_TimeEqual(t *testing.T) {
	a := NewAssert(t)
	now := time.Now()
	later := now.Add(50 * time.Millisecond)
	a.TimeEqual(now, later, 100*time.Millisecond)
}

func TestAssert_TimeAfter(t *testing.T) {
	a := NewAssert(t)
	now := time.Now()
	later := now.Add(time.Hour)
	a.TimeAfter(later, now)
}

func TestAssert_TimeBefore(t *testing.T) {
	a := NewAssert(t)
	now := time.Now()
	earlier := now.Add(-time.Hour)
	a.TimeBefore(earlier, now)
}

func TestAssert_TimeBetween(t *testing.T) {
	a := NewAssert(t)
	now := time.Now()
	start := now.Add(-time.Hour)
	end := now.Add(time.Hour)
	a.TimeBetween(now, start, end)
}

func TestAssert_Panic(t *testing.T) {
	a := NewAssert(t)
	a.Panic(func() {
		panic("expected panic")
	})
}

func TestAssert_NoPanic(t *testing.T) {
	a := NewAssert(t)
	a.NoPanic(func() {
		// No panic here
	})
}

func TestAssert_Eventually(t *testing.T) {
	a := NewAssert(t)
	counter := 0
	a.Eventually(func() bool {
		counter++
		return counter >= 3
	}, 100*time.Millisecond, 10*time.Millisecond)
}

func TestAssert_Never(t *testing.T) {
	a := NewAssert(t)
	a.Never(func() bool {
		return false
	}, 50*time.Millisecond, 10*time.Millisecond)
}
