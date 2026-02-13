// Package lifecycle provides application lifecycle management.
// It coordinates startup and shutdown of application components.
package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the lifecycle state
type State int

const (
	StateStopped State = iota
	StateStarting
	StateStarted
	StateStopping
)

func (s State) String() string {
	switch s {
	case StateStopped:
		return "stopped"
	case StateStarting:
		return "starting"
	case StateStarted:
		return "started"
	case StateStopping:
		return "stopping"
	default:
		return "unknown"
	}
}

// Hook represents a lifecycle hook with start and stop callbacks
type Hook struct {
	// Name identifies this hook (for logging/debugging)
	Name string

	// OnStart is called when the application starts.
	// Hooks are started in the order they were appended.
	OnStart func(ctx context.Context) error

	// OnStop is called when the application stops.
	// Hooks are stopped in reverse order (LIFO).
	OnStop func(ctx context.Context) error
}

// HookRecord records hook execution details
type HookRecord struct {
	Name     string
	Duration time.Duration
	Err      error
}

// Event represents a lifecycle event
type Event interface {
	event()
}

// Events
type (
	OnStarting struct{ Name string }
	OnStarted  struct {
		Name     string
		Duration time.Duration
		Err      error
	}
	OnStopping struct{ Name string }
	OnStopped  struct {
		Name     string
		Duration time.Duration
		Err      error
	}
	Starting    struct{}
	Started     struct{ Err error }
	Stopping    struct{}
	Stopped     struct{ Err error }
	RollingBack struct{ Err error }
	RolledBack  struct{ Err error }
)

func (*OnStarting) event()  {}
func (*OnStarted) event()   {}
func (*OnStopping) event()  {}
func (*OnStopped) event()   {}
func (*Starting) event()    {}
func (*Started) event()     {}
func (*Stopping) event()    {}
func (*Stopped) event()     {}
func (*RollingBack) event() {}
func (*RolledBack) event()  {}

// Logger logs lifecycle events
type Logger interface {
	LogEvent(Event)
}

// DefaultLogger logs to stdout
type DefaultLogger struct{}

func (l *DefaultLogger) LogEvent(e Event) {
	switch ev := e.(type) {
	case *Starting:
		fmt.Println("[lifecycle] Starting application...")
	case *Started:
		if ev.Err != nil {
			fmt.Printf("[lifecycle] Failed to start: %v\n", ev.Err)
		} else {
			fmt.Println("[lifecycle] Application started successfully")
		}
	case *Stopping:
		fmt.Println("[lifecycle] Stopping application...")
	case *Stopped:
		if ev.Err != nil {
			fmt.Printf("[lifecycle] Stopped with errors: %v\n", ev.Err)
		} else {
			fmt.Println("[lifecycle] Application stopped successfully")
		}
	case *OnStarting:
		fmt.Printf("[lifecycle] Starting %s...\n", ev.Name)
	case *OnStarted:
		if ev.Err != nil {
			fmt.Printf("[lifecycle] ✗ %s failed (%v): %v\n", ev.Name, ev.Duration, ev.Err)
		} else {
			fmt.Printf("[lifecycle] ✓ %s started (%v)\n", ev.Name, ev.Duration)
		}
	case *OnStopping:
		fmt.Printf("[lifecycle] Stopping %s...\n", ev.Name)
	case *OnStopped:
		if ev.Err != nil {
			fmt.Printf("[lifecycle] ✗ %s stop failed (%v): %v\n", ev.Name, ev.Duration, ev.Err)
		} else {
			fmt.Printf("[lifecycle] ✓ %s stopped (%v)\n", ev.Name, ev.Duration)
		}
	case *RollingBack:
		fmt.Printf("[lifecycle] Rolling back due to error: %v\n", ev.Err)
	case *RolledBack:
		if ev.Err != nil {
			fmt.Printf("[lifecycle] Rollback completed with errors: %v\n", ev.Err)
		} else {
			fmt.Println("[lifecycle] Rollback completed")
		}
	}
}

// NopLogger discards all events
type NopLogger struct{}

func (l *NopLogger) LogEvent(Event) {}

// Lifecycle manages application lifecycle hooks
type Lifecycle struct {
	mu           sync.Mutex
	state        State
	hooks        []Hook
	numStarted   int
	logger       Logger
	startRecords []HookRecord
	stopRecords  []HookRecord
}

// New creates a new Lifecycle instance
func New(opts ...Option) *Lifecycle {
	lc := &Lifecycle{
		state:  StateStopped,
		logger: &DefaultLogger{},
	}
	for _, opt := range opts {
		opt(lc)
	}
	return lc
}

// Option configures the Lifecycle
type Option func(*Lifecycle)

// WithLogger sets a custom logger
func WithLogger(l Logger) Option {
	return func(lc *Lifecycle) {
		lc.logger = l
	}
}

