package models

import (
	"time"
)

type Keys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

type UserNotify struct {
	Model
	// User
	Endpoint   string `json:"endpoint" gorm:"not null,unique"`
	Keys       `json:"keys"`
	CreateDate time.Time `json:"createDate" gorm:"default:current_timestamp"`
}
