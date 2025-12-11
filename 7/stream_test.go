package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamWriterAt(t *testing.T) {
	cases := []struct {
		Input string
		Pos   int

		Expected []byte
	}{
		{
			Input:    "hello",
			Pos:      0,
			Expected: []byte("hello"),
		},
		{
			Input:    "hello\n",
			Pos:      5,
			Expected: []byte("\x00\x00\x00\x00\x00hello\n"),
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			s := new(stream)
			n, err := s.WriteAt([]byte(tc.Input), int64(tc.Pos))
			assert.NoError(t, err)
			assert.Equal(t, len(tc.Input), n, "wrote exactly the input size")
			assert.Equalf(t, tc.Expected, s.buf, "is expected output")
		})
	}
}

func TestStreamReaderAt(t *testing.T) {
	cases := []struct {
		Initial string
		Pos     int
		BufSize int

		Expected  []byte
		expectEOF bool
	}{
		{
			Initial:  "string",
			Pos:      0,
			BufSize:  6,
			Expected: []byte("string"),
		},
		{
			Initial:   "hello\n",
			Pos:       5,
			BufSize:   5,
			Expected:  []byte("\n"),
			expectEOF: true,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			s := new(stream)
			s.buf = []byte(tc.Initial)
			t.Logf("%s", s.buf)

			p := make([]byte, tc.BufSize)
			n, err := s.ReadAt(p, int64(tc.Pos))

			if tc.expectEOF {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, len(tc.Expected), n, "read exactly the input size")
			assert.Equalf(t, tc.Expected, p[:n], "is expected output")
		})
	}
}
