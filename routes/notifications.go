package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/models"
	"github.com/med8bra/moni-api-go/services"
	"github.com/med8bra/moni-api-go/services/webpush"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func notify(subscription *webpush.Subscription, msg string, options *webpush.Options) {
	if err := services.WebPush.Notify(
		subscription,
		msg,
		options,
	); err != nil {
		logrus.Warnf("failed to push notification %v : %s", subscription, err.Error())
	}
}

func notifyUsers() {
	var userNotify models.UserNotify

	rows, err := models.Database.Model(&userNotify).Rows()

	if err != nil {
		logrus.Error("failed to query user subscriptions : ", err.Error())
		return
	}

	for rows.Next() {
		if err := rows.Scan(&userNotify); err != nil {
			logrus.Error("failed to scan user subscription : ", err.Error())
			return
		}

		subscription := webpush.Subscription{
			Endpoint: userNotify.Endpoint,
			Keys:     webpush.Keys(userNotify.Keys),
		}
		msg := `{"title": "Moni"}`
		notify(&subscription, msg, &webpush.Options{})
	}
}

// @route    POST api/notifications/subscribe
// @desc     Send notification to opening app1
// @access   Public
func postNotificationSubscribe(c *gin.Context) {
	var payload struct {
		Endpoint string `binding:"required"`
		Keys     struct {
			Auth   string `binding:"required"`
			P256dh string `binding:"required"`
		}
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	userNotify := models.UserNotify{Endpoint: payload.Endpoint}

	if models.Database.Where(&userNotify).First(&userNotify).RowsAffected == 1 {
		c.JSON(http.StatusOK, "user already subscribed")
		return
	}

	// subscribe
	userNotify.ID = 0
	if err := models.Database.Create(&userNotify).Error; err != nil {
		c.Error(errors.Wrapf(err, "failed to create user notify %v", userNotify))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusCreated, "subscription saved")

	// send notification
	subscription := webpush.Subscription{
		Endpoint: payload.Endpoint,
		Keys: webpush.Keys{
			Auth:   payload.Keys.Auth,
			P256dh: payload.Keys.P256dh,
		},
	}
	msg := `{"title": "Moni"}`
	go notify(&subscription, msg, &webpush.Options{})
}

// @route    POST api/notifications/subscribe
// @desc     Send notifications to all admins about transaction
// @access   Public
func postNotificationNotify(c *gin.Context) {
	var payload struct {
		TTL int `binding:"required"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, "push triggred")

	// -- notify admins async
	go notifyUsers()
}

func notification(r *gin.RouterGroup) {
	r.
		POST("/subscribe", postNotificationSubscribe).
		POST("/notify", postNotificationNotify)
}
