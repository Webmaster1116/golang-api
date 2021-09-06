package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/models"
)

// @route    GET api/newsletter
// @desc     Get all emails subscribed to newsletter
// @access   Public
func getNewsletter(c *gin.Context) {
	var newsletters []models.Newsletter

	if err := models.Database.Find(&newsletters).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	c.JSON(http.StatusOK, newsletters)
}

// @route    POST api/newsletter
// @desc     Save a email to newsletter object
// @access   Public
func postNewsletter(c *gin.Context) {
	var payload struct {
		Email string `binding:"required,email"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	newsletter := models.Newsletter{Email: payload.Email}

	if err := models.Database.Create(&newsletter).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	c.JSON(http.StatusOK, newsletter)
}

func newsletter(r *gin.RouterGroup) {
	r.
		GET("", getNewsletter).
		POST("", postNewsletter)
}
