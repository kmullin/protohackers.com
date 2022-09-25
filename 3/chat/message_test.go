package chat

import "testing"

func TestMessage(t *testing.T) {
	m := message("valid message")
	if !m.IsValid() {
		t.Fatalf("valid message is invalid")
	}
}
