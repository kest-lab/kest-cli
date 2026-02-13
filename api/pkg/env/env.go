package env

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Priority order (highest to lowest):
// 1. System environment variables (Docker, K8s, shell export)
// 2. .env.{APP_ENV}.local (environment-specific local overrides, not committed)
// 3. .env.local (local overrides, not committed)
// 4. .env.{APP_ENV} (environment-specific, e.g., .env.production)
// 5. .env (base configuration)
// 6. Default values in code

var (
	loaded     bool
	loadedOnce sync.Once
	appEnv     string
)

// Load loads environment files respecting priority.
// System environment variables always have highest priority.
// This function is idempotent - calling it multiple times has no effect.
func Load() {
	loadedOnce.Do(func() {
		loadEnvFiles()
		loaded = true
	})
}

// LoadFresh forces reload of environment files.
// Useful for testing.
func LoadFresh() {
	loadedOnce = sync.Once{}
	loaded = false
	Load()
}

func loadEnvFiles() {
	// Capture system environment variables BEFORE loading any .env files
	systemEnv := captureSystemEnv()

	// Determine environment from system env first (highest priority)
	appEnv = systemEnv["APP_ENV"]
	if appEnv == "" {
		appEnv = systemEnv["GO_ENV"]
	}
	if appEnv == "" {
		appEnv = systemEnv["GIN_MODE"]
	}

	// Check for explicit environment file (ZGO_ENV_FILE)
	if envFile := systemEnv["ZGO_ENV_FILE"]; envFile != "" {
		if _, err := os.Stat(envFile); err == nil {
			_ = godotenv.Load(envFile)
			// Allow APP_ENV from this specific file if not set by system env
			if appEnv == "" {
				appEnv = os.Getenv("APP_ENV")
			}
		}
	}

	// Load .env files in order (lowest priority first)
	// godotenv.Load does NOT override existing values
	files := []string{".env"}

	// Load base first, then check APP_ENV from it if not set
	_ = godotenv.Load(".env")

	// Re-check APP_ENV after loading base .env
	if appEnv == "" {
		appEnv = os.Getenv("APP_ENV")
	}
	if appEnv == "" {
		appEnv = os.Getenv("GO_ENV")
	}

	// Environment-specific file
	if appEnv != "" {
		files = append(files, ".env."+appEnv)
	}

	// Local overrides
	files = append(files, ".env.local")

	// Environment-specific local overrides
	if appEnv != "" {
		files = append(files, ".env."+appEnv+".local")
	}

	// Load remaining files (godotenv.Load won't override existing)
	for _, file := range files[1:] {
		if _, err := os.Stat(file); err == nil {
			_ = godotenv.Load(file)
		}
	}

	// Restore system environment variables (highest priority)
	for key, value := range systemEnv {
		os.Setenv(key, value)
	}
}

// captureSystemEnv captures current environment variables before .env loading
func captureSystemEnv() map[string]string {
	result := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// Get returns the value of an environment variable.
// If the variable is not set, returns the default value.
func Get(key string, defaultValue ...string) string {
	Load()
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetOrFail returns the value of an environment variable.
// Panics if the variable is not set.
func GetOrFail(key string) string {
	Load()
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic("Environment variable " + key + " is required but not set")
}

// GetBool returns the boolean value of an environment variable.
// Recognizes: true, false, 1, 0, yes, no, on, off (case-insensitive)
func GetBool(key string, defaultValue ...bool) bool {
	Load()
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off", "":
		return false
	}

	// Try parsing as bool
	b, err := strconv.ParseBool(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return b
}

// GetInt returns the integer value of an environment variable.
func GetInt(key string, defaultValue ...int) int {
	Load()
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	i, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return i
}

// GetInt64 returns the int64 value of an environment variable.
func GetInt64(key string, defaultValue ...int64) int64 {
	Load()
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	i, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return i
}

// GetFloat returns the float64 value of an environment variable.
func GetFloat(key string, defaultValue ...float64) float64 {
	Load()
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	f, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return f
}

// GetDuration returns the time.Duration value of an environment variable.
func GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	Load()
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	d, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return d
}

// GetSlice returns a slice from a comma-separated environment variable.
func GetSlice(key string, defaultValue ...[]string) []string {
	Load()
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// Set sets an environment variable.
// This is useful for testing.
func Set(key, value string) {
	os.Setenv(key, value)
}

// Unset removes an environment variable.
func Unset(key string) {
	os.Unsetenv(key)
}

// AppEnv returns the current application environment.
// Returns "development" if not set.
func AppEnv() string {
	Load()
	if appEnv != "" {
		return appEnv
	}
	env := Get("APP_ENV", "development")
	return env
}

// IsProduction returns true if running in production mode.
func IsProduction() bool {
	env := strings.ToLower(AppEnv())
	return env == "production" || env == "prod" || env == "release"
}

// IsDevelopment returns true if running in development mode.
func IsDevelopment() bool {
	env := strings.ToLower(AppEnv())
	return env == "development" || env == "dev" || env == "local" || env == "debug"
}

// IsTesting returns true if running in test mode.
func IsTesting() bool {
	env := strings.ToLower(AppEnv())
	return env == "testing" || env == "test"
}

// IsLocal returns true if running locally (development or testing).
func IsLocal() bool {
	return IsDevelopment() || IsTesting()
}
