package environment

import (
	"context"
	"encoding/json"
	"fmt"
)

// Service defines the interface for environment business logic
type Service interface {
	CreateEnvironment(ctx context.Context, req *CreateEnvironmentRequest) (*EnvironmentResponse, error)
	GetEnvironment(ctx context.Context, id uint) (*EnvironmentResponse, error)
	ListEnvironments(ctx context.Context, projectID uint) ([]*EnvironmentResponse, error)
	UpdateEnvironment(ctx context.Context, id uint, req *UpdateEnvironmentRequest) (*EnvironmentResponse, error)
	DeleteEnvironment(ctx context.Context, id uint) error
	DuplicateEnvironment(ctx context.Context, id uint, req *DuplicateEnvironmentRequest) (*EnvironmentResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new environment service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateEnvironment creates a new environment
func (s *service) CreateEnvironment(ctx context.Context, req *CreateEnvironmentRequest) (*EnvironmentResponse, error) {
	// Check if environment with same name already exists
	existing, err := s.repo.GetByProjectAndName(ctx, req.ProjectID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing environment: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("environment with name '%s' already exists in this project", req.Name)
	}

	// Convert request to PO
	env, err := req.ToEnvironmentPO()
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	// Create environment
	if err := s.repo.Create(ctx, env); err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	// Convert to response
	resp, err := env.ToResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to response: %w", err)
	}

	return resp, nil
}

// GetEnvironment gets an environment by ID
func (s *service) GetEnvironment(ctx context.Context, id uint) (*EnvironmentResponse, error) {
	env, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}
	if env == nil {
		return nil, fmt.Errorf("environment not found")
	}

	resp, err := env.ToResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to response: %w", err)
	}

	return resp, nil
}

// ListEnvironments lists all environments for a project
func (s *service) ListEnvironments(ctx context.Context, projectID uint) ([]*EnvironmentResponse, error) {
	envs, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	responses, err := ToResponseList(envs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to responses: %w", err)
	}

	return responses, nil
}

// UpdateEnvironment updates an environment
func (s *service) UpdateEnvironment(ctx context.Context, id uint, req *UpdateEnvironmentRequest) (*EnvironmentResponse, error) {
	// Get existing environment
	env, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}
	if env == nil {
		return nil, fmt.Errorf("environment not found")
	}

	// Update fields
	if req.Name != nil {
		// Check if new name conflicts
		if *req.Name != env.Name {
			existing, err := s.repo.GetByProjectAndName(ctx, env.ProjectID, *req.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to check existing environment: %w", err)
			}
			if existing != nil {
				return nil, fmt.Errorf("environment with name '%s' already exists", *req.Name)
			}
			env.Name = *req.Name
		}
	}

	if req.DisplayName != nil {
		env.DisplayName = *req.DisplayName
	}

	if req.BaseURL != nil {
		env.BaseURL = *req.BaseURL
	}

	if req.Variables != nil {
		varsJSON, err := json.Marshal(*req.Variables)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal variables: %w", err)
		}
		env.Variables = string(varsJSON)
	}

	if req.Headers != nil {
		headersJSON, err := json.Marshal(*req.Headers)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal headers: %w", err)
		}
		env.Headers = string(headersJSON)
	}

	// Update in database
	if err := s.repo.Update(ctx, env); err != nil {
		return nil, fmt.Errorf("failed to update environment: %w", err)
	}

	// Convert to response
	resp, err := env.ToResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to response: %w", err)
	}

	return resp, nil
}

// DeleteEnvironment deletes an environment
func (s *service) DeleteEnvironment(ctx context.Context, id uint) error {
	// Check if environment exists
	env, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}
	if env == nil {
		return fmt.Errorf("environment not found")
	}

	// Delete environment
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	return nil
}

// DuplicateEnvironment duplicates an environment
func (s *service) DuplicateEnvironment(ctx context.Context, id uint, req *DuplicateEnvironmentRequest) (*EnvironmentResponse, error) {
	// Get source environment
	source, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get source environment: %w", err)
	}
	if source == nil {
		return nil, fmt.Errorf("source environment not found")
	}

	// Check if new name conflicts
	existing, err := s.repo.GetByProjectAndName(ctx, source.ProjectID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing environment: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("environment with name '%s' already exists", req.Name)
	}

	// Create new environment by copying source
	newEnv := &EnvironmentPO{
		ProjectID:   source.ProjectID,
		Name:        req.Name,
		DisplayName: fmt.Sprintf("Copy of %s", source.DisplayName),
		BaseURL:     source.BaseURL,
		Variables:   source.Variables,
		Headers:     source.Headers,
	}

	// Override variables if provided
	if req.OverrideVars != nil {
		// Parse existing variables
		var vars map[string]interface{}
		if source.Variables != "" {
			if err := json.Unmarshal([]byte(source.Variables), &vars); err != nil {
				return nil, fmt.Errorf("failed to unmarshal source variables: %w", err)
			}
		} else {
			vars = make(map[string]interface{})
		}

		// Merge with override
		for k, v := range *req.OverrideVars {
			vars[k] = v
		}

		// Marshal back
		varsJSON, err := json.Marshal(vars)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal variables: %w", err)
		}
		newEnv.Variables = string(varsJSON)
	}

	// Create new environment
	if err := s.repo.Create(ctx, newEnv); err != nil {
		return nil, fmt.Errorf("failed to create duplicated environment: %w", err)
	}

	// Convert to response
	resp, err := newEnv.ToResponse()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to response: %w", err)
	}

	return resp, nil
}
