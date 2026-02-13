package lifecycle

import "context"

// StartHook creates a Hook with only OnStart
func StartHook(name string, fn func(ctx context.Context) error) Hook {
	return Hook{Name: name, OnStart: fn}
}

// StopHook creates a Hook with only OnStop
func StopHook(name string, fn func(ctx context.Context) error) Hook {
	return Hook{Name: name, OnStop: fn}
}

// StartStopHook creates a Hook with both OnStart and OnStop
func StartStopHook(name string, start, stop func(ctx context.Context) error) Hook {
	return Hook{Name: name, OnStart: start, OnStop: stop}
}

// SimpleHook creates a Hook from simple functions (no context, no error)
func SimpleHook(name string, start, stop func()) Hook {
	var onStart, onStop func(context.Context) error

	if start != nil {
		onStart = func(context.Context) error {
			start()
			return nil
		}
	}
	if stop != nil {
		onStop = func(context.Context) error {
			stop()
			return nil
		}
	}

	return Hook{Name: name, OnStart: onStart, OnStop: onStop}
}
