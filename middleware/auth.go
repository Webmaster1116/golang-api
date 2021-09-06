package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/services"
)

func Auth(c *gin.Context) {
	auth := services.Authenticator
	// get token
	tokenStr := c.GetHeader("x-auth-token")
	if user, err := auth.Verify(tokenStr); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
	} else {
		c.Set("user", user)
	}
}
