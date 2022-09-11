package message

import (
	"errors"
	"time"
)

type Query struct {
	MinTime time.Time
	MaxTime time.Time
}

func (q *Query) UnmarshalBinary(data []byte) error {
	return errors.New("not implemented")
}
