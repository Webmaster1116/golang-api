package models

type Token struct {
	Model
	User   User   `json:"-"`
	UserID uint   `json:"userId" gorm:"not null"`
	Token  string `json:"token" gorm:"not null"`
}
