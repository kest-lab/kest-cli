package project

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/kest-labs/kest/api/internal/modules/member"
	idpkg "github.com/kest-labs/kest/api/pkg/id"
)

// Common errors
var (
	ErrProjectNotFound          = errors.New("project not found")
	ErrSlugAlreadyExists        = errors.New("project slug already exists")
	ErrInvalidCLIToken          = errors.New("invalid CLI token")
	ErrCLITokenExpired          = errors.New("CLI token has expired")
	ErrCLITokenScopeDenied      = errors.New("CLI token does not have the required scope")
	ErrCLITokenProjectMismatch  = errors.New("CLI token is not scoped to this project")
	ErrUnsupportedCLITokenScope = errors.New("unsupported CLI token scope")
)

// Service defines the interface for project business logic
type Service interface {
	Create(ctx context.Context, userID string, req *CreateProjectRequest) (*Project, error)
	GetByID(ctx context.Context, id string) (*Project, error)
	Update(ctx context.Context, id string, req *UpdateProjectRequest) (*Project, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string, page, perPage int) ([]*Project, int64, error)
	GetStats(ctx context.Context, projectID string) (*ProjectStats, error)
	GenerateCLIToken(ctx context.Context, projectID string, createdBy string, req *GenerateProjectCLITokenRequest) (*GenerateProjectCLITokenResponse, error)
	ValidateCLIToken(ctx context.Context, projectID string, rawToken string, requiredScopes []string) (string, string, error)
}

// service implements Service interface
type service struct {
	repo          Repository
	memberService member.Service
}

// NewService creates a new project service
func NewService(repo Repository, memberService member.Service) Service {
	return &service{
		repo:          repo,
		memberService: memberService,
	}
}

func (s *service) Create(ctx context.Context, userID string, req *CreateProjectRequest) (*Project, error) {
	// Generate slug from name if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	// Check if slug already exists
	existing, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSlugAlreadyExists
	}

	// Generate public key
	publicKey, err := generatePublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate public key: %w", err)
	}

	// Create project
	project := &Project{
		Name:      req.Name,
		Slug:      slug,
		Platform:  req.Platform,
		PublicKey: publicKey,
		Status:    1, // Active by default
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}

	// Auto-assign creator as Owner
	_, err = s.memberService.AddMember(ctx, project.ID, &member.AddMemberRequest{
		UserID: idpkg.Compatible(userID),
		Role:   member.RoleOwner,
	})
	if err != nil {
		// Rollback project creation? For now, we just log or return error
		return nil, fmt.Errorf("failed to assign owner: %w", err)
	}

	return project, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

func (s *service) Update(ctx context.Context, id string, req *UpdateProjectRequest) (*Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	// Apply updates
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Platform != "" {
		project.Platform = req.Platform
	}
	if req.Status != nil {
		project.Status = *req.Status
	}

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrProjectNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, userID string, page, perPage int) ([]*Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage
	return s.repo.List(ctx, userID, offset, perPage)
}

func (s *service) GetStats(ctx context.Context, projectID string) (*ProjectStats, error) {
	return s.repo.GetStats(ctx, projectID)
}

