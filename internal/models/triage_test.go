package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriageLevel(t *testing.T) {
	assert.Equal(t, TriageLevel("red"), TriageLevelRed)
	assert.Equal(t, TriageLevel("yellow"), TriageLevelYellow)
	assert.Equal(t, TriageLevel("green"), TriageLevelGreen)
}

func TestTriageStatus(t *testing.T) {
	assert.Equal(t, TriageStatus("queued"), TriageStatusQueued)
	assert.Equal(t, TriageStatus("processing"), TriageStatusProcessing)
	assert.Equal(t, TriageStatus("completed"), TriageStatusCompleted)
	assert.Equal(t, TriageStatus("failed"), TriageStatusFailed)
}

func TestCreateTriageRequest(t *testing.T) {
	t.Run("Valid triage request", func(t *testing.T) {
		req := CreateTriageRequest{
			Symptoms: map[string]interface{}{
				"fever":    true,
				"duration": 3,
			},
			Channel: "web",
		}

		assert.NotNil(t, req.Symptoms)
		assert.Equal(t, "web", req.Channel)
		assert.Nil(t, req.PatientID)
	})
}