// Append adds a hook to the lifecycle.
// Hooks are started in append order and stopped in reverse order.
func (lc *Lifecycle) Append(h Hook) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.hooks = append(lc.hooks, h)
}

// State returns the current lifecycle state
func (lc *Lifecycle) State() State {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.state
}

// Start executes all OnStart hooks in order.
// If any hook fails, it automatically rolls back by calling Stop.
func (lc *Lifecycle) Start(ctx context.Context) error {
	lc.mu.Lock()
	if lc.state != StateStopped {
		lc.mu.Unlock()
		return fmt.Errorf("cannot start: lifecycle is %s", lc.state)
	}
	lc.state = StateStarting
	lc.numStarted = 0
	lc.startRecords = make([]HookRecord, 0, len(lc.hooks))
	lc.mu.Unlock()

	lc.logger.LogEvent(&Starting{})

	var startErr error
	for i, hook := range lc.hooks {
		if err := ctx.Err(); err != nil {
			startErr = err
			break
		}

		record := lc.runStartHook(ctx, hook)
		lc.mu.Lock()
		lc.startRecords = append(lc.startRecords, record)
		if record.Err == nil {
			lc.numStarted = i + 1
		}
		lc.mu.Unlock()

		if record.Err != nil {
			startErr = record.Err
			break
		}
	}

	if startErr != nil {
		// Rollback
		lc.logger.LogEvent(&RollingBack{Err: startErr})
		rollbackErr := lc.doStop(ctx)
		lc.logger.LogEvent(&RolledBack{Err: rollbackErr})

		lc.mu.Lock()
		lc.state = StateStopped
		lc.mu.Unlock()

		lc.logger.LogEvent(&Started{Err: startErr})
		if rollbackErr != nil {
			return errors.Join(startErr, rollbackErr)
		}
		return startErr
	}

	lc.mu.Lock()
	lc.state = StateStarted
	lc.mu.Unlock()

	lc.logger.LogEvent(&Started{})
	return nil
}

func (lc *Lifecycle) runStartHook(ctx context.Context, hook Hook) HookRecord {
	name := hook.Name
	if name == "" {
		name = "unnamed"
	}

	lc.logger.LogEvent(&OnStarting{Name: name})

	start := time.Now()
	var err error
	if hook.OnStart != nil {
		err = hook.OnStart(ctx)
	}
	duration := time.Since(start)

	lc.logger.LogEvent(&OnStarted{Name: name, Duration: duration, Err: err})

	return HookRecord{Name: name, Duration: duration, Err: err}
}

// Stop executes all OnStop hooks in reverse order.
// It continues even if some hooks fail, collecting all errors.
func (lc *Lifecycle) Stop(ctx context.Context) error {
	lc.mu.Lock()
	if lc.state != StateStarted && lc.state != StateStarting {
		lc.mu.Unlock()
		return nil
	}
	lc.state = StateStopping
	lc.mu.Unlock()

	lc.logger.LogEvent(&Stopping{})

	err := lc.doStop(ctx)

	lc.mu.Lock()
	lc.state = StateStopped
	lc.mu.Unlock()

	lc.logger.LogEvent(&Stopped{Err: err})
	return err
}

func (lc *Lifecycle) doStop(ctx context.Context) error {
	lc.mu.Lock()
	numStarted := lc.numStarted
	hooks := lc.hooks[:numStarted]
	lc.stopRecords = make([]HookRecord, 0, numStarted)
	lc.mu.Unlock()

	var errs []error

	// Stop in reverse order
	for i := len(hooks) - 1; i >= 0; i-- {
		if err := ctx.Err(); err != nil {
			errs = append(errs, err)
			break
		}

		hook := hooks[i]
		record := lc.runStopHook(ctx, hook)

		lc.mu.Lock()
		lc.stopRecords = append(lc.stopRecords, record)
		lc.mu.Unlock()

		if record.Err != nil {
			errs = append(errs, record.Err)
			// Continue stopping other hooks (best effort)
		}
	}

	return errors.Join(errs...)
}

func (lc *Lifecycle) runStopHook(ctx context.Context, hook Hook) HookRecord {
	name := hook.Name
	if name == "" {
		name = "unnamed"
	}

	lc.logger.LogEvent(&OnStopping{Name: name})

	start := time.Now()
	var err error
	if hook.OnStop != nil {
		err = hook.OnStop(ctx)
	}
	duration := time.Since(start)

	lc.logger.LogEvent(&OnStopped{Name: name, Duration: duration, Err: err})

	return HookRecord{Name: name, Duration: duration, Err: err}
}

// StartRecords returns the start hook execution records
func (lc *Lifecycle) StartRecords() []HookRecord {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return append([]HookRecord{}, lc.startRecords...)
}

// StopRecords returns the stop hook execution records
func (lc *Lifecycle) StopRecords() []HookRecord {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return append([]HookRecord{}, lc.stopRecords...)
}