func (s *service) GenerateCLIToken(ctx context.Context, projectID string, createdBy string, req *GenerateProjectCLITokenRequest) (*GenerateProjectCLITokenResponse, error) {
	if req == nil {
		req = &GenerateProjectCLITokenRequest{}
	}

	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	scopes, err := normalizeCLITokenScopes(req.Scopes)
	if err != nil {
		return nil, err
	}

	rawToken, tokenPrefix, tokenHash, err := generateCLITokenMaterial()
	if err != nil {
		return nil, fmt.Errorf("failed to generate CLI token: %w", err)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = fmt.Sprintf("%s CLI token", project.Name)
	}

	token := &ProjectCLIToken{
		ProjectID:   projectID,
		CreatedBy:   createdBy,
		Name:        name,
		TokenPrefix: tokenPrefix,
		Scopes:      scopes,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := s.repo.CreateCLIToken(ctx, token, tokenHash); err != nil {
		return nil, err
	}

	return &GenerateProjectCLITokenResponse{
		Token:     rawToken,
		TokenType: "bearer",
		ProjectID: projectID,
		TokenInfo: toProjectCLITokenResponse(token),
	}, nil
}

func (s *service) ValidateCLIToken(ctx context.Context, projectID string, rawToken string, requiredScopes []string) (string, string, error) {
	tokenHash := hashCLIToken(strings.TrimSpace(rawToken))
	if tokenHash == "" {
		return "", "", ErrInvalidCLIToken
	}

	token, err := s.repo.GetCLITokenByHash(ctx, tokenHash)
	if err != nil {
		return "", "", err
	}
	if token == nil {
		return "", "", ErrInvalidCLIToken
	}
	if token.ProjectID != projectID {
		return "", "", ErrCLITokenProjectMismatch
	}
	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return "", "", ErrCLITokenExpired
	}

	scopes, err := normalizeRequiredCLITokenScopes(requiredScopes)
	if err != nil {
		return "", "", err
	}
	if !hasRequiredScopes(token.Scopes, scopes) {
		return "", "", ErrCLITokenScopeDenied
	}

	if err := s.repo.TouchCLIToken(ctx, token.ID, time.Now().UTC()); err != nil {
		return "", "", err
	}

	return token.ID, token.CreatedBy, nil
}

// generateSlug creates a URL-safe slug from a name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	re := regexp.MustCompile(`[^a-z0-9-]`)
	slug = re.ReplaceAllString(slug, "")

	// Remove duplicate hyphens
	re = regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim hyphens from ends
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 50 {
		slug = slug[:50]
	}

	// Add random suffix if empty
	if slug == "" {
		slug = fmt.Sprintf("project-%d", time.Now().UnixNano()%100000000)
	}

	return slug
}

// generatePublicKey generates a random public key for the project
func generatePublicKey() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 32 hex characters (fits in varchar(64))
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateCLITokenMaterial() (rawToken, tokenPrefix, tokenHash string, err error) {
	bytes := make([]byte, 24)
	if _, err = rand.Read(bytes); err != nil {
		return "", "", "", err
	}

	rawToken = "kest_pat_" + hex.EncodeToString(bytes)
	tokenPrefix = rawToken
	if len(tokenPrefix) > 18 {
		tokenPrefix = tokenPrefix[:18]
	}

	return rawToken, tokenPrefix, hashCLIToken(rawToken), nil
}

func hashCLIToken(rawToken string) string {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return ""
	}

	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func normalizeCLITokenScopes(scopes []string) ([]string, error) {
	if len(scopes) == 0 {
		return []string{CLITokenScopeSpecWrite}, nil
	}

	seen := make(map[string]struct{}, len(scopes))
	normalized := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := supportedCLITokenScopes[scope]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedCLITokenScope, scope)
		}
		if _, exists := seen[scope]; exists {
			continue
		}
		seen[scope] = struct{}{}
		normalized = append(normalized, scope)
	}

	if len(normalized) == 0 {
		return []string{CLITokenScopeSpecWrite}, nil
	}

	sort.Strings(normalized)
	return normalized, nil
}

func normalizeRequiredCLITokenScopes(scopes []string) ([]string, error) {
	if len(scopes) == 0 {
		return nil, nil
	}

	return normalizeCLITokenScopes(scopes)
}

func hasRequiredScopes(actual, required []string) bool {
	if len(required) == 0 {
		return true
	}

	available := make(map[string]struct{}, len(actual))
	for _, scope := range actual {
		available[scope] = struct{}{}
	}

	for _, scope := range required {
		if _, ok := available[scope]; !ok {
			return false
		}
	}

	return true
}
