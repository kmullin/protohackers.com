package message

import (
	"time"
)

type insert struct {
	Timestamp time.Time
	Price     int32
}

func (i insert) Type() msgType {
	return InsertType
}
