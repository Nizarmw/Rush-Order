package session

import (
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

func CreateSession(w http.ResponseWriter, r *http.Request, nama string, meja int) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	customerData := CustomerSession{
		Nama: nama,
		Meja: meja,
		Cart: make(map[string]CartItem),
	}

	jsonData, err := json.Marshal(customerData)
	if err != nil {
		return err
	}
	session.Values[SessionKey] = jsonData
	return session.Save(r, w)
}

func GetSession(r *http.Request) (*CustomerSession, error) {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}

	sessionData, ok := session.Values[SessionKey]

	if !ok {
		return nil, nil
	}
	var customerData CustomerSession
	err = json.Unmarshal([]byte(sessionData.(string)), &customerData)
	if err != nil {
		return nil, err
	}
	return &customerData, nil
}

func ClearCustomerSession(r *http.Request, w http.ResponseWriter) error {
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	delete(session.Values, SessionKey)
	return session.Save(r, w)
}
