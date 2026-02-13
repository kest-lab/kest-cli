package pipeline

import (
	"context"
)

// Pipe is a single processing stage
type Pipe[T any] func(ctx context.Context, passable T, next func(T) T) T

// Pipeline provides a fluent interface for chaining operations
type Pipeline[T any] struct {
	passable T
	pipes    []Pipe[T]
	finally  func(T) T
}

// New creates a new pipeline
func New[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		pipes: make([]Pipe[T], 0),
	}
}

// Send sets the object being sent through the pipeline
func (p *Pipeline[T]) Send(passable T) *Pipeline[T] {
	p.passable = passable
	return p
}

// Through adds pipes to the pipeline
func (p *Pipeline[T]) Through(pipes ...Pipe[T]) *Pipeline[T] {
	p.pipes = append(p.pipes, pipes...)
	return p
}

// Pipe adds a single pipe to the pipeline
func (p *Pipeline[T]) Pipe(pipe Pipe[T]) *Pipeline[T] {
	p.pipes = append(p.pipes, pipe)
	return p
}

// Finally sets the final callback
func (p *Pipeline[T]) Finally(fn func(T) T) *Pipeline[T] {
	p.finally = fn
	return p
}

// Then runs the pipeline with a final destination
func (p *Pipeline[T]) Then(destination func(T) T) T {
	return p.run(context.Background(), destination)
}

// ThenReturn runs the pipeline and returns the result
func (p *Pipeline[T]) ThenReturn() T {
	return p.run(context.Background(), func(t T) T { return t })
}

// ThenWithContext runs the pipeline with context
func (p *Pipeline[T]) ThenWithContext(ctx context.Context, destination func(T) T) T {
	return p.run(ctx, destination)
}

func (p *Pipeline[T]) run(ctx context.Context, destination func(T) T) T {
	// Build the pipeline from end to start
	pipeline := destination

	// Reverse iterate to build the chain
	for i := len(p.pipes) - 1; i >= 0; i-- {
		pipe := p.pipes[i]
		next := pipeline
		pipeline = func(passable T) T {
			return pipe(ctx, passable, next)
		}
	}

	result := pipeline(p.passable)

	if p.finally != nil {
		result = p.finally(result)
	}

	return result
}

// --- Simple functional pipeline for common use cases ---

// Process runs a value through a series of transformers
func Process[T any](value T, transformers ...func(T) T) T {
	for _, fn := range transformers {
		value = fn(value)
	}
	return value
}

// ProcessWithError runs a value through transformers that can return errors
func ProcessWithError[T any](value T, transformers ...func(T) (T, error)) (T, error) {
	var err error
	for _, fn := range transformers {
		value, err = fn(value)
		if err != nil {
			return value, err
		}
	}
	return value, nil
}

// Filter filters a slice based on a predicate
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms a slice using a mapper function
func Map[T any, R any](items []T, mapper func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = mapper(item)
	}
	return result
}

// Reduce reduces a slice to a single value
func Reduce[T any, R any](items []T, initial R, reducer func(R, T) R) R {
	result := initial
	for _, item := range items {
		result = reducer(result, item)
	}
	return result
}
