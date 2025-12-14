package message

import (
	"strconv"

	"github.com/rs/zerolog"
)

type Data struct {
	SessionID int
	Pos       int
	Data      []byte
}

func (d *Data) Marshal() (data []byte) {
	return joinMsg(
		msgData,
		[]byte(strconv.Itoa(d.SessionID)),
		[]byte(strconv.Itoa(d.Pos)),
		d.Data,
	)
}

func (d *Data) MarshalZerologObject(e *zerolog.Event) {
	e.Int("session", d.SessionID).
		Int("pos", d.Pos).
		Int("len", len(d.Data)).
		Bytes("data", d.Data)
}
