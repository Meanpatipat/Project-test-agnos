package repository_test

import (
	"testing"

	"hospital-middleware/models"
	"hospital-middleware/repository"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// MockPatientRepository Tests
// ============================================================

func TestMockPatientRepo_FindByID_ByNationalID(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	patient, err := repo.FindByID("1234567890123", 1) // Hospital A
	assert.NoError(t, err)
	assert.NotNil(t, patient)
	assert.Equal(t, "Somchai", patient.FirstNameEN)
}

func TestMockPatientRepo_FindByID_ByPassportID(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	patient, err := repo.FindByID("CC5551234", 1) // Hospital A
	assert.NoError(t, err)
	assert.NotNil(t, patient)
	assert.Equal(t, "John", patient.FirstNameEN)
}

func TestMockPatientRepo_FindByID_NotFound(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	patient, err := repo.FindByID("9999999999999", 1)
	assert.Error(t, err)
	assert.Nil(t, patient)
	assert.Contains(t, err.Error(), "patient not found")
}

func TestMockPatientRepo_FindByID_WrongHospital(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	// national_id 1234567890123 belongs to Hospital A (ID=1), not Hospital B (ID=2)
	patient, err := repo.FindByID("1234567890123", 2)
	assert.Error(t, err)
	assert.Nil(t, patient)
}

func TestMockPatientRepo_Search_NoFilters(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{}

	patients, err := repo.Search(req, 1) // Hospital A
	assert.NoError(t, err)
	assert.Equal(t, 3, len(patients), "Hospital A should have 3 patients")
}

func TestMockPatientRepo_Search_NoFilters_HospitalB(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{}

	patients, err := repo.Search(req, 2) // Hospital B
	assert.NoError(t, err)
	assert.Equal(t, 2, len(patients), "Hospital B should have 2 patients")
}

func TestMockPatientRepo_Search_NoFilters_NonExistentHospital(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{}

	patients, err := repo.Search(req, 999)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(patients))
}

func TestMockPatientRepo_Search_ByNationalID(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{NationalID: "9876543210987"}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "Somying", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByPassportID(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{PassportID: "DD1111111"}

	patients, err := repo.Search(req, 2) // Hospital B
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "Mana", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByFirstName_CaseInsensitive(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{FirstName: "SOMCHAI"} // uppercase

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "Somchai", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByFirstName_PartialMatch(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{FirstName: "som"} // partial

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	// Should match "Somchai" and "Somying"
	assert.Equal(t, 2, len(patients))
}

func TestMockPatientRepo_Search_ByMiddleName(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{MiddleName: "james"}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "John", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByLastName(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{LastName: "smith"}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "John", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByDateOfBirth(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{DateOfBirth: "1985-06-20"}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "Somying", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByPhoneNumber(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{PhoneNumber: "0856789012"}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "John", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByEmail(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{Email: "somchai.j@email.com"}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "Somchai", patients[0].FirstNameEN)
}

func TestMockPatientRepo_Search_ByEmail_CaseInsensitive(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{Email: "SOMCHAI.J@EMAIL.COM"}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
}

func TestMockPatientRepo_Search_MultipleFilters(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{
		FirstName:   "Somchai",
		DateOfBirth: "1990-01-15",
	}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "1234567890123", patients[0].NationalID)
}

func TestMockPatientRepo_Search_MultipleFilters_NoMatch(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	req := models.PatientSearchRequest{
		FirstName:   "Somchai",
		DateOfBirth: "2000-01-01", // wrong DOB for Somchai
	}

	patients, err := repo.Search(req, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(patients))
}

func TestMockPatientRepo_Search_CrossHospitalIsolation(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	// Search for Hospital A patient data using Hospital B scope
	req := models.PatientSearchRequest{NationalID: "1234567890123"}
	patients, err := repo.Search(req, 2)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(patients), "hospital B should not see hospital A patients")
}

// ============================================================
// MockStaffRepository Tests
// ============================================================

func TestMockStaffRepo_Create_Success(t *testing.T) {
	repo := repository.NewMockStaffRepository()
	staff := &models.Staff{
		Username:     "nurse_a",
		PasswordHash: "hashed_password",
		HospitalID:   1,
	}

	err := repo.Create(staff)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), staff.ID, "should auto-assign ID")
}

func TestMockStaffRepo_Create_AutoIncrementID(t *testing.T) {
	repo := repository.NewMockStaffRepository()
	staff1 := &models.Staff{Username: "user1", PasswordHash: "hash", HospitalID: 1}
	staff2 := &models.Staff{Username: "user2", PasswordHash: "hash", HospitalID: 1}

	repo.Create(staff1)
	repo.Create(staff2)

	assert.Equal(t, uint(1), staff1.ID)
	assert.Equal(t, uint(2), staff2.ID)
}

func TestMockStaffRepo_Create_DuplicateUsername_SameHospital(t *testing.T) {
	repo := repository.NewMockStaffRepository()
	staff1 := &models.Staff{Username: "nurse", PasswordHash: "hash1", HospitalID: 1}
	staff2 := &models.Staff{Username: "nurse", PasswordHash: "hash2", HospitalID: 1}

	err1 := repo.Create(staff1)
	err2 := repo.Create(staff2)

	assert.NoError(t, err1)
	assert.Error(t, err2, "duplicate username in same hospital should fail")
	assert.Contains(t, err2.Error(), "UNIQUE constraint failed")
}

func TestMockStaffRepo_Create_SameUsername_DifferentHospitals(t *testing.T) {
	repo := repository.NewMockStaffRepository()
	staff1 := &models.Staff{Username: "nurse", PasswordHash: "hash1", HospitalID: 1}
	staff2 := &models.Staff{Username: "nurse", PasswordHash: "hash2", HospitalID: 2}

	err1 := repo.Create(staff1)
	err2 := repo.Create(staff2)

	assert.NoError(t, err1)
	assert.NoError(t, err2, "same username in different hospitals should succeed")
}

func TestMockStaffRepo_FindByUsernameAndHospital_Success(t *testing.T) {
	repo := repository.NewMockStaffRepository()
	repo.Create(&models.Staff{Username: "doc", PasswordHash: "hash", HospitalID: 1})

	found, err := repo.FindByUsernameAndHospital("doc", 1)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "doc", found.Username)
}

func TestMockStaffRepo_FindByUsernameAndHospital_NotFound(t *testing.T) {
	repo := repository.NewMockStaffRepository()

	found, err := repo.FindByUsernameAndHospital("ghost", 1)
	assert.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "staff not found")
}

