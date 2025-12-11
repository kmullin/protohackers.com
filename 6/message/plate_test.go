package message

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlateMsg(t *testing.T) {
	cases := []struct {
		Msg      []byte
		Expected *Plate
	}{
		{
			Msg: []byte{0x20, 0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x00, 0x03, 0xe8},
			Expected: &Plate{
				Plate:     "UN1X",
				Timestamp: toTime(1000),
			},
		},
		{
			Msg: []byte{0x20, 0x07, 0x52, 0x45, 0x30, 0x35, 0x42, 0x4b, 0x47, 0x00, 0x01, 0xe2, 0x40},
			Expected: &Plate{
				Plate:     "RE05BKG",
				Timestamp: toTime(123456),
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
