package message

import (
	"strconv"

	"github.com/rs/zerolog"
)

type Close struct {
	SessionID int
}

func (c *Close) Marshal() (data []byte) {
	return joinMsg(
		msgClose,
		[]byte(strconv.Itoa(c.SessionID)),
	)
}

func (c *Close) MarshalZerologObject(e *zerolog.Event) {}
