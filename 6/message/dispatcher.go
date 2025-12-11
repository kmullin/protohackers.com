package message

import (
	"bytes"
	"encoding/binary"
	"io"
)

type IAmDispatcher struct {
	Roads []RoadID
}

func (iad *IAmDispatcher) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypeIAmDispatcher); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, byteOrder, uint8(len(iad.Roads))); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, byteOrder, iad.Roads); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (iad *IAmDispatcher) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, iad)
}

func readIAmDispatcherMsg(r io.Reader) (*IAmDispatcher, error) {
	var numRoads uint8
	if err := binary.Read(r, byteOrder, &numRoads); err != nil {
		return nil, err
	}

	var iad IAmDispatcher
	iad.Roads = make([]RoadID, numRoads)
	if err := binary.Read(r, byteOrder, &iad.Roads); err != nil {
		return nil, err
	}
	return &iad, nil
}
