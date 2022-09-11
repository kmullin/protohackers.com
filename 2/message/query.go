package message

import (
	"time"
)

type Query struct {
	MinTime time.Time
	MaxTime time.Time
}

func (q Query) Type() msgType {
	return QueryType
}
