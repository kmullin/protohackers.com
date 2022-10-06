package database

import (
	"sync"
)

type Db struct {
	keys map[string]string
	mu   *sync.RWMutex
}

func NewDB() *Db {
	return &Db{make(map[string]string), new(sync.RWMutex)}
}

func (db *Db) Status() map[string]string {
	m := make(map[string]string)
	db.mu.RLock()
	defer db.mu.RUnlock()
	for k, v := range db.keys {
		m[k] = v
	}
	return m
}

func (db *Db) Insert(k, v string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.keys[k] = v
}

func (db *Db) Retrieve(k string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	v, ok := db.keys[k]
	return v, ok
}
