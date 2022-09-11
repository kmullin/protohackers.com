package message

import (
	"errors"
	"time"
)

type Insert struct {
	Timestamp time.Time
	Price     int32
}

func (i *Insert) UnmarshalBinary(data []byte) error {
	return errors.New("not implemented")
}
