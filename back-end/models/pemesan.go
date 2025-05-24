package models

type Pemesan struct {
	IDPemesan string `gorm:"column:id_pemesan;primaryKey;size:20" json:"id_pemesan"`
	Nama      string `gorm:"column:nama;size:100" json:"nama"`
	Meja      int    `gorm:"column:meja" json:"meja"`
}
