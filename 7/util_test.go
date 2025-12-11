package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscaping(t *testing.T) {
	cases := []struct {
		EscapedData   []byte
		EscapedLen    int
		UnescapedData []byte
		UnescapedLen  int
	}{
		{
			EscapedData:   []byte("foo\\/bar\\\\baz"),
			EscapedLen:    13,
			UnescapedData: []byte("foo/bar\\baz"),
			UnescapedLen:  11,
		},
		{
			EscapedData:   []byte("\\/"),
			EscapedLen:    2,
			UnescapedData: []byte("/"),
			UnescapedLen:  1,
		},
		{
			EscapedData:   []byte("hello\n"),
			EscapedLen:    6,
			UnescapedData: []byte("hello\n"),
			UnescapedLen:  6,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("unescape %v", i), func(t *testing.T) {
			output := unescapeData(tc.EscapedData)
			assert.Equal(t, tc.UnescapedData, output)
			assert.Len(t, output, tc.UnescapedLen)
		})
		t.Run(fmt.Sprintf("escape %v", i), func(t *testing.T) {
			output := escapeData(tc.UnescapedData)
			assert.Equal(t, tc.EscapedData, output)
			assert.Len(t, output, tc.EscapedLen)
		})
	}
}

func TestReverseBytes(t *testing.T) {
	cases := []struct {
		Input    []byte
		Expected []byte
	}{
		// deal with newlines correctly
		{Input: []byte("hello\n"), Expected: []byte("olleh\n")},
		// without newlines
		{Input: []byte("hello"), Expected: []byte("olleh")},
		// just a newline
		{Input: []byte("\n"), Expected: []byte("\n")},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			output := reverseBytes(tc.Input)
			assert.Equal(t, tc.Expected, output)
			assert.Len(t, output, len(tc.Input))
		})
	}
}
