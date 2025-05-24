package session

import (
	"encoding/json"
	"net/http"
)

func UpdateSessionCart(r *http.Request, item CartItem) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	sessionData, ok := session.Values[SessionKey]
	if !ok {
		return nil
	}

	var customerData CustomerSession
	err = json.Unmarshal([]byte(sessionData.(string)), &customerData)
	if err != nil {
		return err
	}

	customerData.Cart[item.IDProduk] = item

	jsonData, err := json.Marshal(customerData)
	if err != nil {
		return err
	}
	session.Values[SessionKey] = jsonData
	return session.Save(r, nil)
}

func UpdateSessionCartItem(r *http.Request, w http.ResponseWriter, item CartItem) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	sessionData, ok := session.Values[SessionKey]
	if !ok {
		return nil
	}

	var customerData CustomerSession
	err = json.Unmarshal([]byte(sessionData.(string)), &customerData)
	if err != nil {
		return err
	}

	if existingItem, exists := customerData.Cart[item.IDProduk]; exists {
		existingItem.Jumlah += item.Jumlah
		customerData.Cart[item.IDProduk] = existingItem
	} else {
		customerData.Cart[item.IDProduk] = item
	}

	jsonData, err := json.Marshal(customerData)
	if err != nil {
		return err
	}
	session.Values[SessionKey] = jsonData
	return session.Save(r, w)
}

func RemoveFromCart(r *http.Request, w http.ResponseWriter, idProduk string) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	sessionData, ok := session.Values[SessionKey]
	if !ok {
		return nil
	}

	var customerData CustomerSession
	err = json.Unmarshal([]byte(sessionData.(string)), &customerData)
	if err != nil {
		return err
	}

	delete(customerData.Cart, idProduk)

	jsonData, err := json.Marshal(customerData)
	if err != nil {
		return err
	}
	session.Values[SessionKey] = jsonData
	return session.Save(r, w)
}
func GetCartTotal(r *http.Request) (int, error) {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return 0, err
	}

	sessionData, ok := session.Values[SessionKey]
	if !ok {
		return 0, nil
	}

	var customerData CustomerSession
	err = json.Unmarshal([]byte(sessionData.(string)), &customerData)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, item := range customerData.Cart {
		total += item.Jumlah * item.Harga
	}
	return total, nil
}
