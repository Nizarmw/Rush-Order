package session

import "time"

type CustomerSession struct {
	ID        string              `json:"id"`
	Nama      string              `json:"nama"`
	Meja      int                 `json:"meja"`
	Cart      map[string]CartItem `json:"cart"`
	Total     int                 `json:"total"`
	CreatedAt time.Time           `json:"created_at"`
	ExpiredAt time.Time           `json:"expired_at"`
}

type CartItem struct {
	IDProduk   string `json:"id_produk"`
	NamaProduk string `json:"nama_produk"`
	Jumlah     int    `json:"jumlah"`
	Harga      int    `json:"harga"`
	Subtotal   int    `json:"subtotal"`
}
