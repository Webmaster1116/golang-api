package models

import (
	"time"
)

type ExchangeElement struct {
	Compra float64 `json:"compra"`
	Venta  float64 `json:"venta"`
}

type Exchange struct {
	Model
	Current     ExchangeElement `json:"current" gorm:"embedded;embeddedPrefix:current_"`
	Sunat       ExchangeElement `json:"sunat" gorm:"embedded;embeddedPrefix:sunat_"`
	Paralelo    ExchangeElement `json:"paralelo" gorm:"embedded;embeddedPrefix:paralelo_"`
	DollarHouse ExchangeElement `json:"dollarHouse" gorm:"embedded;embeddedPrefix:dollar_house_"`
	Cambix      ExchangeElement `json:"cambix" gorm:"embedded;embeddedPrefix:cambix_"`
	Acomo       ExchangeElement `json:"acomo" gorm:"embedded;embeddedPrefix:acomo_"`
	Bcp         ExchangeElement `json:"bcp" gorm:"embedded;embeddedPrefix:bcp_"`
	// CurrentVenta float64         `json:"CurrentVenta"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp"`
}
