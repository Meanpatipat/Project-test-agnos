package models

import "time"

// Patient represents a patient record scoped to a specific hospital
type Patient struct {
	ID           uint      `json:"-"              gorm:"primaryKey"`
	HospitalID   uint      `json:"-"              gorm:"not null;index"`
	FirstNameTH  string    `json:"first_name_th"  gorm:"type:varchar(255);default:''"`
	MiddleNameTH string    `json:"middle_name_th" gorm:"type:varchar(255);default:''"`
	LastNameTH   string    `json:"last_name_th"   gorm:"type:varchar(255);default:''"`
	FirstNameEN  string    `json:"first_name_en"  gorm:"type:varchar(255);default:''"`
	MiddleNameEN string    `json:"middle_name_en" gorm:"type:varchar(255);default:''"`
	LastNameEN   string    `json:"last_name_en"   gorm:"type:varchar(255);default:''"`
	DateOfBirth  string    `json:"date_of_birth"  gorm:"type:date"`
	PatientHN    string    `json:"patient_hn"     gorm:"type:varchar(50);not null"`
	NationalID   string    `json:"national_id"    gorm:"type:varchar(13);default:'';index"`
	PassportID   string    `json:"passport_id"    gorm:"type:varchar(20);default:'';index"`
	PhoneNumber  string    `json:"phone_number"   gorm:"type:varchar(20);default:'';index"`
	Email        string    `json:"email"          gorm:"type:varchar(255);default:'';index"`
	Gender       string    `json:"gender"         gorm:"type:char(1);check:gender IN ('M','F')"`
	CreatedAt    time.Time `json:"-"              gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"-"              gorm:"autoUpdateTime"`
}

// PatientSearchRequest holds optional search filters for patient queries
type PatientSearchRequest struct {
	NationalID  string `json:"national_id"    form:"national_id"`
	PassportID  string `json:"passport_id"    form:"passport_id"`
	FirstName   string `json:"first_name"     form:"first_name"`
	MiddleName  string `json:"middle_name"    form:"middle_name"`
	LastName    string `json:"last_name"      form:"last_name"`
	DateOfBirth string `json:"date_of_birth"  form:"date_of_birth"`
	PhoneNumber string `json:"phone_number"   form:"phone_number"`
	Email       string `json:"email"          form:"email"`
}
