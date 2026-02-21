package migrations

import (
	"github.com/kest-labs/kest/api/internal/infra/migration"
	"github.com/kest-labs/kest/api/internal/modules/apispec"
	"gorm.io/gorm"
)

func init() {
	register("2026_02_21_000000_add_fields_to_api_examples", &addFieldsToAPIExamples{})
}

type addFieldsToAPIExamples struct {
	migration.BaseMigration
}

func (m *addFieldsToAPIExamples) Up(db *gorm.DB) error {
	migrator := db.Migrator()
	if !migrator.HasColumn("api_examples", "path") {
		if err := migrator.AddColumn(&apispec.APIExamplePO{}, "Path"); err != nil {
			return err
		}
	}
	if !migrator.HasColumn("api_examples", "method") {
		if err := migrator.AddColumn(&apispec.APIExamplePO{}, "Method"); err != nil {
			return err
		}
	}
	if !migrator.HasColumn("api_examples", "description") {
		if err := migrator.AddColumn(&apispec.APIExamplePO{}, "Description"); err != nil {
			return err
		}
	}
	if !migrator.HasColumn("api_examples", "response_headers") {
		if err := migrator.AddColumn(&apispec.APIExamplePO{}, "ResponseHeaders"); err != nil {
			return err
		}
	}
	return nil
}

func (m *addFieldsToAPIExamples) Down(db *gorm.DB) error {
	migrator := db.Migrator()
	if migrator.HasColumn("api_examples", "response_headers") {
		if err := migrator.DropColumn("api_examples", "response_headers"); err != nil {
			return err
		}
	}
	if migrator.HasColumn("api_examples", "description") {
		if err := migrator.DropColumn("api_examples", "description"); err != nil {
			return err
		}
	}
	if migrator.HasColumn("api_examples", "method") {
		if err := migrator.DropColumn("api_examples", "method"); err != nil {
			return err
		}
	}
	if migrator.HasColumn("api_examples", "path") {
		if err := migrator.DropColumn("api_examples", "path"); err != nil {
			return err
		}
	}
	return nil
}
