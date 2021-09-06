package models

import (
	"time"
)

type Transaction struct {
	AmountToPay   float64 `json:"amountToPay"`
	AmountReceive float64 `json:"amountReceive"`
	CurrencyTo    string  `json:"currencyTo"`
	CurrencyFrom  string  `json:"currencyFrom"`
	Savings       float64 `json:"savings"`
	Exchange      float64 `json:"exchange"`
	Status        string  `json:"status"`
}

type File struct {
	Model
	Content []byte
}

type Operation struct {
	Model
	UserID                  uint         `json:"userId"`
	User                    User         `json:"-"`
	AgentID                 uint         `json:"agentId"`
	Agent                   Agent        `json:"-"`
	AgentAccountID          uint         `json:"agentAccountId"`
	ProfileID               uint         `json:"profileId"`
	AccountID               uint         `json:"accountId"`
	Profile                 Profile      `json:"profileDetails" gorm:"not null"`
	Account                 Account      `json:"bankDetails" gorm:"not null"`
	DestinationBank         string       `json:"destinationBank" gorm:"not null"`
	Transaction             Transaction  `json:"transaction" gorm:"not null;embedded;embeddedPrefix:tx_"`
	TransactionNumber       string       `json:"transactionNumber"`
	AgentTransactionNumber  string       `json:"agentTransactionNumber"`
	AgentName               string       `json:"agentName"`
	AgentEmail              string       `json:"agentEmail"`
	AgentAccount            AgentAccount `json:"agentBank"`
	UserTransactionPhotoId  uint         `json:"userTransactionPhotoId"`
	UserTransactionPhoto    File         `json:"-"`
	AgentTransactionPhotoId uint         `json:"agentTransactionPhotoId"`
	AgentTransactionPhoto   File         `json:"-"`
	Savings                 float64      `json:"savings"`
	Exchange                float64      `json:"exchange"`
	Date                    time.Time    `json:"date" gorm:"default:current_timestamp"`
}
