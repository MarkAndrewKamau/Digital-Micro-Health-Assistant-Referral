package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalEnv := os.Getenv("ENVIRONMENT")
	originalPort := os.Getenv("PORT")

	// Clean up after test
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
		os.Setenv("PORT", originalPort)
	}()

	t.Run("Load with defaults", func(t *testing.T) {
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("PORT")

		cfg := Load()

		assert.Equal(t, "development", cfg.Environment)
		assert.Equal(t, "8080", cfg.Port)
		assert.NotEmpty(t, cfg.DatabaseURL)
		assert.NotEmpty(t, cfg.RedisURL)
	})

	t.Run("Load with custom env vars", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("PORT", "3000")

		cfg := Load()

		assert.Equal(t, "production", cfg.Environment)
		assert.Equal(t, "3000", cfg.Port)
	})
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "Returns env value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "Returns default when env not set",
			key:          "MISSING_KEY",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}