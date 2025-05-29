package models

import "time"

type Order struct {
	IDOrder        string    `gorm:"column:id_order;type:varchar(36);primaryKey" json:"id_order"`
	IDPemesan      string    `gorm:"column:id_pemesan;size:36" json:"id_pemesan"`
	TotalHarga     int       `gorm:"column:total_harga;default:0" json:"total_harga"`
	StatusCustomer string    `gorm:"column:status_customer;size:10;default:'pending'" json:"status_customer"`
	StatusAdmin    string    `gorm:"column:status_admin;size:10;default:''" json:"status_admin"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
