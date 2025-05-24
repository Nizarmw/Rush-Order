package routes

import (
	"RushOrder/controller"

	"github.com/gin-gonic/gin"
)

func SetupSessionRoutes(r *gin.Engine) {
	SessionRoutes := r.Group("/api/sessions")
	{
		SessionRoutes.POST("/login", controller.CustomerLoginHandler)
		SessionRoutes.POST("/logout", controller.LogoutHandler)
		SessionRoutes.GET("/", controller.GetCustomerSessionHandler)
	}
}
