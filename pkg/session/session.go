package session

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type Session struct {
	createdAt    time.Time
	lastActivity time.Time
	id           string
	mu           sync.RWMutex // Add mutex to protect data
	data         map[string]any
}

// This interface can be used by any session store like redis, database, etc. to work with sessions
type SessionStore interface {
	Read(id string) (*Session, error)
	Write(session *Session) error
	Destroy(id string) error
	GC(idleExpiration, absoluteExpiration time.Duration) error
}

// should be instantiated on application startup and read configuration options from config file
type SessionManager struct {
	store              SessionStore
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
	cookieName         string
}

// Generate new session id
func generateSessionId() (string, error) {
	id := make([]byte, 32)

	// Generate random bytes inside the buffer
	_, err := rand.Read(id)

	if err != nil {
		return "", err
	}

	// Encode to base64
	return base64.URLEncoding.EncodeToString(id), nil
}

func NewSession() (*Session, error) {
	id, err := generateSessionId()
	if err != nil {
		return nil, err
	}

	return &Session{
		id:           id,
		data:         make(map[string]any),
		createdAt:    time.Now(),
		lastActivity: time.Now(),
	}, nil
}

func (s *Session) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Session) Get(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key]
}

func (s *Session) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (s *Session) UpdateLastActivity() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastActivity = time.Now()
}
