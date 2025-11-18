package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/repository"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/pkg/response"
)

type PatientHandler struct {
	patientRepo *repository.PatientRepository
}

func NewPatientHandler(patientRepo *repository.PatientRepository) *PatientHandler {
	return &PatientHandler{patientRepo: patientRepo}
}

// CreatePatient handles POST /v1/patients
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req models.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	patient, err := h.patientRepo.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create patient")
		return
	}

	response.Success(c, http.StatusCreated, patient)
}

// GetPatient handles GET /v1/patients/:id
func (h *PatientHandler) GetPatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid patient ID")
		return
	}

	patient, err := h.patientRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Patient not found")
		return
	}

	response.Success(c, http.StatusOK, patient)
}

// UpdatePatient handles PUT /v1/patients/:id
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid patient ID")
		return
	}

	var req models.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	patient, err := h.patientRepo.Update(c.Request.Context(), id, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update patient")
		return
	}

	response.Success(c, http.StatusOK, patient)
}