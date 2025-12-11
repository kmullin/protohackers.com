package message

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIAmDispatcher(t *testing.T) {
	cases := []struct {
		Msg      []byte
		Expected *IAmDispatcher
	}{
		{
			Msg: []byte{0x81, 0x01, 0x00, 0x42},
			Expected: &IAmDispatcher{
				Roads: []RoadID{66},
			},
		},
		{
			Msg: []byte{0x81, 0x03, 0x00, 0x42, 0x01, 0x70, 0x13, 0x88},
			Expected: &IAmDispatcher{
				Roads: []RoadID{66, 368, 5000},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := bytes.NewReader(tc.Msg)
			msg, err := New(r)
			assert.Nil(t, err)
			assert.Equal(t, tc.Expected, msg)
		})
	}
}
