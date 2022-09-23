package chat

import "sync"

type Channel struct {
	Users []User
	mu    *sync.RWMutex
}
