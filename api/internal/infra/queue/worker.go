package queue

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
)

// WorkerConfig holds worker configuration
type WorkerConfig struct {
	// Queue name to process
	Queue string

	// Number of concurrent workers
	Concurrency int

	// Sleep duration when queue is empty
	Sleep time.Duration

	// Maximum job processing time
	Timeout time.Duration

	// Stop after processing this many jobs (0 = unlimited)
	MaxJobs int

	// Handler for failed jobs
	FailedJobHandler func(ctx context.Context, payload *JobPayload, err error)

	// Handler called before job processing
	BeforeJob func(ctx context.Context, payload *JobPayload)

	// Handler called after job processing
	AfterJob func(ctx context.Context, payload *JobPayload, err error)
}

// DefaultWorkerConfig returns default worker configuration
func DefaultWorkerConfig() WorkerConfig {
	return WorkerConfig{
		Queue:       "default",
		Concurrency: 1,
		Sleep:       time.Second,
		Timeout:     60 * time.Second,
		MaxJobs:     0,
	}
}

// Worker processes jobs from a queue
type Worker struct {
	config  WorkerConfig
	manager *Manager
	driver  Driver
	stop    chan struct{}
	wg      sync.WaitGroup
	running bool
	mu      sync.Mutex
}

// NewWorker creates a new worker
func NewWorker(config WorkerConfig) *Worker {
	if config.Concurrency < 1 {
		config.Concurrency = 1
	}
	if config.Sleep == 0 {
		config.Sleep = time.Second
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}

	return &Worker{
		config:  config,
		manager: Global(),
		driver:  Global().DefaultDriver(),
		stop:    make(chan struct{}),
	}
}

// SetDriver sets the driver for the worker
func (w *Worker) SetDriver(driver Driver) {
	w.driver = driver
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return nil
	}
	w.running = true
	w.mu.Unlock()

	for i := 0; i < w.config.Concurrency; i++ {
		w.wg.Add(1)
		go w.work(ctx, i)
	}

	return nil
}

// Stop stops the worker gracefully
func (w *Worker) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	w.mu.Unlock()

	close(w.stop)
	w.wg.Wait()
}

// work is the main worker loop
func (w *Worker) work(ctx context.Context, id int) {
	defer w.wg.Done()

	jobsProcessed := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stop:
			return
		default:
			// Check max jobs
			if w.config.MaxJobs > 0 && jobsProcessed >= w.config.MaxJobs {
				return
			}

			// Try to get a job
			payload, err := w.driver.Pop(ctx, w.config.Queue)
			if err != nil {
				// Queue empty or error, sleep and retry
				time.Sleep(w.config.Sleep)
				continue
			}

			// Process the job
			w.processJob(ctx, payload)
			jobsProcessed++
		}
	}
}

// processJob processes a single job
func (w *Worker) processJob(ctx context.Context, data []byte) {
	var payload JobPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Failed to unmarshal job payload: %v", err)
		return
	}

	// Call before handler
	if w.config.BeforeJob != nil {
		w.config.BeforeJob(ctx, &payload)
	}

	// Create job instance
	job, err := w.manager.createJob(payload.Type)
	if err != nil {
		log.Printf("Failed to create job instance: %v", err)
		w.handleFailedJob(ctx, &payload, err)
		return
	}

	// Unmarshal job data
	if err := json.Unmarshal(payload.Data, job); err != nil {
		log.Printf("Failed to unmarshal job data: %v", err)
		w.handleFailedJob(ctx, &payload, err)
		return
	}

	// Create context with timeout
	jobCtx, cancel := context.WithTimeout(ctx, w.config.Timeout)
	defer cancel()

	// Execute the job
	payload.Attempts++
	err = job.Handle(jobCtx)

	// Call after handler
	if w.config.AfterJob != nil {
		w.config.AfterJob(ctx, &payload, err)
	}

	if err != nil {
		// Check if we should retry
		if payload.Attempts < payload.MaxRetries {
			// Get retry delay
			retryDelay := time.Second * time.Duration(payload.Attempts)
			if jr, ok := job.(JobWithRetry); ok {
				retryDelay = jr.RetryDelay()
			}

			// Re-queue the job
			newPayload, _ := json.Marshal(payload)
			w.driver.PushDelayed(ctx, w.config.Queue, newPayload, retryDelay)
		} else {
			// Max retries exceeded
			w.handleFailedJob(ctx, &payload, err)
		}
	}
}

// handleFailedJob handles a failed job
func (w *Worker) handleFailedJob(ctx context.Context, payload *JobPayload, err error) {
	if w.config.FailedJobHandler != nil {
		w.config.FailedJobHandler(ctx, payload, err)
	} else {
		log.Printf("Job %s failed after %d attempts: %v", payload.Type, payload.Attempts, err)
	}
}

// ProcessNext processes the next job in queue (useful for testing)
func (w *Worker) ProcessNext(ctx context.Context) error {
	payload, err := w.driver.Pop(ctx, w.config.Queue)
	if err != nil {
		return err
	}

	w.processJob(ctx, payload)
	return nil
}

// Daemon starts the worker as a daemon process
func Daemon(ctx context.Context, config WorkerConfig) error {
	worker := NewWorker(config)
	return worker.Start(ctx)
}

// Work processes jobs from the default queue
func Work(ctx context.Context, queue string, concurrency int) error {
	config := DefaultWorkerConfig()
	config.Queue = queue
	config.Concurrency = concurrency

	worker := NewWorker(config)
	return worker.Start(ctx)
}
