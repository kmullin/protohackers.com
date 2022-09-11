package message

import (
	"encoding/binary"
)

const msgSize = 9

const (
	insertByte = byte('I')
	queryByte  = byte('Q')
)

var byteOrder = binary.BigEndian

// clientMessage represents requests from the client
type clientMessage struct {
	Type      byte
	Timestamp int32
	Price     int32
}
