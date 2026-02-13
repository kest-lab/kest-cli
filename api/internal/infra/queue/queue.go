package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Job represents a queueable job
type Job interface {
	Handle(ctx context.Context) error
}

// JobWithQueue allows a job to specify its queue
type JobWithQueue interface {
	Job
	Queue() string
}

// JobWithDelay allows a job to specify a delay
type JobWithDelay interface {
	Job
	Delay() time.Duration
}

// JobWithRetry allows a job to specify retry behavior
type JobWithRetry interface {
	Job
	MaxRetries() int
	RetryDelay() time.Duration
}

// Driver defines the queue driver interface
type Driver interface {
	// Push adds a job to the queue
	Push(ctx context.Context, queue string, payload []byte) error

	// PushDelayed adds a job to the queue with a delay
	PushDelayed(ctx context.Context, queue string, payload []byte, delay time.Duration) error

	// Pop retrieves the next job from the queue
	Pop(ctx context.Context, queue string) ([]byte, error)

	// Size returns the number of jobs in the queue
	Size(ctx context.Context, queue string) (int64, error)

	// Clear removes all jobs from the queue
	Clear(ctx context.Context, queue string) error

	// Close closes the driver connection
	Close() error
}

// JobPayload represents the serialized job data
type JobPayload struct {
	Type       string          `json:"type"`
	Data       json.RawMessage `json:"data"`
	Attempts   int             `json:"attempts"`
	MaxRetries int             `json:"max_retries"`
	CreatedAt  time.Time       `json:"created_at"`
}

// Manager manages queue operations
type Manager struct {
	mu           sync.RWMutex
	drivers      map[string]Driver
	defaultQueue string
	jobRegistry  map[string]reflect.Type
}

var (
	manager *Manager
	once    sync.Once
)

// Global returns the global queue manager
func Global() *Manager {
	once.Do(func() {
		manager = &Manager{
			drivers:      make(map[string]Driver),
			defaultQueue: "default",
			jobRegistry:  make(map[string]reflect.Type),
		}
		// Register default sync driver
		manager.drivers["sync"] = NewSyncDriver()
	})
	return manager
}

// SetDefaultQueue sets the default queue name
func (m *Manager) SetDefaultQueue(queue string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultQueue = queue
}

// RegisterDriver registers a queue driver
func (m *Manager) RegisterDriver(name string, driver Driver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[name] = driver
}

// Driver returns a driver by name
func (m *Manager) Driver(name string) Driver {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.drivers[name]
}

// DefaultDriver returns the default driver
func (m *Manager) DefaultDriver() Driver {
	return m.Driver("sync")
}

// RegisterJob registers a job type for deserialization
func (m *Manager) RegisterJob(job Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t := reflect.TypeOf(job)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	m.jobRegistry[t.Name()] = t
}

// createJob creates a job instance from type name
func (m *Manager) createJob(typeName string) (Job, error) {
	m.mu.RLock()
	t, ok := m.jobRegistry[typeName]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown job type: %s", typeName)
	}

	return reflect.New(t).Interface().(Job), nil
}

// Dispatch dispatches a job to the queue
func (m *Manager) Dispatch(ctx context.Context, job Job) error {
	return m.DispatchTo(ctx, m.defaultQueue, job)
}

// DispatchTo dispatches a job to a specific queue
func (m *Manager) DispatchTo(ctx context.Context, queue string, job Job) error {
	// Check if job specifies its own queue
	if jq, ok := job.(JobWithQueue); ok {
		queue = jq.Queue()
	}

	payload, err := m.serializeJob(job)
	if err != nil {
		return err
	}

	driver := m.DefaultDriver()

	// Check if job has delay
	if jd, ok := job.(JobWithDelay); ok {
		if delay := jd.Delay(); delay > 0 {
			return driver.PushDelayed(ctx, queue, payload, delay)
		}
	}

	return driver.Push(ctx, queue, payload)
}

// Later dispatches a job with a delay
func (m *Manager) Later(ctx context.Context, delay time.Duration, job Job) error {
	return m.LaterTo(ctx, m.defaultQueue, delay, job)
}

// LaterTo dispatches a job to a specific queue with a delay
func (m *Manager) LaterTo(ctx context.Context, queue string, delay time.Duration, job Job) error {
	payload, err := m.serializeJob(job)
	if err != nil {
		return err
	}

	return m.DefaultDriver().PushDelayed(ctx, queue, payload, delay)
}

// serializeJob serializes a job to JSON
func (m *Manager) serializeJob(job Job) ([]byte, error) {
	t := reflect.TypeOf(job)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	data, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	maxRetries := 3
	if jr, ok := job.(JobWithRetry); ok {
		maxRetries = jr.MaxRetries()
	}

	payload := JobPayload{
		Type:       t.Name(),
		Data:       data,
		Attempts:   0,
		MaxRetries: maxRetries,
		CreatedAt:  time.Now(),
	}

	return json.Marshal(payload)
}

// deserializeJob deserializes a job from JSON
func (m *Manager) deserializeJob(data []byte) (Job, *JobPayload, error) {
	var payload JobPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, nil, err
	}

	job, err := m.createJob(payload.Type)
	if err != nil {
		return nil, nil, err
	}

	if err := json.Unmarshal(payload.Data, job); err != nil {
		return nil, nil, err
	}

	return job, &payload, nil
}

// Size returns the size of a queue
func (m *Manager) Size(ctx context.Context, queue string) (int64, error) {
	return m.DefaultDriver().Size(ctx, queue)
}

