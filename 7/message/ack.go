package message

import (
	"strconv"

	"github.com/rs/zerolog"
)

type Ack struct {
	SessionID int
	Length    int
}

func (a *Ack) Marshal() (data []byte) {
	return joinMsg(
		msgAck,
		[]byte(strconv.Itoa(a.SessionID)),
		[]byte(strconv.Itoa(a.Length)),
	)
}

func (a *Ack) MarshalZerologObject(e *zerolog.Event) {
	e.Int("session", a.SessionID).Int("len", a.Length)
}
