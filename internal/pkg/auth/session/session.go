package session

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"web-studio-backend/internal/pkg/strhelp"
)

const TTL = 24 * time.Hour

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	UserID    int32
	CSRFToken string
}

// sessions are used to map 'session_id' cookie to user ID.
var (
	sessions = map[string]Session{}
	m        sync.RWMutex
)

func New(userID int32) (string, string, error) {
	m.Lock()
	defer m.Unlock()

	var sessionID string

	for {
		sid, err := generateSessionID()
		if err != nil {
			return "", "", fmt.Errorf("generating session id: %w", err)
		}

		if _, ok := sessions[sid]; !ok {
			sessionID = sid
			break
		}
	}

	csrfToken, err := generateCSRFToken()
	if err != nil {
		return "", "", fmt.Errorf("generating session id: %w", err)
	}

	sessions[sessionID] = Session{
		UserID:    userID,
		CSRFToken: csrfToken,
	}

	time.AfterFunc(TTL, func() {
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

	return &session, nil
}

// Delete deactivates the session.
func Delete(sessionID string) {
	m.Lock()
	defer m.Unlock()

	delete(sessions, sessionID)
}

// DeleteUserSession deletes session for given user.
func DeleteUserSession(userID int32) {
	m.Lock()
	defer m.Unlock()

	var sessionID string

	for sid, sess := range sessions {
		if sess.UserID == userID {
			sessionID = sid
			break
		}
	}

	delete(sessions, sessionID)
}
