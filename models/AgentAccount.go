package models

import (
	"time"
)

type AgentAccount struct {
	Model
	Agent     Agent     `json:"-"`
	AgentID   uint      `json:"agentId"`
	Nickname  string    `json:"nickname"`
	BankName  string    `json:"bankName" gorm:"not null"`
	BankShort string    `json:"bankShort" gorm:"not null"`
	BankUser  string    `json:"bankUser" gorm:"not null"`
	Type      string    `json:"type" gorm:"default:Ahorros"`
	Currency  string    `json:"currency" gorm:"default:Soles"`
	Purpose   string    `json:"purpose" gorm:"default:Receive"`
	Balance   bool      `json:"balance" gorm:"default:true"`
	Date      time.Time `json:"date" gorm:"default:current_timestamp"`
}
