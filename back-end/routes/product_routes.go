package routes

import (
	"RushOrder/controller"

	"github.com/gin-gonic/gin"
)

func SetupProdukRoutes(r *gin.Engine) {
	produkGroup := r.Group("/api/produk")
	{
		produkGroup.POST("/", controller.CreateProdukHandler)
		produkGroup.GET("/", controller.GetProduk)
		produkGroup.GET("/:id", controller.GetProdukByID)
		produkGroup.PUT("/:id", controller.UpdateProduk)
		produkGroup.DELETE("/:id", controller.DeleteProduk)

		produkGroup.GET("/search", controller.SearchProduk)
	}
}
