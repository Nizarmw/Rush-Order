package models

type Produk struct {
	IDProduk    string `gorm:"column:id_produk;primaryKey;size:20" json:"id_produk"`
	NamaProduk  string `gorm:"column:nama_produk;size:100" json:"nama_produk"`
	HargaProduk int    `gorm:"column:harga_produk" json:"harga_produk"`
}
