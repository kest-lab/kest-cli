package run

import (
	"context"
	"errors"
	"time"

	"github.com/kest-labs/kest/api/internal/modules/request"
	"github.com/kest-labs/kest/api/internal/modules/runner"
	"github.com/kest-labs/kest/api/internal/modules/variable"
)

var (
	ErrRunNotFound     = errors.New("run not found")
	ErrInvalidRunInput = errors.New("invalid run input")
)

type Service interface {
	RunRequest(ctx context.Context, workspaceID string, reqModel *request.Request, triggeredBy string, req *RunRequest) (*runner.Response, error)
	RecordRun(ctx context.Context, workspaceID string, triggeredBy string, req *CreateRunRequest) (*Run, error)
	GetRunByID(ctx context.Context, id string) (*Run, error)
	ListRuns(ctx context.Context, workspaceID string, sourceType string, sourceID string, page, perPage int) ([]*Run, int64, error)
}

type service struct {
	repo   Repository
	runner runner.Runner
}

func NewService(repo Repository, runner runner.Runner) Service {
	return &service{
		repo:   repo,
		runner: runner,
	}
}

func (s *service) RunRequest(ctx context.Context, workspaceID string, reqModel *request.Request, triggeredBy string, req *RunRequest) (*runner.Response, error) {
	vars := variable.Variables{}
	if req != nil && req.Variables != nil {
		vars = req.Variables
	}

	resp, err := s.runner.Run(reqModel, vars)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) RecordRun(ctx context.Context, workspaceID string, triggeredBy string, req *CreateRunRequest) (*Run, error) {
	if req == nil || req.SourceType == "" || req.SourceID == "" {
		return nil, ErrInvalidRunInput
	}

	now := time.Now().UTC()
	run := &Run{
		WorkspaceID:      workspaceID,
		SourceType:       req.SourceType,
		SourceID:         req.SourceID,
		SourceName:       req.SourceName,
		Status:           normalizeRunStatus(req.Status),
		TriggeredBy:      triggeredBy,
		EnvironmentID:    req.EnvironmentID,
		ExecutionMode:    normalizeExecutionMode(req.ExecutionMode),
		DurationMs:       req.DurationMs,
		RequestSnapshot:  req.RequestSnapshot,
		ResponseSnapshot: req.ResponseSnapshot,
		ErrorMessage:     req.ErrorMessage,
		Metadata:         req.Metadata,
		StartedAt:        req.StartedAt,
		FinishedAt:       req.FinishedAt,
	}

	if run.StartedAt == nil {
		run.StartedAt = &now
	}
	if run.FinishedAt == nil && run.Status != RunStatusRunning && run.Status != RunStatusPending {
		run.FinishedAt = &now
	}

	steps := make([]RunStep, 0, len(req.Steps))
	passedSteps := 0
	failedSteps := 0
	for index, item := range req.Steps {
		status := normalizeRunStatus(item.Status)
		if status == RunStatusPassed {
			passedSteps++
		}
		if status == RunStatusFailed || status == RunStatusCanceled {
			failedSteps++
		}

		step := RunStep{
			StepIndex:        index,
			SourceType:       item.SourceType,
			SourceID:         item.SourceID,
			SourceName:       item.SourceName,
			Status:           status,
			DurationMs:       item.DurationMs,
			RequestSnapshot:  item.RequestSnapshot,
			ResponseSnapshot: item.ResponseSnapshot,
			ErrorMessage:     item.ErrorMessage,
			Metadata:         item.Metadata,
			StartedAt:        item.StartedAt,
			FinishedAt:       item.FinishedAt,
		}
		if item.StepIndex > 0 {
			step.StepIndex = item.StepIndex
		}
		steps = append(steps, step)
	}

	run.TotalSteps = len(steps)
	run.PassedSteps = passedSteps
	run.FailedSteps = failedSteps

	if err := s.repo.WithTransaction(ctx, func(txRepo Repository) error {
		if err := txRepo.CreateRun(ctx, run); err != nil {
			return err
		}

		for index := range steps {
			steps[index].RunID = run.ID
		}
		if err := txRepo.CreateRunSteps(ctx, steps); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	run.Steps = steps
	return run, nil
}

func (s *service) GetRunByID(ctx context.Context, id string) (*Run, error) {
	run, err := s.repo.GetRunByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if run == nil {
		return nil, ErrRunNotFound
	}

	steps, err := s.repo.ListRunSteps(ctx, id)
	if err != nil {
		return nil, err
	}

	run.Steps = make([]RunStep, 0, len(steps))
	for _, step := range steps {
		run.Steps = append(run.Steps, *step)
	}

	return run, nil
}

func (s *service) ListRuns(ctx context.Context, workspaceID string, sourceType string, sourceID string, page, perPage int) ([]*Run, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	return s.repo.ListRunsByWorkspace(ctx, workspaceID, sourceType, sourceID, page, perPage)
}

func normalizeRunStatus(status string) string {
	switch status {
	case RunStatusRunning, RunStatusPassed, RunStatusFailed, RunStatusCanceled:
		return status
	case "":
		return RunStatusPending
	default:
		return status
	}
}

func normalizeExecutionMode(mode string) string {
	if mode == "" {
		return "local"
	}
	return mode
}
