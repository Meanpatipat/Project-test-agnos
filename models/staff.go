package models

import "time"

// Staff represents a hospital staff member with login credentials
type Staff struct {
	ID           uint      `json:"id"            gorm:"primaryKey"`
	Username     string    `json:"username"       gorm:"type:varchar(100);not null"`
	PasswordHash string    `json:"-"              gorm:"type:varchar(255);not null"` // never expose in JSON
	HospitalID   uint      `json:"hospital_id"    gorm:"not null;index"`
	Hospital     Hospital  `json:"hospital"       gorm:"foreignKey:HospitalID"`
	CreatedAt    time.Time `json:"created_at"     gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at"     gorm:"autoUpdateTime"`
}

// StaffCreateRequest is the input for creating a new staff member
type StaffCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	Hospital string `json:"hospital" binding:"required"` // hospital code (e.g. "HOSP_A")
}

// StaffLoginRequest is the input for staff login
type StaffLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"` // hospital code
}

// StaffLoginResponse is returned on successful login
type StaffLoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Hospital string `json:"hospital"`
}
