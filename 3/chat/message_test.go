package chat

import (
	"bytes"
	"testing"
)

func TestMessage(t *testing.T) {
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
