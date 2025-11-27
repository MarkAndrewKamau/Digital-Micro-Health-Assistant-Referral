package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
)

// TriageRepositoryInterface defines the interface for triage operations
type TriageRepositoryInterface interface {
	Create(ctx context.Context, req *models.CreateTriageRequest) (*models.TriageSession, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.TriageSession, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.TriageStatus) error
	UpdateTriageResult(ctx context.Context, id uuid.UUID, level models.TriageLevel, code string, confidence float64, action string, llmResponse map[string]interface{}) error
	GetByPatientID(ctx context.Context, patientID uuid.UUID, limit int) ([]*models.TriageSession, error)
}

type TriageRepository struct {
	db *pgxpool.Pool
}

func NewTriageRepository(db *pgxpool.Pool) *TriageRepository {
	return &TriageRepository{db: db}
}

func (r *TriageRepository) Create(ctx context.Context, req *models.CreateTriageRequest) (*models.TriageSession, error) {
	symptomsJSON, err := json.Marshal(req.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal symptoms: %w", err)
	}

	query := `
		INSERT INTO triage_sessions (patient_id, symptoms, channel)
		VALUES ($1, $2, $3)
		RETURNING id, patient_id, symptoms, summary_text, triage_level, triage_code, 
				  confidence, recommended_action, llm_response, channel, created_at, updated_at
	`

	var session models.TriageSession
	var symptomsRaw, llmResponseRaw []byte
	var triageLevelStr *string

	err = r.db.QueryRow(ctx, query, req.PatientID, symptomsJSON, req.Channel).Scan(
		&session.ID,
		&session.PatientID,
		&symptomsRaw,
		&session.SummaryText,
		&triageLevelStr,
		&session.TriageCode,
		&session.Confidence,
		&session.RecommendedAction,
		&llmResponseRaw,
		&session.Channel,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating triage session: %v", err)
		return nil, fmt.Errorf("failed to create triage session: %w", err)
	}

	// Parse JSON fields
	if err := json.Unmarshal(symptomsRaw, &session.Symptoms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
	}

	if llmResponseRaw != nil && string(llmResponseRaw) != "null" {
		if err := json.Unmarshal(llmResponseRaw, &session.LLMResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal llm_response: %w", err)
		}
	}

	if triageLevelStr != nil {
		level := models.TriageLevel(*triageLevelStr)
		session.TriageLevel = &level
	}

	// Set default status as queued
	session.Status = models.TriageStatusQueued

	return &session, nil
}

func (r *TriageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TriageSession, error) {
	query := `
		SELECT id, patient_id, symptoms, summary_text, triage_level, triage_code,
			   confidence, recommended_action, llm_response, channel, created_at, updated_at
		FROM triage_sessions
		WHERE id = $1
	`

	var session models.TriageSession
	var symptomsRaw, llmResponseRaw []byte
	var triageLevelStr *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.PatientID,
		&symptomsRaw,
		&session.SummaryText,
		&triageLevelStr,
		&session.TriageCode,
		&session.Confidence,
		&session.RecommendedAction,
		&llmResponseRaw,
		&session.Channel,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("triage session not found")
	}
	if err != nil {
		log.Printf("Error getting triage session: %v", err)
		return nil, fmt.Errorf("failed to get triage session: %w", err)
	}

	// Parse JSON fields
	if err := json.Unmarshal(symptomsRaw, &session.Symptoms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
	}

	if llmResponseRaw != nil && string(llmResponseRaw) != "null" {
		if err := json.Unmarshal(llmResponseRaw, &session.LLMResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal llm_response: %w", err)
		}
	}

	if triageLevelStr != nil {
		level := models.TriageLevel(*triageLevelStr)
		session.TriageLevel = &level
	}

	session.Status = models.TriageStatusQueued

	return &session, nil
}

func (r *TriageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TriageStatus) error {
	query := `UPDATE triage_sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error updating triage session status: %v", err)
		return fmt.Errorf("failed to update triage session status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("triage session not found")
	}

	return nil
}

func (r *TriageRepository) UpdateTriageResult(ctx context.Context, id uuid.UUID, level models.TriageLevel, code string, confidence float64, action string, llmResponse map[string]interface{}) error {
	llmResponseJSON, err := json.Marshal(llmResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal llm_response: %w", err)
	}

	query := `
		UPDATE triage_sessions
		SET triage_level = $1, triage_code = $2, confidence = $3, 
		    recommended_action = $4, llm_response = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	result, err := r.db.Exec(ctx, query, level, code, confidence, action, llmResponseJSON, id)
	if err != nil {
		log.Printf("Error updating triage result: %v", err)
		return fmt.Errorf("failed to update triage result: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("triage session not found")
	}

	return nil
}

func (r *TriageRepository) GetByPatientID(ctx context.Context, patientID uuid.UUID, limit int) ([]*models.TriageSession, error) {
	query := `
		SELECT id, patient_id, symptoms, summary_text, triage_level, triage_code,
			   confidence, recommended_action, llm_response, channel, created_at, updated_at
		FROM triage_sessions
		WHERE patient_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, patientID, limit)
	if err != nil {
		log.Printf("Error getting patient triage sessions: %v", err)
		return nil, fmt.Errorf("failed to get patient triage sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.TriageSession

	for rows.Next() {
		var session models.TriageSession
		var symptomsRaw, llmResponseRaw []byte
		var triageLevelStr *string

		err := rows.Scan(
			&session.ID,
			&session.PatientID,
			&symptomsRaw,
			&session.SummaryText,
			&triageLevelStr,
			&session.TriageCode,
			&session.Confidence,
			&session.RecommendedAction,
			&llmResponseRaw,
			&session.Channel,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error scanning triage session: %v", err)
			return nil, fmt.Errorf("failed to scan triage session: %w", err)
		}

		// Parse JSON fields
		if err := json.Unmarshal(symptomsRaw, &session.Symptoms); err != nil {
			return nil, fmt.Errorf("failed to unmarshal symptoms: %w", err)
		}

		if llmResponseRaw != nil && string(llmResponseRaw) != "null" {
			if err := json.Unmarshal(llmResponseRaw, &session.LLMResponse); err != nil {
				return nil, fmt.Errorf("failed to unmarshal llm_response: %w", err)
			}
		}

		if triageLevelStr != nil {
			level := models.TriageLevel(*triageLevelStr)
			session.TriageLevel = &level
		}

		session.Status = models.TriageStatusQueued

		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating triage sessions: %w", err)
	}

	return sessions, nil
}