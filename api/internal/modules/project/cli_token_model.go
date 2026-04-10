package project

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const (
	CLITokenScopeSpecWrite = "spec:write"
	CLITokenScopeRunWrite  = "run:write"
)

var supportedCLITokenScopes = map[string]struct{}{
	CLITokenScopeSpecWrite: {},
	CLITokenScopeRunWrite:  {},
}

// ProjectCLITokenPO persists project-scoped CLI tokens.
type ProjectCLITokenPO struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ProjectID   uint           `gorm:"not null;index"`
	CreatedBy   uint           `gorm:"not null;index"`
	Name        string         `gorm:"size:100;not null"`
	TokenPrefix string         `gorm:"size:32;not null;index"`
	TokenHash   string         `gorm:"size:64;not null;uniqueIndex"`
	Scopes      string         `gorm:"type:text"`
	LastUsedAt  *time.Time
	ExpiresAt   *time.Time
}

func (ProjectCLITokenPO) TableName() string {
	return "project_cli_tokens"
}

// ProjectCLIToken is the service-layer representation of a CLI token.
type ProjectCLIToken struct {
	ID          uint       `json:"id"`
	ProjectID   uint       `json:"project_id"`
	CreatedBy   uint       `json:"created_by"`
	Name        string     `json:"name"`
	TokenPrefix string     `json:"token_prefix"`
	Scopes      []string   `json:"scopes"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (po *ProjectCLITokenPO) toDomain() *ProjectCLIToken {
	if po == nil {
		return nil
	}

	token := &ProjectCLIToken{
		ID:          po.ID,
		ProjectID:   po.ProjectID,
		CreatedBy:   po.CreatedBy,
		Name:        po.Name,
		TokenPrefix: po.TokenPrefix,
		LastUsedAt:  po.LastUsedAt,
		ExpiresAt:   po.ExpiresAt,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}

	if po.Scopes != "" {
		_ = json.Unmarshal([]byte(po.Scopes), &token.Scopes)
	}

	return token
}

func newProjectCLITokenPO(token *ProjectCLIToken, tokenHash string) *ProjectCLITokenPO {
	if token == nil {
		return nil
	}

	scopesJSON, _ := json.Marshal(token.Scopes)

	return &ProjectCLITokenPO{
		ID:          token.ID,
		ProjectID:   token.ProjectID,
		CreatedBy:   token.CreatedBy,
		Name:        token.Name,
		TokenPrefix: token.TokenPrefix,
		TokenHash:   tokenHash,
		Scopes:      string(scopesJSON),
		LastUsedAt:  token.LastUsedAt,
		ExpiresAt:   token.ExpiresAt,
	}
}
