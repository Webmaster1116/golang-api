package models

import (
	"time"
)

type Agent struct {
	Model
	Name              string      `json:"name" gorm:"not null"`
	Email             string      `json:"email" gorm:"unique;not null"`
	Password          string      `json:"-" gorm:"not null"`
	Type              string      `json:"type" gorm:"default:agent"`
	Online            bool        `json:"online" gorm:"default:true"`
	Enable            bool        `json:"enable" gorm:"default:true"`
	CommissionSoles   int         `json:"commissionSoles" gorm:"default:0"`
	CommissionDollars int         `json:"commissionDollars" gorm:"default:0"`
	Operations        []Operation `json:"operations" gorm:"foreignKey:UserID"`
	Date              time.Time   `json:"date" gorm:"default:current_timestamp"`
}
