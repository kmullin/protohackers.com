package message

import (
	"strconv"

	"github.com/rs/zerolog"
)

type Connect struct {
	SessionID int
}

func (c *Connect) Marshal() (data []byte) {
	return joinMsg(
		msgConnect,
		[]byte(strconv.Itoa(c.SessionID)),
	)
}

func (c *Connect) MarshalZerologObject(e *zerolog.Event) {
	e.Int("session", c.SessionID)
}
