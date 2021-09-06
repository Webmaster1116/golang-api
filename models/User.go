package models

import (
	"time"
)

type User struct {
	Model
	Name                 string    `json:"name" gorm:"not null"`
	Email                string    `json:"email" gorm:"unique;not null"`
	Password             string    `json:"password" gorm:"not null"`
	Verified             bool      `json:"verified" gorm:"default:false"`
	Date                 time.Time `json:"date" gorm:"default:current_timestamp"`
	ResetPasswordToken   string    `json:"resetPasswordToken"`
	ResetPasswordExpires time.Time `json:"resetPasswordExpires"`
}
