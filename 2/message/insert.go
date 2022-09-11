package message

import (
	"time"
)

type Insert struct {
	Timestamp time.Time
	Price     int32
}

func (i Insert) Type() Type {
	return InsertType
}
