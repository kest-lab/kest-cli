package testing

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Fixture provides test setup and cleanup utilities
type Fixture struct {
	t        *testing.T
	cleanups []func()
	db       *gorm.DB
	router   *gin.Engine
}

// FixtureOption configures a fixture
type FixtureOption func(*Fixture)

// NewFixture creates a new test fixture
func NewFixture(t *testing.T, opts ...FixtureOption) *Fixture {
	f := &Fixture{
		t:        t,
		cleanups: make([]func(), 0),
	}

	for _, opt := range opts {
		opt(f)
	}

	// Register cleanup to run after test
	t.Cleanup(func() {
		f.cleanup()
	})

	return f
}

// WithDatabase configures an in-memory SQLite database
func WithDatabase(models ...interface{}) FixtureOption {
	return func(f *Fixture) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			f.t.Fatalf("Failed to open database: %v", err)
		}

		// Auto-migrate models
		if len(models) > 0 {
			if err := db.AutoMigrate(models...); err != nil {
				f.t.Fatalf("Failed to migrate: %v", err)
			}
		}

		f.db = db
		f.OnCleanup(func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		})
	}
}

// WithDatabaseConfig configures database with custom GORM config
func WithDatabaseConfig(config *gorm.Config, models ...interface{}) FixtureOption {
	return func(f *Fixture) {
		db, err := gorm.Open(sqlite.Open(":memory:"), config)
		if err != nil {
			f.t.Fatalf("Failed to open database: %v", err)
		}

		if len(models) > 0 {
			if err := db.AutoMigrate(models...); err != nil {
				f.t.Fatalf("Failed to migrate: %v", err)
			}
		}

		f.db = db
		f.OnCleanup(func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		})
	}
}

// WithRouter configures a Gin router
func WithRouter(setup func(*gin.Engine)) FixtureOption {
	return func(f *Fixture) {
		gin.SetMode(gin.TestMode)
		f.router = gin.New()
		if setup != nil {
			setup(f.router)
		}
	}
}

// WithMigrations runs migrations on the database
func WithMigrations(migrate func(*gorm.DB) error) FixtureOption {
	return func(f *Fixture) {
		if f.db == nil {
			f.t.Fatal("Database must be configured before migrations")
		}
		if err := migrate(f.db); err != nil {
			f.t.Fatalf("Failed to run migrations: %v", err)
		}
	}
}

// OnCleanup registers a cleanup function (LIFO order)
func (f *Fixture) OnCleanup(fn func()) {
	f.cleanups = append(f.cleanups, fn)
}

// cleanup runs all cleanup functions in LIFO order
func (f *Fixture) cleanup() {
	for i := len(f.cleanups) - 1; i >= 0; i-- {
		f.cleanups[i]()
	}
}

// DB returns the database connection
func (f *Fixture) DB() *gorm.DB {
	if f.db == nil {
		f.t.Fatal("Database not configured. Use WithDatabase option.")
	}
	return f.db
}

// Router returns the Gin router
func (f *Fixture) Router() *gin.Engine {
	if f.router == nil {
		f.t.Fatal("Router not configured. Use WithRouter option.")
	}
	return f.router
}

// Seed runs a seeder function
func (f *Fixture) Seed(seeder func(db *gorm.DB) error) *Fixture {
	if f.db == nil {
		f.t.Fatal("Database not configured")
	}
	if err := seeder(f.db); err != nil {
		f.t.Fatalf("Failed to seed: %v", err)
	}
	return f
}

// SeedData inserts data into the database
func (f *Fixture) SeedData(data ...interface{}) *Fixture {
	if f.db == nil {
		f.t.Fatal("Database not configured")
	}
	for _, d := range data {
		if err := f.db.Create(d).Error; err != nil {
			f.t.Fatalf("Failed to seed data: %v", err)
		}
	}
	return f
}

// Transaction runs a function within a transaction that's rolled back
func (f *Fixture) Transaction(fn func(tx *gorm.DB)) {
	if f.db == nil {
		f.t.Fatal("Database not configured")
	}
	tx := f.db.Begin()
	defer tx.Rollback()
	fn(tx)
}

// TransactionScope creates a scoped transaction for the test
func (f *Fixture) TransactionScope() *gorm.DB {
	if f.db == nil {
		f.t.Fatal("Database not configured")
	}
	tx := f.db.Begin()
	f.OnCleanup(func() {
		tx.Rollback()
	})
	return tx
}

// HTTP returns an HTTP test case builder
func (f *Fixture) HTTP() *TestCase {
	if f.router == nil {
		f.t.Fatal("Router not configured")
	}
	return NewTestCase(f.t, f.router)
}

// Database returns a database test case
func (f *Fixture) Database() *DatabaseTestCase {
	if f.db == nil {
		f.t.Fatal("Database not configured")
	}
	return NewDatabaseTestCase(f.t, f.db)
}

// T returns the testing.T instance
func (f *Fixture) T() *testing.T {
	return f.t
}

// Factory provides test data factory utilities
type Factory[T any] struct {
	t       *testing.T
	db      *gorm.DB
	builder func(int) T
}

// NewFactory creates a new factory for generating test data
func NewFactory[T any](t *testing.T, db *gorm.DB, builder func(index int) T) *Factory[T] {
	return &Factory[T]{
		t:       t,
		db:      db,
		builder: builder,
	}
}

// Create creates and persists a single instance
func (f *Factory[T]) Create() T {
	instance := f.builder(0)
	if err := f.db.Create(&instance).Error; err != nil {
		f.t.Fatalf("Failed to create: %v", err)
	}
	return instance
}

// CreateMany creates and persists multiple instances
func (f *Factory[T]) CreateMany(count int) []T {
	instances := make([]T, count)
	for i := 0; i < count; i++ {
		instances[i] = f.builder(i)
		if err := f.db.Create(&instances[i]).Error; err != nil {
			f.t.Fatalf("Failed to create: %v", err)
		}
	}
	return instances
}

// Make creates an instance without persisting
func (f *Factory[T]) Make() T {
	return f.builder(0)
}

// MakeMany creates multiple instances without persisting
func (f *Factory[T]) MakeMany(count int) []T {
	instances := make([]T, count)
	for i := 0; i < count; i++ {
		instances[i] = f.builder(i)
	}
	return instances
}
