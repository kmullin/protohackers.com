package message

import (
	"encoding/binary"
	"io"
)

type IAmCamera struct {
	Road  int
	Mile  int
	Limit int
}

func readIAmCameraMsg(r io.Reader) (*IAmCamera, error) {
	var msg struct {
		Road, Mile, Limit uint16
	}

	if err := binary.Read(r, byteOrder, &msg); err != nil {
		return nil, err
	}

	return &IAmCamera{
		int(msg.Road),
		int(msg.Mile),
		int(msg.Limit),
	}, nil
}
