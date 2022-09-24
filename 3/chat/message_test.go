package chat

/*

import (
	"bytes"
	"testing"
)

func TestReadMessage(t *testing.T) {
	cases := []struct {
		Name  string
		Msg   []byte
		Valid bool
	}{
		{"valid message", []byte("this is a test message\r\n"), true},
		{"invalid message", []byte("this is a test message\r"), false},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			r := bytes.NewReader(c.Msg)
			_, err := ReadMessage(r)
			// if valid expect no err
			if c.Valid && err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}

			// if invalid expect err
			if !c.Valid && err == nil {
				t.Fatalf("expected err, got nil")
			}
		})
	}
}

/*
func TestWriteMessage(t *testing.T) {
	var b bytes.Buffer
	orig := "this is a response message"
	m := Message(orig)
	m.Write(&b)

	// written string should always have a single trailing newline
	if b.String() != fmt.Sprintf("%s\n", orig) {
		t.Fatalf("expected %q, got %q", orig, b.String())
	}
}
*/

/*
var result Message

func BenchmarkReadMessage(b *testing.B) {
	var m Message
	r := bytes.NewReader([]byte("this is a test message...\r\n"))
	for i := 0; i < b.N; i++ {
		m, _ = ReadMessage(r)
		r.Seek(0, 0)
	}
	result = m
}
*/
