package message

import (
	"encoding/binary"
	"io"
	"time"
)

type WantHeartbeat struct {
	Interval time.Duration
}

func readWantHeartbeatMsg(r io.Reader) (*WantHeartbeat, error) {
	var i uint32

	if err := binary.Read(r, byteOrder, &i); err != nil {
		return nil, err
	}

	return &WantHeartbeat{fromDeci(i)}, nil
}

type Heartbeat struct{}

func (hb Heartbeat) MarshalBinary() ([]byte, error) {
	return []byte{MsgTypeHeartbeat}, nil
}
