package routes

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/models"
	"github.com/med8bra/moni-api-go/services"
	"github.com/med8bra/moni-api-go/services/mailer"
	"github.com/med8bra/moni-api-go/util/password"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func sendResetPasswordEmail(user *models.User, recoveryLink string) error {
	resetTemplate, found := services.TemplateManager.Get("resetPassword")
	if !found {
		return ErrResetPasswordTemplateNotFound
	}
	var msg bytes.Buffer
	if err := resetTemplate.Execute(&msg, map[string]interface{}{
		"user": &user,
		"link": &recoveryLink,
	}); err != nil {
		return err
	}
	// send in async
	go func() {
		if err := services.Mailer.Send(&mailer.Email{
			Subject: "Link to Reset Password",
			From:    "Moni Platform",
			Text:    msg.Bytes(),
			To:      user.Email,
		}); err != nil {
			logrus.Error("failed to send verification email: ", err.Error())
		}
	}()
	return nil
}

// @route    POST api/forget/forgetpassword
// @desc     Request Password reset link
// @access   Public
func postForgetPassword(c *gin.Context) {
	var payload struct {
		Email        string `binding:"required,email"`
		RecoveryLink string `binding:"required"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := models.User{Email: payload.Email}
	if err := models.Database.Where(&user).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}
	tokenBuf := make([]byte, 20)
	if _, err := rand.Read(tokenBuf); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	tokenHex := hex.EncodeToString(tokenBuf)

	err := models.Database.Transaction(func(tx *gorm.DB) error {
		user.ResetPasswordToken = tokenHex
		user.ResetPasswordExpires = time.Now().Add(time.Hour)

		if err := tx.Save(&user).Error; err != nil {
			return errors.Wrap(err, "failed to save user")
		}
		if err := sendResetPasswordEmail(&user, payload.RecoveryLink+tokenHex); err != nil {
			return errors.Wrap(err, "failed to send reset email")
		}
		return nil
	})

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	c.JSON(http.StatusOK, "Reset email is sent")
}

// @route    GET api/forget/reset
// @desc     Get user info for reset token
// @access   Public
func getForgetReset(c *gin.Context) {
	tokenHex, found := c.GetQuery("resetPasswordToken")
	if !found {
		c.JSON(http.StatusBadRequest, "missing resetPasswordToken")
		return
	}
	var user struct {
		ID                   uint
		Name                 string
		Email                string
		ResetPasswordExpires time.Time `json:"-"`
	}

	if err := models.Database.Model(&models.User{}).Where(&models.User{ResetPasswordToken: tokenHex}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, "invalid reset token")
		return
	}

	if user.ResetPasswordExpires.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, "token expired")
		return
	}

	c.JSON(http.StatusOK, &user)
}

// @route    POST api/forget/updatePassword
// @desc     Update user password using reset token
// @access   Public
func postUpdatePassword(c *gin.Context) {
	var payload struct {
		Password           string `binding:"required,min=6"`
		ResetPasswordToken string `binding:"required,min=20"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := models.User{ResetPasswordToken: payload.ResetPasswordToken}

	if err := models.Database.Where(&user).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, "invalid token")
		return
	}
	if user.ResetPasswordExpires.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, "token expired")
		return
	}

	if hash, err := password.Hash(payload.Password); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		user.Password = hash
		user.ResetPasswordToken = "."
	}
	if err := models.Database.Updates(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "user password updated")
}

func forget(r *gin.RouterGroup) {
	r.
		POST("/forgetpassword", postForgetPassword).
		GET("/reset", getForgetReset).
		POST("/updatepassword", postUpdatePassword)
}
