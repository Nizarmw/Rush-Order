package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Produk struct {
	IDProduk    uuid.UUID `gorm:"type:char(36);primaryKey" json:"id_produk"`
	NamaProduk  string    `gorm:"type:varchar(100)" json:"nama_produk"`
	HargaProduk int       `json:"harga_produk"`
}

func (p *Produk) BeforeCreate(tx *gorm.DB) (err error) {
	p.IDProduk = uuid.New()
	return
}
