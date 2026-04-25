package router

import (
	"hospital-middleware/handler"
	"hospital-middleware/middleware"
	"hospital-middleware/repository"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router with all routes
func SetupRouter(
	patientRepo repository.PatientRepository,
	staffRepo repository.StaffRepository,
	hospitalRepo repository.HospitalRepository,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "hospital-middleware"})
	})

	// Staff routes (public)
	staffHandler := handler.NewStaffHandler(staffRepo, hospitalRepo, jwtSecret)
	r.POST("/staff/create", staffHandler.CreateStaff)
	r.POST("/staff/login", staffHandler.LoginStaff)

	// Patient routes (protected)
	patientHandler := handler.NewPatientHandler(patientRepo)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		auth.GET("/patient/search", patientHandler.SearchPatient)
	}

	return r
}
