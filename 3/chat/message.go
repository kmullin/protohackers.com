package chat

import (
	"unicode"
)

// message is a single line of ASCII text
type message struct {
	msg     string
	session *Session
}

func (m *message) String() string {
	return m.msg
}

func (m *message) IsValid() bool {
	if len(m.msg) == 0 {
		return false
	}
	for i := 0; i < len(m.msg); i++ {
		if m.msg[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
