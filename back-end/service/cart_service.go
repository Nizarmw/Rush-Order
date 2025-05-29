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

	sessionData, ok := sess.Values[SessionKey]
	if !ok {
		return errors.New("session tidak ditemukan")
	}

	var customer session.CustomerSession
	if err := json.Unmarshal([]byte(sessionData.(string)), &customer); err != nil {
		return err
	}

	if existingItem, exists := customer.Cart[item.IDProduk]; exists {
		existingItem.Jumlah += item.Jumlah
		existingItem.Subtotal = existingItem.Harga * existingItem.Jumlah
		customer.Cart[item.IDProduk] = existingItem
	} else {
		item.Subtotal = item.Harga * item.Jumlah
		customer.Cart[item.IDProduk] = item
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

	// Generate order ID
	orderID := fmt.Sprintf("ORD%d", time.Now().Unix())

	// Create order (tanpa Items)
	order := models.Order{
		IDOrder:    orderID,
		IDPemesan:  customer.ID,
		TotalHarga: customer.Total,
		// Field lain sesuai kebutuhan
	}

	// Save to database using transaction
	tx := config.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// Create and save order items satu per satu
	for _, item := range customer.Cart {
		orderItem := models.OrderItem{
			IDOrder:  orderID,
			IDProduk: item.IDProduk,
			Jumlah:   item.Jumlah,
			Subtotal: item.Subtotal,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return "", err
		}
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
