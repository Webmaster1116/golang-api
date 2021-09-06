package models

import (
	"time"
)

type Invoice struct {
	Model
	AgentID           uint      `json:"agentId"`
	Name              string    `json:"name" gorm:"not null"`
	CommissionSoles   float64   `json:"commissionSoles" gorm:"default:0"`
	CommissionDollars float64   `json:"commissionDollars" gorm:"default:0"`
	StartDay          time.Time `json:"startDay" gorm:"default:current_timestamp"`
	EndDay            time.Time `json:"endDay" gorm:"default:current_timestamp"`
	Date              time.Time `json:"date" gorm:"default:current_timestamp"`
}
