package message

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

type Ticket struct {
	Plate      string
	Road       RoadID
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
		Road       RoadID
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

func readTicketMsg(r io.Reader) (*Ticket, error) {
	var err error
	var t Ticket

	t.Plate, err = readString(r)
	if err != nil {
		return nil, err
	}

	msg := struct {
		Road       RoadID
		Mile1      uint16
		Timestamp1 uint32
		Mile2      uint16
		Timestamp2 uint32
		Speed      uint16
	}{}

	if err := binary.Read(r, byteOrder, &msg); err != nil {
		return nil, err
	}

	t.Road = msg.Road
	t.Mile1 = msg.Mile1
	t.Timestamp1 = toTime(msg.Timestamp1)
	t.Mile2 = msg.Mile2
	t.Timestamp2 = toTime(msg.Timestamp2)
	t.Speed = msg.Speed
	return &t, nil
}
