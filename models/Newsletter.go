package models

import (
	"time"
)

type Newsletter struct {
	Model
	Email string    `json:"email" gorm:"not null"`
	Date  time.Time `json:"date" gorm:"default:current_timestamp"`
}
