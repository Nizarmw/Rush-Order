package service

import (
	"RushOrder/config"
	"RushOrder/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/veritrans/go-midtrans"
)

// Payment model untuk menyimpan data pembayaran
type Payment struct {
	IDPayment     string    `gorm:"column:id_payment;primaryKey;size:50" json:"id_payment"`
	IDOrder       string    `gorm:"column:id_order;size:20;not null" json:"id_order"`
	Amount        int       `gorm:"column:amount;not null" json:"amount"`
	SnapToken     string    `gorm:"column:snap_token;size:255" json:"snap_token"`
	TransactionID string    `gorm:"column:transaction_id;size:100" json:"transaction_id"`
	Status        string    `gorm:"column:status;size:20;default:'pending'" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// Payment status constants
const (
	PaymentStatusPending = "pending"
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
	PaymentStatusCancel  = "cancel"
	PaymentStatusExpired = "expired"
)

// CreateSnapToken creates a Snap token for payment
func CreateSnapToken(orderID string) (string, error) {
	var order models.Order

	// Get order with items
	if err := config.DB.
		Preload("Items").
		Where("id_order = ?", orderID).
		First(&order).Error; err != nil {
		return "", fmt.Errorf("order not found: %v", err)
	}

	// Calculate total amount from order
	var totalAmount int = order.TotalHarga

	// Initialize Midtrans client
	midtransClient := midtrans.NewClient()
	midtransClient.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtransClient.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	midtransClient.APIEnvType = midtrans.Sandbox // Change to Production for live

	snapGateway := midtrans.SnapGateway{Client: midtransClient}

	// Create Snap request
	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(totalAmount),
		},
	}

	// Get Snap token
	snapResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", fmt.Errorf("failed to create snap token: %v", err)
	}

	// Save payment record
	payment := Payment{
		IDPayment: fmt.Sprintf("PAY%d", time.Now().Unix()),
		IDOrder:   orderID,
		Amount:    totalAmount,
		SnapToken: snapResp.Token,
		Status:    PaymentStatusPending,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return "", fmt.Errorf("failed to save payment record: %v", err)
	}

	return snapResp.Token, nil
}

// UpdatePaymentStatus updates payment and order status based on Midtrans notification
func UpdatePaymentStatus(orderID, transactionID, midtransStatus string) error {
	log.Printf("Starting UpdatePaymentStatus for orderID=%s, status=%s", orderID, midtransStatus)

	// Begin transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic recovered in UpdatePaymentStatus: %v", r)
		}
	}()

	// Map Midtrans status to our status
	var paymentStatus, statusCustomer, statusAdmin string
	switch midtransStatus {
	case "settlement", "capture":
		paymentStatus = PaymentStatusSuccess
		statusCustomer = "paid"
		statusAdmin = "process"
	case "cancel", "deny":
		paymentStatus = PaymentStatusCancel
		statusCustomer = "cancelled"
		statusAdmin = "cancelled"
	case "expire":
		paymentStatus = PaymentStatusExpired
		statusCustomer = "expired"
		statusAdmin = "cancelled"
	case "pending":
		paymentStatus = PaymentStatusPending
		statusCustomer = "pending"
		statusAdmin = "process"
	default:
		paymentStatus = PaymentStatusFailed
		statusCustomer = "failed"
		statusAdmin = "cancelled"
	}

	var payment Payment
	if err := tx.Where("id_order = ?", orderID).First(&payment).Error; err != nil {
		log.Printf("Payment not found for orderID=%s: %v", orderID, err)
		tx.Rollback()
		return fmt.Errorf("payment not found: %v", err)
	}

	log.Printf("Found payment for orderID=%s, updating status to %s", orderID, paymentStatus)
	if err := tx.Model(&payment).Updates(map[string]interface{}{
		"status":         paymentStatus,
		"transaction_id": transactionID,
	}).Error; err != nil {
		log.Printf("Error updating payment: %v", err)
		tx.Rollback()
		return fmt.Errorf("failed to update payment: %v", err)
	}

	var order models.Order
	if err := tx.Where("id_order = ?", orderID).First(&order).Error; err != nil {
		log.Printf("Order not found for orderID=%s: %v", orderID, err)
		tx.Rollback()
		return fmt.Errorf("order not found: %v", err)
	}

	log.Printf("Found order for orderID=%s, updating status to %s", orderID, statusCustomer)
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status_customer": statusCustomer,
		"status_admin":    statusAdmin,
	}).Error; err != nil {
		log.Printf("Error updating order status: %v", err)
		tx.Rollback()
		return fmt.Errorf("failed to update order: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully updated payment status for orderID=%s", orderID)
	return nil
}

func GetPaymentByOrderID(orderID string) (*Payment, error) {
	var payment Payment
	if err := config.DB.Where("id_order = ?", orderID).First(&payment).Error; err != nil {
		return nil, fmt.Errorf("payment not found: %v", err)
	}
	return &payment, nil
}

func CreatePaymentFromCart(w http.ResponseWriter, r *http.Request) (string, string, error) {
	orderID, err := CheckoutCart(w, r)
	if err != nil {
		return "", "", fmt.Errorf("failed to checkout cart: %v", err)
	}

	snapToken, err := CreateSnapToken(orderID)
	if err != nil {
		return "", "", fmt.Errorf("failed to create payment: %v", err)
	}

	return orderID, snapToken, nil
}

func GetOrderStatus(orderID string) (*models.Order, error) {
	var order models.Order
	if err := config.DB.
		Preload("Items").
		Where("id_order = ?", orderID).
		First(&order).Error; err != nil {
		return nil, fmt.Errorf("order not found: %v", err)
	}
	return &order, nil
}
