package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Handler defines the interface for log handlers
type Handler interface {
	Handle(ctx context.Context, entry *Entry) error
	Close() error
}

// Entry represents a log entry
type Entry struct {
	Level     Level
	Message   string
	Context   map[string]any
	Time      time.Time
	Channel   string
	RequestID string
	TraceID   string
}

// ConsoleHandler outputs logs to console with optional colors
type ConsoleHandler struct {
	mu           sync.Mutex
	writer       io.Writer
	level        Level
	colorEnabled bool
	timeFormat   string
}

// NewConsoleHandler creates a new console handler
func NewConsoleHandler(level Level, colorEnabled bool) *ConsoleHandler {
	return &ConsoleHandler{
		writer:       os.Stdout,
		level:        level,
		colorEnabled: colorEnabled,
		timeFormat:   "2006-01-02 15:04:05",
	}
}

func (h *ConsoleHandler) Handle(ctx context.Context, entry *Entry) error {
	if entry.Level < h.level {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	var sb strings.Builder

	// Time
	sb.WriteString(entry.Time.Format(h.timeFormat))
	sb.WriteString(" ")

	// Level with color
	levelStr := fmt.Sprintf("%-9s", entry.Level.String())
	if h.colorEnabled {
		sb.WriteString(levelColors[entry.Level])
		sb.WriteString(levelStr)
		sb.WriteString(colorReset)
	} else {
		sb.WriteString(levelStr)
	}
	sb.WriteString(" ")

	// Channel
	if entry.Channel != "" {
		sb.WriteString("[")
		sb.WriteString(entry.Channel)
		sb.WriteString("] ")
	}

	// Message
	sb.WriteString(entry.Message)

	// Context
	if len(entry.Context) > 0 {
		sb.WriteString(" ")
		contextBytes, _ := json.Marshal(entry.Context)
		sb.Write(contextBytes)
	}

	sb.WriteString("\n")

	_, err := h.writer.Write([]byte(sb.String()))
	return err
}

func (h *ConsoleHandler) Close() error {
	return nil
}

// FileHandler outputs logs to files with rotation support
type FileHandler struct {
	mu         sync.Mutex
	path       string
	file       string
	level      Level
	maxSize    int64
	maxAge     int
	maxBackups int
	compress   bool
	timeFormat string
	json       bool

	currentFile *os.File
	currentSize int64
	currentDate string
}

// NewFileHandler creates a new file handler
func NewFileHandler(cfg *Config) *FileHandler {
	return &FileHandler{
		path:       cfg.Path,
		file:       cfg.File,
		level:      cfg.Level,
		maxSize:    int64(cfg.MaxSize) * 1024 * 1024,
		maxAge:     cfg.MaxAge,
		maxBackups: cfg.MaxBackups,
		compress:   cfg.Compress,
		timeFormat: cfg.TimeFormat,
		json:       cfg.JSON,
	}
}

func (h *FileHandler) Handle(ctx context.Context, entry *Entry) error {
	if entry.Level < h.level {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.ensureFile(entry.Time); err != nil {
		return err
	}

	var data []byte
	var err error

	if h.json {
		data, err = h.formatJSON(entry)
	} else {
		data = h.formatText(entry)
	}

	if err != nil {
		return err
	}

	n, err := h.currentFile.Write(data)
	h.currentSize += int64(n)

	return err
}

func (h *FileHandler) formatJSON(entry *Entry) ([]byte, error) {
	record := map[string]any{
		"time":    entry.Time.Format(h.timeFormat),
		"level":   entry.Level.String(),
		"message": entry.Message,
	}
	if entry.Channel != "" {
		record["channel"] = entry.Channel
	}
	if entry.RequestID != "" {
		record["request_id"] = entry.RequestID
	}
	if len(entry.Context) > 0 {
		record["context"] = entry.Context
	}

	data, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func (h *FileHandler) formatText(entry *Entry) []byte {
	var sb strings.Builder

	sb.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	sb.WriteString("] ")

	levelStr := fmt.Sprintf("%-9s", entry.Level.String())
	sb.WriteString(levelStr)
	sb.WriteString(" ")

	if entry.Channel != "" {
		sb.WriteString("[")
		sb.WriteString(entry.Channel)
		sb.WriteString("] ")
	}

	sb.WriteString(entry.Message)

	if len(entry.Context) > 0 {
		sb.WriteString(" ")
		contextBytes, _ := json.Marshal(entry.Context)
		sb.Write(contextBytes)
	}

	sb.WriteString("\n")
	return []byte(sb.String())
}

func (h *FileHandler) ensureFile(t time.Time) error {
	date := t.Format("2006-01-02")

	// Check if we need to rotate
	needRotate := h.currentFile == nil ||
		h.currentDate != date ||
		(h.maxSize > 0 && h.currentSize >= h.maxSize)

	if !needRotate {
		return nil
	}

	// Close current file
	if h.currentFile != nil {
		h.currentFile.Close()
	}

	// Create directory if not exists
	if err := os.MkdirAll(h.path, 0755); err != nil {
		return err
	}

	// Generate filename
	filename := h.generateFilename(t)
	fullPath := filepath.Join(h.path, filename)

	// Open file
	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Get current size
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	h.currentFile = f
	h.currentSize = info.Size()
	h.currentDate = date

	return nil
}

func (h *FileHandler) generateFilename(t time.Time) string {
	filename := h.file
	filename = strings.ReplaceAll(filename, "{Y}", t.Format("2006"))
	filename = strings.ReplaceAll(filename, "{m}", t.Format("01"))
	filename = strings.ReplaceAll(filename, "{d}", t.Format("02"))
	filename = strings.ReplaceAll(filename, "{H}", t.Format("15"))
	return filename
}

func (h *FileHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.currentFile != nil {
		return h.currentFile.Close()
	}
	return nil
}

// DailyHandler is a convenience handler for daily log rotation
func NewDailyHandler(path string, level Level, days int) *FileHandler {
	return &FileHandler{
		path:       path,
		file:       "{Y}-{m}-{d}.log",
		level:      level,
		maxAge:     days,
		maxBackups: days,
		timeFormat: time.RFC3339,
		json:       false,
	}
}
