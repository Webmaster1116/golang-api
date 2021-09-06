package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/pkg/errors"
)

// @route    GET api/accounts
// @desc     Get all accounts
// @access   Private
func getAccounts(c *gin.Context) {
	var accounts []models.Account
	user := c.MustGet("user").(*models.AuthUser)
	if err := models.Database.Where(&models.Account{UserID: user.ID}).
		Order("date DESC").
		Find(&accounts).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to query accounts"))
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, accounts)
}

// @route    POST api/accounts
// @desc     Save an account Data
// @access   Private
func postAccount(c *gin.Context) {
	user := c.MustGet("user").(*models.AuthUser)
	var account models.Account
	if err := c.BindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	account.User = models.User{Model: models.Model{ID: user.ID}}
	if err := models.Database.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, &account)
}

// @route    PUT api/accounts/:id
// @desc     Update an account
// @access   Private
func putAccount(c *gin.Context) {
	userId := c.MustGet("user").(*models.AuthUser).ID
	accountId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Account ID")
		return
	}
	var payload struct {
		BankName      string
		BankShort     string
		BankUser      string
		AccountNumber string
		Type          string `binding:"oneof=Ahorros Corriente"`
		Nickname      string
		Currency      string
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	accountUpdates := models.Account{
		BankName:      payload.BankName,
		BankShort:     payload.BankShort,
		BankUser:      payload.BankUser,
		AccountNumber: payload.AccountNumber,
		Type:          payload.Type,
		Nickname:      payload.Nickname,
		Currency:      payload.Currency,
	}
	if err := models.Database.Model(&models.Account{}).
		Where(&models.Account{Model: models.Model{ID: uint(accountId)}, UserID: userId}).
		Updates(&accountUpdates).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, &accountUpdates)
}

// @route    DELETE api/accounts/:id
// @desc     Delete a account
// @access   Private
func deleteAccount(c *gin.Context) {
	user := c.MustGet("user").(*models.AuthUser)
	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Account ID")
		return
	}
	account := models.Account{Model: models.Model{ID: uint(accountID)}, UserID: user.ID}

	if err := models.Database.First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, "Account not found")
		return
	}

	if err := models.Database.Delete(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func accounts(r *gin.RouterGroup) {
	r.Use(middleware.Auth).
		GET("", getAccounts).
		POST("", postAccount).
		PUT("/:id", putAccount).
		DELETE("/:id", deleteAccount)
}
