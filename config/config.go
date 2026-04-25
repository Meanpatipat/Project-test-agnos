package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	// Server
	Port string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret string

	// External APIs
	HospitalABaseURL string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "hospital_admin"),
		DBPassword:       getEnv("DB_PASSWORD", "hospital_secret_2026"),
		DBName:           getEnv("DB_NAME", "hospital_middleware"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		JWTSecret:        getEnv("JWT_SECRET", "hospital-middleware-jwt-secret-key-2026"),
		HospitalABaseURL: getEnv("HOSPITAL_A_BASE_URL", "https://hospital-a.api.co.th"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
