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
		Message string
		Key     string
		Value   string
	}{
		{"foo=bar", "foo", "bar"},
		{"foo=bar=baz", "foo", "bar=baz"},
		{"foo=", "foo", ""},
		{"foo===", "foo", "=="},
		{"=foo", "", "foo"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			m := msg(tc.Message)
			k, v := m.KV()

			if k != tc.Key || v != tc.Value {
				t.Errorf("unexpected msg: %q, key: %q, value: %q", tc.Message, k, v)
			}
		})
	}
}
