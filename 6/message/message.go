package message

import (
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

var ByteOrder = binary.BigEndian

var ErrLargeMsg = errors.New("msg is too large")

// MaxMsgLen is the maximum decimal value of a uint8
const MaxMsgLen = int(^uint8(0))

type serverMessage interface {
	encoding.BinaryMarshaler
}

type clientMessage interface {
	encoding.BinaryUnmarshaler
}

func New(r io.Reader) (clientMessage, error) {
	var t uint8

	// find out which message we receive
	err := binary.Read(r, ByteOrder, &t)
	if err != nil {
		return nil, err
	}

	switch t {
	case 0x20: // Plate
		var p Plate

		b, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		err = p.UnmarshalBinary(b)
		return &p, err
	case 0x40:
	default:
		return nil, errors.New("unknown message received")
	}

	return nil, nil
}

// readString will read a length prefixed string from r
func readString(r io.Reader) (string, error) {
	var l uint8
	if err := binary.Read(r, ByteOrder, &l); err != nil {
		return "", err
	}

	buf := make([]byte, l)
	if err := binary.Read(r, ByteOrder, &buf); err != nil {
		return "", err
	}

	return string(buf), nil
}

// readTime will read the timestamp from r
func readTime(r io.Reader) (time.Time, error) {
	var ts uint32
	if err := binary.Read(r, ByteOrder, &ts); err != nil {
		return time.Unix(-1, 0).UTC(), err
	}

	return toTime(ts), nil
}

// toTime takes the timestamps from the raw input type and converts it into time.Time
func toTime(t uint32) time.Time {
	return time.Unix(int64(t), 0).UTC()
}
