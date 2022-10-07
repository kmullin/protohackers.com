package chat

import (
	"fmt"
	"testing"

	"github.com/kmullin/protohackers.com/test"
)

func TestSession(t *testing.T) {
	cases := []struct {
		Name    string
		Payload string
		Valid   bool
	}{
		{"missing newline", "foo", false},
		{"valid", "foo\n", true},
		{"invalid single char", "\n", false},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			client, server := test.Conn(t)

			// write some valid name without a trailing newline '\n'
			fmt.Fprintf(client, c.Payload)
			client.Close()

			session, err := NewSession(server)
			if c.Valid && err != nil {
				t.Fatalf("expected nil, got err: %v", err)
			}
			if !c.Valid && err == nil {
				t.Fatalf("expected err, got session user: %+v", session.User)
			}
		})
	}
}
