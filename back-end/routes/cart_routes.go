package routes

import (
	"RushOrder/controller"

	"github.com/gin-gonic/gin"
)

func SetupCartRoutes(r *gin.Engine) {
	CartRoutes := r.Group("/api/carts")
	{
		CartRoutes.GET("/", controller.GetCartHandler)
		CartRoutes.POST("/", controller.AddToCartHandler)
		CartRoutes.PUT("/", controller.UpdateCartItemHandler)
		CartRoutes.DELETE("/", controller.RemoveFromCartHandler)
		CartRoutes.DELETE("/clear", controller.ClearCartHandler)
	}
}
