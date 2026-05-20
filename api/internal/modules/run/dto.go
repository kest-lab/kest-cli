package run

import "time"

type RunRequest struct {
	EnvironmentID *string           `json:"environment_id"`
	Variables     map[string]string `json:"variables"`
}

type CreateRunRequest struct {
	SourceType       string                 `json:"source_type" binding:"required"`
	SourceID         string                 `json:"source_id" binding:"required"`
	SourceName       string                 `json:"source_name"`
	Status           string                 `json:"status"`
	EnvironmentID    *string                `json:"environment_id,omitempty"`
	ExecutionMode    string                 `json:"execution_mode"`
	DurationMs       int64                  `json:"duration_ms"`
	RequestSnapshot  map[string]any         `json:"request_snapshot"`
	ResponseSnapshot map[string]any         `json:"response_snapshot"`
	ErrorMessage     string                 `json:"error_message"`
	Metadata         map[string]any         `json:"metadata"`
	StartedAt        *time.Time             `json:"started_at,omitempty"`
	FinishedAt       *time.Time             `json:"finished_at,omitempty"`
	Steps            []CreateRunStepRequest `json:"steps"`
}

type CreateRunStepRequest struct {
	StepIndex        int             `json:"step_index"`
	SourceType       string          `json:"source_type" binding:"required"`
	SourceID         string          `json:"source_id" binding:"required"`
	SourceName       string          `json:"source_name"`
	Status           string          `json:"status"`
	DurationMs       int64           `json:"duration_ms"`
	RequestSnapshot  map[string]any  `json:"request_snapshot"`
	ResponseSnapshot map[string]any  `json:"response_snapshot"`
	ErrorMessage     string          `json:"error_message"`
	Metadata         map[string]any  `json:"metadata"`
	StartedAt        *time.Time      `json:"started_at,omitempty"`
	FinishedAt       *time.Time      `json:"finished_at,omitempty"`
}

type RunResponse struct {
	ID               string            `json:"id"`
	WorkspaceID      string            `json:"workspace_id"`
	SourceType       string            `json:"source_type"`
	SourceID         string            `json:"source_id"`
	SourceName       string            `json:"source_name"`
	Status           string            `json:"status"`
	TriggeredBy      string            `json:"triggered_by"`
	EnvironmentID    *string           `json:"environment_id,omitempty"`
	ExecutionMode    string            `json:"execution_mode"`
	TotalSteps       int               `json:"total_steps"`
	PassedSteps      int               `json:"passed_steps"`
	FailedSteps      int               `json:"failed_steps"`
	DurationMs       int64             `json:"duration_ms"`
	RequestSnapshot  map[string]any    `json:"request_snapshot,omitempty"`
	ResponseSnapshot map[string]any    `json:"response_snapshot,omitempty"`
	ErrorMessage     string            `json:"error_message,omitempty"`
	Metadata         map[string]any    `json:"metadata,omitempty"`
	StartedAt        *time.Time        `json:"started_at,omitempty"`
	FinishedAt       *time.Time        `json:"finished_at,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Steps            []RunStepResponse `json:"steps,omitempty"`
}

type RunStepResponse struct {
	ID               string         `json:"id"`
	RunID            string         `json:"run_id"`
	StepIndex        int            `json:"step_index"`
	SourceType       string         `json:"source_type"`
	SourceID         string         `json:"source_id"`
	SourceName       string         `json:"source_name"`
	Status           string         `json:"status"`
	DurationMs       int64          `json:"duration_ms"`
	RequestSnapshot  map[string]any `json:"request_snapshot,omitempty"`
	ResponseSnapshot map[string]any `json:"response_snapshot,omitempty"`
	ErrorMessage     string         `json:"error_message,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	StartedAt        *time.Time     `json:"started_at,omitempty"`
	FinishedAt       *time.Time     `json:"finished_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func toRunResponse(run *Run) *RunResponse {
	if run == nil {
		return nil
	}

	resp := &RunResponse{
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
		RequestSnapshot:  run.RequestSnapshot,
		ResponseSnapshot: run.ResponseSnapshot,
		ErrorMessage:     run.ErrorMessage,
		Metadata:         run.Metadata,
		StartedAt:        run.StartedAt,
		FinishedAt:       run.FinishedAt,
		CreatedAt:        run.CreatedAt,
		UpdatedAt:        run.UpdatedAt,
	}

	if len(run.Steps) > 0 {
		resp.Steps = make([]RunStepResponse, 0, len(run.Steps))
		for _, step := range run.Steps {
			resp.Steps = append(resp.Steps, RunStepResponse{
				ID:               step.ID,
				RunID:            step.RunID,
				StepIndex:        step.StepIndex,
				SourceType:       step.SourceType,
				SourceID:         step.SourceID,
				SourceName:       step.SourceName,
				Status:           step.Status,
				DurationMs:       step.DurationMs,
				RequestSnapshot:  step.RequestSnapshot,
				ResponseSnapshot: step.ResponseSnapshot,
				ErrorMessage:     step.ErrorMessage,
				Metadata:         step.Metadata,
				StartedAt:        step.StartedAt,
				FinishedAt:       step.FinishedAt,
				CreatedAt:        step.CreatedAt,
				UpdatedAt:        step.UpdatedAt,
			})
		}
	}

	return resp
}

func toRunResponseList(runs []*Run) []*RunResponse {
	result := make([]*RunResponse, 0, len(runs))
	for _, item := range runs {
		result = append(result, toRunResponse(item))
	}
	return result
}
