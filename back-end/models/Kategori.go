package models

type KategoriProduk string

const (
	KategoriMakanan KategoriProduk = "makanan"
	KategoriMinuman KategoriProduk = "minuman"
	KategoriSnack   KategoriProduk = "snack"
)

func (k KategoriProduk) IsValid() bool {
	switch k {
	case KategoriMakanan, KategoriMinuman, KategoriSnack:
		return true
	}
	return false
}

func GetValidKategoriProduk() []KategoriProduk {
	return []KategoriProduk{
		KategoriMakanan,
		KategoriMinuman,
		KategoriSnack,
	}
}
