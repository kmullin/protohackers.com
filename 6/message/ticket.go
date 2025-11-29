package message

import (
	"bytes"
	"encoding/binary"
	"time"
)

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
