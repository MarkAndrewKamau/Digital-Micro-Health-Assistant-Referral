package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/database"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/repository"
)

type AuthService struct {
	redis    *database.Redis
	userRepo *repository.UserRepository
}

func NewAuthService(redis *database.Redis, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		redis:    redis,
		userRepo: userRepo,
	}
}

// Register creates a new user
func (s *AuthService) Register(ctx context.Context, phone string, role models.UserRole) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user already exists with this phone number")
	}

	// Create new user
	user, err := s.userRepo.Create(ctx, phone, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates user and returns user object (does NOT create user)
func (s *AuthService) Login(ctx context.Context, phone string) (*models.User, error) {
	// Get user - must exist
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found - please register first")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}

	return user, nil
}

// CreateSession creates a new session for the user
func (s *AuthService) CreateSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*models.Session, error) {
	return s.userRepo.CreateSession(ctx, userID, userAgent, ipAddress)
}

// ValidateSession validates a session token
func (s *AuthService) ValidateSession(ctx context.Context, sessionToken string) (*models.User, error) {
	session, err := s.userRepo.GetSessionByToken(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	// Check if session is expired
	if session.IsExpired() {
		return nil, fmt.Errorf("session expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user is inactive")
	}

	return user, nil
}

// DeleteSession deletes a session (logout)
func (s *AuthService) DeleteSession(ctx context.Context, sessionToken string) error {
	return s.userRepo.DeleteSession(ctx, sessionToken)
}