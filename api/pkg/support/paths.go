package support

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	basePath     string
	basePathOnce sync.Once
)

// SetBasePath sets the application base path
func SetBasePath(path string) {
	basePath = path
}

// getBasePath returns the application base path
func getBasePath() string {
	basePathOnce.Do(func() {
		if basePath == "" {
			// Try to get from current working directory
			if wd, err := os.Getwd(); err == nil {
				basePath = wd
			}
		}
	})
	return basePath
}

// BasePath returns the path to the base of the install
func BasePath(paths ...string) string {
	base := getBasePath()
	if len(paths) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, paths...)...)
}

// AppPath returns the path to the app directory
func AppPath(paths ...string) string {
	return BasePath(append([]string{"app"}, paths...)...)
}

// ConfigPath returns the path to the config directory
func ConfigPath(paths ...string) string {
	return BasePath(append([]string{"config"}, paths...)...)
}

// DatabasePath returns the path to the database directory
func DatabasePath(paths ...string) string {
	return BasePath(append([]string{"database"}, paths...)...)
}

// PublicPath returns the path to the public directory
func PublicPath(paths ...string) string {
	return BasePath(append([]string{"public"}, paths...)...)
}

// ResourcePath returns the path to the resources directory
func ResourcePath(paths ...string) string {
	return BasePath(append([]string{"resources"}, paths...)...)
}

// StoragePath returns the path to the storage directory
func StoragePath(paths ...string) string {
	return BasePath(append([]string{"storage"}, paths...)...)
}

// StorageAppPath returns the path to storage/app
func StorageAppPath(paths ...string) string {
	return StoragePath(append([]string{"app"}, paths...)...)
}

// StorageLogsPath returns the path to storage/logs
func StorageLogsPath(paths ...string) string {
	return StoragePath(append([]string{"logs"}, paths...)...)
}

// StorageCachePath returns the path to storage/cache
func StorageCachePath(paths ...string) string {
	return StoragePath(append([]string{"cache"}, paths...)...)
}

// LangPath returns the path to the lang directory
func LangPath(paths ...string) string {
	return BasePath(append([]string{"lang"}, paths...)...)
}

// TestsPath returns the path to the tests directory
func TestsPath(paths ...string) string {
	return BasePath(append([]string{"tests"}, paths...)...)
}

// EnsureDirectoryExists ensures a directory exists, creates it if not
func EnsureDirectoryExists(path string, perm ...os.FileMode) error {
	p := os.FileMode(0755)
	if len(perm) > 0 {
		p = perm[0]
	}
	return os.MkdirAll(path, p)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirectoryExists checks if a directory exists
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
