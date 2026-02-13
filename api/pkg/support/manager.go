package support

import (
	"fmt"
	"sync"
)

// Manager provides a driver-based management pattern.
// It allows managing multiple driver implementations with lazy loading and caching.
type Manager[T any] struct {
	mu             sync.RWMutex
	drivers        map[string]T
	customCreators map[string]func() T
	defaultDriver  string
	createDriver   func(name string) (T, error)
}

// ManagerConfig holds configuration for creating a Manager.
type ManagerConfig[T any] struct {
	DefaultDriver string
	CreateDriver  func(name string) (T, error)
}

// NewManager creates a new Manager instance.
//
// Example:
//
//	cacheManager := NewManager(ManagerConfig[Cache]{
//	    DefaultDriver: "redis",
//	    CreateDriver: func(name string) (Cache, error) {
//	        switch name {
//	        case "redis":
//	            return NewRedisCache(), nil
//	        case "memory":
//	            return NewMemoryCache(), nil
//	        default:
//	            return nil, fmt.Errorf("driver [%s] not supported", name)
//	        }
//	    },
//	})
func NewManager[T any](config ManagerConfig[T]) *Manager[T] {
	return &Manager[T]{
		drivers:        make(map[string]T),
		customCreators: make(map[string]func() T),
		defaultDriver:  config.DefaultDriver,
		createDriver:   config.CreateDriver,
	}
}

// Driver returns a driver instance by name.
// If no name is provided, returns the default driver.
func (m *Manager[T]) Driver(name ...string) (T, error) {
	driverName := m.defaultDriver
	if len(name) > 0 && name[0] != "" {
		driverName = name[0]
	}

	if driverName == "" {
		var zero T
		return zero, fmt.Errorf("no driver specified and no default driver configured")
	}

	// Check if already created
	m.mu.RLock()
	if driver, ok := m.drivers[driverName]; ok {
		m.mu.RUnlock()
		return driver, nil
	}
	m.mu.RUnlock()

	// Create new driver
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if driver, ok := m.drivers[driverName]; ok {
		return driver, nil
	}

	driver, err := m.resolveDriver(driverName)
	if err != nil {
		var zero T
		return zero, err
	}

	m.drivers[driverName] = driver
	return driver, nil
}

// MustDriver returns a driver instance or panics on error.
func (m *Manager[T]) MustDriver(name ...string) T {
	driver, err := m.Driver(name...)
	if err != nil {
		panic(err)
	}
	return driver
}

// resolveDriver creates a driver instance.
func (m *Manager[T]) resolveDriver(name string) (T, error) {
	// Check custom creators first
	if creator, ok := m.customCreators[name]; ok {
		return creator(), nil
	}

	// Use the default create function
	if m.createDriver != nil {
		return m.createDriver(name)
	}

	var zero T
	return zero, fmt.Errorf("driver [%s] not supported", name)
}

// Extend registers a custom driver creator.
//
// Example:
//
//	cacheManager.Extend("custom", func() Cache {
//	    return NewCustomCache()
//	})
func (m *Manager[T]) Extend(name string, creator func() T) *Manager[T] {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.customCreators[name] = creator
	return m
}

// SetDefaultDriver sets the default driver name.
func (m *Manager[T]) SetDefaultDriver(name string) *Manager[T] {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultDriver = name
	return m
}

// GetDefaultDriver returns the default driver name.
func (m *Manager[T]) GetDefaultDriver() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.defaultDriver
}

// GetDrivers returns all created driver instances.
func (m *Manager[T]) GetDrivers() map[string]T {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]T, len(m.drivers))
	for k, v := range m.drivers {
		result[k] = v
	}
	return result
}

// ForgetDriver removes a driver from the cache.
func (m *Manager[T]) ForgetDriver(name string) *Manager[T] {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.drivers, name)
	return m
}

// ForgetDrivers removes all drivers from the cache.
func (m *Manager[T]) ForgetDrivers() *Manager[T] {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers = make(map[string]T)
	return m
}

// HasDriver checks if a driver has been created.
func (m *Manager[T]) HasDriver(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.drivers[name]
	return ok
}

// HasCustomCreator checks if a custom creator exists for the driver.
func (m *Manager[T]) HasCustomCreator(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.customCreators[name]
	return ok
}
