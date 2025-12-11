package message

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Error struct {
	Msg string
}

func (e *Error) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, MsgTypeError); err != nil {
		return nil, err
	}

	if err := writeString(&buf, e.Msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *Error) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, e)
}

func readErrorMsg(r io.Reader) (*Error, error) {
	var err error
	var e Error

	e.Msg, err = readString(r)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
