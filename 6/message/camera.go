package message

import (
	"bytes"
	"encoding/binary"
)

type IAmCamera struct {
	Road  int
	Mile  int
	Limit int
}

func (iac *IAmCamera) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	msg := struct {
		Road, Mile, Limit uint16
	}{}

	if err := binary.Read(r, byteOrder, &msg); err != nil {
		return err
	}

	iac.Road = int(msg.Road)
	iac.Mile = int(msg.Mile)
	iac.Limit = int(msg.Limit)

	return nil
}
