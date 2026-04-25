package models

import "time"

// Hospital represents a registered hospital in the system
type Hospital struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Name      string    `json:"name"       gorm:"type:varchar(255);not null"`
	Code      string    `json:"code"       gorm:"type:varchar(50);not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
