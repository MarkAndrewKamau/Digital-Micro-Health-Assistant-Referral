package models

import (
	"time"

	"github.com/google/uuid"
)

type TriageLevel string
type TriageStatus string

const (
	TriageLevelRed    TriageLevel = "red"
	TriageLevelYellow TriageLevel = "yellow"
	TriageLevelGreen  TriageLevel = "green"
)

const (
	TriageStatusQueued     TriageStatus = "queued"
	TriageStatusProcessing TriageStatus = "processing"
	TriageStatusCompleted  TriageStatus = "completed"
	TriageStatusFailed     TriageStatus = "failed"
)

type TriageSession struct {
	ID                uuid.UUID              `json:"id"`
	PatientID         *uuid.UUID             `json:"patient_id,omitempty"`
	Symptoms          map[string]interface{} `json:"symptoms"`
	SummaryText       *string                `json:"summary_text,omitempty"`
	TriageLevel       *TriageLevel           `json:"triage_level,omitempty"`
	TriageCode        *string                `json:"triage_code,omitempty"`
	Confidence        *float64               `json:"confidence,omitempty"`
	RecommendedAction *string                `json:"recommended_action,omitempty"`
	LLMResponse       map[string]interface{} `json:"llm_response,omitempty"`
	Channel           string                 `json:"channel"`
	Status            TriageStatus           `json:"status"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type CreateTriageRequest struct {
	PatientID *uuid.UUID             `json:"patient_id"`
	Symptoms  map[string]interface{} `json:"symptoms" binding:"required"`
	Channel   string                 `json:"channel" binding:"required,oneof=sms ussd web"`
	Context   map[string]interface{} `json:"context"`
}

type TriageResponse struct {
	SessionID         uuid.UUID              `json:"session_id"`
	Status            TriageStatus           `json:"status"`
	TriageLevel       *TriageLevel           `json:"triage_level,omitempty"`
	RecommendedAction *string                `json:"recommended_action,omitempty"`
	Message           string                 `json:"message"`
	CreatedAt         time.Time              `json:"created_at"`
}