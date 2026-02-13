package support

import (
	"math/rand"
	"sort"
)

// Arr provides collection and array manipulation utilities.

// Contains checks if a slice contains a value
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// ContainsFunc checks if any element matches the predicate
func ContainsFunc[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// First returns the first element that matches the predicate, or zero value if none
func First[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// FirstOr returns the first element that matches, or the default value
func FirstOr[T any](slice []T, predicate func(T) bool, defaultValue T) T {
	if v, ok := First(slice, predicate); ok {
		return v
	}
	return defaultValue
}

// Last returns the last element that matches the predicate
func Last[T any](slice []T, predicate func(T) bool) (T, bool) {
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return slice[i], true
		}
	}
	var zero T
	return zero, false
}

// Filter returns elements that match the predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reject returns elements that don't match the predicate
func Reject[T any](slice []T, predicate func(T) bool) []T {
	return Filter(slice, func(v T) bool { return !predicate(v) })
}

// Map transforms each element using the mapper function
func Map[T any, R any](slice []T, mapper func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = mapper(v)
	}
	return result
}

// MapWithIndex transforms each element with its index
func MapWithIndex[T any, R any](slice []T, mapper func(int, T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = mapper(i, v)
	}
	return result
}

// FlatMap maps and flattens the result
func FlatMap[T any, R any](slice []T, mapper func(T) []R) []R {
	result := make([]R, 0)
	for _, v := range slice {
		result = append(result, mapper(v)...)
	}
	return result
}

// Reduce reduces a slice to a single value
func Reduce[T any, R any](slice []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// Each iterates over elements
func Each[T any](slice []T, fn func(T)) {
	for _, v := range slice {
		fn(v)
	}
}

// EachWithIndex iterates over elements with index
func EachWithIndex[T any](slice []T, fn func(int, T)) {
	for i, v := range slice {
		fn(i, v)
	}
}

// Unique returns unique elements
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// UniqueBy returns unique elements based on a key function
func UniqueBy[T any, K comparable](slice []T, keyFn func(T) K) []T {
	seen := make(map[K]bool)
	result := make([]T, 0)
	for _, v := range slice {
		key := keyFn(v)
		if !seen[key] {
			seen[key] = true
			result = append(result, v)
		}
	}
	return result
}

// GroupBy groups elements by a key function
func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range slice {
		key := keyFn(v)
		result[key] = append(result[key], v)
	}
	return result
}

// KeyBy creates a map keyed by a function
func KeyBy[T any, K comparable](slice []T, keyFn func(T) K) map[K]T {
	result := make(map[K]T)
	for _, v := range slice {
		result[keyFn(v)] = v
	}
	return result
}

// Pluck extracts a value from each element
func Pluck[T any, R any](slice []T, getter func(T) R) []R {
	return Map(slice, getter)
}

// Chunk splits a slice into chunks of the given size
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	result := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result
}

// Flatten flattens a 2D slice into a 1D slice
func Flatten[T any](slices [][]T) []T {
	result := make([]T, 0)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// Reverse reverses a slice
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// Shuffle randomly shuffles a slice
func Shuffle[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

// Take returns the first n elements
func Take[T any](slice []T, n int) []T {
	if n > len(slice) {
		n = len(slice)
	}
	if n < 0 {
		n = 0
	}
	return slice[:n]
}

// TakeLast returns the last n elements
func TakeLast[T any](slice []T, n int) []T {
	if n > len(slice) {
		n = len(slice)
	}
	if n < 0 {
		n = 0
	}
	return slice[len(slice)-n:]
}

// Skip skips the first n elements
func Skip[T any](slice []T, n int) []T {
	if n > len(slice) {
		return []T{}
	}
	if n < 0 {
		n = 0
	}
	return slice[n:]
}

// SortBy sorts elements by a comparable key
func SortBy[T any, K any](slice []T, keyFn func(T) K, less func(a, b K) bool) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	sort.Slice(result, func(i, j int) bool {
		return less(keyFn(result[i]), keyFn(result[j]))
	})
	return result
}

// Partition splits a slice into two based on a predicate
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	pass := make([]T, 0)
	fail := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			pass = append(pass, v)
		} else {
			fail = append(fail, v)
		}
	}
	return pass, fail
}

// Zip combines two slices into pairs
func Zip[T any, U any](a []T, b []U) [][2]any {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}
	result := make([][2]any, length)
	for i := 0; i < length; i++ {
		result[i] = [2]any{a[i], b[i]}
	}
	return result
}

// Diff returns elements in a that are not in b
func Diff[T comparable](a, b []T) []T {
	bSet := make(map[T]bool)
	for _, v := range b {
		bSet[v] = true
	}
	result := make([]T, 0)
	for _, v := range a {
		if !bSet[v] {
			result = append(result, v)
		}
	}
	return result
}

// Intersect returns elements that exist in both slices
func Intersect[T comparable](a, b []T) []T {
	bSet := make(map[T]bool)
	for _, v := range b {
		bSet[v] = true
	}
	result := make([]T, 0)
	for _, v := range a {
		if bSet[v] {
			result = append(result, v)
		}
	}
	return Unique(result)
}

// Sum returns the sum of all elements
func Sum[T int | int64 | float64](slice []T) T {
	var sum T
	for _, v := range slice {
		sum += v
	}
	return sum
}

// Avg returns the average of all elements
func Avg[T int | int64 | float64](slice []T) float64 {
	if len(slice) == 0 {
		return 0
	}
	return float64(Sum(slice)) / float64(len(slice))
}

// Min returns the minimum value
func Min[T int | int64 | float64 | string](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	min := slice[0]
	for _, v := range slice[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum value
func Max[T int | int64 | float64 | string](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// IsEmpty checks if a slice is empty
func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}

// IsNotEmpty checks if a slice is not empty
func IsNotEmpty[T any](slice []T) bool {
	return len(slice) > 0
}

// Random returns a random element from the slice
func Random[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[rand.Intn(len(slice))], true
}

// Wrap wraps a value in a slice if it's not already a slice
func Wrap[T any](value T) []T {
	return []T{value}
}
