package routes

import (
	"RushOrder/controller"

	"github.com/gin-gonic/gin"
)

func SetupCartRoutes(r *gin.Engine) {
	CartRoutes := r.Group("/api/carts")
	{
		CartRoutes.GET("/", controller.GetCartHandler)
		CartRoutes.POST("/:idpemesan", controller.AddToCartHandler)
		CartRoutes.DELETE("/clear", controller.ClearCartHandler)
		CartRoutes.DELETE("/:id", controller.RemoveFromCartHandler)
	}
}
