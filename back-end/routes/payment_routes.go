package routes

import (
	"RushOrder/controller"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(r *gin.Engine) {
	payment := r.Group("/api/payment")
	{
		payment.POST("/", controller.CreatePaymentHandler)
		payment.POST("/webhook", controller.MidtransWebhookHandler)
		payment.GET("/:order_id", controller.GetPaymentHandler)
		payment.POST("/checkout", controller.CheckoutAndPayHandler)
		payment.POST("/simulate/:order_id", controller.SimulatePaymentSuccessHandler)
	}

	order := r.Group("/api/order")
	{
		order.GET("/:order_id/status", controller.GetOrderStatusHandler)
	}
}
