package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, phone string, role models.UserRole) (*models.User, error) {
	query := `
		INSERT INTO users (phone, role)
		VALUES ($1, $2)
		RETURNING id, phone, name, email, role, patient_id, clinician_id, is_active, last_login_at, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, phone, role).Scan(
		&user.ID,
		&user.Phone,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.PatientID,
		&user.ClinicianID,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, phone, name, email, role, patient_id, clinician_id, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.PatientID,
		&user.ClinicianID,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, phone, name, email, role, patient_id, clinician_id, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE phone = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.Phone,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.PatientID,
		&user.ClinicianID,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // Not found, not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET last_login_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// GenerateSessionToken generates a secure random session token
func (r *UserRepository) GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (r *UserRepository) CreateSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*models.Session, error) {
	sessionToken, err := r.GenerateSessionToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days

	query := `
		INSERT INTO sessions (user_id, session_token, expires_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, session_token, expires_at, user_agent, ip_address, created_at
	`

	var session models.Session
	err = r.db.QueryRow(ctx, query, userID, sessionToken, expiresAt, userAgent, ipAddress).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionToken,
		&session.ExpiresAt,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &session, nil
}

func (r *UserRepository) GetSessionByToken(ctx context.Context, sessionToken string) (*models.Session, error) {
	query := `
		SELECT id, user_id, session_token, expires_at, user_agent, ip_address, created_at
		FROM sessions
		WHERE session_token = $1
	`

	var session models.Session
	err := r.db.QueryRow(ctx, query, sessionToken).Scan(
		&session.ID,
		&session.UserID,
		&session.SessionToken,
		&session.ExpiresAt,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *UserRepository) DeleteSession(ctx context.Context, sessionToken string) error {
	query := `DELETE FROM sessions WHERE session_token = $1`

	_, err := r.db.Exec(ctx, query, sessionToken)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}