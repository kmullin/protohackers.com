package message

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

var byteOrder = binary.BigEndian

var ErrLargeMsg = errors.New("msg is too large")

// MaxMsgLen is the maximum decimal value of a uint8
const MaxMsgLen = int(^uint8(0)) // 255

const (
	MsgTypeError         = uint8(0x10)
	MsgTypePlate         = uint8(0x20)
	MsgTypeTicket        = uint8(0x21)
	MsgTypeWantHeartbeat = uint8(0x40)
	MsgTypeHeartbeat     = uint8(0x41)
	MsgTypeIAmCamera     = uint8(0x80)
	MsgTypeIAmDispatcher = uint8(0x81)
)

type Message any

type stringWriter interface {
	io.StringWriter
	io.Writer
}

func New(r io.Reader) (Message, error) {
	var t uint8

	// find out which message we receive
	err := binary.Read(r, byteOrder, &t)
	if err != nil {
		return nil, err
	}

	switch t {
	case MsgTypePlate:
		return readPlateMsg(r)
	case MsgTypeWantHeartbeat:
		return readWantHeartbeatMsg(r)
	case MsgTypeIAmCamera:
		return readIAmCameraMsg(r)
	case MsgTypeIAmDispatcher:
		return readIAmDispatcherMsg(r)
	default:
		return nil, errors.New("unknown message received")
	}
}

func writeString(w stringWriter, s string) error {
	if len(s) > MaxMsgLen {
		return ErrLargeMsg
	}

	l := uint8(len(s))

	if err := binary.Write(w, byteOrder, &l); err != nil {
		return err
	}

	if _, err := w.WriteString(s); err != nil {
		return err
	}
	return nil
}

// readString will read a length prefixed string from r
func readString(r io.Reader) (string, error) {
	var l uint8
	if err := binary.Read(r, byteOrder, &l); err != nil {
		return "", err
	}

	buf := make([]byte, l)
	if err := binary.Read(r, byteOrder, &buf); err != nil {
		return "", err
	}

	return string(buf), nil
}

// readTime will read the timestamp from r
func readTime(r io.Reader) (time.Time, error) {
	var ts uint32
	if err := binary.Read(r, byteOrder, &ts); err != nil {
		return time.Unix(-1, 0).UTC(), err
	}

	return toTime(ts), nil
}

// toTime takes the timestamps from the raw input type and converts it into time.Time
func toTime(t uint32) time.Time {
	return time.Unix(int64(t), 0).UTC()
}

// fromDeci converts Decisecond interval into a usable time.Duration
func fromDeci(t uint32) time.Duration {
	return time.Duration(t) * time.Second / 10
}
