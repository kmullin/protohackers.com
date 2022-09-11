package message

import (
	"time"
)

type query struct {
	MinTime time.Time
	MaxTime time.Time
}

func (q query) Type() msgType {
	return QueryType
}
