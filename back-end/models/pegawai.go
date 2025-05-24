package models

type Pegawai struct {
	IDPegawai int    `gorm:"column:id_pegawai;primaryKey;autoIncrement" json:"id_pegawai"`
	Username  string `gorm:"column:username;size:50;unique;not null" json:"username"`
	Password  string `gorm:"column:password;size:255;not null" json:"-"`
}
