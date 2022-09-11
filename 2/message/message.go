package message

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

var Unknown = errors.New("unknown message")

type msgType int

const (
	QueryType msgType = iota
	InsertType
)

type Message interface {
	Type() msgType
}

func New(r io.Reader) (Message, error) {
	var m clientMessage
	err := binary.Read(r, byteOrder, &m)
	if err != nil {
		return nil, err
	}

	switch m.Type {
	case insertByte:
		return insert{
			Timestamp: unixTime(m.Timestamp),
			Price:     m.Price,
		}, nil
	case queryByte:
		return query{
			MinTime: unixTime(m.Timestamp),
			MaxTime: unixTime(m.Price),
		}, nil
	default:
		return nil, Unknown
	}
}

func unixTime(i int32) time.Time {
	return time.Unix(int64(i), 0).UTC()
}
