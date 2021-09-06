package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/pkg/errors"
)

// @route    GET api/agentprofile/id
// @desc     Get agent profile with agent id
// @access   Private
func getAgentProfile(c *gin.Context) {
	agentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}
	agentProfile := models.AgentProfile{AgentID: uint(agentId)}

	if err := models.Database.Where(&agentProfile).First(&agentProfile).Error; err != nil {
		c.JSON(http.StatusBadRequest, "agent profile not found")
		return
	}
	c.JSON(http.StatusOK, agentProfile)
}

// @route    GET api/agentprofile
// @desc     Get user's profile
// @access   Private
func getUserProfile(c *gin.Context) {
	agentId := c.MustGet("user").(*models.AuthUser).ID
	profile := models.AgentProfile{AgentID: agentId}

	if err := models.Database.Where(&profile).Order("date DESC").First(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "agent profile not found")
		return
	}
	c.JSON(http.StatusOK, profile)
}

// @route    POST api/agentprofiles
// @desc     Save a agent's profile
// @access   Private
func postAgentProfiles(c *gin.Context) {
	userID := c.MustGet("user").(*models.AuthUser).ID
	var payload struct {
		FirstName   string `binding:"required"`
		LastName    string `binding:"required"`
		Email       string `binding:"required,email"`
		Phone       string `binding:"required"`
		Address     string `binding:"required"`
		Dni         string `binding:"required"`
		AccountType string `binding:"required,oneof=Personal Empresa"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	profile := models.AgentProfile{
		AgentID:     userID,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Phone:       payload.Phone,
		Dni:         payload.Dni,
		AccountType: payload.AccountType,
	}

	if err := models.Database.Create(&profile).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to create user profile"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, profile)
}

// @route    Put api/agentprofile/:id
// @desc     Edit Agent Profile Details
// @access   Private
func putAgentProfile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	agentId := uint(id)
	var payload models.AgentProfile
	profile := models.AgentProfile{AgentID: agentId}
	if err := models.Database.Where(&profile).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, "profile not found")
		return
	}

	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// update
	if err := models.Database.Model(&profile).Updates(&payload).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to update profile"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, profile)
}

func agentprofile(r *gin.RouterGroup) {
	r.Use(middleware.Auth).
		GET("/:id", getAgentProfile).
		GET("", getUserProfile).
		POST("", postAgentProfiles).
		PUT("/:id", putAgentProfile)
}
