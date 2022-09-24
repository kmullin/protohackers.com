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

func ScanMessage(r io.Reader) (Message, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitFunc)
	if !scanner.Scan() {
		return Message(""), fmt.Errorf("err scanning")
	}

	m := Message(scanner.Text())
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

func splitFunc(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF {
		// if we're already at EOF, we dont want any remaining data
		return 0, nil, nil
	}
	return bufio.ScanLines(data, atEOF)
}
