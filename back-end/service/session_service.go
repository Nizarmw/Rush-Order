package service

import (
	"RushOrder/session"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	Store      *sessions.CookieStore
	SessionKey string
)

const (
	SessionName   = "customer_session"
	SessionMaxAge = 3600 * 2
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Error loading .env file")
	}
	SessionKey = os.Getenv("SESSION_KEY")
}

func InitSessionStore(secretKey string) {
	Store = sessions.NewCookieStore([]byte(secretKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   SessionMaxAge,
		HttpOnly: true,
		Secure:   false,
	}
}

func CreateSession(w http.ResponseWriter, r *http.Request, data session.CustomerSession) error {
	sess, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
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
