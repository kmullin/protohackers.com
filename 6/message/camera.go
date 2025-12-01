package message

import (
	"bytes"
	"encoding/binary"
	"io"
)

type IAmCamera struct {
	Road, Mile, Limit uint16
}

func (iac *IAmCamera) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypeIAmCamera); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, byteOrder, iac); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (iac *IAmCamera) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, iac)
}

func readIAmCameraMsg(r io.Reader) (*IAmCamera, error) {
	var msg IAmCamera

	if err := binary.Read(r, byteOrder, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
