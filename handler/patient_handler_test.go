package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"hospital-middleware/handler"
	"hospital-middleware/models"
	"hospital-middleware/repository"

	"github.com/stretchr/testify/assert"
)

var testCounter uint64

// helper: create staff + login, return JWT token
func getAuthToken(t *testing.T, tsURL, hospital string) string {
	t.Helper()
	n := atomic.AddUint64(&testCounter, 1)
	username := fmt.Sprintf("testuser_%s_%d", hospital, n)

	postJSON(tsURL+"/staff/create", map[string]string{
		"username": username,
		"password": "testpass123",
		"hospital": hospital,
	})

	resp, _ := postJSON(tsURL+"/staff/login", map[string]string{
		"username": username,
		"password": "testpass123",
		"hospital": hospital,
	})

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var lr models.StaffLoginResponse
	json.Unmarshal(data, &lr)
	return lr.Token
}

func getWithAuth(url, token string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

// ============================================================
// Patient Search – Positive Tests
// ============================================================

func TestSearchPatient_NoFilters_ReturnsHospitalPatients(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search", token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	// Hospital A has 3 patients in mock
	assert.Equal(t, 3, len(patients))
}

func TestSearchPatient_ByNationalID(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search?national_id=1234567890123", token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "Somchai", patients[0].FirstNameEN)
}

func TestSearchPatient_ByPassportID(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search?passport_id=CC5551234", token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "John", patients[0].FirstNameEN)
}

func TestSearchPatient_ByFirstName(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search?first_name=somchai", token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	assert.Equal(t, 1, len(patients))
}

func TestSearchPatient_ByEmail(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search?email=somying.r@email.com", token)
	assert.NoError(t, err)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "F", patients[0].Gender)
}

func TestSearchPatient_ByDateOfBirth(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search?date_of_birth=1990-01-15", token)
	assert.NoError(t, err)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	assert.Equal(t, 1, len(patients))
	assert.Equal(t, "Somchai", patients[0].FirstNameEN)
}

// ============================================================
// Patient Search – Negative / Security Tests
// ============================================================

func TestSearchPatient_NoAuthToken(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/patient/search")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSearchPatient_InvalidToken(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := getWithAuth(ts.URL+"/patient/search", "invalid.token.here")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSearchPatient_HospitalIsolation(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	// Login as Hospital B staff
	tokenB := getAuthToken(t, ts.URL, "HOSP_B")

	// Search for a Hospital A patient by national_id — should return 0
	resp, err := getWithAuth(ts.URL+"/patient/search?national_id=1234567890123", tokenB)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	// Hospital B staff should NOT see Hospital A patients
	assert.Equal(t, 0, len(patients))
}

func TestSearchPatient_HospitalB_SeesOwnPatients(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	tokenB := getAuthToken(t, ts.URL, "HOSP_B")
	resp, err := getWithAuth(ts.URL+"/patient/search", tokenB)
	assert.NoError(t, err)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	// Hospital B has 2 patients in mock
	assert.Equal(t, 2, len(patients))
}

func TestSearchPatient_NoMatch(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search?national_id=0000000000000", token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&result)
	data, _ := json.Marshal(result.Data)
	var patients []models.Patient
	json.Unmarshal(data, &patients)

	assert.Equal(t, 0, len(patients))
}

// ============================================================
// Health Check
// ============================================================

func TestHealthCheck(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSearchPatient_MissingContextVariables(t *testing.T) {
	// Setup just the handler to skip auth middleware
	repo := repository.NewMockPatientRepository()
	h := handler.NewPatientHandler(repo)
	
	r := gin.Default()
	r.GET("/patient/search/nomid", h.SearchPatient)
	
	ts := httptest.NewServer(r)
	defer ts.Close()
	
	resp, err := http.Get(ts.URL + "/patient/search/nomid")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestSearchPatient_WrongTypeContextVariable(t *testing.T) {
	repo := repository.NewMockPatientRepository()
	h := handler.NewPatientHandler(repo)
	
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("hospital_id", "not_a_uint")
		c.Next()
	})
	r.GET("/patient/search/wrongtype", h.SearchPatient)
	
	ts := httptest.NewServer(r)
	defer ts.Close()
	
	resp, err := http.Get(ts.URL + "/patient/search/wrongtype")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestSearchPatient_RepoError(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	token := getAuthToken(t, ts.URL, "HOSP_A")
	resp, err := getWithAuth(ts.URL+"/patient/search?national_id=ERROR", token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}




