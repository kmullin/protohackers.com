package message

import (
	"bytes"
	"encoding/binary"
)

type Error struct {
	Msg string
}

func (e Error) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypeError); err != nil {
		return nil, err
	}

	if err := writeString(&buf, e.Msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
