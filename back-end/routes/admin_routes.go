package routes

import (
	"RushOrder/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAdminRoutes(r *gin.Engine, db *gorm.DB) {
	admin := r.Group("/api/admin")
	admin.POST("/login", controller.AdminLoginHandler(db))
	admin.POST("/register", controller.AdminRegisterHandler(db))
	admin.POST("/logout", controller.AdminLogoutHandler())
	admin.GET("/orders", controller.GetAdminOrdersHandler)
	admin.PUT("/orders/status", controller.UpdateAdminStatusHandler)

	// Legacy route for compatibility
	admin.GET("/order", controller.GetOrdersAdminHandler(db))
}
