package logger

import "strings"

// Level represents logging level
type Level int

const (
	// LevelDebug detailed debug information (only in development)
	LevelDebug Level = iota
	// LevelInfo interesting events (user login, SQL logs)
	LevelInfo
	// LevelNotice normal but significant events
	LevelNotice
	// LevelWarning exceptional occurrences that are not errors
	LevelWarning
	// LevelError runtime errors
	LevelError
	// LevelCritical critical conditions
	LevelCritical
	// LevelAlert action must be taken immediately
	LevelAlert
	// LevelEmergency system is unusable
	LevelEmergency
	// LevelFatal critical error that causes exit
	LevelFatal
)

var levelNames = map[Level]string{
	LevelDebug:     "DEBUG",
	LevelInfo:      "INFO",
	LevelNotice:    "NOTICE",
	LevelWarning:   "WARNING",
	LevelError:     "ERROR",
	LevelCritical:  "CRITICAL",
	LevelAlert:     "ALERT",
	LevelEmergency: "EMERGENCY",
	LevelFatal:     "FATAL",
}

var levelColors = map[Level]string{
	LevelDebug:     "\033[37m",    // white
	LevelInfo:      "\033[32m",    // green
	LevelNotice:    "\033[36m",    // cyan
	LevelWarning:   "\033[33m",    // yellow
	LevelError:     "\033[31m",    // red
	LevelCritical:  "\033[35m",    // magenta
	LevelAlert:     "\033[91m",    // bright red
	LevelEmergency: "\033[41;97m", // white on red bg
	LevelFatal:     "\033[41;97m", // white on red bg (same as emergency)
}

const colorReset = "\033[0m"

// String returns the string representation of the level
func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return "UNKNOWN"
}

// ParseLevel parses a level string into Level
func ParseLevel(s string) Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "NOTICE":
		return LevelNotice
	case "WARNING", "WARN":
		return LevelWarning
	case "ERROR", "ERR":
		return LevelError
	case "CRITICAL", "CRIT":
		return LevelCritical
	case "ALERT":
		return LevelAlert
	case "EMERGENCY", "EMERG":
		return LevelEmergency
	default:
		return LevelDebug
	}
}
