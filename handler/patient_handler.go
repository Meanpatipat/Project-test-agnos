package handler

import (
	"net/http"

	"hospital-middleware/models"
	"hospital-middleware/repository"

	"github.com/gin-gonic/gin"
)

// PatientHandler handles patient-related HTTP requests
type PatientHandler struct {
	repo repository.PatientRepository
}

// NewPatientHandler creates a new PatientHandler
func NewPatientHandler(repo repository.PatientRepository) *PatientHandler {
	return &PatientHandler{repo: repo}
}

// SearchPatient handles GET /patient/search
// Requires authentication. Searches only within the staff member's hospital.
// All query parameters are optional; if none are provided, returns all patients in the hospital.
func (h *PatientHandler) SearchPatient(c *gin.Context) {
	// Get hospital_id from JWT context (set by auth middleware)
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Authentication required.",
		})
		return
	}

	hID, ok := hospitalID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid session data.",
		})
		return
	}

	// Bind optional search parameters
	var req models.PatientSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid search parameters.",
		})
		return
	}

	// Search patients scoped to the staff's hospital
	patients, err := h.repo.Search(req, hID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to search patients.",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Search completed successfully.",
		Data:    patients,
	})
}
