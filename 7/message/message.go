package message

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
)

var separator = []byte("/")

// Numeric field values must be smaller than 2147483648
const maxMsgInt = 2147483648

// LRCP messages must be smaller than 1000 bytes. You might have to break up data into multiple data messages in order to fit it below this limit.
const MaxSize = 999

// our message types
var (
	msgConnect = []byte("connect")
	msgData    = []byte("data")
	msgAck     = []byte("ack")
	msgClose   = []byte("close")
)

type Msg interface {
	Marshal() []byte // kinda like encoding.BinaryMarshaler, but without the error
	zerolog.LogObjectMarshaler
}

func joinMsg(parts ...[]byte) []byte {
	joined := bytes.Join(parts, separator)
	return append(append(separator, joined...), separator...)
}

func New(b []byte) (Msg, error) {
	msg, err := isValidMsg(b)
	if err != nil {
		return nil, fmt.Errorf("invalid msg: %w", err)
	}
	return msg, nil
}

// isValid checks to see if the payload received is valid
// Packet contents must begin with a forward slash, end with a forward slash, have a valid message type, and have the correct number of fields for the message type.
// Numeric field values must be smaller than 2147483648. This means sessions are limited to 2 billion bytes of data transferred in each direction.
// LRCP messages must be smaller than 1000 bytes. You might have to break up data into multiple data messages in order to fit it below this limit.
func isValidMsg(b []byte) (Msg, error) {
	var msg Msg

	// split it we want at max 5 sub slices
	// data msg might contain the separator key so treat it as in individual field
	// this allows us to check the split for validity
	bb := bytes.SplitN(b, separator, 5)
	// if there are enough fields to work with
	if len(bb) < 4 {
		return nil, fmt.Errorf("not enough fields: %v", len(bb))
	}

	// each message has session id as the second field
	sessionID, err := isValidNumeric(bb[2])
	if err != nil {
		return nil, fmt.Errorf("not valid session id: %w", err)
	}

	// first field designates the message type
	msgType := bb[1]
	switch {
	case bytes.Equal(msgType, msgConnect):
		if len(bb) != 4 {
			return nil, fmt.Errorf("connect msg doesnt have 4 fields: %v", len(bb))
		}

		msg = &Connect{SessionID: sessionID}
	case bytes.Equal(msgType, msgData):
		if len(bb) != 5 { // 5 is our splitN
			return nil, fmt.Errorf("data msg doesnt have enough fields: %v", len(bb))
		}

		pos, err := isValidNumeric(bb[3])
		if err != nil {
			return nil, fmt.Errorf("not valid pos: %w", err)
		}
		msg = &Data{
			SessionID: sessionID,
			Pos:       pos,
			Data:      bytes.TrimSuffix(bb[4], separator),
		}
	case bytes.Equal(msgType, msgAck):
		if len(bb) != 5 {
			return nil, fmt.Errorf("ack msg doesnt have enough fields: %v", len(bb))
		}

		l, err := isValidNumeric(bb[3])
		if err != nil {
			return nil, fmt.Errorf("not valid length: %w", err)
		}
		msg = &Ack{
			SessionID: sessionID,
			Length:    l,
		}
	case bytes.Equal(msgType, msgClose):
		if len(bb) != 4 {
			return nil, fmt.Errorf("close msg doesnt have enough fields: %v", len(bb))
		}
		msg = &Close{SessionID: sessionID}
	}

	return msg, nil
}

func isValidNumeric(b []byte) (int, error) {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, fmt.Errorf("err converting to int: %w", err)
	}

	// Numeric field values must be smaller than 2147483648.
	// This means sessions are limited to 2 billion bytes of data transferred in each direction.
	// SESSION field must be a non-negative integer.
	// POS field must be a non-negative integer.
	// LENGTH field must be a non-negative integer.
	if i < 0 || i >= maxMsgInt {
		return 0, fmt.Errorf("number not in range 0-%v: %v", maxMsgInt-1, i)
	}

	return i, nil
}
