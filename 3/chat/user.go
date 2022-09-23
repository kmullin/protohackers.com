package chat

import "time"

type User struct {
	Name     string
	lastSeen time.Time
}
