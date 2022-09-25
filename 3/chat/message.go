package chat

import (
	"unicode"
)

// message is a single line of ASCII text
type message string

func (m message) String() string {
	return string(m)
}

func (m message) IsValid() bool {
	if len(m) == 0 {
		return false
	}
	for i := 0; i < len(m); i++ {
		if m[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
