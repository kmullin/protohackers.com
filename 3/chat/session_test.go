package chat

import (
	"fmt"
	"testing"

	"git.kpmullin.com/kmullin/protocolhackers.com/test"
)

func TestSession(t *testing.T) {
	client, server := test.Conn(t)

	// write some valid name without a trailing newline '\n'
	fmt.Fprintf(client, "foo")
	client.Close()

	session, err := NewSession(server)
	if err == nil {
		t.Fatalf("expected err, got session user: %+v", session.User)
	}
}
