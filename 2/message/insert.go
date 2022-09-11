package message

import (
	"time"
)

type Insert struct {
	Timestamp time.Time
	Price     int32
}
