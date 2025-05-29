package service

import (
	"RushOrder/config"
	"RushOrder/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/veritrans/go-midtrans"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Error loading .env file")
	}

	if os.Getenv("MIDTRANS_SERVER_KEY") == "" || os.Getenv("MIDTRANS_CLIENT_KEY") == "" {
		panic("Midtrans environment variables are not set")
	}
}

func CreateSnapToken(orderID string) (string, error) {
	var order models.Order

	if err := config.DB.
		Preload("Items").
		Where("id_order = ?", orderID).
		First(&order).Error; err != nil {
		return "", fmt.Errorf("order not found: %v", err)
	}

	if order.StatusCustomer != models.CustomerStatusPending {
		return "", fmt.Errorf("payment already processed for this order")
	}

	midtransClient := midtrans.NewClient()
	midtransClient.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtransClient.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	midtransClient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{Client: midtransClient}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(order.TotalHarga),
		},
	}

	snapResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", fmt.Errorf("failed to create snap token: %v", err)
	}

	var existingPayment models.Payment
	if err := config.DB.Where("id_order = ?", orderID).First(&existingPayment).Error; err == nil {
		existingPayment.SnapToken = snapResp.Token
		if err := config.DB.Save(&existingPayment).Error; err != nil {
			return "", fmt.Errorf("failed to update payment record: %v", err)
		}
		return snapResp.Token, nil
	}

	payment := models.Payment{
		IDPayment: fmt.Sprintf("PAY%d", time.Now().Unix()),
		IDOrder:   orderID,
		Amount:    order.TotalHarga,
		SnapToken: snapResp.Token,
		Status:    models.PaymentStatusPending,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return "", fmt.Errorf("failed to save payment record: %v", err)
	}

	return snapResp.Token, nil
}

func UpdatePaymentStatus(orderID, transactionID, midtransStatus string) error {
	log.Printf("Starting UpdatePaymentStatus for orderID=%s, status=%s", orderID, midtransStatus)
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic recovered in UpdatePaymentStatus: %v", r)
		}
	}()

	paymentStatus, customerStatus, adminStatus := models.GetStatusFromMidtrans(midtransStatus)

	log.Printf("Status mapping - Payment: %s, Customer: %s, Admin: %s", paymentStatus, customerStatus, adminStatus)

	var payment models.Payment
	if err := tx.Where("id_order = ?", orderID).First(&payment).Error; err != nil {
		log.Printf("Payment not found for orderID=%s: %v", orderID, err)
		tx.Rollback()
		return fmt.Errorf("payment not found: %v", err)
	}

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

	updateData := map[string]interface{}{
		"status_customer": customerStatus,
	}

	if adminStatus != "" {
		updateData["status_admin"] = adminStatus
	} else {
		updateData["status_admin"] = ""
	}

	if err := tx.Model(&order).Updates(updateData).Error; err != nil {
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

func UpdateAdminStatus(orderID, newAdminStatus string) error {
	log.Printf("Updating admin status for orderID=%s to %s", orderID, newAdminStatus)

	if newAdminStatus != models.AdminStatusProcess && newAdminStatus != models.AdminStatusCompleted {
		return fmt.Errorf("invalid admin status: %s", newAdminStatus)
	}

	var order models.Order
	if err := config.DB.Where("id_order = ?", orderID).First(&order).Error; err != nil {
		return fmt.Errorf("order not found: %v", err)
	}

	if order.StatusCustomer != models.CustomerStatusSuccess {
		return fmt.Errorf("cannot update admin status: customer payment not completed")
	}

	if order.StatusAdmin == "" {
		return fmt.Errorf("cannot update admin status: order not in admin queue")
	}

	if err := config.DB.Model(&order).Update("status_admin", newAdminStatus).Error; err != nil {
		return fmt.Errorf("failed to update admin status: %v", err)
	}

	log.Printf("Successfully updated admin status for orderID=%s to %s", orderID, newAdminStatus)
	return nil
}

func GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
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

func GetAdminOrders(status string) ([]models.Order, error) {
	var orders []models.Order
	query := config.DB.Preload("Items").Where("status_customer = ?", models.CustomerStatusSuccess)

	if status != "" {
		query = query.Where("status_admin = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin orders: %v", err)
	}

	return orders, nil
}
