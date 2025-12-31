package main

import "sync"

type AckTracker struct {
	m  map[int]int // session : last read pos
	mu sync.RWMutex
}

func newAckTracker() *AckTracker {
	return &AckTracker{
		m: make(map[int]int),
	}
}

func (at *AckTracker) Record(sessionId, pos int) {
	at.mu.Lock()
	defer at.mu.Unlock()
	at.m[sessionId] = pos
}

func (at *AckTracker) LastPos(sessionId int) int {
	at.mu.RLock()
	defer at.mu.RUnlock()
	return at.m[sessionId]
}
