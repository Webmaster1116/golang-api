package models

import (
	"time"
)

type AgentProfile struct {
	Model
	Agent       Agent     `json:"-"`
	AgentID     uint      `json:"agentId"`
	FirstName   string    `json:"firstName" gorm:"not null"`
	LastName    string    `json:"lastName" gorm:"not null"`
	Address     string    `json:"address" gorm:"not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Phone       string    `json:"phone" gorm:"not null"`
	Dni         string    `json:"dni" gorm:"not null"`
	CompanyName string    `json:"companyName"`
	Ruc         string    `json:"ruc"`
	AccountType string    `json:"accountType" gorm:"default:Personal"`
	Date        time.Time `json:"date" gorm:"default:current_timestamp"`
}
