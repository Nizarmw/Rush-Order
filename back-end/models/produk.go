package models

type Produk struct {
	IDProduk    string         `gorm:"column:id_produk;primaryKey;size:20" json:"id_produk"`
	NamaProduk  string         `gorm:"column:nama_produk;size:100" json:"nama_produk"`
	Deskripsi   string         `gorm:"column:deskripsi;size:255" json:"deskripsi"`
	HargaProduk int            `gorm:"column:harga_produk" json:"harga_produk"`
	ImageURL    string         `gorm:"type:text" json:"image_url"`
	Kategori    KategoriProduk `gorm:"column:kategori;size:20;not null" json:"kategori"`
}
