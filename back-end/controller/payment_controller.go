package controller

import (
	"RushOrder/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaymentRequest untuk request pembayaran berdasarkan order ID
type PaymentRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

// CreatePaymentHandler creates payment for existing order
func CreatePaymentHandler(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	snapToken, err := service.CreateSnapToken(req.OrderID)
	if err != nil {
		log.Printf("Error creating payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "payment created successfully",
		"order_id":   req.OrderID,
		"snap_token": snapToken,
	})
}

// CheckoutAndPayHandler handles checkout cart and create payment in one step
func CheckoutAndPayHandler(c *gin.Context) {
	orderID, snapToken, err := service.CreatePaymentFromCart(c.Writer, c.Request)
	if err != nil {
		log.Printf("Error in checkout and pay: %v", err)
		switch err.Error() {
		case "failed to checkout cart: session tidak ditemukan":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session tidak ditemukan"})
		case "failed to checkout cart: keranjang kosong":
			c.JSON(http.StatusBadRequest, gin.H{"error": "keranjang kosong"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat pembayaran"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "checkout and payment created successfully",
		"order_id":   orderID,
		"snap_token": snapToken,
	})
}

// MidtransWebhookHandler handles Midtrans webhook notification
func MidtransWebhookHandler(c *gin.Context) {
	var notif map[string]interface{}
	if err := c.ShouldBindJSON(&notif); err != nil {
		log.Printf("Invalid webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook payload"})
		return
	}

	// Extract required fields from notification
	orderID, ok1 := notif["order_id"].(string)
	transactionID, ok2 := notif["transaction_id"].(string)
	status, ok3 := notif["transaction_status"].(string)

	if !ok1 || !ok2 || !ok3 {
		log.Printf("Missing required fields in webhook: %+v", notif)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	log.Printf("Received webhook - OrderID: %s, TransactionID: %s, Status: %s", orderID, transactionID, status)

	// Update payment status
	if err := service.UpdatePaymentStatus(orderID, transactionID, status); err != nil {
		log.Printf("Error updating payment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update status pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment status updated successfully"})
}

// GetPaymentHandler retrieves payment information by order ID
func GetPaymentHandler(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
		return
	}

	payment, err := service.GetPaymentByOrderID(orderID)
	if err != nil {
		log.Printf("Payment not found for order %s: %v", orderID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id_payment":     payment.IDPayment,
		"id_order":       payment.IDOrder,
		"amount":         payment.Amount,
		"snap_token":     payment.SnapToken,
		"transaction_id": payment.TransactionID,
		"status":         payment.Status,
		"created_at":     payment.CreatedAt,
		"updated_at":     payment.UpdatedAt,
	})
}

// GetOrderStatusHandler retrieves order status with payment info
func GetOrderStatusHandler(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
		return
	}

	order, err := service.GetOrderStatus(orderID)
	if err != nil {
		log.Printf("Order not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Get payment info
	payment, err := service.GetPaymentByOrderID(orderID)
	if err != nil {
		log.Printf("Payment not found for order %s: %v", orderID, err)
		// Return order without payment info
		c.JSON(http.StatusOK, gin.H{
			"id_order":        order.IDOrder,
			"id_pemesan":      order.IDPemesan,
			"total_harga":     order.TotalHarga,
			"status_customer": order.StatusCustomer,
			"status_admin":    order.StatusAdmin,
			"items":           order.Items,
			"payment":         nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id_order":        order.IDOrder,
		"id_pemesan":      order.IDPemesan,
		"total_harga":     order.TotalHarga,
		"status_customer": order.StatusCustomer,
		"status_admin":    order.StatusAdmin,
		"items":           order.Items,
		"payment": gin.H{
			"id_payment":     payment.IDPayment,
			"amount":         payment.Amount,
			"snap_token":     payment.SnapToken,
			"transaction_id": payment.TransactionID,
			"status":         payment.Status,
			"created_at":     payment.CreatedAt,
		},
	})
}
