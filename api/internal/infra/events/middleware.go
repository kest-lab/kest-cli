// Package events provides event middleware for cross-cutting concerns
package events

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Logger interface for logging middleware
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// LogLevel represents logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// LoggingMiddleware creates a middleware that logs event handling
func LoggingMiddleware(logger Logger, level LogLevel) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			start := time.Now()
			meta := event.Metadata()

			logMsg := fmt.Sprintf("handling event: %s", event.EventName())
			logArgs := []any{
				"event_id", meta.ID,
				"event_name", event.EventName(),
				"correlation_id", meta.CorrelationID,
				"source", meta.Source,
			}

			// Log start based on level
			switch level {
			case LogLevelDebug:
				logger.Debug(logMsg, logArgs...)
			case LogLevelInfo:
				logger.Info(logMsg, logArgs...)
			}

			err := next(ctx, event)
			duration := time.Since(start)

			logArgs = append(logArgs, "duration_ms", duration.Milliseconds())

			if err != nil {
				logArgs = append(logArgs, "error", err.Error())
				logger.Error("event handling failed: "+event.EventName(), logArgs...)
			} else {
				switch level {
				case LogLevelDebug:
					logger.Debug("event handled: "+event.EventName(), logArgs...)
				case LogLevelInfo:
					logger.Info("event handled: "+event.EventName(), logArgs...)
				}
			}

			return err
		}
	}
}

// SlogAdapter adapts slog.Logger to Logger interface
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter creates a new slog adapter
func NewSlogAdapter(logger *slog.Logger) *SlogAdapter {
	return &SlogAdapter{logger: logger}
}

func (a *SlogAdapter) Debug(msg string, args ...any) { a.logger.Debug(msg, args...) }
func (a *SlogAdapter) Info(msg string, args ...any)  { a.logger.Info(msg, args...) }
func (a *SlogAdapter) Warn(msg string, args ...any)  { a.logger.Warn(msg, args...) }
func (a *SlogAdapter) Error(msg string, args ...any) { a.logger.Error(msg, args...) }

// TracingMiddleware creates a middleware that adds OpenTelemetry tracing
func TracingMiddleware(tracer trace.Tracer) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			meta := event.Metadata()

			ctx, span := tracer.Start(ctx, "event.handle."+event.EventName(),
				trace.WithAttributes(
					attribute.String("event.id", meta.ID),
					attribute.String("event.name", event.EventName()),
					attribute.String("event.correlation_id", meta.CorrelationID),
					attribute.String("event.causation_id", meta.CausationID),
					attribute.String("event.source", meta.Source),
					attribute.String("event.timestamp", meta.Timestamp.Format(time.RFC3339)),
				),
			)
			defer span.End()

			err := next(ctx, event)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}

			return err
		}
	}
}

// RecoveryMiddleware creates a middleware that recovers from panics
func RecoveryMiddleware(onPanic func(event Event, recovered interface{}, stack []byte)) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) (err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := debug.Stack()
					if onPanic != nil {
						onPanic(event, r, stack)
					}
					err = fmt.Errorf("panic recovered in event handler: %v", r)
				}
			}()
			return next(ctx, event)
		}
	}
}

// RetryMiddleware creates a middleware that retries failed handlers
func RetryMiddleware(maxRetries int, backoff time.Duration) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			var lastErr error

			for attempt := 0; attempt <= maxRetries; attempt++ {
				// Check context cancellation
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				lastErr = next(ctx, event)
				if lastErr == nil {
					return nil
				}

				// Don't sleep after the last attempt
				if attempt < maxRetries {
					// Exponential backoff
					sleepDuration := backoff * time.Duration(1<<uint(attempt))
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(sleepDuration):
					}
				}
			}

			return fmt.Errorf("event handler failed after %d retries: %w", maxRetries+1, lastErr)
		}
	}
}

// TimeoutMiddleware creates a middleware that enforces a timeout on handlers
func TimeoutMiddleware(timeout time.Duration) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				done <- next(ctx, event)
			}()

			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				return fmt.Errorf("event handler timed out after %v: %w", timeout, ctx.Err())
			}
		}
	}
}

// MetricsMiddleware creates a middleware that records metrics
type MetricsRecorder interface {
	RecordEventHandled(eventName string, duration time.Duration, success bool)
}

func MetricsMiddleware(recorder MetricsRecorder) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			start := time.Now()
			err := next(ctx, event)
			duration := time.Since(start)

			recorder.RecordEventHandled(event.EventName(), duration, err == nil)
			return err
		}
	}
}

// FilterMiddleware creates a middleware that filters events based on a predicate
func FilterMiddleware(predicate func(Event) bool) EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			if !predicate(event) {
				return nil // Skip this event
			}
			return next(ctx, event)
		}
	}
}

// CorrelationMiddleware ensures events have correlation IDs
func CorrelationMiddleware() EventMiddleware {
	return func(next EventHandler) EventHandler {
		return func(ctx context.Context, event Event) error {
			meta := event.Metadata()

			// If no correlation ID, use the event ID as correlation ID
			if meta.CorrelationID == "" {
				if baseEvent, ok := event.(interface{ SetCorrelationID(string) }); ok {
					baseEvent.SetCorrelationID(meta.ID)
				}
			}

			return next(ctx, event)
		}
	}
}
