package models

import (
	"time"
)

type Account struct {
	Model
	User          User      `json:"-"`
	UserID        uint      `json:"userId"`
	BankName      string    `json:"bankName" gorm:"not null"`
	BankShort     string    `json:"bankShort" gorm:"not null"`
	BankUser      string    `json:"bankUser" gorm:"not null"`
	AccountNumber string    `json:"accountNumber" gorm:"not null"`
	Type          string    `json:"type" gorm:"default:Ahorros"`
	Currency      string    `json:"currency" gorm:"default:Soles"`
	Nickname      string    `json:"nickname"`
	Date          time.Time `json:"date" gorm:"default:current_timestamp"`
}
