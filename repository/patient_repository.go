package repository

import (
	"hospital-middleware/models"
)

// PatientRepository defines the interface for patient data access
type PatientRepository interface {
	// FindByID searches for a patient by national_id or passport_id within a specific hospital
	FindByID(id string, hospitalID uint) (*models.Patient, error)

	// Search searches for patients matching optional filters within a specific hospital
	Search(req models.PatientSearchRequest, hospitalID uint) ([]models.Patient, error)
}

// StaffRepository defines the interface for staff data access
type StaffRepository interface {
	// Create creates a new staff member
	Create(staff *models.Staff) error

	// FindByUsernameAndHospital finds a staff member by username and hospital ID
	FindByUsernameAndHospital(username string, hospitalID uint) (*models.Staff, error)
}

// HospitalRepository defines the interface for hospital data access
type HospitalRepository interface {
	// FindByCode finds a hospital by its unique code
	FindByCode(code string) (*models.Hospital, error)
}
