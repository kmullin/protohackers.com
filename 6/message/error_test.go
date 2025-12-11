package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMsg(t *testing.T) {
	cases := []struct {
		Msg      string
		Expected []byte
	}{
		{"bad", []byte{0x10, 0x03, 0x62, 0x61, 0x64}},
		{"illegal msg", []byte{0x10, 0x0b, 0x69, 0x6c, 0x6c, 0x65,
			0x67, 0x61, 0x6c, 0x20, 0x6d, 0x73, 0x67}},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			e := Error{tc.Msg}
			b, err := e.MarshalBinary()
			assert.Nil(t, err)
			assert.Equal(t, tc.Expected, b)
		})
	}
}
