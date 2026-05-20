package run

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestRepositoryCreateAndLoadRunWithSteps(t *testing.T) {
	db := newRunTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	startedAt := time.Now().UTC().Add(-2 * time.Second)
	finishedAt := time.Now().UTC()
	run := &Run{
		WorkspaceID:   "workspace-1",
		SourceType:    RunSourceRequest,
		SourceID:      "request-1",
		SourceName:    "GET /health",
		Status:        RunStatusPassed,
		TriggeredBy:   "user-1",
		ExecutionMode: "local",
		TotalSteps:    1,
		PassedSteps:   1,
		DurationMs:    120,
		RequestSnapshot: map[string]any{
			"method": "GET",
			"url":    "https://api.example.com/health",
		},
		ResponseSnapshot: map[string]any{
			"status": 200,
		},
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
	}

	if err := repo.CreateRun(ctx, run); err != nil {
		t.Fatalf("failed to create run: %v", err)
	}

	steps := []RunStep{
		{
			RunID:       run.ID,
			StepIndex:   0,
			SourceType:  RunSourceRequest,
			SourceID:    "request-1",
			SourceName:  "GET /health",
			Status:      RunStatusPassed,
			DurationMs:  120,
			StartedAt:   &startedAt,
			FinishedAt:  &finishedAt,
			RequestSnapshot: map[string]any{
				"url": "https://api.example.com/health",
			},
			ResponseSnapshot: map[string]any{
				"status": 200,
			},
		},
	}

	if err := repo.CreateRunSteps(ctx, steps); err != nil {
		t.Fatalf("failed to create run steps: %v", err)
	}

	loadedRun, err := repo.GetRunByID(ctx, run.ID)
	if err != nil {
		t.Fatalf("failed to load run: %v", err)
	}
	if loadedRun == nil {
		t.Fatal("expected run, got nil")
	}
	if loadedRun.WorkspaceID != "workspace-1" {
		t.Fatalf("expected workspace-1, got %q", loadedRun.WorkspaceID)
	}
	if loadedRun.SourceType != RunSourceRequest {
		t.Fatalf("expected source type %q, got %q", RunSourceRequest, loadedRun.SourceType)
	}

	loadedSteps, err := repo.ListRunSteps(ctx, run.ID)
	if err != nil {
		t.Fatalf("failed to list run steps: %v", err)
	}
	if len(loadedSteps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(loadedSteps))
	}
	if loadedSteps[0].StepIndex != 0 {
		t.Fatalf("expected step_index 0, got %d", loadedSteps[0].StepIndex)
	}
}

func TestRepositoryListRunsByWorkspaceScopesResults(t *testing.T) {
	db := newRunTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	first := &Run{
		WorkspaceID:   "workspace-1",
		SourceType:    RunSourceRequest,
		SourceID:      "request-1",
		Status:        RunStatusPassed,
		TriggeredBy:   "user-1",
		ExecutionMode: "local",
	}
	second := &Run{
		WorkspaceID:   "workspace-2",
		SourceType:    RunSourceCollection,
		SourceID:      "collection-1",
		Status:        RunStatusFailed,
		TriggeredBy:   "user-2",
		ExecutionMode: "local",
	}

	if err := repo.CreateRun(ctx, first); err != nil {
		t.Fatalf("failed to create first run: %v", err)
	}
	if err := repo.CreateRun(ctx, second); err != nil {
		t.Fatalf("failed to create second run: %v", err)
	}

	runs, total, err := repo.ListRunsByWorkspace(ctx, "workspace-1", "", "", 1, 20)
	if err != nil {
		t.Fatalf("failed to list runs: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	if runs[0].WorkspaceID != "workspace-1" {
		t.Fatalf("expected workspace-1 run, got %q", runs[0].WorkspaceID)
	}
}

func newRunTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.Callback().Create().Before("gorm:before_create").Register("test:assign_uuid_primary_key", func(tx *gorm.DB) {
		if tx == nil || tx.Statement == nil || tx.Statement.Schema == nil {
			return
		}

		idField := tx.Statement.Schema.LookUpField("ID")
		if idField == nil || idField.FieldType.Kind() != reflect.String {
			return
		}

		value := tx.Statement.ReflectValue
		for value.Kind() == reflect.Ptr {
			if value.IsNil() {
				return
			}
			value = value.Elem()
		}
		if value.Kind() != reflect.Struct {
			return
		}

		_, isZero := idField.ValueOf(tx.Statement.Context, value)
		if isZero {
			_ = idField.Set(tx.Statement.Context, value, uuid.NewString())
		}
	}); err != nil {
		t.Fatalf("failed to register id callback: %v", err)
	}

	if err := db.AutoMigrate(&RunPO{}, &RunStepPO{}); err != nil {
		t.Fatalf("failed to migrate run schema: %v", err)
	}

	return db
}
