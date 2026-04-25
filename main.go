package main

import (
	"fmt"
	"log"
	"os"

	"hospital-middleware/config"
	"hospital-middleware/database"
	"hospital-middleware/repository"
	"hospital-middleware/router"
)

func main() {
	cfg := config.LoadConfig()

	// Determine mode: if DB_HOST is set, use PostgreSQL; otherwise use mock
	var (
		patientRepo  repository.PatientRepository
		staffRepo    repository.StaffRepository
		hospitalRepo repository.HospitalRepository
	)

	if os.Getenv("DB_HOST") != "" || cfg.DBHost != "localhost" {
		// Production / Docker mode — connect to PostgreSQL
		db := database.Connect(cfg)
		patientRepo = repository.NewPostgresPatientRepository(db)
		staffRepo = repository.NewPostgresStaffRepository(db)
		hospitalRepo = repository.NewPostgresHospitalRepository(db)
		log.Println("Running with PostgreSQL database")
	} else {
		// Development mode — use mock data
		patientRepo = repository.NewMockPatientRepository()
		staffRepo = repository.NewMockStaffRepository()
		hospitalRepo = repository.NewMockHospitalRepository()
		log.Println("Running with mock repositories")
	}

	r := router.SetupRouter(patientRepo, staffRepo, hospitalRepo, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Hospital Middleware API starting on %s", addr)
	log.Printf("Endpoints:")
	log.Printf("  POST /staff/create")
	log.Printf("  POST /staff/login")
	log.Printf("  GET  /patient/search  (requires JWT)")
	log.Printf("  GET  /health")

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
