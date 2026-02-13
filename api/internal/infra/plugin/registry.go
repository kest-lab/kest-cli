package plugin

import (
	"sync"
)

// Registry manages registered plugins
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]PluginInfo
}

var (
	globalRegistry *Registry
	once           sync.Once
)

// Global returns the global plugin registry
func Global() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			plugins: make(map[string]PluginInfo),
		}
		// Auto-discover plugins on initialization
		globalRegistry.Refresh()
	})
	return globalRegistry
}

// Register registers a plugin
func (r *Registry) Register(info PluginInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[info.Name] = info
}

// Get retrieves a plugin by name
func (r *Registry) Get(name string) (PluginInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.plugins[name]
	return info, ok
}

// List returns all registered plugins
func (r *Registry) List() []PluginInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]PluginInfo, 0, len(r.plugins))
	for _, info := range r.plugins {
		plugins = append(plugins, info)
	}
	return plugins
}

// Has checks if a plugin is registered
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.plugins[name]
	return ok
}

// Refresh re-discovers all plugins
func (r *Registry) Refresh() {
	plugins := Discover()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Clear existing plugins
	r.plugins = make(map[string]PluginInfo)

	// Register discovered plugins
	for _, plugin := range plugins {
		r.plugins[plugin.Name] = plugin
	}
}

// Count returns the number of registered plugins
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.plugins)
}
