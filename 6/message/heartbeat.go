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
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypeHeartbeat); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
