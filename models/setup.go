package models

import (
	"github.com/med8bra/moni-api-go/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Init(c *config.DB) (err error) {
	if Database, err = gorm.Open(mysql.Open(c.MYSQL_DSN), &gorm.Config{}); err != nil {
		return err
	}
	// migrate models
	return Database.AutoMigrate(
		&Account{},
		&Admin{},
		&AdminNotify{},
		&Agent{},
		&AgentAccount{},
		&AgentProfile{},
		&Exchange{},
		&Invoice{},
		&Newsletter{},
		&Operation{},
		&Profile{},
		&Token{},
		&User{},
		&UserNotify{},
	)
}
