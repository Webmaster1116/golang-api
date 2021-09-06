package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/models"
)

// @route    GET api/adminnotify
// @desc     Create AdminNotify
// @access   Public
func postAdminNotify(c *gin.Context) {
	var payload struct {
		Endpoint string      `binding:"required"`
		Keys     models.Keys `binding:"required"`
	}

	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	adminNotify := models.AdminNotify{Endpoint: payload.Endpoint, Keys: payload.Keys}
	if err := models.Database.Create(&adminNotify).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "Server error")
		return
	}
	c.JSON(http.StatusOK, "AdminNotify saved")
}

func adminnotify(r *gin.RouterGroup) {
	r.
		POST("", postAdminNotify)
}
