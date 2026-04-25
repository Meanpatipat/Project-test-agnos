package config_test

import (
	"os"
	"testing"

	"hospital-middleware/config"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Config – Defaults
// ============================================================

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any environment variables that might interfere
	envVars := []string{"PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "JWT_SECRET", "HOSPITAL_A_BASE_URL"}
	savedValues := make(map[string]string)

	for _, key := range envVars {
		if val, exists := os.LookupEnv(key); exists {
			savedValues[key] = val
			os.Unsetenv(key)
		}
	}
	defer func() {
		for key, val := range savedValues {
			os.Setenv(key, val)
		}
	}()

	cfg := config.LoadConfig()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "hospital_admin", cfg.DBUser)
	assert.Equal(t, "hospital_secret_2026", cfg.DBPassword)
	assert.Equal(t, "hospital_middleware", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, "hospital-middleware-jwt-secret-key-2026", cfg.JWTSecret)
	assert.Equal(t, "https://hospital-a.api.co.th", cfg.HospitalABaseURL)
}

// ============================================================
// Config – Environment Variable Override
// ============================================================

func TestLoadConfig_EnvOverrides(t *testing.T) {
	// Set custom env vars
	os.Setenv("PORT", "9090")
	os.Setenv("DB_HOST", "custom-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "custom_user")
	os.Setenv("DB_PASSWORD", "custom_pass")
	os.Setenv("DB_NAME", "custom_db")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("JWT_SECRET", "custom-jwt-secret")
	os.Setenv("HOSPITAL_A_BASE_URL", "https://custom.api.com")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("HOSPITAL_A_BASE_URL")
	}()

	cfg := config.LoadConfig()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "custom-host", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "custom_user", cfg.DBUser)
	assert.Equal(t, "custom_pass", cfg.DBPassword)
	assert.Equal(t, "custom_db", cfg.DBName)
	assert.Equal(t, "require", cfg.DBSSLMode)
	assert.Equal(t, "custom-jwt-secret", cfg.JWTSecret)
	assert.Equal(t, "https://custom.api.com", cfg.HospitalABaseURL)
}

func TestLoadConfig_PartialEnvOverride(t *testing.T) {
	// Only override PORT, rest should use defaults
	os.Setenv("PORT", "3000")
	defer os.Unsetenv("PORT")

	// Clear other vars to ensure defaults
	envVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "JWT_SECRET", "HOSPITAL_A_BASE_URL"}
	savedValues := make(map[string]string)
	for _, key := range envVars {
		if val, exists := os.LookupEnv(key); exists {
			savedValues[key] = val
			os.Unsetenv(key)
		}
	}
	defer func() {
		for key, val := range savedValues {
			os.Setenv(key, val)
		}
	}()

	cfg := config.LoadConfig()

	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "localhost", cfg.DBHost, "DBHost should use default")
	assert.Equal(t, "5432", cfg.DBPort, "DBPort should use default")
}

func TestLoadConfig_EmptyEnvVar(t *testing.T) {
	// An empty string is still a valid value from os.Getenv perspective
	os.Setenv("PORT", "")
	defer os.Unsetenv("PORT")

	cfg := config.LoadConfig()
	assert.Equal(t, "", cfg.Port, "empty string env var should override default")
}

func TestLoadConfig_ReturnsNewInstance(t *testing.T) {
	cfg1 := config.LoadConfig()
	cfg2 := config.LoadConfig()

	assert.NotSame(t, cfg1, cfg2, "each call should return a new Config instance")
}
