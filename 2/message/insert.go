package message

import (
	"errors"
	"time"
)

type Insert struct {
	Timestamp time.Time
	Price     int32
}

func (i *Insert) MarshalBinary() ([]byte, error) {
	return nil, errors.New("not implemented")
}
