package message

import (
	"io"
	"time"
)

type Plate struct {
	Plate     string
	Timestamp time.Time
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
