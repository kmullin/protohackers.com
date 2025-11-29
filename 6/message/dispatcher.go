package message

import (
	"bytes"
	"encoding/binary"
)

type IAmDispatcher struct {
	Roads []int
}

func (iad *IAmDispatcher) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	var numRoads uint8

	if err := binary.Read(r, byteOrder, &numRoads); err != nil {
		return err
	}

	var road uint16
	for i := 0; i < int(numRoads); i++ {
		if err := binary.Read(r, byteOrder, &road); err != nil {
			return err
		}

		iad.Roads = append(iad.Roads, int(road))
	}
	return nil
}
