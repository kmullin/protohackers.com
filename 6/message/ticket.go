package message

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

type Ticket struct {
	Plate      string
	Road       uint16
	Mile1      uint16
	Timestamp1 time.Time
	Mile2      uint16
	Timestamp2 time.Time
	Speed      uint16
}

func (t *Ticket) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypeTicket); err != nil {
		return nil, err
	}

	if err := writeString(&buf, t.Plate); err != nil {
		return nil, err
	}

	msg := struct {
		Road       uint16
		Mile1      uint16
		Timestamp1 uint32
		Mile2      uint16
		Timestamp2 uint32
		Speed      uint16
	}{
		t.Road,
		t.Mile1,
		uint32(t.Timestamp1.Unix()),
		t.Mile2,
		uint32(t.Timestamp2.Unix()),
		t.Speed * 100, // the inferred average speed of the car multiplied by 100
	}

	if err := binary.Write(&buf, byteOrder, &msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t *Ticket) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, t)
}
