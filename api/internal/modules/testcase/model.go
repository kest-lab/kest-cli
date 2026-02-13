package testcase

import (
	"time"

	"gorm.io/gorm"
)

// TestCasePO represents a test case in the database
type TestCasePO struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	APISpecID   uint   `gorm:"not null;index:idx_testcase_api_spec" json:"api_spec_id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Env         string `gorm:"size:50" json:"env"` // dev, test, staging, prod

	// Request configuration (JSON encoded)
	Headers     string `gorm:"type:text" json:"headers"`      // map[string]string
	QueryParams string `gorm:"type:text" json:"query_params"` // map[string]string
	PathParams  string `gorm:"type:text" json:"path_params"`  // map[string]string
	RequestBody string `gorm:"type:text" json:"request_body"` // any

	// Scripts
	PreScript  string `gorm:"type:text" json:"pre_script"`  // JavaScript
	PostScript string `gorm:"type:text" json:"post_script"` // JavaScript

	// Assertions (JSON array)
	Assertions string `gorm:"type:text" json:"assertions"` // []Assertion

	// Variable extraction (JSON array)
	ExtractVars string `gorm:"type:text" json:"extract_vars"` // []ExtractVar

	CreatedBy uint           `gorm:"index" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for TestCasePO
func (TestCasePO) TableName() string {
	return "test_cases"
}

// Assertion represents a test assertion
type Assertion struct {
	Type     string `json:"type"`     // status, json_path, response_time, header, body_contains
	Path     string `json:"path"`     // for json_path type
	Operator string `json:"operator"` // equals, not_equals, exists, not_exists, contains, below, above
	Expect   any    `json:"expect"`   // expected value
	Message  string `json:"message"`  // custom assertion message
}

// ExtractVar represents a variable to extract from response
type ExtractVar struct {
	Name   string `json:"name"`   // variable name
	Source string `json:"source"` // body, header, cookie
	Path   string `json:"path"`   // JSON path or header name
}
