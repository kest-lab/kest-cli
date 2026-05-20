package migrations

import (
	"time"

	"gorm.io/gorm"

	"github.com/kest-labs/kest/api/internal/infra/migration"
)

func init() {
	register("2026_05_20_000032_create_unified_runs_tables", &createUnifiedRunsTables{
		BaseMigration: migration.BaseMigration{
			UseTransaction: true,
		},
	})
}

type createUnifiedRunsTables struct {
	migration.BaseMigration
}

func (m *createUnifiedRunsTables) Up(db *gorm.DB) error {
	if err := db.AutoMigrate(&runMigrationModel{}, &runStepMigrationModel{}); err != nil {
		return err
	}

	return nil
}

func (m *createUnifiedRunsTables) Down(db *gorm.DB) error {
	if db.Migrator().HasTable("run_steps") {
		if err := db.Migrator().DropTable("run_steps"); err != nil {
			return err
		}
	}
	if db.Migrator().HasTable("runs") {
		if err := db.Migrator().DropTable("runs"); err != nil {
			return err
		}
	}

	return nil
}

type runMigrationModel struct {
	ID               string `gorm:"primaryKey"`
	WorkspaceID      string `gorm:"not null;index:idx_runs_workspace_created"`
	SourceType       string `gorm:"size:32;not null;index:idx_runs_source"`
	SourceID         string `gorm:"not null;index:idx_runs_source"`
	SourceName       string `gorm:"size:255"`
	Status           string `gorm:"size:20;not null;default:'pending';index:idx_runs_status"`
	TriggeredBy      string `gorm:"not null;index"`
	EnvironmentID    *string
	ExecutionMode    string `gorm:"size:20;not null;default:'local'"`
	TotalSteps       int    `gorm:"default:0"`
	PassedSteps      int    `gorm:"default:0"`
	FailedSteps      int    `gorm:"default:0"`
	DurationMs       int64  `gorm:"default:0"`
	RequestSnapshot  string `gorm:"type:text"`
	ResponseSnapshot string `gorm:"type:text"`
	ErrorMessage     string `gorm:"type:text"`
	Metadata         string `gorm:"type:text"`
	StartedAt        *time.Time
	FinishedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (runMigrationModel) TableName() string {
	return "runs"
}

type runStepMigrationModel struct {
	ID               string `gorm:"primaryKey"`
	RunID            string `gorm:"not null;index:idx_run_steps_run_sort"`
	StepIndex        int    `gorm:"not null;index:idx_run_steps_run_sort"`
	SourceType       string `gorm:"size:32;not null"`
	SourceID         string `gorm:"not null"`
	SourceName       string `gorm:"size:255"`
	Status           string `gorm:"size:20;not null;default:'pending';index:idx_run_steps_status"`
	DurationMs       int64  `gorm:"default:0"`
	RequestSnapshot  string `gorm:"type:text"`
	ResponseSnapshot string `gorm:"type:text"`
	ErrorMessage     string `gorm:"type:text"`
	Metadata         string `gorm:"type:text"`
	StartedAt        *time.Time
	FinishedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (runStepMigrationModel) TableName() string {
	return "run_steps"
}
