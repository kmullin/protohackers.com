package message

import (
	"bytes"
	"encoding"
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

func (q Query) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	cm := clientMessage{queryByte, int32(q.MinTime.Unix()), int32(q.MaxTime.Unix())}
	err := binary.Write(&buf, ByteOrder, &cm)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i Insert) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	cm := clientMessage{insertByte, int32(i.Timestamp.Unix()), i.Price}
	err := binary.Write(&buf, ByteOrder, &cm)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Message interface {
	encoding.BinaryMarshaler
}

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

func Write(w io.Writer, a any) error {
	return binary.Write(w, ByteOrder, a)
}

func unixTime(i int32) time.Time {
	return time.Unix(int64(i), 0).UTC()
}
