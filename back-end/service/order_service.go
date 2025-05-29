package service

import (
	"RushOrder/config" // Keep models import if other parts of the service use it directly
	"log"
)

// OrderItemDetail struct to include product details, now flattened
type OrderItemDetail struct {
	// Fields from order_items table
	IDItem   int    `json:"id_item"`
	IDOrder  string `json:"id_order"`
	IDProduk string `json:"id_produk"`
	Jumlah   int    `json:"jumlah"`   // Quantity of the item
	Subtotal int    `json:"subtotal"` // Subtotal for this item (Jumlah * HargaProduk)

	// Fields from produks table
	NamaProduk  string `json:"nama_produk"`  // Name of the product
	HargaProduk int    `json:"harga_produk"` // Unit price of the product
}

func GetOrderItems(orderID string) ([]OrderItemDetail, error) {
	var items []OrderItemDetail
	// Join OrderItem with Produk to get product name and price
	// The select statement maps directly to the fields in the flattened OrderItemDetail struct
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
			// Log should now show a flat structure for item
			log.Printf("Fetched item for order %s: %+v", orderID, item)
		}
	}

	return items, err
}
