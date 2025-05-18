package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pegawai struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"not null" json:"password"`
	Role     string    `gorm:"default:admin" json:"role"`
}

func (u *Pegawai) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
