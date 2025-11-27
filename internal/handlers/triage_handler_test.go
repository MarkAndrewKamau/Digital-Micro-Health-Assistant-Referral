package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
)

// Mock TriageRepository
type MockTriageRepository struct {
	mock.Mock
}

func (m *MockTriageRepository) Create(ctx context.Context, req *models.CreateTriageRequest) (*models.TriageSession, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TriageSession), args.Error(1)
}

func (m *MockTriageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TriageSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TriageSession), args.Error(1)
}

func (m *MockTriageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TriageStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTriageRepository) UpdateTriageResult(ctx context.Context, id uuid.UUID, level models.TriageLevel, code string, confidence float64, action string, llmResponse map[string]interface{}) error {
	args := m.Called(ctx, id, level, code, confidence, action, llmResponse)
	return args.Error(0)
}

func (m *MockTriageRepository) GetByPatientID(ctx context.Context, patientID uuid.UUID, limit int) ([]*models.TriageSession, error) {
	args := m.Called(ctx, patientID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TriageSession), args.Error(1)
}

func TestCreateTriage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Create triage session", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		sessionID := uuid.New()
		expectedSession := &models.TriageSession{
			ID:        sessionID,
			Symptoms:  map[string]interface{}{"fever": true},
			Channel:   "web",
			Status:    models.TriageStatusQueued,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.CreateTriageRequest")).
			Return(expectedSession, nil)

		router := gin.New()
		router.POST("/triage", handler.CreateTriage)

		reqBody := models.CreateTriageRequest{
			Symptoms: map[string]interface{}{"fever": true},
			Channel:  "web",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/triage", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, sessionID.String(), data["session_id"])
		assert.Equal(t, "queued", data["status"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Invalid request body", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		router := gin.New()
		router.POST("/triage", handler.CreateTriage)

		req, _ := http.NewRequest("POST", "/triage", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Fail - Empty symptoms", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		router := gin.New()
		router.POST("/triage", handler.CreateTriage)

		reqBody := models.CreateTriageRequest{
			Symptoms: map[string]interface{}{},
			Channel:  "web",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/triage", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))

		errorData := response["error"].(map[string]interface{})
		assert.Equal(t, "INVALID_SYMPTOMS", errorData["code"])
	})

	t.Run("Fail - Invalid channel", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		router := gin.New()
		router.POST("/triage", handler.CreateTriage)

		reqBody := map[string]interface{}{
			"symptoms": map[string]interface{}{"fever": true},
			"channel":  "invalid_channel",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/triage", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetTriage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Get triage session", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		sessionID := uuid.New()
		expectedSession := &models.TriageSession{
			ID:        sessionID,
			Symptoms:  map[string]interface{}{"fever": true},
			Channel:   "web",
			Status:    models.TriageStatusQueued,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("GetByID", mock.Anything, sessionID).Return(expectedSession, nil)

		router := gin.New()
		router.GET("/triage/:id", handler.GetTriage)

		req, _ := http.NewRequest("GET", "/triage/"+sessionID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Invalid UUID", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		router := gin.New()
		router.GET("/triage/:id", handler.GetTriage)

		req, _ := http.NewRequest("GET", "/triage/invalid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Fail - Session not found", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		sessionID := uuid.New()
		mockRepo.On("GetByID", mock.Anything, sessionID).Return(nil, assert.AnError)

		router := gin.New()
		router.GET("/triage/:id", handler.GetTriage)

		req, _ := http.NewRequest("GET", "/triage/"+sessionID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetPatientTriages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Get patient triage sessions", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		patientID := uuid.New()
		sessions := []*models.TriageSession{
			{
				ID:        uuid.New(),
				PatientID: &patientID,
				Symptoms:  map[string]interface{}{"fever": true},
				Channel:   "web",
				Status:    models.TriageStatusQueued,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockRepo.On("GetByPatientID", mock.Anything, patientID, 10).Return(sessions, nil)

		router := gin.New()
		router.GET("/triage/patient/:patient_id", handler.GetPatientTriages)

		req, _ := http.NewRequest("GET", "/triage/patient/"+patientID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["count"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail - Invalid patient ID", func(t *testing.T) {
		mockRepo := new(MockTriageRepository)
		handler := NewTriageHandler(mockRepo)

		router := gin.New()
		router.GET("/triage/patient/:patient_id", handler.GetPatientTriages)

		req, _ := http.NewRequest("GET", "/triage/patient/invalid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}