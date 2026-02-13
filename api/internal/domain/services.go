package domain

import (
	"context"
	"errors"
)

// Domain Services encapsulate business logic that doesn't naturally fit within an entity.
// They operate on domain objects and enforce business rules.

// PasswordHasher defines the contract for password hashing
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

// TokenGenerator defines the contract for token generation
type TokenGenerator interface {
	Generate(userID uint, username string) (string, error)
	Validate(token string) (userID uint, username string, err error)
}

// AuthenticationService handles user authentication logic
type AuthenticationService struct {
	userRepo UserRepository
	hasher   PasswordHasher
	tokens   TokenGenerator
}

// NewAuthenticationService creates a new authentication service
func NewAuthenticationService(
	userRepo UserRepository,
	hasher PasswordHasher,
	tokens TokenGenerator,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo: userRepo,
		hasher:   hasher,
		tokens:   tokens,
	}
}

// AuthResult represents the result of authentication
type AuthResult struct {
	User        *User
	AccessToken string
}

// Authenticate authenticates a user by username/email and password
func (s *AuthenticationService) Authenticate(ctx context.Context, identifier, password string) (*AuthResult, error) {
	// Try to find user by username first, then by email
	user, err := s.userRepo.FindByUsername(ctx, identifier)
	if err != nil {
		user, err = s.userRepo.FindByEmail(ctx, identifier)
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	// Check if account is active
	if !user.IsActive() {
		return nil, ErrAccountDisabled
	}

	// Verify password
	if !s.hasher.Verify(password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.tokens.Generate(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:        user,
		AccessToken: token,
	}, nil
}

// ValidateToken validates a token and returns the user
func (s *AuthenticationService) ValidateToken(ctx context.Context, token string) (*User, error) {
	userID, _, err := s.tokens.Validate(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive() {
		return nil, ErrAccountDisabled
	}

	return user, nil
}

// RegistrationService handles user registration logic
type RegistrationService struct {
	userRepo UserRepository
	hasher   PasswordHasher
}

// NewRegistrationService creates a new registration service
func NewRegistrationService(userRepo UserRepository, hasher PasswordHasher) *RegistrationService {
	return &RegistrationService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

// RegistrationRequest represents a registration request
type RegistrationRequest struct {
	Username string
	Email    string
	Password string
	Nickname string
}

// Register registers a new user
func (s *RegistrationService) Register(ctx context.Context, req *RegistrationRequest) (*User, error) {
	// Validate email
	email, err := NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Validate username
	username, err := NewUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(ctx, email.String())
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Username: username.String(),
		Email:    email.String(),
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   int(UserStatusActive),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// PasswordService handles password-related operations
type PasswordService struct {
	userRepo UserRepository
	hasher   PasswordHasher
}

// NewPasswordService creates a new password service
func NewPasswordService(userRepo UserRepository, hasher PasswordHasher) *PasswordService {
	return &PasswordService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

// ChangePassword changes a user's password
func (s *PasswordService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if !s.hasher.Verify(oldPassword, user.Password) {
		return ErrInvalidCredentials
	}

	// Validate new password
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Hash new password
	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user)
}

// ResetPassword resets a user's password (admin operation)
func (s *PasswordService) ResetPassword(ctx context.Context, userID uint, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user)
}

// Specification pattern for complex queries

// Specification defines a business rule that can be combined with others
type Specification[T any] interface {
	IsSatisfiedBy(item T) bool
}

// AndSpecification combines two specifications with AND
type AndSpecification[T any] struct {
	left  Specification[T]
	right Specification[T]
}

// IsSatisfiedBy checks if both specifications are satisfied
func (s *AndSpecification[T]) IsSatisfiedBy(item T) bool {
	return s.left.IsSatisfiedBy(item) && s.right.IsSatisfiedBy(item)
}

// OrSpecification combines two specifications with OR
type OrSpecification[T any] struct {
	left  Specification[T]
	right Specification[T]
}

// IsSatisfiedBy checks if either specification is satisfied
func (s *OrSpecification[T]) IsSatisfiedBy(item T) bool {
	return s.left.IsSatisfiedBy(item) || s.right.IsSatisfiedBy(item)
}

// NotSpecification negates a specification
type NotSpecification[T any] struct {
	spec Specification[T]
}

// IsSatisfiedBy checks if the specification is not satisfied
func (s *NotSpecification[T]) IsSatisfiedBy(item T) bool {
	return !s.spec.IsSatisfiedBy(item)
}

// And combines this specification with another using AND
func And[T any](left, right Specification[T]) Specification[T] {
	return &AndSpecification[T]{left: left, right: right}
}

// Or combines this specification with another using OR
func Or[T any](left, right Specification[T]) Specification[T] {
	return &OrSpecification[T]{left: left, right: right}
}

// Not negates a specification
func Not[T any](spec Specification[T]) Specification[T] {
	return &NotSpecification[T]{spec: spec}
}

// User Specifications

// ActiveUserSpec checks if a user is active
type ActiveUserSpec struct{}

// IsSatisfiedBy checks if the user is active
func (s *ActiveUserSpec) IsSatisfiedBy(user *User) bool {
	return user.IsActive()
}

// EmailDomainSpec checks if a user's email belongs to a specific domain
type EmailDomainSpec struct {
	domain string
}

// NewEmailDomainSpec creates a new email domain specification
func NewEmailDomainSpec(domain string) *EmailDomainSpec {
	return &EmailDomainSpec{domain: domain}
}

// IsSatisfiedBy checks if the user's email belongs to the domain
func (s *EmailDomainSpec) IsSatisfiedBy(user *User) bool {
	email, err := NewEmail(user.Email)
	if err != nil {
		return false
	}
	return email.Domain() == s.domain
}
