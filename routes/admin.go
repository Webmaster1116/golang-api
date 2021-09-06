package routes

import (
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
	"github.com/med8bra/moni-api-go/services"
	"github.com/med8bra/moni-api-go/util/password"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// @route    POST api/admin/new
// @desc     Register Admin
// @access   Public
func postAdminNew(c *gin.Context) {
	var payload struct {
		Name     string `binding:"required"`
		Email    string `binding:"required,email"`
		Password string `binding:"required,min=6"`
	}

	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	admin := models.Admin{Email: payload.Email}
	// check admin exists
	if models.Database.Where(&admin).First(&admin).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, "admin already exists")
		return
	}
	// hash password
	if hash, err := password.Hash(payload.Password); err != nil {
		c.Error(errors.Wrap(err, "failed to hash admin password"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	} else {
		admin.Password = hash
	}
	// save admin
	admin.Name = payload.Name
	if err := models.Database.Create(&admin).Error; err != nil {
		c.Error(errors.Wrap(err, "failed to create admin entity"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}

	if token, err := services.Authenticator.Sign(admin.ID, time.Hour); err != nil {
		c.Error(errors.Wrap(err, "failed to sign admin token"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

// @route    GET api/admin/auth
// @desc     Get logged user
// @access   Private
func getAdminAuth(c *gin.Context) {
	adminID := c.MustGet("user").(*models.AuthUser).ID
	var admin struct {
		ID    uint
		Name  string
		Email string
	}
	if err := models.Database.Model(&models.Admin{}).First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, admin)
}

// @route    POST api/admin/auth
// @desc     Authenticate Admin & get token
// @access   Public
func postAdminAuth(c *gin.Context) {
	var payload struct {
		Email    string `binding:"required,email"`
		Password string `binding:"required,min=6"`
	}
	auth := services.Authenticator
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	admin := models.Admin{Email: payload.Email}
	if err := models.Database.Where(&admin).First(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, "admin not found")
		return
	}
	hash, _ := hex.DecodeString(admin.Password)
	if err := bcrypt.CompareHashAndPassword(hash, []byte(payload.Password)); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Credentials")
		return
	}

	token, err := auth.Sign(admin.ID, 24*time.Hour)
	if err != nil {
		c.Error(errors.Wrap(err, "failed to sign admin token"))
		c.JSON(http.StatusInternalServerError, "cannot authenticate admin")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func admin(r *gin.RouterGroup) {
	// private
	r.GET("/auth", middleware.Auth, getAdminAuth)

	// public
	r.
		POST("/new", postAdminNew).
		POST("/auth", postAdminAuth)
}
