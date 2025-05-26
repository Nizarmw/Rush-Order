package routes

import (
	"RushOrder/controller"
	"RushOrder/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAdminRoutes(r *gin.Engine, db *gorm.DB) {
	admin := r.Group("/api/admin")
	admin.POST("/login", controller.AdminLoginHandler(db))
	admin.POST("/register", controller.AdminRegisterHandler(db))
    admin.POST("/logout", controller.AdminLogoutHandler())
	admin.Use(middleware.AdminAuthMiddleware(db))

	// 	admin.GET("/orders", controller.AdminGetOrdersHandler())

}
