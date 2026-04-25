package middleware

import (
	"net/http"
	"strings"
	"time"

	"hospital-middleware/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the custom claims stored in JWT tokens
type JWTClaims struct {
	StaffID    uint   `json:"staff_id"`
	Username   string `json:"username"`
	HospitalID uint   `json:"hospital_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a staff member
func GenerateToken(staff *models.Staff, secret string) (string, error) {
	claims := JWTClaims{
		StaffID:    staff.ID,
		Username:   staff.Username,
		HospitalID: staff.HospitalID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hospital-middleware",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// AuthMiddleware validates JWT tokens and sets staff context
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Authorization header is required.",
			})
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid authorization format. Use: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid or expired token.",
			})
			return
		}

		// Store staff info in context for use by handlers
		c.Set("staff_id", claims.StaffID)
		c.Set("username", claims.Username)
		c.Set("hospital_id", claims.HospitalID)

		c.Next()
	}
}
