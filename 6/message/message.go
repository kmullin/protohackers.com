package message

import (
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

var byteOrder = binary.BigEndian

var ErrLargeMsg = errors.New("msg is too large")

// MaxMsgLen is the maximum decimal value of a uint8
const MaxMsgLen = int(^uint8(0)) // 255

// client messages
const (
	MsgTypePlate         = uint8(0x20)
	MsgTypeWantHeartbeat = uint8(0x40)
	MsgTypeIAmCamera     = uint8(0x80)
	MsgTypeIAmDispatcher = uint8(0x81)
)

// server messages
const (
	MsgTypeError     = uint8(0x10)
	MsgTypeTicket    = uint8(0x21)
	MsgTypeHeartbeat = uint8(0x41)
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
	// client messages
	case MsgTypePlate:
		return readPlateMsg(r)
	case MsgTypeWantHeartbeat:
		return readWantHeartbeatMsg(r)
	case MsgTypeIAmCamera:
		return readIAmCameraMsg(r)
	case MsgTypeIAmDispatcher:
		return readIAmDispatcherMsg(r)
	// server messges
	case MsgTypeError:
		return readErrorMsg(r)
	case MsgTypeTicket:
		return readTicketMsg(r)
	case MsgTypeHeartbeat:
		return &Heartbeat{}, nil
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

// writeTo is a generic helper for client messages
func writeTo(w io.Writer, msg encoding.BinaryMarshaler) (int64, error) {
	b, err := msg.MarshalBinary()
	if err != nil {
		return 0, err
	}

	if err := binary.Write(w, byteOrder, b); err != nil {
		return 0, err
	}

	return int64(binary.Size(b)), nil
}

// writeTime will write to w the binary representation of the given time
func writeTime(w io.Writer, t time.Time) error {
	return binary.Write(w, byteOrder, uint32(t.Unix()))
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

// toDeci returns deciseconds for a given duration
func toDeci(d time.Duration) uint32 {
	return uint32(d / (time.Second / 10))
}
