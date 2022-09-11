package message

import (
	"encoding/binary"
)

const (
	insertByte = byte('I')
	queryByte  = byte('Q')
)

var byteOrder = binary.BigEndian

// clientMessage represents requests from the client
type clientMessage struct {
	Type byte
	N1   int32
	N2   int32
}
