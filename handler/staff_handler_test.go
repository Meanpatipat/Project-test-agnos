package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hospital-middleware/models"
	"hospital-middleware/repository"
	"hospital-middleware/router"

	"github.com/stretchr/testify/assert"
)

const testJWTSecret = "test-secret-key"

func setupTestServer() *httptest.Server {
	patientRepo := repository.NewMockPatientRepository()
	staffRepo := repository.NewMockStaffRepository()
	hospitalRepo := repository.NewMockHospitalRepository()
	r := router.SetupRouter(patientRepo, staffRepo, hospitalRepo, testJWTSecret)
	return httptest.NewServer(r)
}

func postJSON(url string, body interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(body)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}

// ============================================================
// Staff Create Tests
// ============================================================

func TestCreateStaff_Success(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "nurse_a",
		"password": "secret123",
		"hospital": "HOSP_A",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Staff member created successfully.", result.Message)
}

func TestCreateStaff_DuplicateUsername(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	body := map[string]string{
		"username": "nurse_dup",
		"password": "secret123",
		"hospital": "HOSP_A",
	}
	postJSON(ts.URL+"/staff/create", body)

	// Second creation with same username+hospital should fail
	resp, err := postJSON(ts.URL+"/staff/create", body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestCreateStaff_InvalidHospital(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "nurse_x",
		"password": "secret123",
		"hospital": "INVALID_HOSP",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateStaff_MissingFields(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "nurse_y",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateStaff_ShortPassword(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "nurse_z",
		"password": "123",
		"hospital": "HOSP_A",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateStaff_ShortUsername(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "ab",
		"password": "secret123",
		"hospital": "HOSP_A",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateStaff_PasswordTooLong(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// bcrypt fails if password is > 72 bytes
	longPass := ""
	for i := 0; i < 80; i++ {
		longPass += "a"
	}

	resp, err := postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "nurse_longpass",
		"password": longPass,
		"hospital": "HOSP_A",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}


// ============================================================
// Staff Login Tests
// ============================================================

func TestLoginStaff_Success(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// Create staff first
	postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "doc_login",
		"password": "mypassword",
		"hospital": "HOSP_A",
	})

	// Login
	resp, err := postJSON(ts.URL+"/staff/login", map[string]string{
		"username": "doc_login",
		"password": "mypassword",
		"hospital": "HOSP_A",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Login successful.", result.Message)

	// Verify token is present
	data, _ := json.Marshal(result.Data)
	var loginResp models.StaffLoginResponse
	json.Unmarshal(data, &loginResp)
	assert.NotEmpty(t, loginResp.Token)
	assert.Equal(t, "doc_login", loginResp.Username)
	assert.Equal(t, "HOSP_A", loginResp.Hospital)
}

func TestLoginStaff_WrongPassword(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "doc_wp",
		"password": "correct_pass",
		"hospital": "HOSP_A",
	})

	resp, err := postJSON(ts.URL+"/staff/login", map[string]string{
		"username": "doc_wp",
		"password": "wrong_pass",
		"hospital": "HOSP_A",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestLoginStaff_NonExistentUser(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := postJSON(ts.URL+"/staff/login", map[string]string{
		"username": "ghost_user",
		"password": "password",
		"hospital": "HOSP_A",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestLoginStaff_WrongHospital(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	postJSON(ts.URL+"/staff/create", map[string]string{
		"username": "doc_wh",
		"password": "password",
		"hospital": "HOSP_A",
	})

	resp, err := postJSON(ts.URL+"/staff/login", map[string]string{
		"username": "doc_wh",
		"password": "password",
		"hospital": "HOSP_B",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestLoginStaff_MissingFields(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := postJSON(ts.URL+"/staff/login", map[string]string{
		"username": "doc_mf",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
