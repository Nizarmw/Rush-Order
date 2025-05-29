package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID       string    `gorm:"uniqueIndex;not null" json:"order_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	SnapToken     string    `gorm:"not null" json:"snap_token"`
	TransactionID string    `json:"transaction_id"`
	Status        string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Order Order `gorm:"foreignKey:OrderID;references:ID" json:"-"`
}
