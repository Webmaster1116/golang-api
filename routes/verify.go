package routes

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/models"
	"github.com/med8bra/moni-api-go/services"
	"github.com/med8bra/moni-api-go/services/mailer"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	verificationTokenExpires = time.Hour
)

func sendVerifyEmail(tx *gorm.DB, user *models.User, host string) error {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return err
	}
	tokenHex := hex.EncodeToString(token)
	verificationToken := models.Token{
		UserID: user.ID,
		Token:  tokenHex,
	}
	if err := tx.Create(&verificationToken).Error; err != nil {
		return err
	}
	verificationTemplate, found := services.TemplateManager.Get("verification")
	if !found {
		return ErrVerificationTemplateNotFound
	}
	var msg bytes.Buffer
	verificationLink := fmt.Sprintf("http://%s/api/verify/confirmation?token=3D%s", host, tokenHex)
	if err := verificationTemplate.Execute(&msg, verificationLink); err != nil {
		return err
	}
	// send in async
	go func() {
		if err := services.Mailer.Send(&mailer.Email{
			Subject: "Account Verification Token",
			From:    "Moni Platfrom",
			HTML:    msg.Bytes(),
			To:      user.Email,
		}); err != nil {
			logrus.Error("Failed to send verification Email: ", err.Error())
		}
	}()
	return nil
}

// @route    GET api/verify/confirmation
// @desc     Confirm a User's email
// @access   Public
func getVerifyConfirmation(c *gin.Context) {
	tokenHex, found := c.GetQuery("token")
	if !found {
		c.JSON(http.StatusBadRequest, "missing token")
		return
	}
	var token models.Token
	// find token
	if err := models.Database.Where("token = ?", tokenHex).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, "invalid token")
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Server error")
		}
		return
	}

	// check expired
	if time.Since(token.CreatedAt) > verificationTokenExpires {
		models.Database.Delete(&token)
		c.JSON(http.StatusBadRequest, "token expired")
		return
	}

	// update user
	var user models.User
	if err := models.Database.First(&user, token.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, "invalid token")
		} else {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, "Server error")
		}
		return
	}

	user.Verified = true
	err := models.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&token).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}

	c.JSON(http.StatusOK, "The account has been verified. Please log in.")
}

// @route    POST api/verify/resend
// @desc     Resend Confirmation Mail
// @access   Public
func postVerifyResend(c *gin.Context) {
	var payload struct {
		Email string `binding:"required,email"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// get user
	user := models.User{Email: payload.Email}
	if err := models.Database.Where(&user).First(&user).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, "We were unable to find a user with that email")
		return
	}

	// already verified
	if user.Verified {
		c.JSON(http.StatusBadRequest, "This account has already been verified. Please log in.")
		return
	}

	if err := sendVerifyEmail(models.Database, &user, c.Request.Host); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	c.JSON(http.StatusOK, "A verification email has been sent to "+payload.Email)
}

func verify(r *gin.RouterGroup) {
	r.
		GET("/confirmation", getVerifyConfirmation).
		POST("/resend", postVerifyResend)
}
