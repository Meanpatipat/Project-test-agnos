package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hospital-middleware/middleware"
	"hospital-middleware/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-jwt-secret"

func init() {
	gin.SetMode(gin.TestMode)
}

// setupAuthRouter creates a Gin router with the auth middleware protecting a test endpoint.
func setupAuthRouter(secret string) *gin.Engine {
	r := gin.New()
	r.Use(middleware.AuthMiddleware(secret))
	r.GET("/protected", func(c *gin.Context) {
		staffID, _ := c.Get("staff_id")
		username, _ := c.Get("username")
		hospitalID, _ := c.Get("hospital_id")
		c.JSON(http.StatusOK, gin.H{
			"staff_id":    staffID,
			"username":    username,
			"hospital_id": hospitalID,
		})
	})
	return r
}

// ============================================================
// GenerateToken Tests
// ============================================================

func TestGenerateToken_Success(t *testing.T) {
	staff := &models.Staff{
		ID:         1,
		Username:   "nurse_a",
		HospitalID: 10,
	}

	token, err := middleware.GenerateToken(staff, testSecret)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse and validate the token
	claims := &middleware.JWTClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsed.Valid)
	assert.Equal(t, uint(1), claims.StaffID)
	assert.Equal(t, "nurse_a", claims.Username)
	assert.Equal(t, uint(10), claims.HospitalID)
	assert.Equal(t, "hospital-middleware", claims.Issuer)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, 5*time.Second)
}

func TestGenerateToken_DifferentStaff(t *testing.T) {
	staff1 := &models.Staff{ID: 1, Username: "user1", HospitalID: 1}
	staff2 := &models.Staff{ID: 2, Username: "user2", HospitalID: 2}

	token1, err1 := middleware.GenerateToken(staff1, testSecret)
	token2, err2 := middleware.GenerateToken(staff2, testSecret)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, token1, token2, "different staff should produce different tokens")
}

func TestGenerateToken_DifferentSecrets(t *testing.T) {
	staff := &models.Staff{ID: 1, Username: "user1", HospitalID: 1}

	token, err := middleware.GenerateToken(staff, "secret-A")
	assert.NoError(t, err)

	// Token should NOT validate with a different secret
	claims := &middleware.JWTClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret-B"), nil
	})
	assert.Error(t, err, "token should not validate with wrong secret")
}

// ============================================================
// AuthMiddleware Tests – Positive
// ============================================================

func TestAuthMiddleware_ValidToken(t *testing.T) {
	router := setupAuthRouter(testSecret)

	staff := &models.Staff{ID: 5, Username: "doc_test", HospitalID: 3}
	token, _ := middleware.GenerateToken(staff, testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, float64(5), body["staff_id"])
	assert.Equal(t, "doc_test", body["username"])
	assert.Equal(t, float64(3), body["hospital_id"])
}

func TestAuthMiddleware_ValidToken_CaseInsensitiveBearer(t *testing.T) {
	router := setupAuthRouter(testSecret)

	staff := &models.Staff{ID: 1, Username: "user", HospitalID: 1}
	token, _ := middleware.GenerateToken(staff, testSecret)

	// "bearer" in lowercase should also work
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============================================================
// AuthMiddleware Tests – Negative
// ============================================================

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errResp models.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Contains(t, errResp.Message, "Authorization header is required")
}

func TestAuthMiddleware_EmptyAuthHeader(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_MissingBearerPrefix(t *testing.T) {
	router := setupAuthRouter(testSecret)

	staff := &models.Staff{ID: 1, Username: "user", HospitalID: 1}
	token, _ := middleware.GenerateToken(staff, testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", token) // no "Bearer " prefix
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errResp models.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Contains(t, errResp.Message, "Invalid authorization format")
}

func TestAuthMiddleware_WrongPrefix(t *testing.T) {
	router := setupAuthRouter(testSecret)

	staff := &models.Staff{ID: 1, Username: "user", HospitalID: 1}
	token, _ := middleware.GenerateToken(staff, testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Basic "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_MalformedToken(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errResp models.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Contains(t, errResp.Message, "Invalid or expired token")
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	router := setupAuthRouter(testSecret)

	staff := &models.Staff{ID: 1, Username: "user", HospitalID: 1}
	token, _ := middleware.GenerateToken(staff, "different-secret")

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	router := setupAuthRouter(testSecret)

	// Manually create an already-expired token
	claims := middleware.JWTClaims{
		StaffID:    1,
		Username:   "expired_user",
		HospitalID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
			Issuer:    "hospital-middleware",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(testSecret))

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errResp models.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Contains(t, errResp.Message, "Invalid or expired token")
}

func TestAuthMiddleware_TokenWithDifferentSigningMethod(t *testing.T) {
	router := setupAuthRouter(testSecret)

	// Create a token with "none" algorithm (unsigned)
	claims := middleware.JWTClaims{
		StaffID:    1,
		Username:   "hacker",
		HospitalID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signedToken, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_BearerOnly_NoToken(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
