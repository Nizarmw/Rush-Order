package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pemesan struct {
	IDPemesan uuid.UUID `gorm:"type:char(36);primaryKey" json:"id_pemesan"`
	Nama      string    `gorm:"type:varchar(100)" json:"nama"`
	Meja      int       `json:"meja"`
}

func (p *Pemesan) BeforeCreate(tx *gorm.DB) (err error) {
	p.IDPemesan = uuid.New()
	return
}
