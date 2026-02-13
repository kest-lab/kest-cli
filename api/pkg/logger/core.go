package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger is the main logger struct that manages handlers
type Logger struct {
	mu       sync.RWMutex
	handlers []Handler
	config   *Config
	channel  string
	context  map[string]any
	defaultH Handler
}

var defaultLogger *Logger
var defaultMu sync.RWMutex

// New creates a new Logger instance
func New(cfg *Config) *Logger {
	l := &Logger{
		config:  cfg,
		channel: ChannelApp,
		context: make(map[string]any),
	}

	// Create default handler based on config
	if cfg.Stack {
		// Stack handler (multi-channel)
		// For simplicity, we just add a file handler and console if debug
		l.inputHandlers(cfg)
	} else {
		l.inputHandlers(cfg)
	}

	return l
}

func (l *Logger) inputHandlers(cfg *Config) {
	// Add File Handler
	fileHandler := NewFileHandler(cfg)
	l.AddHandler(fileHandler)

	// Add Console Handler if enabled
	if cfg.StdoutPrint {
		consoleHandler := NewConsoleHandler(cfg.Level, cfg.ColorEnabled)
		l.AddHandler(consoleHandler)
	}
}

// Default returns the default logger instance
func Default() *Logger {
	defaultMu.RLock()
	defer defaultMu.RUnlock()

	if defaultLogger == nil {
		// Fallback to simple default if not initialized
		defaultLogger = New(DefaultConfig())
	}
	return defaultLogger
}

// SetDefault sets the default logger instance
func SetDefault(l *Logger) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultLogger = l
}

// Channel returns a new logger instance with the specified channel
func (l *Logger) Channel(name string) *Logger {
	newLogger := l.clone()
	newLogger.channel = name
	return newLogger
}

// WithContext returns a new logger instance with the specified context
func (l *Logger) WithContext(ctx map[string]any) *Logger {
	newLogger := l.clone()
	for k, v := range ctx {
		newLogger.context[k] = v
	}
	return newLogger
}

func (l *Logger) clone() *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return &Logger{
		handlers: l.handlers, // Share handlers
		config:   l.config,
		channel:  l.channel,
		context:  copyMap(l.context),
	}
}

func copyMap(m map[string]any) map[string]any {
	cp := make(map[string]any)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// AddHandler adds a handler to the logger
func (l *Logger) AddHandler(h Handler) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.handlers = append(l.handlers, h)
}

// Log logs a message with the specified level
func (l *Logger) Log(level Level, msg string, ctx map[string]any) {
	entry := &Entry{
		Level:   level,
		Message: msg,
		Context: l.mergeContext(ctx),
		Time:    time.Now(),
		Channel: l.channel,
	}

	l.mu.RLock()
	handlers := l.handlers
	l.mu.RUnlock()

	for _, h := range handlers {
		_ = h.Handle(context.Background(), entry)
	}
}

func (l *Logger) log(ctx context.Context, level Level, msg string, logCtx map[string]any) {
	entry := &Entry{
		Level:   level,
		Message: msg,
		Context: l.mergeContext(logCtx),
		Time:    time.Now(),
		Channel: l.channel,
	}

	l.mu.RLock()
	handlers := l.handlers
	l.mu.RUnlock()

	for _, h := range handlers {
		_ = h.Handle(ctx, entry)
	}
}

func (l *Logger) mergeContext(ctx map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range l.context {
		merged[k] = v
	}
	for k, v := range ctx {
		merged[k] = v
	}
	return merged
}

// Helper methods for facade match
func getContext(ctx []map[string]any) map[string]any {
	if len(ctx) > 0 {
		return ctx[0]
	}
	return nil
}

func (l *Logger) Debug(msg string, ctx ...map[string]any) {
	l.Log(LevelDebug, msg, getContext(ctx))
}

func (l *Logger) Info(msg string, ctx ...map[string]any) {
	l.Log(LevelInfo, msg, getContext(ctx))
}

func (l *Logger) Warning(msg string, ctx ...map[string]any) {
	l.Log(LevelWarning, msg, getContext(ctx))
}

func (l *Logger) Error(msg string, ctx ...map[string]any) {
	l.Log(LevelError, msg, getContext(ctx))
}

func (l *Logger) Fatal(msg string, ctx ...map[string]any) {
	l.Log(LevelFatal, msg, getContext(ctx))
	os.Exit(1)
}

func (l *Logger) Notice(msg string, ctx ...map[string]any) {
	l.Log(LevelNotice, msg, getContext(ctx))
}

func (l *Logger) Critical(msg string, ctx ...map[string]any) {
	l.Log(LevelCritical, msg, getContext(ctx))
}

func (l *Logger) Alert(msg string, ctx ...map[string]any) {
	l.Log(LevelAlert, msg, getContext(ctx))
}

func (l *Logger) Emergency(msg string, ctx ...map[string]any) {
	l.Log(LevelEmergency, msg, getContext(ctx))
}

// Formatted methods
func (l *Logger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warningf(format string, args ...any) {
	l.Warning(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.Fatal(fmt.Sprintf(format, args...))
}

func (l *Logger) Noticef(format string, args ...any) {
	l.Notice(fmt.Sprintf(format, args...))
}

func (l *Logger) Criticalf(format string, args ...any) {
	l.Critical(fmt.Sprintf(format, args...))
}

func (l *Logger) Alertf(format string, args ...any) {
	l.Alert(fmt.Sprintf(format, args...))
}

func (l *Logger) Emergencyf(format string, args ...any) {
	l.Emergency(fmt.Sprintf(format, args...))
}

func (l *Logger) Sync() error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	// In a real impl, loop handlers and close/sync
	return nil
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, h := range l.handlers {
		h.Close()
	}
	return nil
}
