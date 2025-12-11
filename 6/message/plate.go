package message

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

type Plate struct {
	Plate     string
	Timestamp time.Time
}

// func (p *Plate) MarshalZerologObject(e *zerolog.Event) {
// 	e.Str("plate", p.Plate).
// 		Time("timestamp", p.Timestamp)
// }

func (p *Plate) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypePlate); err != nil {
		return nil, err
	}
	if err := writeString(&buf, p.Plate); err != nil {
		return nil, err
	}
	if err := writeTime(&buf, p.Timestamp); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Plate) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, p)
}

func readPlateMsg(r io.Reader) (*Plate, error) {
	var err error
	var p Plate

	p.Plate, err = readString(r)
	if err != nil {
		return nil, err
	}

	p.Timestamp, err = readTime(r)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
