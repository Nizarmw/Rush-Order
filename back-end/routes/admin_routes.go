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

	admin.GET("/order", controller.GetOrdersAdminHandler(db))
}
