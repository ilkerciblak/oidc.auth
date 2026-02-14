package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type StateManager struct {
	mu    *sync.RWMutex
	store map[string]stateEntry
}

func NewStateManager() *StateManager{
	return &StateManager{
		mu:    &sync.RWMutex{},
		store: map[string]stateEntry{},
	}
}

type stateEntry struct {
	expirety time.Time
}

func newStateEntry(duration time.Duration) stateEntry {
	return stateEntry{
		expirety: time.Now().Add(duration),
	}
}

func (s *StateManager) GenerateState() (string, error) {
	state, err := generateState()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[state] = newStateEntry(time.Minute * 5)

	return state, nil
}

func (s *StateManager) ValidateState(state_string string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.store[state_string]
	if !exists {
		return fmt.Errorf("[Invalid State]: state not found in store")
	}

	if entry.expirety.Before(time.Now()) {
		return fmt.Errorf("[Invalid State]: state has expired")
	}

	return nil
}

func generateState() (string, error) {
	randomBytes := make([]byte, 32)

	// generated random bytes using crypto/rand
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("[rand.Read]:  %v", err)
	}

	state := base64.URLEncoding.EncodeToString(randomBytes)
	return state, nil
}
func (sm *StateManager) Delete(state string) {
	sm.mu.Lock()
	delete(sm.store, state)
	sm.mu.Unlock()
}


