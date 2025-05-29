package models

type Order struct {
	IDOrder        string      `gorm:"column:id_order;primaryKey;size:20" json:"id_order"`
	IDPemesan      string      `gorm:"column:id_pemesan;size:20" json:"id_pemesan"`
	TotalHarga     int         `gorm:"column:total_harga;default:0" json:"total_harga"`
	StatusCustomer string      `gorm:"column:status_customer;size:10;default:'pending'" json:"status_customer"`
	StatusAdmin    string      `gorm:"column:status_admin;size:10;default:'process'" json:"status_admin"`
	Items          []OrderItem `gorm:"foreignKey:IDOrder;references:IDOrder" json:"items"`
}
