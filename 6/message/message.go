package message

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

var ByteOrder = binary.BigEndian

var ErrLargeMsg = errors.New("msg is too large")

// MaxMsgLen is the maximum decimal value of a uint8
const MaxMsgLen = int(^uint8(0))

type Message interface {
	encoding.BinaryMarshaler
}

type Error struct {
	Msg string
}

func (e Error) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if len(e.Msg) > MaxMsgLen {
		return nil, ErrLargeMsg
	}

	msg := struct {
		Type   uint8
		StrLen uint8
	}{
		0x10, uint8(len(e.Msg)),
	}

	err := binary.Write(&buf, ByteOrder, &msg)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(e.Msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Ticket struct {
	Plate      string
	Road       int
	Mile1      int
	Timestamp1 time.Time
	Mile2      int
	Timestamp2 time.Time
	Speed      int
}

func (t Ticket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	msg := struct {
		Type   uint8
		StrLen uint8
	}{
		0x21, uint8(len(t.Plate)),
	}

	if err := binary.Write(&buf, ByteOrder, &msg); err != nil {
		return nil, err
	}

	if _, err := buf.WriteString(t.Plate); err != nil {
		return nil, err
	}

	msg2 := struct {
		Road       uint16
		Mile1      uint16
		Timestamp1 uint32
		Mile2      uint16
		Timestamp2 uint32
		Speed      uint16
	}{
		uint16(t.Road),
		uint16(t.Mile1),
		uint32(t.Timestamp1.Unix()),
		uint16(t.Mile2),
		uint32(t.Timestamp2.Unix()),
		uint16(t.Speed),
	}

	if err := binary.Write(&buf, ByteOrder, &msg2); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Plate struct {
	Plate     string
	Timestamp time.Time
}

func (p Plate) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func New(r io.Reader) (Message, error) {
	var t uint8

	// find out which message we receive
	err := binary.Read(r, ByteOrder, &t)
	if err != nil {
		return nil, err
	}

	switch t {
	case 0x20: // Plate
		var l uint8

		if err := binary.Read(r, ByteOrder, &l); err != nil {
			return nil, err
		}

		buf := make([]byte, l)
		if err := binary.Read(r, ByteOrder, &buf); err != nil {
			return nil, err
		}

		var ts uint32
		if err := binary.Read(r, ByteOrder, &ts); err != nil {
			return nil, err
		}

		return &Plate{Plate: string(buf), Timestamp: time.Unix(int64(ts), 0).UTC()}, nil
	case 0x40:
	default:
		return nil, errors.New("unknown message received")
	}

	return nil, nil
}
