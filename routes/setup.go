package routes

import "github.com/gin-gonic/gin"

type RouterRegister func(*gin.RouterGroup)

var routes = map[string]RouterRegister{
	"/accounts":      accounts,
	"/users":         users,
	"/auth":          auth,
	"/admin":         admin,
	"/verify":        verify,
	"/newsletter":    newsletter,
	"/adminnotify":   adminnotify,
	"/exchange":      exchange,
	"/forget":        forget,
	"/notification":  notification,
	"/agentprofile":  agentprofile,
	"/mail":          mail,
	"/profiles":      profiles,
	"/agentaccounts": agentaccounts,
	"/agent":         agent,
	"/operations":    operations,
}

func Init(r *gin.RouterGroup) error {
	for group, register := range routes {
		register(r.Group(group))
	}
	return nil
}
