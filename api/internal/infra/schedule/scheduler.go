package schedule

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task represents a schedulable task
type Task interface {
	Run(ctx context.Context) error
}

// TaskFunc is a function that implements Task
type TaskFunc func(ctx context.Context) error

func (f TaskFunc) Run(ctx context.Context) error {
	return f(ctx)
}

// Event represents a scheduled event
type Event struct {
	name            string
	task            Task
	expression      string
	timezone        *time.Location
	withoutOverlap  bool
	onOneServer     bool
	runInBackground bool
	mutex           sync.Mutex
	running         bool

	// Schedule fields
	minute     string
	hour       string
	dayOfMonth string
	month      string
	dayOfWeek  string

	// Callbacks
	beforeCallbacks []func()
	afterCallbacks  []func()
	onSuccess       func()
	onFailure       func(error)
}

// NewEvent creates a new scheduled event
func NewEvent(name string, task Task) *Event {
	return &Event{
		name:       name,
		task:       task,
		timezone:   time.Local,
		minute:     "*",
		hour:       "*",
		dayOfMonth: "*",
		month:      "*",
		dayOfWeek:  "*",
	}
}

// Call creates a new event with a function
func Call(name string, fn func(ctx context.Context) error) *Event {
	return NewEvent(name, TaskFunc(fn))
}

// --- Schedule Methods ---

// Cron sets a custom cron expression
func (e *Event) Cron(expression string) *Event {
	e.expression = expression
	return e
}

// EveryMinute runs the task every minute
func (e *Event) EveryMinute() *Event {
	e.minute = "*"
	return e
}

// EveryTwoMinutes runs the task every two minutes
func (e *Event) EveryTwoMinutes() *Event {
	e.minute = "*/2"
	return e
}

// EveryFiveMinutes runs the task every five minutes
func (e *Event) EveryFiveMinutes() *Event {
	e.minute = "*/5"
	return e
}

// EveryTenMinutes runs the task every ten minutes
func (e *Event) EveryTenMinutes() *Event {
	e.minute = "*/10"
	return e
}

// EveryFifteenMinutes runs the task every fifteen minutes
func (e *Event) EveryFifteenMinutes() *Event {
	e.minute = "*/15"
	return e
}

// EveryThirtyMinutes runs the task every thirty minutes
func (e *Event) EveryThirtyMinutes() *Event {
	e.minute = "*/30"
	return e
}

// Hourly runs the task every hour
func (e *Event) Hourly() *Event {
	e.minute = "0"
	return e
}

// HourlyAt runs the task every hour at the given minute
func (e *Event) HourlyAt(minute int) *Event {
	e.minute = itoa(minute)
	return e
}

// Daily runs the task every day at midnight
func (e *Event) Daily() *Event {
	e.minute = "0"
	e.hour = "0"
	return e
}

// DailyAt runs the task every day at the given time
func (e *Event) DailyAt(hour, minute int) *Event {
	e.minute = itoa(minute)
	e.hour = itoa(hour)
	return e
}

// At is an alias for DailyAt
func (e *Event) At(hour, minute int) *Event {
	return e.DailyAt(hour, minute)
}

// TwiceDaily runs the task twice daily at the given hours
func (e *Event) TwiceDaily(hour1, hour2 int) *Event {
	e.minute = "0"
	e.hour = itoa(hour1) + "," + itoa(hour2)
	return e
}

// Weekly runs the task every week on Sunday at midnight
func (e *Event) Weekly() *Event {
	e.minute = "0"
	e.hour = "0"
	e.dayOfWeek = "0"
	return e
}

// WeeklyOn runs the task every week on the given day and time
func (e *Event) WeeklyOn(day time.Weekday, timeStr string) *Event {
	e.dayOfWeek = itoa(int(day))
	hour, minute := parseTime(timeStr)
	e.hour = itoa(hour)
	e.minute = itoa(minute)
	return e
}

// Monthly runs the task on the first day of every month
func (e *Event) Monthly() *Event {
	e.minute = "0"
	e.hour = "0"
	e.dayOfMonth = "1"
	return e
}

// MonthlyOn runs the task on the given day of each month
func (e *Event) MonthlyOn(day, hour, minute int) *Event {
	e.dayOfMonth = itoa(day)
	e.hour = itoa(hour)
	e.minute = itoa(minute)
	return e
}

