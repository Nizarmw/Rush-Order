package service

import (
	"RushOrder/config"
	"log"
)

type OrderItemDetail struct {
	IDItem   int    `json:"id_item"`
	IDOrder  string `json:"id_order"`
	IDProduk string `json:"id_produk"`
	Jumlah   int    `json:"jumlah"`
	Subtotal int    `json:"subtotal"`

	NamaProduk  string `json:"nama_produk"`
	HargaProduk int    `json:"harga_produk"`
}

func GetOrderItems(orderID string) ([]OrderItemDetail, error) {
	var items []OrderItemDetail
	err := config.DB.Table("order_items").
		Select("order_items.id_item, order_items.id_order, order_items.id_produk, order_items.jumlah, order_items.subtotal, produks.nama_produk, produks.harga_produk").
		Joins("join produks on produks.id_produk = order_items.id_produk").
		Where("order_items.id_order = ?", orderID).
		Find(&items).Error

	if err != nil {
		log.Printf("Error fetching order items for order %s: %v", orderID, err)
		return nil, err
	}
	if len(items) == 0 {
		log.Printf("No items found for order %s", orderID)
	} else {
		for _, item := range items {
			log.Printf("Fetched item for order %s: %+v", orderID, item)
		}
	}

	return items, err
}