// Clear clears a queue
func (m *Manager) Clear(ctx context.Context, queue string) error {
	return m.DefaultDriver().Clear(ctx, queue)
}

// --- Convenience functions ---

// Dispatch dispatches a job using the global manager
func Dispatch(ctx context.Context, job Job) error {
	return Global().Dispatch(ctx, job)
}

// DispatchTo dispatches a job to a specific queue
func DispatchTo(ctx context.Context, queue string, job Job) error {
	return Global().DispatchTo(ctx, queue, job)
}

// Later dispatches a job with a delay
func Later(ctx context.Context, delay time.Duration, job Job) error {
	return Global().Later(ctx, delay, job)
}

// LaterTo dispatches a job to a specific queue with a delay
func LaterTo(ctx context.Context, queue string, delay time.Duration, job Job) error {
	return Global().LaterTo(ctx, queue, delay, job)
}

// RegisterJob registers a job type
func RegisterJob(job Job) {
	Global().RegisterJob(job)
}

// Size returns the size of a queue
func Size(ctx context.Context, queue string) (int64, error) {
	return Global().Size(ctx, queue)
}

// Clear clears a queue
func Clear(ctx context.Context, queue string) error {
	return Global().Clear(ctx, queue)
}

// --- Sync Driver (executes jobs immediately) ---

// SyncDriver executes jobs synchronously (for development/testing)
type SyncDriver struct {
	mu     sync.Mutex
	queues map[string][][]byte
}

// NewSyncDriver creates a new sync driver
func NewSyncDriver() *SyncDriver {
	return &SyncDriver{
		queues: make(map[string][][]byte),
	}
}

// Push adds a job and executes it immediately
func (d *SyncDriver) Push(ctx context.Context, queue string, payload []byte) error {
	// For sync driver, we execute immediately
	job, jobPayload, err := Global().deserializeJob(payload)
	if err != nil {
		// If we can't deserialize, just store it
		d.mu.Lock()
		d.queues[queue] = append(d.queues[queue], payload)
		d.mu.Unlock()
		return nil
	}

	// Execute the job
	if err := job.Handle(ctx); err != nil {
		// Check retry
		if jobPayload.Attempts < jobPayload.MaxRetries {
			jobPayload.Attempts++
			newPayload, _ := json.Marshal(jobPayload)
			d.mu.Lock()
			d.queues[queue] = append(d.queues[queue], newPayload)
			d.mu.Unlock()
		}
		return err
	}

	return nil
}

// PushDelayed adds a job with delay (for sync, we just delay then execute)
func (d *SyncDriver) PushDelayed(ctx context.Context, queue string, payload []byte, delay time.Duration) error {
	time.Sleep(delay)
	return d.Push(ctx, queue, payload)
}

// Pop retrieves the next job
func (d *SyncDriver) Pop(ctx context.Context, queue string) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	jobs := d.queues[queue]
	if len(jobs) == 0 {
		return nil, errors.New("queue is empty")
	}

	payload := jobs[0]
	d.queues[queue] = jobs[1:]
	return payload, nil
}

// Size returns queue size
func (d *SyncDriver) Size(ctx context.Context, queue string) (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return int64(len(d.queues[queue])), nil
}

// Clear clears the queue
func (d *SyncDriver) Clear(ctx context.Context, queue string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.queues[queue] = nil
	return nil
}

// Close closes the driver
func (d *SyncDriver) Close() error {
	return nil
}

// --- Memory Driver (stores jobs for worker processing) ---

// MemoryDriver stores jobs in memory for async processing
type MemoryDriver struct {
	mu      sync.Mutex
	queues  map[string]chan []byte
	bufSize int
	closed  bool
}

// NewMemoryDriver creates a new memory driver
func NewMemoryDriver(bufferSize int) *MemoryDriver {
	return &MemoryDriver{
		queues:  make(map[string]chan []byte),
		bufSize: bufferSize,
	}
}

func (d *MemoryDriver) getQueue(name string) chan []byte {
	d.mu.Lock()
	defer d.mu.Unlock()

	if ch, ok := d.queues[name]; ok {
		return ch
	}

	ch := make(chan []byte, d.bufSize)
	d.queues[name] = ch
	return ch
}

// Push adds a job to the queue
func (d *MemoryDriver) Push(ctx context.Context, queue string, payload []byte) error {
	if d.closed {
		return errors.New("driver is closed")
	}

	select {
	case d.getQueue(queue) <- payload:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// PushDelayed adds a job with delay
func (d *MemoryDriver) PushDelayed(ctx context.Context, queue string, payload []byte, delay time.Duration) error {
	go func() {
		time.Sleep(delay)
		d.Push(context.Background(), queue, payload)
	}()
	return nil
}

// Pop retrieves the next job
func (d *MemoryDriver) Pop(ctx context.Context, queue string) ([]byte, error) {
	select {
	case payload := <-d.getQueue(queue):
		return payload, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Size returns approximate queue size
func (d *MemoryDriver) Size(ctx context.Context, queue string) (int64, error) {
	return int64(len(d.getQueue(queue))), nil
}

// Clear clears the queue
func (d *MemoryDriver) Clear(ctx context.Context, queue string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if ch, ok := d.queues[queue]; ok {
		close(ch)
		d.queues[queue] = make(chan []byte, d.bufSize)
	}
	return nil
}

// Close closes all queues
func (d *MemoryDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.closed = true
	for _, ch := range d.queues {
		close(ch)
	}
	return nil
}
