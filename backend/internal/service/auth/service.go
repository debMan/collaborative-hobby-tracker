package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/debMan/collaborative-hobby-tracker/backend/internal/domain/models"
	"github.com/debMan/collaborative-hobby-tracker/backend/internal/repository"
	"github.com/debMan/collaborative-hobby-tracker/backend/pkg/auth"
)

var (
	// ErrEmailAlreadyExists is returned when trying to register with an email that already exists
	ErrEmailAlreadyExists = errors.New("email already exists")
	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidEmail is returned when email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	// ErrPasswordTooShort is returned when password is too short
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	// ErrEmailRequired is returned when email is empty
	ErrEmailRequired = errors.New("email is required")
	// ErrPasswordRequired is returned when password is empty
	ErrPasswordRequired = errors.New("password is required")
	// ErrNameRequired is returned when name is empty
	ErrNameRequired = errors.New("name is required")
)

const minPasswordLength = 8

// emailRegex is a simple regex for basic email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Service provides authentication operations
type Service struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	tokenDuration time.Duration
}

// AuthResult represents the result of authentication operations
type AuthResult struct {
	Token  string
	UserID string
	Email  string
	Name   string
}

// NewService creates a new auth service
func NewService(userRepo repository.UserRepository, jwtSecret string, tokenDuration time.Duration) *Service {
	return &Service{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		tokenDuration: tokenDuration,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, email, password, name string) (*AuthResult, error) {
	// Validate inputs
	if email == "" {
		return nil, ErrEmailRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}
	if name == "" {
		return nil, ErrNameRequired
	}

	// Validate email format
	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	// Validate password length
	if len(password) < minPasswordLength {
		return nil, ErrPasswordTooShort
	}

	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, email, s.jwtSecret, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResult{
		Token:  token,
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Name:   user.Name,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	// Validate inputs
	if email == "" {
		return nil, ErrEmailRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	err = auth.VerifyPassword(user.PasswordHash, password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, email, s.jwtSecret, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResult{
		Token:  token,
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Name:   user.Name,
	}, nil
}
