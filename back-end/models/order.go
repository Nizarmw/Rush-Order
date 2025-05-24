package models

type Order struct {
	IDOrder    string `gorm:"column:id_order;primaryKey;size:20" json:"id_order"`
	IDPemesan  string `gorm:"column:id_pemesan;size:20" json:"id_pemesan"`
	TotalHarga int    `gorm:"column:total_harga;default:0" json:"total_harga"`
	Status     bool   `gorm:"column:status;default:false" json:"status"`
}
