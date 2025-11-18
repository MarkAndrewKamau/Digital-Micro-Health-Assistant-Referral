package models

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID                uuid.UUID         `json:"id"`
	Phone             string            `json:"phone"`
	Name              *string           `json:"name,omitempty"`
	DateOfBirth       *time.Time        `json:"date_of_birth,omitempty"`
	Gender            *string           `json:"gender,omitempty"`
	PreferredLanguage string            `json:"preferred_language"`
	ConsentFlags      map[string]bool   `json:"consent_flags"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type CreatePatientRequest struct {
	Phone             string          `json:"phone" binding:"required"`
	Name              *string         `json:"name"`
	DateOfBirth       *time.Time      `json:"date_of_birth"`
	Gender            *string         `json:"gender"`
	PreferredLanguage string          `json:"preferred_language"`
	ConsentFlags      map[string]bool `json:"consent_flags"`
}