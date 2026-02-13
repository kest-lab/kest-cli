package support

import (
	"errors"
	"testing"
	"time"
)

func TestBlank(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"whitespace string", "   ", true},
		{"non-empty string", "hello", false},
		{"zero int", 0, true},
		{"non-zero int", 1, false},
		{"zero float", 0.0, true},
		{"non-zero float", 1.5, false},
		{"empty slice", []int{}, true},
		{"non-empty slice", []int{1, 2}, false},
		{"empty map", map[string]int{}, true},
		{"non-empty map", map[string]int{"a": 1}, false},
		{"false bool", false, false}, // booleans are never blank
		{"true bool", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Blank(tt.value); got != tt.expected {
				t.Errorf("Blank(%v) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestFilled(t *testing.T) {
	if !Filled("hello") {
		t.Error("Filled('hello') should be true")
	}
	if Filled("") {
		t.Error("Filled('') should be false")
	}
}

func TestTap(t *testing.T) {
	var sideEffect string
	result := Tap("hello", func(s string) {
		sideEffect = s + " world"
	})

	if result != "hello" {
		t.Errorf("Tap should return original value, got %s", result)
	}
	if sideEffect != "hello world" {
		t.Errorf("Tap should execute callback, got %s", sideEffect)
	}
}

func TestWith(t *testing.T) {
	// Without callback
	result := With("hello", nil)
	if result != "hello" {
		t.Errorf("With without callback should return value, got %s", result)
	}

	// With callback
	result = With("hello", func(s string) string {
		return s + " world"
	})
	if result != "hello world" {
		t.Errorf("With callback should transform value, got %s", result)
	}
}

func TestTransform(t *testing.T) {
	// Filled value
	result := Transform("hello", func(s string) string {
		return s + " world"
	}, "default")
	if result != "hello world" {
		t.Errorf("Transform should apply callback, got %s", result)
	}

	// Blank value
	result = Transform("", func(s string) string {
		return s + " world"
	}, "default")
	if result != "default" {
		t.Errorf("Transform should return default for blank, got %s", result)
	}
}

func TestRetry(t *testing.T) {
	// Success on first try
	attempts := 0
	result, err := Retry(3, func(attempt int) (string, error) {
		attempts = attempt
		return "success", nil
	})
	if err != nil || result != "success" || attempts != 1 {
		t.Errorf("Retry should succeed on first try")
	}

	// Success on third try
	attempts = 0
	result, err = Retry(3, func(attempt int) (string, error) {
		attempts = attempt
		if attempt < 3 {
			return "", errors.New("fail")
		}
		return "success", nil
	})
	if err != nil || result != "success" || attempts != 3 {
		t.Errorf("Retry should succeed on third try")
	}

	// All attempts fail
	result, err = Retry(3, func(attempt int) (string, error) {
		return "", errors.New("fail")
	})
	if err == nil {
		t.Error("Retry should return error when all attempts fail")
	}
}

func TestRetryWithDelay(t *testing.T) {
	start := time.Now()
	attempts := 0

	_, _ = RetryWithDelay(3, 10*time.Millisecond, func(attempt int) (string, error) {
		attempts = attempt
		if attempt < 3 {
			return "", errors.New("fail")
		}
		return "success", nil
	})

	elapsed := time.Since(start)
	// Should have 2 delays of 10ms each
	if elapsed < 15*time.Millisecond {
		t.Errorf("RetryWithDelay should have delays, elapsed: %v", elapsed)
	}
}

func TestRetryWhen(t *testing.T) {
	tempErr := errors.New("temporary")
	permErr := errors.New("permanent")

	// Should retry on temporary error
	attempts := 0
	_, err := RetryWhen(3, func(attempt int) (string, error) {
		attempts = attempt
		return "", tempErr
	}, func(err error) bool {
		return err == tempErr
	})
	if attempts != 3 {
		t.Errorf("RetryWhen should retry on matching error, attempts: %d", attempts)
	}

	// Should not retry on permanent error
	attempts = 0
	_, err = RetryWhen(3, func(attempt int) (string, error) {
		attempts = attempt
		return "", permErr
	}, func(err error) bool {
		return err == tempErr
	})
	if attempts != 1 || err != permErr {
		t.Errorf("RetryWhen should not retry on non-matching error")
	}
}
