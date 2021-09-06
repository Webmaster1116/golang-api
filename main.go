package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/config"
	"github.com/med8bra/moni-api-go/models"
	"github.com/med8bra/moni-api-go/routes"
	"github.com/med8bra/moni-api-go/services"
	"github.com/med8bra/moni-api-go/util/password"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func main() {
	// -- load config
	if err := config.Init(); err != nil {
		logrus.Fatal("failed to load config: ", err.Error())
	}
	// -- setup config
	logrus.SetLevel(config.Config.Level)
	// -- init services
	if err := services.Init(&config.Config); err != nil {
		logrus.Fatal("failed to load services: ", err.Error())
	}
	// -- open database
	if err := models.Init(&config.Config.DB); err != nil {
		logrus.Fatal("failed to open database: ", err.Error())
	}
	// -- load templates
	if err := services.TemplateManager.LoadDir("templates"); err != nil {
		logrus.Fatal("failed to load templates: ", err.Error())
	}
	// -- create default admin if not exists
	adminPasswordHash, _ := password.Hash("adminPassword")
	models.Database.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&models.Admin{Name: "admin", Email: "admin@moni.pe", Password: adminPasswordHash})
	// -- insert exchange row if not exists
	models.Database.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&models.Exchange{
			Model:       models.Model{ID: 1},
			Current:     models.ExchangeElement{Compra: 3.48, Venta: 3.52},
			Sunat:       models.ExchangeElement{Compra: 4.06, Venta: 4.069},
			Paralelo:    models.ExchangeElement{Compra: 4.07, Venta: 3.52},
			DollarHouse: models.ExchangeElement{Compra: 4.085, Venta: 4.105},
		})
	// -- create server
	app := gin.Default()
	// -- use CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "x-auth-token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if err := routes.Init(app.Group("/api")); err != nil {
		logrus.Fatal("failed to set up routes: ", err.Error())
	}
	// -- run server
	if err := app.Run(); err != nil {
		logrus.Fatal("server failed : ", err.Error())
	}
}
