package message

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

type WantHeartbeat struct {
	Interval time.Duration
}

func (wh *WantHeartbeat) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypeWantHeartbeat); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, byteOrder, toDeci(wh.Interval)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (wh *WantHeartbeat) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, wh)
}

func readWantHeartbeatMsg(r io.Reader) (*WantHeartbeat, error) {
	var i uint32

	if err := binary.Read(r, byteOrder, &i); err != nil {
		return nil, err
	}

	return &WantHeartbeat{fromDeci(i)}, nil
}

type Heartbeat struct{}

func (hb *Heartbeat) MarshalBinary() ([]byte, error) {
	return []byte{MsgTypeHeartbeat}, nil
}

func (hb *Heartbeat) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, hb)
}
