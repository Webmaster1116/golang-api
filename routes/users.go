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

// @route    GET api/users
// @desc     Get all user's profile
// @access   Private
func getUsers(c *gin.Context) {
	var users []struct {
		ID       uint
		Name     string
		Email    string
		Verified bool
		Date     time.Time
	}
	err := models.Database.
		Model(&models.User{}).
		Order("date DESC").
		Find(&users).
		Error

	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, users)
}

// @route    GET api/users/count
// @desc     Get count of user profiles
// @access   Private
func getUsersCount(c *gin.Context) {
	var usersCount int64
	err := models.Database.
		Model(&models.User{}).
		Count(&usersCount).
		Error

	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, usersCount)
}

// @route    POST api/users
// @desc     Register user
// @access   Public
func postUser(c *gin.Context) {
	var payload struct {
		Name     string `binding:"required"`
		Email    string `binding:"required,email"`
		Password string `binding:"required,min=6"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if models.Database.Where("email = ?", payload.Email).First(&models.User{}).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, "user already exists")
		return
	}

	if hash, err := password.Hash(payload.Password); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		payload.Password = hash
	}
	user := models.User{Name: payload.Name, Email: payload.Email, Password: payload.Password}
	err := models.Database.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&user).Error; err != nil {
			return errors.Wrap(err, "failed to create user")
		}
		if err := sendVerifyEmail(tx, &user, c.Request.Host); err != nil {
			return errors.Wrap(err, "failed to send verification mail")
		}
		return nil
	})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	token, err := services.Authenticator.Sign(user.ID, 24*time.Hour)
	if err != nil {
		c.Error(errors.Wrap(err, "failed to sign user token"))
		c.JSON(http.StatusInternalServerError, "cannot authenticate user")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func users(r *gin.RouterGroup) {
	// private
	r.Group("", middleware.Auth).
		GET("", getUsers).
		GET("/count", getUsersCount)

		// public
	r.POST("", postUser)
}
