package service

import (
	"RushOrder/session"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	Store *sessions.CookieStore
)

const (
	SessionName   = "customer_session"
	SessionMaxAge = 3600 * 2
	SessionKey    = "customer_data"
)

func InitSessionStore(secretKey string) {
	Store = sessions.NewCookieStore([]byte(secretKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   SessionMaxAge,
		HttpOnly: true,
		Secure:   false,
	}
}

func CreateSession(w http.ResponseWriter, r *http.Request, id, nama string, meja int) error {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	customerData := session.CustomerSession{
		ID:    id,
		Nama:  nama,
		Meja:  meja,
		Cart:  make(map[string]session.CartItem),
		Total: 0,
	}

	jsonData, err := json.Marshal(customerData)
	if err != nil {
		return err
	}

	sess.Values[SessionKey] = string(jsonData)
	return sess.Save(r, w)
}

func GetSession(r *http.Request) (*session.CustomerSession, error) {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}

	sessionData, ok := sess.Values[SessionKey]
	if !ok {
		return nil, nil
	}

	var customerData session.CustomerSession
	err = json.Unmarshal([]byte(sessionData.(string)), &customerData)
	if err != nil {
		return nil, err
	}

	return &customerData, nil
}
func ClearCustomerSession(w http.ResponseWriter, r *http.Request) error {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	delete(sess.Values, SessionKey)
	return sess.Save(r, w)
}
