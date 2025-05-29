package controller

import (
	"RushOrder/models"
	"RushOrder/service"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

type AdminStatusRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

func CreatePaymentHandler(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	snapToken, err := service.CreateSnapToken(req.OrderID)
	if err != nil {
		log.Printf("Error creating payment: %v", err)
		if err.Error() == "payment already processed for this order" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pembayaran sudah diproses untuk order ini"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "payment created successfully",
		"order_id":   req.OrderID,
		"snap_token": snapToken,
	})
}

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

func MidtransWebhookHandler(c *gin.Context) {
	var notif map[string]interface{}
	if err := c.ShouldBindJSON(&notif); err != nil {
		log.Printf("Invalid webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid webhook payload"})
		return
	}

	orderID, ok1 := notif["order_id"].(string)
	transactionID, ok2 := notif["transaction_id"].(string)
	status, ok3 := notif["transaction_status"].(string)

	if !ok1 || !ok2 || !ok3 {
		log.Printf("Missing required fields in webhook: %+v", notif)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	log.Printf("Received webhook - OrderID: %s, TransactionID: %s, Status: %s", orderID, transactionID, status)

	if err := service.UpdatePaymentStatus(orderID, transactionID, status); err != nil {
		log.Printf("Error updating payment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update status pembayaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment status updated successfully"})
}

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

	items, _ := service.GetOrderItems(orderID)
	payment, err := service.GetPaymentByOrderID(orderID)
	response := gin.H{
		"id_order":        order.IDOrder,
		"id_pemesan":      order.IDPemesan,
		"total_harga":     order.TotalHarga,
		"status_customer": order.StatusCustomer,
		"status_admin":    order.StatusAdmin,
		"items":           items,
	}

	if err == nil {
		response["payment"] = gin.H{
			"id_payment":     payment.IDPayment,
			"amount":         payment.Amount,
			"snap_token":     payment.SnapToken,
			"transaction_id": payment.TransactionID,
			"status":         payment.Status,
			"created_at":     payment.CreatedAt,
		}
	} else {
		response["payment"] = nil
	}

	c.JSON(http.StatusOK, response)
}

func UpdateAdminStatusHandler(c *gin.Context) {
	var req AdminStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	if req.Status != models.AdminStatusProcess && req.Status != models.AdminStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin status"})
		return
	}

	if err := service.UpdateAdminStatus(req.OrderID, req.Status); err != nil {
		log.Printf("Error updating admin status: %v", err)
		switch err.Error() {
		case "cannot update admin status: customer payment not completed":
			c.JSON(http.StatusBadRequest, gin.H{"error": "customer belum menyelesaikan pembayaran"})
		case "cannot update admin status: order not in admin queue":
			c.JSON(http.StatusBadRequest, gin.H{"error": "order tidak ada dalam antrian admin"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update status admin"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "admin status updated successfully",
		"order_id": req.OrderID,
		"status":   req.Status,
	})
}

func SimulatePaymentSuccessHandler(c *gin.Context) {
	log.Println("--- SimulatePaymentSuccessHandler invoked (restored to full simulation) ---")
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
		return
	}

	transactionID := fmt.Sprintf("SIM_%d", time.Now().Unix())
	err := service.UpdatePaymentStatus(orderID, transactionID, "settlement")
	if err != nil {
		log.Printf("Error simulating payment success for order %s: %v", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simulasi pembayaran"})
		return
	}

	err = service.UpdateAdminStatus(orderID, models.AdminStatusProcess)
	if err != nil {
		log.Printf("Error updating admin status to 'process' after simulating payment for order %s: %v", orderID, err)
	}

	log.Printf("Payment simulated as 'settlement' and admin status set to 'process' for order %s.", orderID)
	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment simulated successfully, and order is set to 'process' for admin.",
		"order_id":       orderID,
		"transaction_id": transactionID,
		"status_after_simulation": gin.H{
			"customer": "success",
			"admin":    models.AdminStatusProcess,
		},
	})
}
