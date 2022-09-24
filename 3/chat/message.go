package chat

import (
	"bufio"
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

func msgSplitFunc(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF {
		// if we're already at EOF, we dont want any remaining data
		return 0, nil, nil
	}
	return bufio.ScanLines(data, atEOF)
}
