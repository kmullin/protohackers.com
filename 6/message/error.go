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

	if len(e.Msg) > MaxMsgLen {
		return nil, ErrLargeMsg
	}

	t := uint8(0x10)
	if err := binary.Write(&buf, ByteOrder, &t); err != nil {
		return nil, err
	}

	if err := writeString(&buf, e.Msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
