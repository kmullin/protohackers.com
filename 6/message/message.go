package message

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
)

var ByteOrder = binary.BigEndian

type Message interface {
	encoding.BinaryMarshaler
}

type Error struct {
	Msg string
}

var ErrLargeMsg = errors.New("msg is too large")

// MaxMsgLen is the maximum decimal value of a uint8
const MaxMsgLen = int(^uint8(0))

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
