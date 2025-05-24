package service

import (
	"RushOrder/session"
	"errors"
	"sync"
	"time"
)

type SessionManager struct {
	session map[string]*session.CustomerSession
	mutex   sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		session: make(map[string]*session.CustomerSession),
	}
}

func (s *SessionManager) CreateSession(id, nama string, meja int) *session.CustomerSession {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newSession := &session.CustomerSession{
		ID:        id,
		Nama:      nama,
		Meja:      meja,
		Cart:      make(map[string]session.CartItem),
		Total:     0,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}
	s.session[id] = newSession

	return newSession
}

func (s *SessionManager) GetSession(id string) (*session.CustomerSession, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sess, ok := s.session[id]
	if !ok {
		return nil, false
	}
	if time.Now().After(sess.ExpiredAt) {
		delete(s.session, id)
		return nil, false
	}
	return sess, true
}

func (s *SessionManager) UpdateSession(id, nama string, meja int) (*session.CustomerSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sess, ok := s.session[id]
	if !ok {
		return nil, errors.New("session not found")
	}
	sess.Nama = nama
	sess.Meja = meja
	return sess, nil
}

func (s *SessionManager) DeleteSession(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.session, id)
	return nil
}
