package migration

// MigratorOptions configures the migrator behavior for run operations.
type MigratorOptions struct {
	// Pretend shows SQL statements without executing them
	Pretend bool

	// Step increments batch number for each migration instead of using
	// a single batch number for all migrations in the run
	Step bool

	// Force allows running migrations in production environment
	Force bool
}

// RollbackOptions configures the migrator behavior for rollback operations.
type RollbackOptions struct {
	// Steps specifies the number of migrations to rollback.
	// If > 0, rolls back exactly N migrations regardless of batch.
	Steps int

	// Batch specifies a specific batch number to rollback.
	// If > 0, rolls back all migrations in that batch.
	Batch int

	// Pretend shows SQL statements without executing them
	Pretend bool
}

// NewMigratorOptions creates default migrator options.
func NewMigratorOptions() MigratorOptions {
	return MigratorOptions{
		Pretend: false,
		Step:    false,
		Force:   false,
	}
}

// NewRollbackOptions creates default rollback options.
func NewRollbackOptions() RollbackOptions {
	return RollbackOptions{
		Steps:   0,
		Batch:   0,
		Pretend: false,
	}
}

// WithPretend returns options with pretend mode enabled.
func (o MigratorOptions) WithPretend() MigratorOptions {
	o.Pretend = true
	return o
}

// WithStep returns options with step mode enabled.
func (o MigratorOptions) WithStep() MigratorOptions {
	o.Step = true
	return o
}

// WithForce returns options with force mode enabled.
func (o MigratorOptions) WithForce() MigratorOptions {
	o.Force = true
	return o
}

// WithSteps returns rollback options with the specified number of steps.
func (o RollbackOptions) WithSteps(steps int) RollbackOptions {
	o.Steps = steps
	return o
}

// WithBatch returns rollback options with the specified batch number.
func (o RollbackOptions) WithBatch(batch int) RollbackOptions {
	o.Batch = batch
	return o
}

// WithPretend returns rollback options with pretend mode enabled.
func (o RollbackOptions) WithPretend() RollbackOptions {
	o.Pretend = true
	return o
}
