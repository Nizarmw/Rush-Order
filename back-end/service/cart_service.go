package service

import (
	"RushOrder/session"
	"errors"
)

func (s *SessionManager) AddToCart(sessionID, IDProduk, NamaProduk string, Jumlah, Harga int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sess, ok := s.session[sessionID]
	if !ok {
		return errors.New("session not found")
	}

	sess.Cart[IDProduk] = session.CartItem{
		IDProduk:   IDProduk,
		NamaProduk: NamaProduk,
		Jumlah:     Jumlah,
		Harga:      Harga,
		Subtotal:   Jumlah * Harga,
	}

	total := 0
	for _, item := range sess.Cart {
		total += item.Subtotal
	}
	sess.Total = total

	return nil
}

func (s *SessionManager) UpdateCartItem(sessionID, IDProduk string, Jumlah int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sess, ok := s.session[sessionID]
	if !ok {
		return errors.New("session not found")
	}

	item, exists := sess.Cart[IDProduk]
	if !exists {
		return errors.New("item not found in cart")
	}
	item.Jumlah = Jumlah
	item.Subtotal = Jumlah * item.Harga
	sess.Cart[IDProduk] = item

	total := 0
	for _, cartItem := range sess.Cart {
		total += cartItem.Subtotal
	}
	sess.Total = total
	return nil
}

func (s *SessionManager) RemoveFromCart(sessionID, IDProduk string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sess, ok := s.session[sessionID]
	if !ok {
		return errors.New("session not found")
	}

	delete(sess.Cart, IDProduk)

	total := 0
	for _, item := range sess.Cart {
		total += item.Subtotal
	}
	sess.Total = total

	return nil
}

func (s *SessionManager) ClearCart(sessionID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sess, ok := s.session[sessionID]
	if !ok {
		return errors.New("session not found")
	}

	sess.Cart = make(map[string]session.CartItem)
	sess.Total = 0

	return nil
}
