package message

import (
	"bytes"
	"time"
)

type Plate struct {
	Plate     string
	Timestamp time.Time
}

func (p *Plate) UnmarshalBinary(data []byte) error {
	var err error

	r := bytes.NewReader(data)

	p.Plate, err = readString(r)
	if err != nil {
		return err
	}

	p.Timestamp, err = readTime(r)
	if err != nil {
		return err
	}

	return nil
}
