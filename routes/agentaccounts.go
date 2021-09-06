package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// @route    GET api/agentaccounts/balance/:balance/:id
// @desc     Make account balance true or false
// @access   Private
func getAgentAccountsBalance(c *gin.Context) {
	var payload struct {
		ID      uint `uri:"id" binding:"required"`
		Balance bool `uri:"balance" binding:"required"`
	}
	if err := c.BindUri(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var agentAccount models.AgentAccount

	if err := models.Database.First(&agentAccount, payload.ID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Error(errors.Wrap(err, "failed to get account"))
		}
		c.JSON(http.StatusBadRequest, "account not found")
		return
	}
	agentAccount.Balance = payload.Balance
	if err := models.Database.Save(&agentAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agentAccount)
}

// @route    GET api/agentaccounts/:id
// @desc     Get agent accounts with agent id
// @access   Private
func getAgentAccountsId(c *gin.Context) {
	var payload struct {
		AgentID uint `uri:"id" binding:"required"`
	}
	if err := c.BindUri(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var agentAccounts []models.AgentAccount

	if err := models.Database.
		Where(&models.AgentAccount{AgentID: payload.AgentID}).
		Find(&agentAccounts).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agentAccounts)
}

// @route    GET api/agentaccounts
// @desc     Get all accounts of agent
// @access   Private
func getAgentAccounts(c *gin.Context) {
	userId := c.MustGet("user").(*models.User).ID
	var agentAccounts []models.AgentAccount

	if err := models.Database.
		Where(&models.AgentAccount{AgentID: userId}).
		Order("date DESC").
		Find(&agentAccounts).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agentAccounts)
}

// @route    POST api/agentaccounts
// @desc     Save an agent account Data
// @access   Private
func postAgentAccounts(c *gin.Context) {
	userId := c.MustGet("user").(*models.User).ID
	var payload struct {
		BankName  string `binding:"required"`
		BankShort string `binding:"required"`
		BankUser  string `binding:"required"`
		Type      string `binding:"required, oneof=Ahorros Corriente"`
		Currency  string `binding:"required"`
		Purpose   string `binding:"required"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	agentAccount := models.AgentAccount{
		BankName:  payload.BankName,
		BankShort: payload.BankShort,
		BankUser:  payload.BankUser,
		Type:      payload.Type,
		Currency:  payload.Currency,
		Purpose:   payload.Purpose,
		AgentID:   userId,
	}

	if err := models.Database.Create(&agentAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agentAccount)
}

// @route    PUT api/agentaccounts/:id
// @desc     Update an agent account
// @access   Private
func putAgentAccounts(c *gin.Context) {
	userId := c.MustGet("user").(*models.User).ID
	var payload struct {
		ID        uint   `uri:"id" binding:"required"`
		BankName  string `binding:"required"`
		BankShort string `binding:"required"`
		BankUser  string `binding:"required"`
		Type      string `binding:"required, oneof=Ahorros Corriente"`
		Currency  string `binding:"required"`
		Purpose   string `binding:"required"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	agentAccount := models.AgentAccount{
		Model:     models.Model{ID: payload.ID},
		BankName:  payload.BankName,
		BankShort: payload.BankShort,
		BankUser:  payload.BankUser,
		Type:      payload.Type,
		Currency:  payload.Currency,
		Purpose:   payload.Purpose,
		AgentID:   userId,
	}

	if err := models.Database.Updates(&agentAccount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, "account not found")
		} else {
			c.JSON(http.StatusInternalServerError, "server error")
		}
		return
	}
	c.JSON(http.StatusOK, agentAccount)
}

// @route    DELETE api/agentaccounts/:id
// @desc     Delete a account
// @access   Private
func deleteAgentAccounts(c *gin.Context) {
	var payload struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := models.Database.Delete(&models.AgentAccount{}, payload.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, "account not found")
		} else {
			c.JSON(http.StatusInternalServerError, "server error")
		}
		return
	}
	c.JSON(http.StatusOK, "account removed")
}

func agentaccounts(r *gin.RouterGroup) {
	r.Use(middleware.Auth).
		GET("/balance/:balance/:id", getAgentAccountsBalance).
		GET("/:id", getAgentAccountsId).
		GET("", getAgentAccounts).
		POST("", postAgentAccounts).
		PUT("/:id", putAgentAccounts).
		DELETE("/:id", deleteAgentAccounts)
}
