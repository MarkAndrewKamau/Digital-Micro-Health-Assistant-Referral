package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRolePatient   UserRole = "patient"
	UserRoleCHV       UserRole = "chv"
	UserRoleClinician UserRole = "clinician"
	UserRoleAdmin     UserRole = "admin"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Phone        string     `json:"phone"`
	Name         *string    `json:"name,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Role         UserRole   `json:"role"`
	PatientID    *uuid.UUID `json:"patient_id,omitempty"`
	ClinicianID  *uuid.UUID `json:"clinician_id,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserAgent    *string   `json:"user_agent,omitempty"`
	IPAddress    *string   `json:"ip_address,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type OTPCode struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	Code      string    `json:"code"`
	Purpose   string    `json:"purpose"`
	Verified  bool      `json:"verified"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}