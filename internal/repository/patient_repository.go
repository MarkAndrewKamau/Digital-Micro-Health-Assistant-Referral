package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/digital-health-assistant/internal/models"
)

type PatientRepository struct {
	db *pgxpool.Pool
}

func NewPatientRepository(db *pgxpool.Pool) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Create(ctx context.Context, req *models.CreatePatientRequest) (*models.Patient, error) {
	consentFlagsJSON, err := json.Marshal(req.ConsentFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal consent flags: %w", err)
	}

	query := `
		INSERT INTO patients (phone, name, date_of_birth, gender, preferred_language, consent_flags)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, phone, name, date_of_birth, gender, preferred_language, consent_flags, created_at, updated_at
	`

	var patient models.Patient
	var consentFlagsRaw []byte

	err = r.db.QueryRow(ctx, query,
		req.Phone,
		req.Name,
		req.DateOfBirth,
		req.Gender,
		req.PreferredLanguage,
		consentFlagsJSON,
	).Scan(
		&patient.ID,
		&patient.Phone,
		&patient.Name,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.PreferredLanguage,
		&consentFlagsRaw,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	if err := json.Unmarshal(consentFlagsRaw, &patient.ConsentFlags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal consent flags: %w", err)
	}

	return &patient, nil
}

func (r *PatientRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Patient, error) {
	query := `
		SELECT id, phone, name, date_of_birth, gender, preferred_language, consent_flags, created_at, updated_at
		FROM patients
		WHERE id = $1
	`

	var patient models.Patient
	var consentFlagsRaw []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&patient.ID,
		&patient.Phone,
		&patient.Name,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.PreferredLanguage,
		&consentFlagsRaw,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("patient not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	if err := json.Unmarshal(consentFlagsRaw, &patient.ConsentFlags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal consent flags: %w", err)
	}

	return &patient, nil
}

func (r *PatientRepository) GetByPhone(ctx context.Context, phone string) (*models.Patient, error) {
	query := `
		SELECT id, phone, name, date_of_birth, gender, preferred_language, consent_flags, created_at, updated_at
		FROM patients
		WHERE phone = $1
	`

	var patient models.Patient
	var consentFlagsRaw []byte

	err := r.db.QueryRow(ctx, query, phone).Scan(
		&patient.ID,
		&patient.Phone,
		&patient.Name,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.PreferredLanguage,
		&consentFlagsRaw,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // Return nil, nil for not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get patient by phone: %w", err)
	}

	if err := json.Unmarshal(consentFlagsRaw, &patient.ConsentFlags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal consent flags: %w", err)
	}

	return &patient, nil
}

func (r *PatientRepository) Update(ctx context.Context, id uuid.UUID, req *models.CreatePatientRequest) (*models.Patient, error) {
	consentFlagsJSON, err := json.Marshal(req.ConsentFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal consent flags: %w", err)
	}

	query := `
		UPDATE patients
		SET name = $1, date_of_birth = $2, gender = $3, preferred_language = $4, consent_flags = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, phone, name, date_of_birth, gender, preferred_language, consent_flags, created_at, updated_at
	`

	var patient models.Patient
	var consentFlagsRaw []byte

	err = r.db.QueryRow(ctx, query,
		req.Name,
		req.DateOfBirth,
		req.Gender,
		req.PreferredLanguage,
		consentFlagsJSON,
		id,
	).Scan(
		&patient.ID,
		&patient.Phone,
		&patient.Name,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.PreferredLanguage,
		&consentFlagsRaw,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("patient not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}

	if err := json.Unmarshal(consentFlagsRaw, &patient.ConsentFlags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal consent flags: %w", err)
	}

	return &patient, nil
}