func TestMockStaffRepo_FindByUsernameAndHospital_WrongHospital(t *testing.T) {
	repo := repository.NewMockStaffRepository()
	repo.Create(&models.Staff{Username: "doc", PasswordHash: "hash", HospitalID: 1})

	found, err := repo.FindByUsernameAndHospital("doc", 2) // exists in hospital 1, not 2
	assert.Error(t, err)
	assert.Nil(t, found)
}

// ============================================================
// MockHospitalRepository Tests
// ============================================================

func TestMockHospitalRepo_FindByCode_Success(t *testing.T) {
	repo := repository.NewMockHospitalRepository()

	hospital, err := repo.FindByCode("HOSP_A")
	assert.NoError(t, err)
	assert.NotNil(t, hospital)
	assert.Equal(t, uint(1), hospital.ID)
	assert.Equal(t, "HOSP_A", hospital.Code)
}

func TestMockHospitalRepo_FindByCode_HospitalB(t *testing.T) {
	repo := repository.NewMockHospitalRepository()

	hospital, err := repo.FindByCode("HOSP_B")
	assert.NoError(t, err)
	assert.NotNil(t, hospital)
	assert.Equal(t, uint(2), hospital.ID)
}

func TestMockHospitalRepo_FindByCode_NotFound(t *testing.T) {
	repo := repository.NewMockHospitalRepository()

	hospital, err := repo.FindByCode("NON_EXISTENT")
	assert.Error(t, err)
	assert.Nil(t, hospital)
	assert.Contains(t, err.Error(), "hospital not found")
}

func TestMockHospitalRepo_FindByCode_EmptyCode(t *testing.T) {
	repo := repository.NewMockHospitalRepository()

	hospital, err := repo.FindByCode("")
	assert.Error(t, err)
	assert.Nil(t, hospital)
}

func TestMockHospitalRepo_FindByCode_CaseSensitive(t *testing.T) {
	repo := repository.NewMockHospitalRepository()

	// Codes are stored as "HOSP_A", searching with lowercase should not match
	hospital, err := repo.FindByCode("hosp_a")
	assert.Error(t, err)
	assert.Nil(t, hospital, "hospital code lookup should be case-sensitive")
}
