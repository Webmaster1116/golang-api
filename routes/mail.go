package routes

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/services"
	"github.com/med8bra/moni-api-go/services/mailer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ProfileDetails struct {
	FirstName string `binding:"required"`
	LastName  string `binding:"required"`
}

type BankDetails struct {
	BankName      string `binding:"required"`
	AccountNumber string `binding:"required"`
	Type          string `binding:"required"`
}

// @route    POST api/mail/verify
// @desc     Send payment confirmation
// @access   Private
func postMailVerify(c *gin.Context) {
	var payload struct {
		Email string    `binding:"required,email"`
		Date  time.Time `binding:"required"`
		ProfileDetails
		BankDetails
		Transaction struct {
			AmountToPay   float64 `binding:"required"`
			AmountRecieve float64 `binding:"required"`
			CurrenyTo     string  `binding:"required"`
			CurrenyFrom   string  `binding:"required"`
		}
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	paymentConfirmedTemplate, found := services.TemplateManager.Get("paymentConfirmed")
	if !found {
		c.Error(errors.New("failed to get paymentConfirmed template"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	var msg bytes.Buffer
	if err := paymentConfirmedTemplate.Execute(&msg, &payload); err != nil {
		c.Error(errors.Wrap(err, "failed to render payment confirmed template"))
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}

	c.JSON(http.StatusOK, "notification triggred")

	go func() {
		// TODO: add content-type HTML
		if err := services.Mailer.Send(&mailer.Email{
			Subject: "Payment Added",
			From:    "Moni Platform",
			HTML:    msg.Bytes(),
			To:      payload.Email,
		}); err != nil {
			logrus.Error("failed to send payment confirmation email : ", err.Error())
		}
	}()
}

// @route    POST api/mail/new
// @desc     Send payment confirmation
// @access   Private
func postMailNew(c *gin.Context) {
	var payload struct {
		ID    uint      `binding:"required"`
		Email string    `binding:"required,email"`
		Date  time.Time `binding:"required"`
		ProfileDetails
		BankDetails struct {
			ReceiveBank BankDetails
		}
		Transaction struct {
			AmountToPay   float64 `binding:"required"`
			AmountRecieve float64 `binding:"required"`
			CurrenyTo     string  `binding:"required"`
			CurrenyFrom   string  `binding:"required"`
		}
		Agent      ProfileDetails
		AgentBank  BankDetails
		AgentEmail string `binding:"required,email"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "sending payment notifications")
	// async
	go func() {
		var msg bytes.Buffer
		// -- send to user
		template, found := services.TemplateManager.Get("newOperation")
		if !found {
			logrus.Error("failed to get new operation template")
			return
		}
		if err := template.Execute(&msg, payload); err != nil {
			logrus.Error("failed to render new operation template : ", err.Error())
			return
		}
		if err := services.Mailer.Send(&mailer.Email{
			Subject: "Nueva operación iniciada",
			From:    "Moni Platform",
			HTML:    msg.Bytes(),
			To:      payload.Email,
		}); err != nil {
			logrus.Error("failed to send new operation mail : ", err.Error())
			return
		}
		// -- send to user
		template, found = services.TemplateManager.Get("newOperationAssigned")
		if !found {
			logrus.Error("failed to get new operation assigned template")
			return
		}
		msg.Reset()
		if err := template.Execute(&msg, payload); err != nil {
			logrus.Error("failed to render new operation assigned template : ", err.Error())
			return
		}
		if err := services.Mailer.Send(&mailer.Email{
			Subject: "Nueva operación iniciada",
			From:    "Moni Platform",
			HTML:    msg.Bytes(),
			To:      payload.AgentEmail,
		}); err != nil {
			logrus.Error("failed to send new operation assigned mail : ", err.Error())
			return
		}
	}()
}

// @route    POST api/mail/invoice
// @desc     Send invoice
// @access   Public
func postMailInvoice(c *gin.Context) {
	var payload struct {
		ID                uint    `binding:"required"`
		Name              string  `binding:"required"`
		Email             string  `binding:"required,email"`
		commissionDollars float64 `binding:"required"`
		commissionSoles   float64 `binding:"required"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "sending invoice")
	// -- async
	go func() {
		var msg bytes.Buffer
		// -- send to user
		template, found := services.TemplateManager.Get("invoice")
		if !found {
			logrus.Error("failed to get invoice template")
			return
		}
		if err := template.Execute(&msg, payload); err != nil {
			logrus.Error("failed to render invoice template : ", err.Error())
			return
		}
		if err := services.Mailer.Send(&mailer.Email{
			Subject: "Factura",
			From:    "Moni Platform",
			HTML:    msg.Bytes(),
			To:      payload.Email,
		}); err != nil {
			logrus.Error("failed to send invoice mail : ", err.Error())
			return
		}
	}()
}

// @route    POST api/mail/success
// @desc     Notify operation success
// @access   Private
func postMailSuccess(c *gin.Context) {
	var payload struct {
		Email string `binding:"required,email"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "sending status")
	// -- async
	go func() {
		var msg bytes.Buffer
		// -- send to user
		template, found := services.TemplateManager.Get("statusChange")
		if !found {
			logrus.Error("failed to get status change template")
			return
		}
		if err := template.Execute(&msg, payload); err != nil {
			logrus.Error("failed to render status change template : ", err.Error())
			return
		}
		if err := services.Mailer.Send(&mailer.Email{
			Subject: "Operación completada",
			From:    "Moni Platform",
			HTML:    msg.Bytes(),
			To:      payload.Email,
		}); err != nil {
			logrus.Error("failed to send status change mail : ", err.Error())
			return
		}
	}()
}

func mail(r *gin.RouterGroup) {
	// private
	r.Group("", middleware.Auth).
		POST("/verify", postMailVerify).
		POST("/new", postMailNew).
		POST("/success", postMailSuccess)

	// public
	r.POST("/invoice", postMailInvoice)
}
