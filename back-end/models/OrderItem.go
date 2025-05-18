package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	IDItem   uint      `gorm:"primaryKey;autoIncrement" json:"id_item"`
	IDOrder  uuid.UUID `gorm:"type:char(36)" json:"id_order"`
	IDProduk uuid.UUID `gorm:"type:char(36)" json:"id_produk"`
	Jumlah   int       `json:"jumlah"`
	Subtotal int       `json:"subtotal"`

	Order  Order  `gorm:"foreignKey:IDOrder;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Produk Produk `gorm:"foreignKey:IDProduk;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
