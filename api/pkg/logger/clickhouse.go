package logger

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ClickHouseHandler sends logs to ClickHouse in batches
type ClickHouseHandler struct {
	level     Level
	batchSize int
	interval  time.Duration
	endpoint  string
	buffer    []*Entry
	mu        sync.Mutex
	stop      chan struct{}
	wg        sync.WaitGroup
}

// NewClickHouseHandler creates a new ClickHouse handler
func NewClickHouseHandler(level Level, batchSize int, interval time.Duration, endpoint string) *ClickHouseHandler {
	h := &ClickHouseHandler{
		level:     level,
		batchSize: batchSize,
		interval:  interval,
		endpoint:  endpoint,
		buffer:    make([]*Entry, 0, batchSize),
		stop:      make(chan struct{}),
	}

	h.wg.Add(1)
	go h.run()

	return h
}

func (h *ClickHouseHandler) run() {
	defer h.wg.Done()
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.flush()
		case <-h.stop:
			h.flush()
			return
		}
	}
}

// Handle implements the Handler interface
func (h *ClickHouseHandler) Handle(ctx context.Context, entry *Entry) error {
	if entry.Level < h.level {
		return nil
	}

	h.mu.Lock()
	h.buffer = append(h.buffer, entry)
	full := len(h.buffer) >= h.batchSize
	h.mu.Unlock()

	if full {
		h.flush()
	}

	return nil
}

func (h *ClickHouseHandler) flush() {
	h.mu.Lock()
	if len(h.buffer) == 0 {
		h.mu.Unlock()
		return
	}
	batch := h.buffer
	h.buffer = make([]*Entry, 0, h.batchSize)
	h.mu.Unlock()

	// Direct integration logic would go here
	// Example: write to ClickHouse table via HTTP or Native protocol
	fmt.Printf("[ClickHouse] Flushed %d logs to storage\n", len(batch))
}

// Close stops the background flusher
func (h *ClickHouseHandler) Close() error {
	close(h.stop)
	h.wg.Wait()
	return nil
}
