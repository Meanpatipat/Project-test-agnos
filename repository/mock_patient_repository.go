package repository

import (
	"fmt"
	"strings"

	"hospital-middleware/models"
)

// ============================================================
// Mock Patient Repository
// ============================================================

// MockPatientRepository provides mock patient data for testing
type MockPatientRepository struct {
	patients []models.Patient
}

// NewMockPatientRepository creates a new MockPatientRepository with sample data
func NewMockPatientRepository() *MockPatientRepository {
	repo := &MockPatientRepository{}
	repo.seedData()
	return repo
}

// FindByID searches for a patient by national_id or passport_id within a hospital
func (r *MockPatientRepository) FindByID(id string, hospitalID uint) (*models.Patient, error) {
	for _, p := range r.patients {
		if p.HospitalID == hospitalID && (p.NationalID == id || p.PassportID == id) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("patient not found with id: %s", id)
}

// Search searches for patients matching optional filters within a hospital
func (r *MockPatientRepository) Search(req models.PatientSearchRequest, hospitalID uint) ([]models.Patient, error) {
	var results []models.Patient

	for _, p := range r.patients {
		if p.HospitalID != hospitalID {
			continue
		}

		if req.NationalID == "ERROR" {
			return nil, fmt.Errorf("simulated repository error")
		}

		if req.NationalID != "" && p.NationalID != req.NationalID {
			continue
		}
		if req.PassportID != "" && p.PassportID != req.PassportID {
			continue
		}
		if req.FirstName != "" {
			lf := strings.ToLower(req.FirstName)
			if !strings.Contains(strings.ToLower(p.FirstNameEN), lf) &&
				!strings.Contains(strings.ToLower(p.FirstNameTH), lf) {
				continue
			}
		}
		if req.MiddleName != "" {
			lm := strings.ToLower(req.MiddleName)
			if !strings.Contains(strings.ToLower(p.MiddleNameEN), lm) &&
				!strings.Contains(strings.ToLower(p.MiddleNameTH), lm) {
				continue
			}
		}
		if req.LastName != "" {
			ll := strings.ToLower(req.LastName)
			if !strings.Contains(strings.ToLower(p.LastNameEN), ll) &&
				!strings.Contains(strings.ToLower(p.LastNameTH), ll) {
				continue
			}
		}
		if req.DateOfBirth != "" && p.DateOfBirth != req.DateOfBirth {
			continue
		}
		if req.PhoneNumber != "" && p.PhoneNumber != req.PhoneNumber {
			continue
		}
		if req.Email != "" && !strings.EqualFold(p.Email, req.Email) {
			continue
		}

		results = append(results, p)
	}

	return results, nil
}

func (r *MockPatientRepository) seedData() {
	r.patients = []models.Patient{
		// Hospital A (ID = 1)
		{
			ID: 1, HospitalID: 1,
			FirstNameTH: "สมชาย", MiddleNameTH: "", LastNameTH: "ใจดี",
			FirstNameEN: "Somchai", MiddleNameEN: "", LastNameEN: "Jaidee",
			DateOfBirth: "1990-01-15", PatientHN: "HN-000001",
			NationalID: "1234567890123", PassportID: "AA1234567",
			PhoneNumber: "0812345678", Email: "somchai.j@email.com", Gender: "M",
		},
		{
			ID: 2, HospitalID: 1,
			FirstNameTH: "สมหญิง", MiddleNameTH: "", LastNameTH: "รักสุข",
			FirstNameEN: "Somying", MiddleNameEN: "", LastNameEN: "Raksuk",
			DateOfBirth: "1985-06-20", PatientHN: "HN-000002",
			NationalID: "9876543210987", PassportID: "BB9876543",
			PhoneNumber: "0898765432", Email: "somying.r@email.com", Gender: "F",
		},
		{
			ID: 3, HospitalID: 1,
			FirstNameTH: "จอห์น", MiddleNameTH: "เจมส์", LastNameTH: "สมิธ",
			FirstNameEN: "John", MiddleNameEN: "James", LastNameEN: "Smith",
			DateOfBirth: "1978-11-30", PatientHN: "HN-000003",
			NationalID: "", PassportID: "CC5551234",
			PhoneNumber: "0856789012", Email: "john.smith@email.com", Gender: "M",
		},
		// Hospital B (ID = 2)
		{
			ID: 4, HospitalID: 2,
			FirstNameTH: "มานะ", MiddleNameTH: "", LastNameTH: "ทำดี",
			FirstNameEN: "Mana", MiddleNameEN: "", LastNameEN: "Thamdee",
			DateOfBirth: "1992-03-10", PatientHN: "HN-B00001",
			NationalID: "1111111111111", PassportID: "DD1111111",
			PhoneNumber: "0811111111", Email: "mana.t@email.com", Gender: "M",
		},
		{
			ID: 5, HospitalID: 2,
			FirstNameTH: "วิไล", MiddleNameTH: "", LastNameTH: "สุขใจ",
			FirstNameEN: "Wilai", MiddleNameEN: "", LastNameEN: "Sukjai",
			DateOfBirth: "1988-12-05", PatientHN: "HN-B00002",
			NationalID: "2222222222222", PassportID: "EE2222222",
			PhoneNumber: "0822222222", Email: "wilai.s@email.com", Gender: "F",
		},
	}
}

// ============================================================
// Mock Staff Repository
// ============================================================

// MockStaffRepository provides mock staff data for testing
type MockStaffRepository struct {
	staffs []models.Staff
	nextID uint
}

// NewMockStaffRepository creates a new MockStaffRepository
func NewMockStaffRepository() *MockStaffRepository {
	return &MockStaffRepository{
		staffs: []models.Staff{},
		nextID: 1,
	}
}

// Create creates a new staff member
func (r *MockStaffRepository) Create(staff *models.Staff) error {
	// Check uniqueness (username + hospital_id)
	for _, s := range r.staffs {
		if s.Username == staff.Username && s.HospitalID == staff.HospitalID {
			return fmt.Errorf("UNIQUE constraint failed: staff.username, staff.hospital_id")
		}
	}
	staff.ID = r.nextID
	r.nextID++
	r.staffs = append(r.staffs, *staff)
	return nil
}

// FindByUsernameAndHospital finds a staff member by username and hospital ID
func (r *MockStaffRepository) FindByUsernameAndHospital(username string, hospitalID uint) (*models.Staff, error) {
	for _, s := range r.staffs {
		if s.Username == username && s.HospitalID == hospitalID {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("staff not found")
}

// ============================================================
// Mock Hospital Repository
// ============================================================

// MockHospitalRepository provides mock hospital data for testing
type MockHospitalRepository struct {
	hospitals []models.Hospital
}

// NewMockHospitalRepository creates a new MockHospitalRepository
func NewMockHospitalRepository() *MockHospitalRepository {
	return &MockHospitalRepository{
		hospitals: []models.Hospital{
			{ID: 1, Name: "โรงพยาบาล A", Code: "HOSP_A"},
			{ID: 2, Name: "โรงพยาบาล B", Code: "HOSP_B"},
		},
	}
}

// FindByCode finds a hospital by its unique code
func (r *MockHospitalRepository) FindByCode(code string) (*models.Hospital, error) {
	for _, h := range r.hospitals {
		if h.Code == code {
			return &h, nil
		}
	}
	return nil, fmt.Errorf("hospital not found with code: %s", code)
}
