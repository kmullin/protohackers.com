package message

import (
	"bytes"
	"encoding/binary"
	"time"
)

type WantHeartbeat struct {
	Interval time.Duration
}

func (hb *WantHeartbeat) UnmarshalBinary(data []byte) error {
	var i uint32

	r := bytes.NewReader(data)

	if err := binary.Read(r, byteOrder, &i); err != nil {
		return err
	}

	hb.Interval = fromDeci(i)

	return nil
}

type Heartbeat struct{}

func (hb Heartbeat) MarshalBinary() ([]byte, error) {
	return []byte{MsgTypeHeartbeat}, nil
}
