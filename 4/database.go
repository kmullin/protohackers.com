package main

import "sync"

type db struct {
	keys map[string]string
	mu   *sync.RWMutex
}

func NewDB() *db {
	return &db{make(map[string]string), new(sync.RWMutex)}
}

func (db *db) Insert(k, v string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.keys[k] = v
}

func (db *db) Retrieve(k string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	v, ok := db.keys[k]
	return v, ok
}
