package run

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const (
	RunSourceRequest    = "request"
	RunSourceCollection = "collection"
	RunSourceTestCase   = "test_case"
	RunSourceFlow       = "flow"
)

const (
	RunStatusPending  = "pending"
	RunStatusRunning  = "running"
	RunStatusPassed   = "passed"
	RunStatusFailed   = "failed"
	RunStatusCanceled = "canceled"
)

type RunPO struct {
	ID               string         `gorm:"primaryKey" json:"id"`
	WorkspaceID      string         `gorm:"not null;index:idx_runs_workspace_created" json:"workspace_id"`
	SourceType       string         `gorm:"size:32;not null;index:idx_runs_source" json:"source_type"`
	SourceID         string         `gorm:"not null;index:idx_runs_source" json:"source_id"`
	SourceName       string         `gorm:"size:255" json:"source_name"`
	Status           string         `gorm:"size:20;not null;default:'pending';index:idx_runs_status" json:"status"`
	TriggeredBy      string         `gorm:"not null;index" json:"triggered_by"`
	EnvironmentID    *string        `gorm:"index" json:"environment_id,omitempty"`
	ExecutionMode    string         `gorm:"size:20;not null;default:'local'" json:"execution_mode"`
	TotalSteps       int            `gorm:"default:0" json:"total_steps"`
	PassedSteps      int            `gorm:"default:0" json:"passed_steps"`
	FailedSteps      int            `gorm:"default:0" json:"failed_steps"`
	DurationMs       int64          `gorm:"default:0" json:"duration_ms"`
	RequestSnapshot  string         `gorm:"type:text" json:"request_snapshot"`
	ResponseSnapshot string         `gorm:"type:text" json:"response_snapshot"`
	ErrorMessage     string         `gorm:"type:text" json:"error_message"`
	Metadata         string         `gorm:"type:text" json:"metadata"`
	StartedAt        *time.Time     `json:"started_at"`
	FinishedAt       *time.Time     `json:"finished_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RunPO) TableName() string {
	return "runs"
}

type RunStepPO struct {
	ID               string         `gorm:"primaryKey" json:"id"`
	RunID            string         `gorm:"not null;index:idx_run_steps_run_sort" json:"run_id"`
	StepIndex        int            `gorm:"not null;index:idx_run_steps_run_sort" json:"step_index"`
	SourceType       string         `gorm:"size:32;not null" json:"source_type"`
	SourceID         string         `gorm:"not null" json:"source_id"`
	SourceName       string         `gorm:"size:255" json:"source_name"`
	Status           string         `gorm:"size:20;not null;default:'pending';index:idx_run_steps_status" json:"status"`
	DurationMs       int64          `gorm:"default:0" json:"duration_ms"`
	RequestSnapshot  string         `gorm:"type:text" json:"request_snapshot"`
	ResponseSnapshot string         `gorm:"type:text" json:"response_snapshot"`
	ErrorMessage     string         `gorm:"type:text" json:"error_message"`
	Metadata         string         `gorm:"type:text" json:"metadata"`
	StartedAt        *time.Time     `json:"started_at"`
	FinishedAt       *time.Time     `json:"finished_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RunStepPO) TableName() string {
	return "run_steps"
}

