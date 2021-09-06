package routes

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// var exchangeId = "5f0b10586057690017293870"
// var url = "https://cuantoestaeldolar.pe/cambio-de-dolar-online"

// @route    GET api/exchange
// @desc     Get exchange rate
// @access   Private
func getExchange(c *gin.Context) {
	var exchange models.Exchange

	if err := models.Database.
		Order("created_at ASC").
		First(&exchange).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, "No exchange in server")
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          exchange.ID,
		"current":     exchange.Current,
		"PEN":         exchange.Current.Compra,
		"USD":         exchange.Current.Venta,
		"sunat":       exchange.Sunat,
		"paralelo":    exchange.Paralelo,
		"dollarHouse": exchange.DollarHouse,
		"cambix":      exchange.Cambix,
		"acomo":       exchange.Acomo,
		"bcp":         exchange.Bcp,
	})
}

// @route    GET api/exchange/calculate
// @desc     Get exchange rate
// @access   Public
func getExchangeCalculate(c *gin.Context) {
	var payload struct {
		OriginCurrency      string  `form:"originCurrency" binding:"required"`
		DestinationCurrency string  `form:"destinationCurrency" binding:"required"`
		Amount              float64 `form:"amount" binding:"required"`
	}

	if err := c.BindQuery(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var exchange models.Exchange

	if err := models.Database.Order("created_at ASC").First(&exchange).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, "No exchange in server")
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Server error")
		}
		return
	}

	rate := exchange.Current.Venta
	ex := (1 / rate) * payload.Amount
	savings := math.Abs(exchange.Bcp.Venta-exchange.Current.Venta) * payload.Amount
	if payload.OriginCurrency == "USD" {
		rate = exchange.Current.Compra
		ex = rate * payload.Amount
		savings = math.Abs(exchange.Bcp.Compra-exchange.Current.Compra) * payload.Amount
	}
	c.JSON(http.StatusOK, gin.H{
		"rate":     rate,
		"exchange": ex,
		"ahorros":  savings,
	})
}

// @route    PUT api/exchange/:id
// @desc     Update current value
// @access   Public
func putExchange(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Compra float64
		Venta  float64
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var exchange models.Exchange
	if err := models.Database.First(&exchange, id).Error; err != nil {
		c.JSON(http.StatusNotFound, "Database exchange not found")
		return
	}
	logrus.Infof("updating exchange(%v) : %v\n", id, payload)

	if payload.Compra != 0 {
		exchange.Current.Compra = payload.Compra
	}
	if payload.Venta != 0 {
		exchange.Current.Venta = payload.Venta
	}

	if err := models.Database.Save(&exchange).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	c.JSON(http.StatusOK, &exchange)
}

func exchange(r *gin.RouterGroup) {
	r.
		GET("", middleware.Auth, getExchange).
		GET("/calculate", getExchangeCalculate).
		PUT("/:id", putExchange)
}
