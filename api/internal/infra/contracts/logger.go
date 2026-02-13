package contracts

import "context"

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, ctx ...map[string]any)
	Debugf(format string, args ...any)
	Info(msg string, ctx ...map[string]any)
	Infof(format string, args ...any)
	Warning(msg string, ctx ...map[string]any)
	Warningf(format string, args ...any)
	Error(msg string, ctx ...map[string]any)
	Errorf(format string, args ...any)
	Channel(name string) Logger
	WithContext(ctx map[string]any) Logger
}

// ContextLogger defines context-aware logging
type ContextLogger interface {
	Debug(ctx context.Context, msg string, data ...map[string]any)
	Info(ctx context.Context, msg string, data ...map[string]any)
	Warning(ctx context.Context, msg string, data ...map[string]any)
	Error(ctx context.Context, msg string, data ...map[string]any)
}
