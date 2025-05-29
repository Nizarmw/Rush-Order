package service

import (
	"RushOrder/config"
	"RushOrder/models"
	"RushOrder/session"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func AddToCart(w http.ResponseWriter, r *http.Request, item session.CartItem) error {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	var customer session.CustomerSession

	// Coba ambil data session
	sessionData, ok := sess.Values[SessionKey]
	if ok {
		// Jika ada, decode JSON ke struct
		if err := json.Unmarshal([]byte(sessionData.(string)), &customer); err != nil {
			return err
		}
	} else {
		// Jika belum ada, buat struct kosong
		customer = session.CustomerSession{
			Cart:  make(map[string]session.CartItem),
			Total: 0,
		}
	}

	// Update atau tambah item
	if existingItem, exists := customer.Cart[item.IDProduk]; exists {
		existingItem.Jumlah += item.Jumlah
		existingItem.Subtotal = existingItem.Harga * existingItem.Jumlah
		customer.Cart[item.IDProduk] = existingItem
	} else {
		item.Subtotal = item.Harga * item.Jumlah
		customer.Cart[item.IDProduk] = item
	}

	// Hitung ulang total
	total := 0
	for _, val := range customer.Cart {
		total += val.Subtotal
	}
	customer.Total = total

	// Simpan kembali ke session
	jsonData, err := json.Marshal(customer)
	if err != nil {
		return err
	}
	sess.Values[SessionKey] = string(jsonData)
	return sess.Save(r, w)
}

func GetCart(r *http.Request) (map[string]session.CartItem, int, error) {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return nil, 0, err
	}

	sessionData, ok := sess.Values[SessionKey]
	if !ok {
		return nil, 0, errors.New("session tidak ditemukan")
	}

	var customer session.CustomerSession
	if err := json.Unmarshal([]byte(sessionData.(string)), &customer); err != nil {
		return nil, 0, err
	}

	return customer.Cart, customer.Total, nil
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request, idProduk string) error {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	sessionData, ok := sess.Values[SessionKey]
	if !ok {
		return errors.New("session tidak ditemukan")
	}

	var customer session.CustomerSession
	if err := json.Unmarshal([]byte(sessionData.(string)), &customer); err != nil {
		return err
	}

	delete(customer.Cart, idProduk)

	total := 0
	for _, val := range customer.Cart {
		total += val.Subtotal
	}
	customer.Total = total

	jsonData, err := json.Marshal(customer)
	if err != nil {
		return err
	}
	sess.Values[SessionKey] = string(jsonData)
	return sess.Save(r, w)
}

func ClearCart(w http.ResponseWriter, r *http.Request) error {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	sessionData, ok := sess.Values[SessionKey]
	if !ok {
		return errors.New("session tidak ditemukan")
	}

	var customer session.CustomerSession
	if err := json.Unmarshal([]byte(sessionData.(string)), &customer); err != nil {
		return err
	}

	customer.Cart = make(map[string]session.CartItem)
	customer.Total = 0

	jsonData, err := json.Marshal(customer)
	if err != nil {
		return err
	}
	sess.Values[SessionKey] = string(jsonData)
	return sess.Save(r, w)
}

func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request, idProduk string, jumlah int) error {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	sessionData, ok := sess.Values[SessionKey]
	if !ok {
		return errors.New("session tidak ditemukan")
	}

	var customer session.CustomerSession
	if err := json.Unmarshal([]byte(sessionData.(string)), &customer); err != nil {
		return err
	}

	item, exists := customer.Cart[idProduk]
	if !exists {
		return errors.New("item tidak ditemukan di keranjang")
	}

	if jumlah <= 0 {
		delete(customer.Cart, idProduk)
	} else {
		item.Jumlah = jumlah
		item.Subtotal = item.Harga * item.Jumlah
		customer.Cart[idProduk] = item
	}

	total := 0
	for _, val := range customer.Cart {
		total += val.Subtotal
	}
	customer.Total = total
	jsonData, err := json.Marshal(customer)
	if err != nil {
		return err
	}
	sess.Values[SessionKey] = string(jsonData)
	return sess.Save(r, w)
}

// CheckoutCart converts the cart items to an order and saves it to the database
func CheckoutCart(w http.ResponseWriter, r *http.Request) (string, error) {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return "", err
	}

	sessionData, ok := sess.Values[SessionKey]
	if !ok {
		return "", errors.New("session tidak ditemukan")
	}

	var customer session.CustomerSession
	if err := json.Unmarshal([]byte(sessionData.(string)), &customer); err != nil {
		return "", err
	}

	// Check if cart is empty
	if len(customer.Cart) == 0 {
		return "", errors.New("keranjang kosong")
	}

	// Generate order ID (you can replace this with your own ID generation logic)
	orderID := fmt.Sprintf("ORD%d", time.Now().Unix())

	// Create order
	order := models.Order{
		IDOrder:    orderID,
		IDPemesan:  customer.ID,
		TotalHarga: customer.Total,

		Items: []models.OrderItem{},
	}

	// Create order items
	for _, item := range customer.Cart {
		orderItem := models.OrderItem{
			IDOrder:  orderID,
			IDProduk: item.IDProduk,
			Jumlah:   item.Jumlah,
			Subtotal: item.Subtotal,
		}
		order.Items = append(order.Items, orderItem)
	}

	// Save to database using transaction
	tx := config.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// Clear the cart after successful checkout
	customer.Cart = make(map[string]session.CartItem)
	customer.Total = 0

	jsonData, err := json.Marshal(customer)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Save the updated session
	sess.Values[SessionKey] = string(jsonData)
	if err := sess.Save(r, w); err != nil {
		tx.Rollback()
		return "", err
	}

	// Commit transaction
	tx.Commit()

	return orderID, nil
}
