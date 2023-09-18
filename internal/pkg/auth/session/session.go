package session

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"web-studio-backend/internal/pkg/strhelp"
)

const Timeout = 24 * time.Hour

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	ID        string
	UserID    int16
	CSRFToken string
}

// sessions are used to map 'session_id' cookie to user ID.
var (
	sessions = map[string]*Session{}
	m        sync.RWMutex
)

func New(userID int16) (string, string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", "", fmt.Errorf("generating session id: %w", err)
	}

	csrfToken, err := generateCSRFToken()
	if err != nil {
		return "", "", fmt.Errorf("generating session id: %w", err)
	}

	m.Lock()
	defer m.Unlock()
	sessions[sessionID] = &Session{
		UserID:    userID,
		CSRFToken: csrfToken,
	}

	time.AfterFunc(Timeout, func() {
		Delete(sessionID)
	})

	return sessionID, csrfToken, nil
}

func generateSessionID() (string, error) {
	// Session identifiers should be at least 128 bits long to prevent brute-force session guessing attacks
	return strhelp.GenerateRandomString(32)
}

func generateCSRFToken() (string, error) {
	// Session identifiers should be at least 128 bits long to prevent brute-force session guessing attacks
	return strhelp.GenerateRandomString(32)
}

// GetSession returns session by sessionID.
func GetSession(sessionID string) (*Session, error) {
	m.RLock()
	defer m.RUnlock()

	session, ok := sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

// Delete deactivates the session.
func Delete(sessionID string) {
	m.Lock()
	defer m.Unlock()

	delete(sessions, sessionID)
}
