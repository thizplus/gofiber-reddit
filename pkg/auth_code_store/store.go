package auth_code_store

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"gofiber-template/domain/dto"
)

// AuthCodeData stores the data associated with an authorization code
type AuthCodeData struct {
	Token     string
	User      dto.UserResponse
	IsNewUser bool
	State     string
	ExpiresAt time.Time
}

// Store manages authorization codes
type Store struct {
	mu    sync.RWMutex
	codes map[string]*AuthCodeData
}

var (
	instance *Store
	once     sync.Once
)

// GetInstance returns singleton instance of the store
func GetInstance() *Store {
	once.Do(func() {
		instance = &Store{
			codes: make(map[string]*AuthCodeData),
		}
		// Start cleanup goroutine
		go instance.cleanupExpiredCodes()
	})
	return instance
}

// GenerateCode creates a new authorization code and stores the data
func (s *Store) GenerateCode(token string, user dto.UserResponse, isNewUser bool, state string) (string, error) {
	// Generate random code
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := base64.URLEncoding.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Store code with 5 minute expiration
	s.codes[code] = &AuthCodeData{
		Token:     token,
		User:      user,
		IsNewUser: isNewUser,
		State:     state,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return code, nil
}

// ExchangeCode retrieves and deletes the data for a given code
func (s *Store) ExchangeCode(code string, state string) (*AuthCodeData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.codes[code]
	if !exists {
		return nil, false
	}

	// Check if code is expired
	if time.Now().After(data.ExpiresAt) {
		delete(s.codes, code)
		return nil, false
	}

	// Validate state if provided
	if state != "" && data.State != state {
		return nil, false
	}

	// Delete code after use (one-time use)
	delete(s.codes, code)

	return data, true
}

// cleanupExpiredCodes removes expired codes every minute
func (s *Store) cleanupExpiredCodes() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for code, data := range s.codes {
			if now.After(data.ExpiresAt) {
				delete(s.codes, code)
			}
		}
		s.mu.Unlock()
	}
}
