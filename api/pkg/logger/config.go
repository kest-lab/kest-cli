package logger

import "time"

// Config holds logger configuration
type Config struct {
	// Channel is the default log channel
	Channel string

	// Level is the minimum logging level
	Level Level

	// Path is the directory for log files
	Path string

	// File is the log file name pattern (supports date: {Y}-{m}-{d}.log)
	File string

	// MaxSize is the maximum size in MB before rotation
	MaxSize int

	// MaxAge is the maximum days to retain old logs
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int

	// Compress determines if rotated files should be compressed
	Compress bool

	// StdoutPrint outputs to stdout (default: true in debug mode)
	StdoutPrint bool

	// ColorEnabled enables colored output for console
	ColorEnabled bool

	// TimeFormat is the time format for log entries
	TimeFormat string

	// Stack enables stack driver logging (multi-channel)
	Stack bool

	// JSON outputs logs in JSON format
	JSON bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Channel:      "stack",
		Level:        LevelDebug,
		Path:         "storage/logs",
		File:         "{Y}-{m}-{d}.log",
		MaxSize:      100,
		MaxAge:       14,
		MaxBackups:   7,
		Compress:     true,
		StdoutPrint:  true,
		ColorEnabled: true,
		TimeFormat:   time.RFC3339,
		JSON:         false,
	}
}

// ProductionConfig returns configuration optimized for production
func ProductionConfig() *Config {
	cfg := DefaultConfig()
	cfg.Level = LevelWarning
	cfg.StdoutPrint = false
	cfg.ColorEnabled = false
	cfg.JSON = true
	return cfg
}
