package repository

import (
	"fmt"
	"strings"

	"hospital-middleware/models"

	"gorm.io/gorm"
)

// PostgresPatientRepository implements PatientRepository with PostgreSQL via GORM
type PostgresPatientRepository struct {
	db *gorm.DB
}

// NewPostgresPatientRepository creates a new PostgresPatientRepository
func NewPostgresPatientRepository(db *gorm.DB) *PostgresPatientRepository {
	return &PostgresPatientRepository{db: db}
}

// FindByID searches for a patient by national_id or passport_id within a hospital
func (r *PostgresPatientRepository) FindByID(id string, hospitalID uint) (*models.Patient, error) {
	var patient models.Patient
	result := r.db.Where(
		"hospital_id = ? AND (national_id = ? OR passport_id = ?)",
		hospitalID, id, id,
	).First(&patient)

	if result.Error != nil {
		return nil, fmt.Errorf("patient not found with id: %s", id)
	}
	return &patient, nil
}

// Search searches for patients matching optional filters within a hospital
func (r *PostgresPatientRepository) Search(req models.PatientSearchRequest, hospitalID uint) ([]models.Patient, error) {
	query := r.db.Where("hospital_id = ?", hospitalID)

	if req.NationalID != "" {
		query = query.Where("national_id = ?", req.NationalID)
	}
	if req.PassportID != "" {
		query = query.Where("passport_id = ?", req.PassportID)
	}
	if req.FirstName != "" {
		name := "%" + strings.ToLower(req.FirstName) + "%"
		query = query.Where("(LOWER(first_name_en) LIKE ? OR LOWER(first_name_th) LIKE ?)", name, name)
	}
	if req.MiddleName != "" {
		name := "%" + strings.ToLower(req.MiddleName) + "%"
		query = query.Where("(LOWER(middle_name_en) LIKE ? OR LOWER(middle_name_th) LIKE ?)", name, name)
	}
	if req.LastName != "" {
		name := "%" + strings.ToLower(req.LastName) + "%"
		query = query.Where("(LOWER(last_name_en) LIKE ? OR LOWER(last_name_th) LIKE ?)", name, name)
	}
	if req.DateOfBirth != "" {
		query = query.Where("date_of_birth = ?", req.DateOfBirth)
	}
	if req.PhoneNumber != "" {
		query = query.Where("phone_number = ?", req.PhoneNumber)
	}
	if req.Email != "" {
		query = query.Where("LOWER(email) = ?", strings.ToLower(req.Email))
	}

	var patients []models.Patient
	result := query.Find(&patients)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to search patients: %w", result.Error)
	}

	return patients, nil
}

// PostgresStaffRepository implements StaffRepository with PostgreSQL via GORM
type PostgresStaffRepository struct {
	db *gorm.DB
}

// NewPostgresStaffRepository creates a new PostgresStaffRepository
func NewPostgresStaffRepository(db *gorm.DB) *PostgresStaffRepository {
	return &PostgresStaffRepository{db: db}
}

// Create creates a new staff member in the database
func (r *PostgresStaffRepository) Create(staff *models.Staff) error {
	result := r.db.Create(staff)
	return result.Error
}

// FindByUsernameAndHospital finds a staff member by username and hospital ID
func (r *PostgresStaffRepository) FindByUsernameAndHospital(username string, hospitalID uint) (*models.Staff, error) {
	var staff models.Staff
	result := r.db.Preload("Hospital").Where("username = ? AND hospital_id = ?", username, hospitalID).First(&staff)
	if result.Error != nil {
		return nil, fmt.Errorf("staff not found")
	}
	return &staff, nil
}

// PostgresHospitalRepository implements HospitalRepository with PostgreSQL via GORM
type PostgresHospitalRepository struct {
	db *gorm.DB
}

// NewPostgresHospitalRepository creates a new PostgresHospitalRepository
func NewPostgresHospitalRepository(db *gorm.DB) *PostgresHospitalRepository {
	return &PostgresHospitalRepository{db: db}
}

// FindByCode finds a hospital by its unique code
func (r *PostgresHospitalRepository) FindByCode(code string) (*models.Hospital, error) {
	var hospital models.Hospital
	result := r.db.Where("code = ?", code).First(&hospital)
	if result.Error != nil {
		return nil, fmt.Errorf("hospital not found with code: %s", code)
	}
	return &hospital, nil
}
