package models

import (
	"time"
)

type Profile struct {
	Model
	User           User      `json:"-"`
	UserID         uint      `json:"userId"`
	FirstName      string    `json:"firstName" gorm:"not null"`
	SecondName     string    `json:"secondName" gorm:"not null"`
	LastName       string    `json:"lastName" gorm:"not null"`
	BirthDate      time.Time `json:"birthDate" gorm:"not null"`
	Email          string    `json:"email" gorm:"not null"`
	Phone          string    `json:"phone" gorm:"not null"`
	Nationality    string    `json:"nationality" gorm:"not null"`
	Doctype        string    `json:"doctype" gorm:"not null"`
	Dni            string    `json:"dni" gorm:"not null"`
	MotherLastName string    `json:"motherLastName"`
	AccountType    string    `json:"accountType" gorm:"default:Personal"`
	Date           time.Time `json:"date" gorm:"default:current_timestamp"`
}
