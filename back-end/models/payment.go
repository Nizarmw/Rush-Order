package models

import "time"

type Payment struct {
	IDPayment     string    `gorm:"column:id_payment;primaryKey;size:50" json:"id_payment"`
	IDOrder       string    `gorm:"column:id_order;size:20;not null;uniqueIndex" json:"id_order"`
	Amount        int       `gorm:"column:amount;not null" json:"amount"`
	SnapToken     string    `gorm:"column:snap_token;size:255" json:"snap_token"`
	TransactionID string    `gorm:"column:transaction_id;size:100" json:"transaction_id"`
	Status        string    `gorm:"column:status;size:20;default:'pending'" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Order Order `gorm:"foreignKey:IDOrder;references:IDOrder" json:"-"`
}

const (
	PaymentStatusPending = "pending"
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
	PaymentStatusExpired = "expired"
)

const (
	CustomerStatusPending = "pending"
	CustomerStatusSuccess = "success"
	CustomerStatusFailed  = "failed"
)

const (
	AdminStatusProcess   = "process"
	AdminStatusCompleted = "completed"
)

func GetStatusFromMidtrans(midstransStatus string) (paymentStatus, customerStatus, adminStatus string) {
	switch midstransStatus {
	case "settlement", "capture":
		return PaymentStatusSuccess, CustomerStatusSuccess, AdminStatusProcess
	case "pending":
		return PaymentStatusPending, CustomerStatusPending, ""
	case "cancel", "deny":
		return PaymentStatusFailed, CustomerStatusFailed, ""
	case "expire":
		return PaymentStatusExpired, CustomerStatusFailed, ""
	default:
		return PaymentStatusFailed, CustomerStatusFailed, ""
	}
}
