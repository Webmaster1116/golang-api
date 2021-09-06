package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/pkg/errors"
)

// @route    GET api/profiles/all
// @desc     Get All user's profile
// @access   Private
func getProfilesAll(c *gin.Context) {
	userId := c.MustGet("user").(*models.AuthUser).ID
	var page struct {
		Page  int `form:"page" binding:"required,min=1"`
		Limit int `form:"limit" binding:"required,min=1"`
	}
	if err := c.BindQuery(&page); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	skip := (page.Page - 1) * page.Limit
	var profiles []models.Profile

	if err := models.Database.
		Limit(page.Limit).
		Offset(skip).
		Order("date DESC").
		Where(&models.Profile{UserID: userId}).
		Find(&profiles).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to query profiles page"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, profiles)
}

// @route    GET api/profiles
// @desc     Get user's profile
// @access   Private
func getProfiles(c *gin.Context) {
	userId := c.MustGet("user").(*models.AuthUser).ID
	var profiles []models.Profile
	if err := models.Database.
		Where(&models.Profile{UserID: userId}).
		Find(&profiles).Error; err != nil {

		c.Error(errors.Wrap(err, "failed to query user profiles"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, profiles)
}

// @route    POST api/profiles
// @desc     Save a user's profile
// @access   Private
func postProfiles(c *gin.Context) {
	userID := c.MustGet("user").(*models.AuthUser).ID
	var payload struct {
		FirstName   string    `binding:"required"`
		LastName    string    `binding:"required"`
		BirthDate   time.Time `binding:"required"`
		Email       string    `binding:"required,email"`
		Phone       string    `binding:"required"`
		Nationality string    `binding:"required"`
		Doctype     string    `binding:"required"`
		Dni         string    `binding:"required"`
		AccountType string    `binding:"required,oneof=Personal Empresa"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// check exists
	profile := models.Profile{Email: payload.Email}
	var count int64
	if err := models.Database.Model(&profile).
		Where(&profile).
		Count(&count).Error; err != nil || count > 0 {
		if err != nil {
			c.Error(err)
		}
		c.JSON(http.StatusBadRequest, "profile alread exists")
		return
	}

	profile = models.Profile{
		UserID:      userID,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		BirthDate:   payload.BirthDate,
		Email:       payload.Email,
		Phone:       payload.Phone,
		Nationality: payload.Nationality,
		Doctype:     payload.Doctype,
		Dni:         payload.Dni,
		AccountType: payload.AccountType,
	}
	// create
	if err := models.Database.Create(&profile).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to create profile"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, profile)
}

// @route    Put api/profiles
// @desc     Edit Profile Details
// @access   Private
func putProfiles(c *gin.Context) {
	userID := c.MustGet("user").(*models.AuthUser).ID
	var payload struct {
		FirstName   string
		LastName    string
		BirthDate   time.Time
		Email       string `binding:"email"`
		Phone       string
		Nationality string
		Doctype     string
		Dni         string
		AccountType string `binding:"oneof=Personal Empresa"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// find
	profile := models.Profile{UserID: userID}
	profileUpdates := models.Profile{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		BirthDate:   payload.BirthDate,
		Email:       payload.Email,
		Phone:       payload.Phone,
		Nationality: payload.Nationality,
		Doctype:     payload.Doctype,
		Dni:         payload.Dni,
		AccountType: payload.AccountType,
	}
	if err := models.Database.
		Where(&profile).
		First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, "profile not found")
		return
	}
	// update
	if err := models.Database.
		Model(&profile).
		Updates(&profileUpdates).Error; err != nil {
		if err != nil {
			c.Error(err)
		}
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, payload)
}

func profiles(r *gin.RouterGroup) {
	r.Use(middleware.Auth).
		GET("/all", getProfilesAll).
		GET("", getProfiles).
		POST("", postProfiles).
		PUT("", putProfiles)
}
