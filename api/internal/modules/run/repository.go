package run

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/kest-labs/kest/api/pkg/dbutil"
)

type Repository interface {
	CreateRun(ctx context.Context, run *Run) error
	UpdateRun(ctx context.Context, run *Run) error
	GetRunByID(ctx context.Context, id string) (*Run, error)
	ListRunsByWorkspace(ctx context.Context, workspaceID string, sourceType string, sourceID string, page, perPage int) ([]*Run, int64, error)
	CreateRunSteps(ctx context.Context, steps []RunStep) error
	ReplaceRunSteps(ctx context.Context, runID string, steps []RunStep) error
	ListRunSteps(ctx context.Context, runID string) ([]*RunStep, error)
	WithTransaction(ctx context.Context, fn func(Repository) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRun(ctx context.Context, run *Run) error {
	po := newRunPO(run)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}

	run.ID = po.ID
	run.CreatedAt = po.CreatedAt
	run.UpdatedAt = po.UpdatedAt
	return nil
}

func (r *repository) UpdateRun(ctx context.Context, run *Run) error {
	po := newRunPO(run)
	return r.db.WithContext(ctx).
		Model(&RunPO{}).
		Where("id = ?", run.ID).
		Updates(po).Error
}

func (r *repository) GetRunByID(ctx context.Context, id string) (*Run, error) {
	var po RunPO
	if err := dbutil.ByID(r.db.WithContext(ctx), id).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return po.toDomain(), nil
}

func (r *repository) ListRunsByWorkspace(ctx context.Context, workspaceID string, sourceType string, sourceID string, page, perPage int) ([]*Run, int64, error) {
	var poList []*RunPO
	var total int64

	query := r.db.WithContext(ctx).Model(&RunPO{}).Where("workspace_id = ?", workspaceID)
	if sourceType != "" {
		query = query.Where("source_type = ?", sourceType)
	}
	if sourceID != "" {
		query = query.Where("source_id = ?", sourceID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&poList).Error; err != nil {
		return nil, 0, err
	}

	runs := make([]*Run, 0, len(poList))
	for _, po := range poList {
		runs = append(runs, po.toDomain())
	}

	return runs, total, nil
}

func (r *repository) CreateRunSteps(ctx context.Context, steps []RunStep) error {
	if len(steps) == 0 {
		return nil
	}

	poList := make([]*RunStepPO, 0, len(steps))
	for i := range steps {
		poList = append(poList, newRunStepPO(&steps[i]))
	}

	return r.db.WithContext(ctx).Create(&poList).Error
}

func (r *repository) ReplaceRunSteps(ctx context.Context, runID string, steps []RunStep) error {
	if err := r.db.WithContext(ctx).Where("run_id = ?", runID).Delete(&RunStepPO{}).Error; err != nil {
		return err
	}

	return r.CreateRunSteps(ctx, steps)
}

func (r *repository) ListRunSteps(ctx context.Context, runID string) ([]*RunStep, error) {
	var poList []*RunStepPO
	if err := r.db.WithContext(ctx).
		Where("run_id = ?", runID).
		Order("step_index ASC, created_at ASC").
		Find(&poList).Error; err != nil {
		return nil, err
	}

	steps := make([]*RunStep, 0, len(poList))
	for _, po := range poList {
		steps = append(steps, po.toDomain())
	}

	return steps, nil
}

func (r *repository) WithTransaction(ctx context.Context, fn func(Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repository{db: tx})
	})
}
