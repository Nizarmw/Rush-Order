package models

type OrderItem struct {
	IDItem   int    `gorm:"column:id_item;primaryKey;autoIncrement" json:"id_item"`
	IDOrder  string `gorm:"column:id_order;type:varchar(20);not null" json:"id_order"`
	IDProduk string `gorm:"column:id_produk;type:varchar(20);not null" json:"id_produk"`
	Jumlah   int    `gorm:"column:jumlah" json:"jumlah"`
	Subtotal int    `gorm:"column:subtotal" json:"subtotal"`
	Order    Order  `gorm:"foreignKey:IDOrder;references:IDOrder" json:"-"`
}
