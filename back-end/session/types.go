package session

type CustomerSession struct {
	Nama string              `json:"nama"`
	Meja int                 `json:"meja"`
	Cart map[string]CartItem `json:"cart"`
}

type CartItem struct {
	IDProduk string `json:"id_produk"`
	Jumlah   int    `json:"jumlah"`
	Harga    int    `json:"harga"`
}
