package message

import (
	"encoding/binary"
	"io"
)

type IAmDispatcher struct {
	Roads []int
}

func readIAmDispatcherMsg(r io.Reader) (*IAmDispatcher, error) {
	var numRoads uint8

	if err := binary.Read(r, byteOrder, &numRoads); err != nil {
		return nil, err
	}

	var road uint16
	var roads []int
	for i := 0; i < int(numRoads); i++ {
		if err := binary.Read(r, byteOrder, &road); err != nil {
			return nil, err
		}

		roads = append(roads, int(road))
	}
	return &IAmDispatcher{roads}, nil
}