type Run struct {
	ID               string                 `json:"id"`
	WorkspaceID      string                 `json:"workspace_id"`
	SourceType       string                 `json:"source_type"`
	SourceID         string                 `json:"source_id"`
	SourceName       string                 `json:"source_name"`
	Status           string                 `json:"status"`
	TriggeredBy      string                 `json:"triggered_by"`
	EnvironmentID    *string                `json:"environment_id,omitempty"`
	ExecutionMode    string                 `json:"execution_mode"`
	TotalSteps       int                    `json:"total_steps"`
	PassedSteps      int                    `json:"passed_steps"`
	FailedSteps      int                    `json:"failed_steps"`
	DurationMs       int64                  `json:"duration_ms"`
	RequestSnapshot  map[string]any         `json:"request_snapshot,omitempty"`
	ResponseSnapshot map[string]any         `json:"response_snapshot,omitempty"`
	ErrorMessage     string                 `json:"error_message,omitempty"`
	Metadata         map[string]any         `json:"metadata,omitempty"`
	StartedAt        *time.Time             `json:"started_at,omitempty"`
	FinishedAt       *time.Time             `json:"finished_at,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Steps            []RunStep              `json:"steps,omitempty"`
}

type RunStep struct {
	ID               string                 `json:"id"`
	RunID            string                 `json:"run_id"`
	StepIndex        int                    `json:"step_index"`
	SourceType       string                 `json:"source_type"`
	SourceID         string                 `json:"source_id"`
	SourceName       string                 `json:"source_name"`
	Status           string                 `json:"status"`
	DurationMs       int64                  `json:"duration_ms"`
	RequestSnapshot  map[string]any         `json:"request_snapshot,omitempty"`
	ResponseSnapshot map[string]any         `json:"response_snapshot,omitempty"`
	ErrorMessage     string                 `json:"error_message,omitempty"`
	Metadata         map[string]any         `json:"metadata,omitempty"`
	StartedAt        *time.Time             `json:"started_at,omitempty"`
	FinishedAt       *time.Time             `json:"finished_at,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

func (po *RunPO) toDomain() *Run {
	if po == nil {
		return nil
	}

	return &Run{
		ID:               po.ID,
		WorkspaceID:      po.WorkspaceID,
		SourceType:       po.SourceType,
		SourceID:         po.SourceID,
		SourceName:       po.SourceName,
		Status:           po.Status,
		TriggeredBy:      po.TriggeredBy,
		EnvironmentID:    po.EnvironmentID,
		ExecutionMode:    po.ExecutionMode,
		TotalSteps:       po.TotalSteps,
		PassedSteps:      po.PassedSteps,
		FailedSteps:      po.FailedSteps,
		DurationMs:       po.DurationMs,
		RequestSnapshot:  decodeRunJSON(po.RequestSnapshot),
		ResponseSnapshot: decodeRunJSON(po.ResponseSnapshot),
		ErrorMessage:     po.ErrorMessage,
		Metadata:         decodeRunJSON(po.Metadata),
		StartedAt:        po.StartedAt,
		FinishedAt:       po.FinishedAt,
		CreatedAt:        po.CreatedAt,
		UpdatedAt:        po.UpdatedAt,
	}
}

func (po *RunStepPO) toDomain() *RunStep {
	if po == nil {
		return nil
	}

	return &RunStep{
		ID:               po.ID,
		RunID:            po.RunID,
		StepIndex:        po.StepIndex,
		SourceType:       po.SourceType,
		SourceID:         po.SourceID,
		SourceName:       po.SourceName,
		Status:           po.Status,
		DurationMs:       po.DurationMs,
		RequestSnapshot:  decodeRunJSON(po.RequestSnapshot),
		ResponseSnapshot: decodeRunJSON(po.ResponseSnapshot),
		ErrorMessage:     po.ErrorMessage,
		Metadata:         decodeRunJSON(po.Metadata),
		StartedAt:        po.StartedAt,
		FinishedAt:       po.FinishedAt,
		CreatedAt:        po.CreatedAt,
		UpdatedAt:        po.UpdatedAt,
	}
}

func newRunPO(run *Run) *RunPO {
	if run == nil {
		return nil
	}

	return &RunPO{
		ID:               run.ID,
		WorkspaceID:      run.WorkspaceID,
		SourceType:       run.SourceType,
		SourceID:         run.SourceID,
		SourceName:       run.SourceName,
		Status:           run.Status,
		TriggeredBy:      run.TriggeredBy,
		EnvironmentID:    run.EnvironmentID,
		ExecutionMode:    run.ExecutionMode,
		TotalSteps:       run.TotalSteps,
		PassedSteps:      run.PassedSteps,
		FailedSteps:      run.FailedSteps,
		DurationMs:       run.DurationMs,
		RequestSnapshot:  encodeRunJSON(run.RequestSnapshot),
		ResponseSnapshot: encodeRunJSON(run.ResponseSnapshot),
		ErrorMessage:     run.ErrorMessage,
		Metadata:         encodeRunJSON(run.Metadata),
		StartedAt:        run.StartedAt,
		FinishedAt:       run.FinishedAt,
	}
}

func newRunStepPO(step *RunStep) *RunStepPO {
	if step == nil {
		return nil
	}

	return &RunStepPO{
		ID:               step.ID,
		RunID:            step.RunID,
		StepIndex:        step.StepIndex,
		SourceType:       step.SourceType,
		SourceID:         step.SourceID,
		SourceName:       step.SourceName,
		Status:           step.Status,
		DurationMs:       step.DurationMs,
		RequestSnapshot:  encodeRunJSON(step.RequestSnapshot),
		ResponseSnapshot: encodeRunJSON(step.ResponseSnapshot),
		ErrorMessage:     step.ErrorMessage,
		Metadata:         encodeRunJSON(step.Metadata),
		StartedAt:        step.StartedAt,
		FinishedAt:       step.FinishedAt,
	}
}

func decodeRunJSON(value string) map[string]any {
	if value == "" {
		return nil
	}

	var out map[string]any
	_ = json.Unmarshal([]byte(value), &out)
	return out
}

func encodeRunJSON(value map[string]any) string {
	if value == nil {
		return ""
	}

	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}
