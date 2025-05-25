package service

import (
	"RushOrder/session"
	"encoding/json"
	"errors"
	"net/http"
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
