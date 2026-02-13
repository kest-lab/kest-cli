package logger

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

// SentryHandler sends logs to Sentry
type SentryHandler struct {
	level Level
}

// NewSentryHandler creates a new Sentry handler
func NewSentryHandler(dsn string, appEnv string) (*SentryHandler, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      appEnv,
		AttachStacktrace: true,
	})
	if err != nil {
		return nil, err
	}

	return &SentryHandler{
		level: LevelError, // Default to Error and above
	}, nil
}

// Handle implements the Handler interface
func (h *SentryHandler) Handle(ctx context.Context, entry *Entry) error {
	if entry.Level < h.level {
		return nil
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(h.mapLevel(entry.Level))
		scope.SetContext("details", entry.Context)
		if entry.Channel != "" {
			scope.SetTag("channel", entry.Channel)
		}
		if entry.RequestID != "" {
			scope.SetTag("request_id", entry.RequestID)
		}

		sentry.CaptureMessage(entry.Message)
	})

	return nil
}

// Close flushes sentry events
func (h *SentryHandler) Close() error {
	sentry.Flush(2 * time.Second)
	return nil
}

func (h *SentryHandler) mapLevel(l Level) sentry.Level {
	switch l {
	case LevelDebug:
		return sentry.LevelDebug
	case LevelInfo, LevelNotice:
		return sentry.LevelInfo
	case LevelWarning:
		return sentry.LevelWarning
	case LevelError:
		return sentry.LevelError
	default:
		return sentry.LevelFatal
	}
}
