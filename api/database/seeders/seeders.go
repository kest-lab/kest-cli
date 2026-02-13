package seeders

import "gorm.io/gorm"

// Seeder interface defines the contract for database seeders
type Seeder interface {
	Run(db *gorm.DB) error
}

var registry []Seeder

// register adds a seeder to the registry
func register(s Seeder) {
	registry = append(registry, s)
}

// All returns all registered seeders
func All() []Seeder {
	return registry
}

// RunAll executes all registered seeders with the given database connection
func RunAll(db *gorm.DB) error {
	for _, seeder := range registry {
		if err := seeder.Run(db); err != nil {
			return err
		}
	}
	return nil
}
