package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/med8bra/moni-api-go/services"
	"github.com/med8bra/moni-api-go/util/password"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// @route    GET api/agent/reset/:id
// @desc     Reset agent commissions
// @access   Private
func getAgentRestId(c *gin.Context) {
	agentId := c.Param("id")
	var agent models.Agent

	if err := models.Database.First(&agent, agentId).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Error(errors.Wrap(err, "failed to get agent"))
		}
		c.JSON(http.StatusBadRequest, "agent not found")
		return
	}
	// reset comissions
	agent.CommissionDollars = 0
	agent.CommissionSoles = 0
	if err := models.Database.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agent)
}

// @route    GET api/agent/enable/:id
// @desc     Make agent enabled
// @access   Private
func getAgentEnableId(c *gin.Context) {
	agentId := c.Param("id")
	var agent models.Agent

	if err := models.Database.First(&agent, agentId).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Error(errors.Wrap(err, "failed to get agent"))
		}
		c.JSON(http.StatusBadRequest, "agent not found")
		return
	}
	// enable agent
	agent.Enable = true
	if err := models.Database.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agent)
}

// @route    GET api/agent/disable/:id
// @desc     Make agent disabled
// @access   Private
func getAgentDisableId(c *gin.Context) {
	agentId := c.Param("id")
	var agent models.Agent

	if err := models.Database.First(&agent, agentId).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Error(errors.Wrap(err, "failed to get agent"))
		}
		c.JSON(http.StatusBadRequest, "agent not found")
		return
	}
	// disable agent
	agent.Enable = false
	if err := models.Database.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agent)
}

// @route    GET api/agent/online/:id
// @desc     Make agent online
// @access   Private
func getAgentOnlineId(c *gin.Context) {
	agentId := c.Param("id")
	var agent models.Agent

	if err := models.Database.First(&agent, agentId).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Error(errors.Wrap(err, "failed to get agent"))
		}
		c.JSON(http.StatusBadRequest, "agent not found")
		return
	}
	// make online
	agent.Online = true
	if err := models.Database.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agent)
}

// @route    GET api/agent/offline/:id
// @desc     Make agent offline
// @access   Private
func getAgentOfflineId(c *gin.Context) {
	agentId := c.Param("id")
	var agent models.Agent

	if err := models.Database.First(&agent, agentId).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.Error(errors.Wrap(err, "failed to get agent"))
		}
		c.JSON(http.StatusBadRequest, "agent not found")
		return
	}
	// make online
	agent.Online = false
	if err := models.Database.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agent)
}

// @route    GET api/agent
// @desc     Get all agent's profile
// @access   Private
func getAgent(c *gin.Context) {
	var agents []models.Agent
	if err := models.Database.Preload("Operations").Model(&models.Agent{}).Find(&agents).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to get agents"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, agents)
}

// @route    GET api/agent/invoice
// @desc     Get every agent's invoice data
// @access   Private
func getAgentInvoice(c *gin.Context) {
	var invoices []models.Invoice

	if err := models.Database.
		Order("date DESC").
		Find(&invoices).
		Error; err != nil {
		c.Error(errors.Wrap(err, "failed to get invoices"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, invoices)
}

// @route    POST api/agent/new
// @desc     Register Agent
// @access   Public
func postAgentNew(c *gin.Context) {
	var payload struct {
		Name     string `binding:"required"`
		Email    string `binding:"required,email"`
		Password string `binding:"required,min=6"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	agent := models.Agent{
		Name:  payload.Name,
		Email: payload.Email,
	}
	// check email address
	if models.Database.Where("email = ?", payload.Email).First(&models.Agent{}).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, "agent already exists")
		return
	}
	// hash password
	if hash, err := password.Hash(payload.Password); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		agent.Password = hash
	}

	if err := models.Database.Create(&agent).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to create agent"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	// generate token
	if token, err := services.Authenticator.Sign(agent.ID, 24*time.Hour); err == nil {
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.Error(errors.Wrap(err, "failed to generate token"))
		c.JSON(http.StatusInternalServerError, "failed to generate token")
	}
}

// @route    GET api/agent/auth
// @desc     Get logged Agent
// @access   Private
func getAgentAuth(c *gin.Context) {
	// copy user
	agentID := c.MustGet("user").(*models.AuthUser).ID
	var agent struct {
		ID    uint
		Name  string
		Email string
	}
	if err := models.Database.Model(&models.Agent{}).First(&agent, agentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, agent)
}

// @route    POST api/agent/auth
// @desc     Authenticate Agent & get token
// @access   Public
func postAgentAuth(c *gin.Context) {
	var payload struct {
		Email    string `binding:"required,email"`
		Password string `binding:"required,min=6"`
	}
	auth := services.Authenticator
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	agent := models.Agent{Email: payload.Email}
	if err := models.Database.Where(&agent).First(&agent).Error; err != nil {
		c.JSON(http.StatusBadRequest, "agent not found")
		return
	}
	if !password.Verify(payload.Password, agent.Password) {
		c.JSON(http.StatusBadRequest, "Invalid Credentials")
		return
	}

	token, err := auth.Sign(agent.ID, 24*time.Hour)
	if err != nil {
		c.Error(errors.Wrap(err, "failed to sign agent token"))
		c.JSON(http.StatusInternalServerError, "cannot authenticate agent")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func agent(r *gin.RouterGroup) {
	// -- private
	r.Group("", middleware.Auth).
		GET("/reset/:id", getAgentRestId).
		GET("/enable/:id", getAgentEnableId).
		GET("/disable/:id", getAgentDisableId).
		GET("/online/:id", getAgentOnlineId).
		GET("/offline/:id", getAgentOfflineId).
		GET("/invoice", getAgentInvoice).
		GET("", getAgent).
		GET("/auth", getAgentAuth)

	// -- public
	r.
		POST("/new", postAgentNew).
		POST("/auth", postAgentAuth)
}
