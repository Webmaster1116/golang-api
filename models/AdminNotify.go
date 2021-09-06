package models

import (
	"time"
)

type AdminNotify struct {
	Model
	// Admin
	Endpoint   string `json:"endpoint" gorm:"not null,unique"`
	Keys       `json:"keys"`
	CreateDate time.Time `json:"createDate" gorm:"default:current_timestamp"`
}
