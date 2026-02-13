package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Plugin represents a plugin interface
type Plugin interface {
	Name() string
	Version() string
	Description() string
}

// Discover automatically discovers installed plugins in PATH
// Looks for executables matching pattern: zgo-*
func Discover() []PluginInfo {
	var plugins []PluginInfo

	// Get PATH environment variable
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return plugins
	}

	// Split PATH into directories
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	// Add current directory to search paths for development convenience
	if wd, err := os.Getwd(); err == nil {
		paths = append([]string{wd}, paths...)
	}

	// Track discovered plugins to avoid duplicates
	discovered := make(map[string]bool)

	for _, dir := range paths {
		// Find all zgo-* executables
		pattern := filepath.Join(dir, "zgo-*")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			// Check if file is executable
			if !isExecutable(match) {
				continue
			}

			// Extract plugin name (remove zgo- prefix)
			name := strings.TrimPrefix(filepath.Base(match), "zgo-")

			// Skip if already discovered
			if discovered[name] {
				continue
			}

			discovered[name] = true

			// Get plugin info
			info := PluginInfo{
				Name:   name,
				Binary: match,
			}

			// Try to get version and description
			if version := getPluginVersion(match); version != "" {
				info.Version = version
			}

			plugins = append(plugins, info)
		}
	}

	return plugins
}

// PluginInfo contains information about a discovered plugin
type PluginInfo struct {
	Name        string
	Binary      string
	Version     string
	Description string
}

// isExecutable checks if a file is executable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if it's a regular file and has execute permission
	mode := info.Mode()
	return mode.IsRegular() && (mode.Perm()&0111 != 0)
}

// getPluginVersion attempts to get plugin version
func getPluginVersion(binary string) string {
	cmd := exec.Command(binary, "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// Execute runs a plugin command
func Execute(pluginName string, args []string) error {
	binary := "zgo-" + pluginName
	cmd := exec.Command(binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// IsInstalled checks if a plugin is installed
func IsInstalled(pluginName string) bool {
	binary := "zgo-" + pluginName
	_, err := exec.LookPath(binary)
	return err == nil
}
