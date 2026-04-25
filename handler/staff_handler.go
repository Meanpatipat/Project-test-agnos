package handler

import (
	"net/http"
	"strings"

	"hospital-middleware/middleware"
	"hospital-middleware/models"
	"hospital-middleware/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// StaffHandler handles staff-related HTTP requests
type StaffHandler struct {
	staffRepo    repository.StaffRepository
	hospitalRepo repository.HospitalRepository
	jwtSecret    string
}

// NewStaffHandler creates a new StaffHandler
func NewStaffHandler(staffRepo repository.StaffRepository, hospitalRepo repository.HospitalRepository, jwtSecret string) *StaffHandler {
	return &StaffHandler{
		staffRepo:    staffRepo,
		hospitalRepo: hospitalRepo,
		jwtSecret:    jwtSecret,
	}
}

// CreateStaff handles POST /staff/create
func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req models.StaffCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid input. Required fields: username (min 3 chars), password (min 6 chars), hospital (hospital code).",
		})
		return
	}

	// Find the hospital by code
	hospital, err := h.hospitalRepo.FindByCode(req.Hospital)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Hospital not found. Please provide a valid hospital code.",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to process password.",
		})
		return
	}

	// Create staff record
	staff := &models.Staff{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		HospitalID:   hospital.ID,
		Hospital:     *hospital,
	}

	if err := h.staffRepo.Create(staff); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Status:  http.StatusConflict,
				Message: "Username already exists in this hospital.",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create staff record.",
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Status:  http.StatusCreated,
		Message: "Staff member created successfully.",
		Data: gin.H{
			"id":       staff.ID,
			"username": staff.Username,
			"hospital": hospital.Code,
		},
	})
}

// LoginStaff handles POST /staff/login
func (h *StaffHandler) LoginStaff(c *gin.Context) {
	var req models.StaffLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid input. Required fields: username, password, hospital.",
		})
		return
	}

	// Find the hospital by code
	hospital, err := h.hospitalRepo.FindByCode(req.Hospital)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalid credentials.",
		})
		return
	}

	// Find staff by username and hospital
	staff, err := h.staffRepo.FindByUsernameAndHospital(req.Username, hospital.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalid credentials.",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalid credentials.",
		})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(staff, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate authentication token.",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Login successful.",
		Data: models.StaffLoginResponse{
			Token:    token,
			Username: staff.Username,
			Hospital: hospital.Code,
		},
	})
}
