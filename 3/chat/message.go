package chat

import (
	"bufio"
	"unicode"
)

// Message is a single line of ASCII text
type Message string

func (m Message) String() string {
	return string(m)
}

func (m Message) isValid() bool {
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

func splitFunc(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF {
		// if we're already at EOF, we dont want any remaining data
		return 0, nil, nil
	}
	return bufio.ScanLines(data, atEOF)
}
