package routes

import (
	"RushOrder/controller"
	"log" // Import log

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(r *gin.Engine) {
	log.Println("--- RegisterPaymentRoutes: START ---")
	payment := r.Group("/api/payment")
	{
		log.Println("--- RegisterPaymentRoutes: Defining /api/payment/ ---")
		payment.POST("/", controller.CreatePaymentHandler)

		log.Println("--- RegisterPaymentRoutes: Defining /api/payment/webhook ---")
		payment.POST("/webhook", controller.MidtransWebhookHandler)

		log.Println("--- RegisterPaymentRoutes: Defining /api/payment/:order_id ---")
		payment.GET("/:order_id", controller.GetPaymentHandler)

		log.Println("--- RegisterPaymentRoutes: Defining /api/payment/checkout ---")
		payment.POST("/checkout", controller.CheckoutAndPayHandler)

		log.Println("--- RegisterPaymentRoutes: Defining /api/payment/simulate/:order_id ---")
		payment.POST("/simulate/:order_id", controller.SimulatePaymentSuccessHandler)
		log.Println("--- RegisterPaymentRoutes: DEFINED /api/payment/simulate/:order_id ---")
	}

	order := r.Group("/api/order")
	{
		log.Println("--- RegisterPaymentRoutes: Defining /api/order/:order_id/status ---")
		order.GET("/:order_id/status", controller.GetOrderStatusHandler)
	}
	log.Println("--- RegisterPaymentRoutes: END ---")
}
