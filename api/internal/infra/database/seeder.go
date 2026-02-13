package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Seeder is the interface for database seeders
type Seeder interface {
	Run(ctx context.Context, db *gorm.DB) error
	Name() string
}

// SeederFunc is a function adapter for Seeder
type SeederFunc struct {
	name string
	fn   func(ctx context.Context, db *gorm.DB) error
}

func (s SeederFunc) Run(ctx context.Context, db *gorm.DB) error {
	return s.fn(ctx, db)
}

func (s SeederFunc) Name() string {
	return s.name
}

// NewSeeder creates a new SeederFunc
func NewSeeder(name string, fn func(ctx context.Context, db *gorm.DB) error) Seeder {
	return SeederFunc{name: name, fn: fn}
}

// SeederManager manages and runs seeders
type SeederManager struct {
	db      *gorm.DB
	seeders []Seeder
	called  map[string]bool
	output  SeederOutput
}

// SeederOutput defines the output interface for seeder progress
type SeederOutput interface {
	Info(format string, args ...interface{})
	Success(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// DefaultSeederOutput provides default console output
type DefaultSeederOutput struct{}

func (o DefaultSeederOutput) Info(format string, args ...interface{}) {
	fmt.Printf("  ℹ "+format+"\n", args...)
}

func (o DefaultSeederOutput) Success(format string, args ...interface{}) {
	fmt.Printf("  ✓ "+format+"\n", args...)
}

func (o DefaultSeederOutput) Error(format string, args ...interface{}) {
	fmt.Printf("  ✗ "+format+"\n", args...)
}

// NewSeederManager creates a new SeederManager
func NewSeederManager(db *gorm.DB) *SeederManager {
	return &SeederManager{
		db:      db,
		seeders: make([]Seeder, 0),
		called:  make(map[string]bool),
		output:  DefaultSeederOutput{},
	}
}

// SetOutput sets the output handler
func (m *SeederManager) SetOutput(output SeederOutput) *SeederManager {
	m.output = output
	return m
}

// Register adds a seeder to the manager
func (m *SeederManager) Register(seeder Seeder) *SeederManager {
	m.seeders = append(m.seeders, seeder)
	return m
}

// RegisterFunc adds a function seeder to the manager
func (m *SeederManager) RegisterFunc(name string, fn func(ctx context.Context, db *gorm.DB) error) *SeederManager {
	return m.Register(NewSeeder(name, fn))
}

// Call runs a specific seeder by name
func (m *SeederManager) Call(ctx context.Context, name string) error {
	for _, seeder := range m.seeders {
		if seeder.Name() == name {
			return m.runSeeder(ctx, seeder, false)
		}
	}
	return fmt.Errorf("seeder '%s' not found", name)
}

// CallOnce runs a seeder only if it hasn't been called before
func (m *SeederManager) CallOnce(ctx context.Context, name string) error {
	if m.called[name] {
		return nil
	}
	return m.Call(ctx, name)
}

// Run executes all registered seeders
func (m *SeederManager) Run(ctx context.Context) error {
	for _, seeder := range m.seeders {
		if err := m.runSeeder(ctx, seeder, false); err != nil {
			return err
		}
	}
	return nil
}

// RunFresh truncates tables and runs all seeders
func (m *SeederManager) RunFresh(ctx context.Context, tables ...string) error {
	// Truncate specified tables
	for _, table := range tables {
		if err := m.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			m.output.Error("Failed to truncate table %s: %v", table, err)
			return err
		}
		m.output.Info("Truncated table: %s", table)
	}

	// Reset called status
	m.called = make(map[string]bool)

	return m.Run(ctx)
}

func (m *SeederManager) runSeeder(ctx context.Context, seeder Seeder, silent bool) error {
	name := seeder.Name()

	if !silent {
		m.output.Info("Running seeder: %s", name)
	}

	startTime := time.Now()

	if err := seeder.Run(ctx, m.db); err != nil {
		m.output.Error("Seeder %s failed: %v", name, err)
		return err
	}

	elapsed := time.Since(startTime)
	m.called[name] = true

	if !silent {
		m.output.Success("Seeder %s completed in %v", name, elapsed.Round(time.Millisecond))
	}

	return nil
}

// List returns all registered seeder names
func (m *SeederManager) List() []string {
	names := make([]string, len(m.seeders))
	for i, s := range m.seeders {
		names[i] = s.Name()
	}
	return names
}

// --- Factory helpers for generating test data ---

// Factory helps create test data
type Factory[T any] struct {
	db       *gorm.DB
	builder  func(index int) T
	count    int
	modifier func(*T)
}

// NewFactory creates a new factory for a model
func NewFactory[T any](db *gorm.DB, builder func(index int) T) *Factory[T] {
	return &Factory[T]{
		db:      db,
		builder: builder,
		count:   1,
	}
}

// Count sets the number of records to create
func (f *Factory[T]) Count(n int) *Factory[T] {
	f.count = n
	return f
}

// State applies a modifier to each created record
func (f *Factory[T]) State(modifier func(*T)) *Factory[T] {
	f.modifier = modifier
	return f
}

// Create inserts records into the database
func (f *Factory[T]) Create() ([]T, error) {
	records := make([]T, f.count)
	for i := 0; i < f.count; i++ {
		record := f.builder(i)
		if f.modifier != nil {
			f.modifier(&record)
		}
		if err := f.db.Create(&record).Error; err != nil {
			return nil, err
		}
		records[i] = record
	}
	return records, nil
}

// Make creates records without inserting them
func (f *Factory[T]) Make() []T {
	records := make([]T, f.count)
	for i := 0; i < f.count; i++ {
		record := f.builder(i)
		if f.modifier != nil {
			f.modifier(&record)
		}
		records[i] = record
	}
	return records
}

// CreateOne inserts a single record
func (f *Factory[T]) CreateOne() (T, error) {
	records, err := f.Count(1).Create()
	if err != nil {
		var zero T
		return zero, err
	}
	return records[0], nil
}

// MakeOne creates a single record without inserting
func (f *Factory[T]) MakeOne() T {
	return f.Count(1).Make()[0]
}
