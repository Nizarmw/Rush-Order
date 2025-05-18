package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	IDOrder    uuid.UUID `gorm:"type:char(36);primaryKey" json:"id_order"`
	IDPemesan  uuid.UUID `gorm:"type:char(36)" json:"id_pemesan"`
	TotalHarga int       `gorm:"default:0" json:"total_harga"`
	Status     bool      `gorm:"default:false" json:"status"`

	Pemesan Pemesan `gorm:"foreignKey:IDPemesan;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.IDOrder = uuid.New()
	return
}
