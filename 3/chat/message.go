package chat

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

// Message is a single line of ASCII text
type Message string

func ReadMessage(r io.Reader) (Message, error) {
	bufReader := bufio.NewReader(r)
	s, err := bufReader.ReadString('\n')
	if err != nil {
		return Message(""), err
	}

	m := Message(s)
	if !m.isValid() {
		return Message(""), fmt.Errorf("invalid msg")
	}
	return m, nil
}

func (m Message) Write(w io.Writer) error {
	_, err := fmt.Fprintln(w, m)
	return err
}

func (m Message) isValid() bool {
	for i := 0; i < len(m); i++ {
		if m[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