// Quarterly runs the task on the first day of each quarter
func (e *Event) Quarterly() *Event {
	e.minute = "0"
	e.hour = "0"
	e.dayOfMonth = "1"
	e.month = "1,4,7,10"
	return e
}

// Yearly runs the task on January 1st
func (e *Event) Yearly() *Event {
	e.minute = "0"
	e.hour = "0"
	e.dayOfMonth = "1"
	e.month = "1"
	return e
}

// Weekdays runs the task only on weekdays
func (e *Event) Weekdays() *Event {
	e.dayOfWeek = "1-5"
	return e
}

// Weekends runs the task only on weekends
func (e *Event) Weekends() *Event {
	e.dayOfWeek = "0,6"
	return e
}

// Mondays runs the task only on Mondays
func (e *Event) Mondays() *Event {
	e.dayOfWeek = "1"
	return e
}

// Tuesdays runs the task only on Tuesdays
func (e *Event) Tuesdays() *Event {
	e.dayOfWeek = "2"
	return e
}

// Wednesdays runs the task only on Wednesdays
func (e *Event) Wednesdays() *Event {
	e.dayOfWeek = "3"
	return e
}

// Thursdays runs the task only on Thursdays
func (e *Event) Thursdays() *Event {
	e.dayOfWeek = "4"
	return e
}

// Fridays runs the task only on Fridays
func (e *Event) Fridays() *Event {
	e.dayOfWeek = "5"
	return e
}

// Saturdays runs the task only on Saturdays
func (e *Event) Saturdays() *Event {
	e.dayOfWeek = "6"
	return e
}

// Sundays runs the task only on Sundays
func (e *Event) Sundays() *Event {
	e.dayOfWeek = "0"
	return e
}

// --- Options ---

// WithoutOverlapping prevents task from running if already running
func (e *Event) WithoutOverlapping() *Event {
	e.withoutOverlap = true
	return e
}

// RunInBackground runs the task in a goroutine
func (e *Event) RunInBackground() *Event {
	e.runInBackground = true
	return e
}

// Timezone sets the timezone for the event
func (e *Event) Timezone(tz string) *Event {
	loc, err := time.LoadLocation(tz)
	if err == nil {
		e.timezone = loc
	}
	return e
}

// Before registers a callback to run before the task
func (e *Event) Before(fn func()) *Event {
	e.beforeCallbacks = append(e.beforeCallbacks, fn)
	return e
}

// After registers a callback to run after the task
func (e *Event) After(fn func()) *Event {
	e.afterCallbacks = append(e.afterCallbacks, fn)
	return e
}

// OnSuccess registers a callback for successful execution
func (e *Event) OnSuccess(fn func()) *Event {
	e.onSuccess = fn
	return e
}

// OnFailure registers a callback for failed execution
func (e *Event) OnFailure(fn func(error)) *Event {
	e.onFailure = fn
	return e
}

// Name returns the event name
func (e *Event) Name() string {
	return e.name
}

// IsDue checks if the event is due to run
func (e *Event) IsDue(t time.Time) bool {
	t = t.In(e.timezone)

	// Check minute
	if !matchField(e.minute, t.Minute()) {
		return false
	}

	// Check hour
	if !matchField(e.hour, t.Hour()) {
		return false
	}

	// Check day of month
	if !matchField(e.dayOfMonth, t.Day()) {
		return false
	}

	// Check month
	if !matchField(e.month, int(t.Month())) {
		return false
	}

	// Check day of week
	if !matchField(e.dayOfWeek, int(t.Weekday())) {
		return false
	}

	return true
}

// Run executes the event
func (e *Event) Run(ctx context.Context) error {
	// Check overlap
	if e.withoutOverlap {
		e.mutex.Lock()
		if e.running {
			e.mutex.Unlock()
			return nil
		}
		e.running = true
		e.mutex.Unlock()
		defer func() {
			e.mutex.Lock()
			e.running = false
			e.mutex.Unlock()
		}()
	}

	// Run before callbacks
	for _, fn := range e.beforeCallbacks {
		fn()
	}

	// Run task
	var err error
	if e.runInBackground {
		go func() {
			if taskErr := e.task.Run(ctx); taskErr != nil {
				if e.onFailure != nil {
					e.onFailure(taskErr)
				}
			} else if e.onSuccess != nil {
				e.onSuccess()
			}
		}()
	} else {
		err = e.task.Run(ctx)
		if err != nil {
			if e.onFailure != nil {
				e.onFailure(err)
			}
		} else if e.onSuccess != nil {
			e.onSuccess()
		}
	}

	// Run after callbacks
	for _, fn := range e.afterCallbacks {
		fn()
	}

	return err
}

