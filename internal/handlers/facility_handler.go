package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/repository"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/pkg/response"
)

type FacilityHandler struct {
	facilityRepo *repository.FacilityRepository
}

func NewFacilityHandler(facilityRepo *repository.FacilityRepository) *FacilityHandler {
	return &FacilityHandler{facilityRepo: facilityRepo}
}

// GetFacility handles GET /v1/facilities/:id
func (h *FacilityHandler) GetFacility(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid facility ID")
		return
	}

	facility, err := h.facilityRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Facility not found")
		return
	}

	response.Success(c, http.StatusOK, facility)
}

// GetNearbyFacilities handles GET /v1/facilities/nearby
func (h *FacilityHandler) GetNearbyFacilities(c *gin.Context) {
	var req models.NearbyFacilitiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid query parameters")
		return
	}

	facilities, err := h.facilityRepo.GetNearby(c.Request.Context(), req.Latitude, req.Longitude, req.RadiusKM)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "QUERY_FAILED", "Failed to query facilities")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"facilities": facilities,
		"count":      len(facilities),
	})
}

// ListFacilities handles GET /v1/facilities
func (h *FacilityHandler) ListFacilities(c *gin.Context) {
	county := c.Query("county")
	facilityTypeStr := c.Query("type")

	var county_ *string
	var facilityType_ *models.FacilityType

	if county != "" {
		county_ = &county
	}

	if facilityTypeStr != "" {
		ft := models.FacilityType(facilityTypeStr)
		facilityType_ = &ft
	}

	facilities, err := h.facilityRepo.List(c.Request.Context(), county_, facilityType_)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "QUERY_FAILED", "Failed to list facilities")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"facilities": facilities,
		"count":      len(facilities),
	})
}