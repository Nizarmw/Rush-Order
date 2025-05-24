package models

type OrderItem struct {
	IDItem   int    `gorm:"column:id_item;primaryKey;autoIncrement" json:"id_item"`
	IDOrder  string `gorm:"column:id_order;size:20" json:"id_order"`
	IDProduk string `gorm:"column:id_produk;size:20" json:"id_produk"`
	Jumlah   int    `gorm:"column:jumlah" json:"jumlah"`
	Subtotal int    `gorm:"column:subtotal" json:"subtotal"`
}

func (OrderItem) TableName() string {
	return "Order_Item"
}
