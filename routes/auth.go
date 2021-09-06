package routes

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/med8bra/moni-api-go/services"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// @route    GET api/auth
// @desc     Get logged user
// @access   Private
func getAuth(c *gin.Context) {
	// copy user
	userID := c.MustGet("user").(*models.AuthUser).ID
	var user struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := models.Database.Model(&models.User{}).First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

// @route    POST api/auth
// @desc     Authenticate user & get token
// @access   Public
func postAuth(c *gin.Context) {
	var payload struct {
		Email    string `binding:"required,email"`
		Password string `binding:"required,min=6"`
	}
	auth := services.Authenticator
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	password := payload.Password
	var user models.User
	if err := models.Database.Where(&models.User{Email: payload.Email}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, "user not found")
		return
	}
	hash, _ := hex.DecodeString(user.Password)
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Credentials")
		return
	}

	if !user.Verified {
		c.JSON(http.StatusUnauthorized, "Your account has not been verified")
		return
	}

	token, err := auth.Sign(user.ID, 24*time.Hour)
	if err != nil {
		c.Error(errors.Wrap(err, "failed to sign user token"))
		c.JSON(http.StatusInternalServerError, "cannot authenticate user")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func auth(r *gin.RouterGroup) {
	r.
		GET("", middleware.Auth, getAuth).
		POST("", postAuth)
}
