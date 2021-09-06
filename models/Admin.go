package models

import (
	"time"
)

type Admin struct {
	Model
	Name     string    `json:"name" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null;unique"`
	Password string    `json:"password" gorm:"not null"`
	Type     string    `json:"type" gorm:"default:admin"`
	Date     time.Time `json:"date" gorm:"default:current_timestamp"`
}
