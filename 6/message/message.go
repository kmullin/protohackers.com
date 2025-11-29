package message

import (
	"bytes"
	"encoding"
	"encoding/binary"
)

var ByteOrder = binary.BigEndian

type Message interface {
	encoding.BinaryMarshaler
}

type Error struct {
	Msg string
}

func (e Error) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

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
