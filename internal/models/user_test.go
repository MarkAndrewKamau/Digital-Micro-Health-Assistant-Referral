package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "Session not expired",
			expiresAt: time.Now().Add(24 * time.Hour),
			want:      false,
		},
		{
			name:      "Session expired",
			expiresAt: time.Now().Add(-24 * time.Hour),
			want:      true,
		},
		{
			name:      "Session expires in 1 minute",
			expiresAt: time.Now().Add(1 * time.Minute),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &Session{
				ExpiresAt: tt.expiresAt,
			}
			got := session.IsExpired()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRole(t *testing.T) {
	assert.Equal(t, UserRole("patient"), UserRolePatient)
	assert.Equal(t, UserRole("chv"), UserRoleCHV)
	assert.Equal(t, UserRole("clinician"), UserRoleClinician)
	assert.Equal(t, UserRole("admin"), UserRoleAdmin)
}