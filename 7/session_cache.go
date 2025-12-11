package main

import (
	"sync"
)

type SessionCache struct {
	c  map[int]*Session
	mu sync.RWMutex
}

func NewSessionCache() *SessionCache {
	return &SessionCache{c: make(map[int]*Session)}
}

// Add adds a new session to the cache
func (sc *SessionCache) Add(s *Session) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.c[s.ID] = s
}

// Expire removes a session from the cache
func (sc *SessionCache) Expire(s *Session) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.c, s.ID)
}

// Get retrieves a session from the cache or returns nil if a session doesn't exist
func (sc *SessionCache) Get(id int) *Session {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	s := sc.c[id]
	return s
}
