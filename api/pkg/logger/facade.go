package logger

import "context"

// Package-level functions that delegate to the default logger
// These provide a static facade pattern for convenient access

// Channel returns a logger for the specified channel
func Channel(name string) *Logger {
	return Default().Channel(name)
}

// WithContext returns a logger with additional context
func WithContext(ctx map[string]any) *Logger {
	return Default().WithContext(ctx)
}

// Debug logs a debug message
func Debug(msg string, ctx ...map[string]any) {
	Default().Debug(msg, ctx...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...any) {
	Default().Debugf(format, args...)
}

// Info logs an info message
func Info(msg string, ctx ...map[string]any) {
	Default().Info(msg, ctx...)
}

// Infof logs a formatted info message
func Infof(format string, args ...any) {
	Default().Infof(format, args...)
}

// Notice logs a notice message
func Notice(msg string, ctx ...map[string]any) {
	Default().Notice(msg, ctx...)
}

// Noticef logs a formatted notice message
func Noticef(format string, args ...any) {
	Default().Noticef(format, args...)
}

// Warning logs a warning message
func Warning(msg string, ctx ...map[string]any) {
	Default().Warning(msg, ctx...)
}

// Warn is an alias for Warning
func Warn(msg string, ctx ...map[string]any) {
	Default().Warning(msg, ctx...)
}

// Warningf logs a formatted warning message
func Warningf(format string, args ...any) {
	Default().Warningf(format, args...)
}

// Error logs an error message
func Error(msg string, ctx ...map[string]any) {
	Default().Error(msg, ctx...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...any) {
	Default().Errorf(format, args...)
}

// Critical logs a critical message
func Critical(msg string, ctx ...map[string]any) {
	Default().Critical(msg, ctx...)
}

// Criticalf logs a formatted critical message
func Criticalf(format string, args ...any) {
	Default().Criticalf(format, args...)
}

// Alert logs an alert message
func Alert(msg string, ctx ...map[string]any) {
	Default().Alert(msg, ctx...)
}

// Alertf logs a formatted alert message
func Alertf(format string, args ...any) {
	Default().Alertf(format, args...)
}

// Emergency logs an emergency message
func Emergency(msg string, ctx ...map[string]any) {
	Default().Emergency(msg, ctx...)
}

// Emergencyf logs a formatted emergency message
func Emergencyf(format string, args ...any) {
	Default().Emergencyf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, ctx ...map[string]any) {
	Default().Fatal(msg, ctx...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...any) {
	Default().Fatalf(format, args...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return Default().Sync()
}

// Close closes all handlers
func Close() error {
	return Default().Close()
}

// ContextLog logs with context.Context support
type ContextLog struct {
	ctx context.Context
	l   *Logger
}

// Ctx returns a context-aware logger
func Ctx(ctx context.Context) *ContextLog {
	return &ContextLog{ctx: ctx, l: Default()}
}

// Debug logs a debug message with context
func (c *ContextLog) Debug(msg string, ctx ...map[string]any) {
	c.l.log(c.ctx, LevelDebug, msg, getContext(ctx))
}

// Info logs an info message with context
func (c *ContextLog) Info(msg string, ctx ...map[string]any) {
	c.l.log(c.ctx, LevelInfo, msg, getContext(ctx))
}

// Warning logs a warning message with context
func (c *ContextLog) Warning(msg string, ctx ...map[string]any) {
	c.l.log(c.ctx, LevelWarning, msg, getContext(ctx))
}

// Error logs an error message with context
func (c *ContextLog) Error(msg string, ctx ...map[string]any) {
	c.l.log(c.ctx, LevelError, msg, getContext(ctx))
}
