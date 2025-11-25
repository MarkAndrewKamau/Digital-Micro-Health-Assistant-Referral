package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/repository"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/pkg/response"
)

type TriageHandler struct {
	triageRepo *repository.TriageRepository
}

func NewTriageHandler(triageRepo *repository.TriageRepository) *TriageHandler {
	return &TriageHandler{triageRepo: triageRepo}
}

// CreateTriage handles POST /v1/triage
func (h *TriageHandler) CreateTriage(c *gin.Context) {
	var req models.CreateTriageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body: "+err.Error())
		return
	}

	// Validate symptoms not empty
	if len(req.Symptoms) == 0 {
		response.Error(c, http.StatusBadRequest, "INVALID_SYMPTOMS", "Symptoms cannot be empty")
		return
	}

	// Create triage session
	session, err := h.triageRepo.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "TRIAGE_CREATE_FAILED", "Failed to create triage session")
		return
	}

	// Build response
	triageResponse := models.TriageResponse{
		SessionID:         session.ID,
		Status:            models.TriageStatusQueued,
		TriageLevel:       session.TriageLevel,
		RecommendedAction: session.RecommendedAction,
		Message:           "Triage session created and queued for processing",
		CreatedAt:         session.CreatedAt,
	}

	response.Success(c, http.StatusCreated, triageResponse)
}

// GetTriage handles GET /v1/triage/:id
func (h *TriageHandler) GetTriage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid triage session ID")
		return
	}

	session, err := h.triageRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Triage session not found")
		return
	}

	response.Success(c, http.StatusOK, session)
}

// GetPatientTriages handles GET /v1/triage/patient/:patient_id
func (h *TriageHandler) GetPatientTriages(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid patient ID")
		return
	}

	limit := 10 // Default limit
	sessions, err := h.triageRepo.GetByPatientID(c.Request.Context(), patientID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "QUERY_FAILED", "Failed to retrieve triage sessions")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}