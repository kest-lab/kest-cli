package support

// Pipeline implements the pipeline pattern (middleware chain).
// It allows passing a value through a series of pipes (handlers).
type Pipeline[T any] struct {
	passable T
	pipes    []Pipe[T]
	finally  func(T)
}

// Pipe is a function that processes a value and calls the next handler.
type Pipe[T any] func(passable T, next func(T) T) T

// NewPipeline creates a new Pipeline instance.
func NewPipeline[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		pipes: make([]Pipe[T], 0),
	}
}

// Send sets the object being sent through the pipeline.
//
// Example:
//
//	result := NewPipeline[*Request]().
//	    Send(request).
//	    Through(authMiddleware, logMiddleware).
//	    Then(func(r *Request) *Request {
//	        return controller.Handle(r)
//	    })
func (p *Pipeline[T]) Send(passable T) *Pipeline[T] {
	p.passable = passable
	return p
}

// Through sets the array of pipes.
func (p *Pipeline[T]) Through(pipes ...Pipe[T]) *Pipeline[T] {
	p.pipes = pipes
	return p
}

// Pipe adds additional pipes to the pipeline.
func (p *Pipeline[T]) Pipe(pipes ...Pipe[T]) *Pipeline[T] {
	p.pipes = append(p.pipes, pipes...)
	return p
}

// Finally sets a callback to be executed after the pipeline completes.
func (p *Pipeline[T]) Finally(callback func(T)) *Pipeline[T] {
	p.finally = callback
	return p
}

// Then runs the pipeline with a final destination callback.
func (p *Pipeline[T]) Then(destination func(T) T) T {
	pipeline := destination

	// Build the pipeline from the end to the beginning (onion model)
	for i := len(p.pipes) - 1; i >= 0; i-- {
		pipe := p.pipes[i]
		next := pipeline
		pipeline = func(passable T) T {
			return pipe(passable, next)
		}
	}

	result := pipeline(p.passable)

	if p.finally != nil {
		p.finally(result)
	}

	return result
}

// ThenReturn runs the pipeline and returns the result.
func (p *Pipeline[T]) ThenReturn() T {
	return p.Then(func(passable T) T {
		return passable
	})
}

// SimplePipe creates a simple pipe that doesn't need to call next explicitly.
// The next handler is called automatically after the callback.
func SimplePipe[T any](callback func(T) T) Pipe[T] {
	return func(passable T, next func(T) T) T {
		result := callback(passable)
		return next(result)
	}
}

// ConditionalPipe creates a pipe that only executes if the condition is true.
func ConditionalPipe[T any](condition bool, pipe Pipe[T]) Pipe[T] {
	return func(passable T, next func(T) T) T {
		if condition {
			return pipe(passable, next)
		}
		return next(passable)
	}
}

// RecoverPipe creates a pipe that recovers from panics.
func RecoverPipe[T any](handler func(T, any) T) Pipe[T] {
	return func(passable T, next func(T) T) T {
		defer func() {
			if r := recover(); r != nil {
				passable = handler(passable, r)
			}
		}()
		return next(passable)
	}
}
