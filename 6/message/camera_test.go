package message

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIAmCamera(t *testing.T) {
	cases := []struct {
		Msg      []byte
		Expected *IAmCamera
	}{
		{
			Msg: []byte{0x80, 0x00, 0x42, 0x00, 0x64, 0x00, 0x3c},
			Expected: &IAmCamera{
				Road:  66,
				Mile:  100,
				Limit: 60,
			},
		},
		{
			Msg: []byte{0x80, 0x01, 0x70, 0x04, 0xd2, 0x00, 0x28},
			Expected: &IAmCamera{
				Road:  368,
				Mile:  1234,
				Limit: 40,
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
