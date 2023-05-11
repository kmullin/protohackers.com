package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var validBogusCoins = []string{
	"7F1u3wSD5RbOHQmupo9nx4TnhQ",
	"7iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX",
	"7LOrwbDlS8NujgjddyogWgIM93MV5N2VR",
	"7adNeSwJkMakpEcln9HEtthSRtxdmEHOT8T",
}

var invalidBogusCoins = []string{
	"9F1u3wSD5RbOHQmupo9nx4TnhQ",
	"9iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX",
	"9LOrwbDlS8NujgjddyogWgIM93MV5N2VR",
	"7adNeSwJkMakpEcln9HEtthSRtxdmEHOT8Tsd89f",
}

func TestBogusCoin(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for _, s := range validBogusCoins {
			assert.Truef(t, IsBogusCoinAddress(s), "%s", s)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		for _, s := range invalidBogusCoins {
			assert.Falsef(t, IsBogusCoinAddress(s), "%s", s)
		}
	})
}

func TestReplaceCoins(t *testing.T) {
	expected := "[TinyDev654] Please pay the ticket price of 15 Boguscoins to one of these addresses: 7YWHMfk9JZe0LM0g1ZauHuiSxhI 7YWHMfk9JZe0LM0g1ZauHuiSxhI 7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	given := "[TinyDev654] Please pay the ticket price of 15 Boguscoins to one of these addresses: 7YWHMfk9JZe0LM0g1ZauHuiSxhI 7FaqH6wRIWnqcJBoDAa0xiNirCBjftwt 77eBgtfG7Q3CqzOizaUuituyl3S"
	msg := ReplaceBogusCoins(given)
	assert.Equal(t, expected, msg)
}