// --- Scheduler ---

// Scheduler manages scheduled events
type Scheduler struct {
	mu      sync.RWMutex
	events  []*Event
	stop    chan struct{}
	running bool
}

var (
	globalScheduler *Scheduler
	once            sync.Once
)

// Global returns the global scheduler instance
func Global() *Scheduler {
	once.Do(func() {
		globalScheduler = New()
	})
	return globalScheduler
}

// New creates a new scheduler
func New() *Scheduler {
	return &Scheduler{
		events: make([]*Event, 0),
		stop:   make(chan struct{}),
	}
}

// Register registers an event with the scheduler
func (s *Scheduler) Register(event *Event) *Scheduler {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
	return s
}

// Call registers a function as a scheduled event
func (s *Scheduler) Call(name string, fn func(ctx context.Context) error) *Event {
	event := Call(name, fn)
	s.Register(event)
	return event
}

// Events returns all registered events
func (s *Scheduler) Events() []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events
}

// DueEvents returns events that are due to run
func (s *Scheduler) DueEvents(t time.Time) []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var due []*Event
	for _, event := range s.events {
		if event.IsDue(t) {
			due = append(due, event)
		}
	}
	return due
}

// Run runs all due events
func (s *Scheduler) Run(ctx context.Context) {
	now := time.Now()
	for _, event := range s.DueEvents(now) {
		if err := event.Run(ctx); err != nil {
			log.Printf("Scheduled task '%s' failed: %v", event.Name(), err)
		}
	}
}

// Start starts the scheduler loop
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// Run immediately for any due tasks
	s.Run(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stop:
			return
		case <-ticker.C:
			s.Run(ctx)
		}
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}
	s.running = false
	close(s.stop)
	s.stop = make(chan struct{})
}

// Clear removes all registered events
func (s *Scheduler) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = make([]*Event, 0)
}

// --- Convenience functions ---

// Register registers an event with the global scheduler
func Register(event *Event) *Scheduler {
	return Global().Register(event)
}

// Schedule creates and registers a new event
func Schedule(name string, fn func(ctx context.Context) error) *Event {
	return Global().Call(name, fn)
}

// Run runs due events on the global scheduler
func Run(ctx context.Context) {
	Global().Run(ctx)
}

// Start starts the global scheduler
func Start(ctx context.Context) {
	Global().Start(ctx)
}

// Stop stops the global scheduler
func Stop() {
	Global().Stop()
}

// --- Helper functions ---

func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	return string(rune('0'+i/10)) + string(rune('0'+i%10))
}

func parseTime(timeStr string) (hour, minute int) {
	// Parse "HH:MM" format
	if len(timeStr) >= 5 && timeStr[2] == ':' {
		hour = int(timeStr[0]-'0')*10 + int(timeStr[1]-'0')
		minute = int(timeStr[3]-'0')*10 + int(timeStr[4]-'0')
	}
	return
}

func matchField(field string, value int) bool {
	if field == "*" {
		return true
	}

	// Handle step values (*/n)
	if len(field) > 2 && field[0] == '*' && field[1] == '/' {
		step := int(field[2] - '0')
		if len(field) > 3 {
			step = step*10 + int(field[3]-'0')
		}
		return value%step == 0
	}

	// Handle ranges (n-m)
	for i := 0; i < len(field); i++ {
		if field[i] == '-' {
			start := parseNum(field[:i])
			end := parseNum(field[i+1:])
			return value >= start && value <= end
		}
	}

	// Handle lists (n,m,...)
	for i := 0; i < len(field); {
		j := i
		for j < len(field) && field[j] != ',' {
			j++
		}
		if parseNum(field[i:j]) == value {
			return true
		}
		i = j + 1
	}

	return false
}

func parseNum(s string) int {
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
