package message

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

var Unknown = errors.New("unknown message")

type Query struct {
	MinTime time.Time
	MaxTime time.Time
}

type Insert struct {
	Timestamp time.Time
	Price     int32
}

type Message interface{}

func New(r io.Reader) (Message, error) {
	var m clientMessage
	err := binary.Read(r, ByteOrder, &m)
	if err != nil {
		return nil, err
	}

	switch m.Type {
	case insertByte:
		return Insert{
			Timestamp: unixTime(m.N1),
			Price:     m.N2,
		}, nil
	case queryByte:
		return Query{
			MinTime: unixTime(m.N1),
			MaxTime: unixTime(m.N2),
		}, nil
	default:
		return nil, Unknown
	}
}

func unixTime(i int32) time.Time {
	return time.Unix(int64(i), 0).UTC()
}
