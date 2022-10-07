package main

import (
	"fmt"
	"testing"
)

func TestMsg(t *testing.T) {
	// foo=bar will insert a key foo with value "bar".
	// foo=bar=baz will insert a key foo with value "bar=baz".
	// foo= will insert a key foo with value "" (i.e. the empty string).
	// foo=== will insert a key foo with value "==".
	// =foo will insert a key of the empty string with value "foo".
	testCases := []struct {
		Message []byte
		Key     string
		Value   string
	}{
		{[]byte("foo=bar"), "foo", "bar"},
		{[]byte("foo=bar=baz"), "foo", "bar=baz"},
		{[]byte("foo="), "foo", ""},
		{[]byte("foo==="), "foo", "=="},
		{[]byte("=foo"), "", "foo"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			m := NewMessage(tc.Message)
			if m.Key != tc.Key || m.Value != tc.Value {
				t.Errorf("unexpected msg: %q, key: %q, value: %q", m, tc.Key, tc.Value)
			}
		})
	}
}
