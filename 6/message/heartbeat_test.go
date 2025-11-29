package message

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWantHeartbeat(t *testing.T) {
	cases := []struct {
		Msg      []byte
		Expected *WantHeartbeat
	}{
		{
			Msg:      []byte{0x40, 0x00, 0x00, 0x00, 0x0a},
			Expected: &WantHeartbeat{fromDeci(10)},
		},
		{
			Msg:      []byte{0x40, 0x00, 0x00, 0x04, 0xdb},
			Expected: &WantHeartbeat{fromDeci(1243)},
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

func TestHeartbeat(t *testing.T) {
	var hb Heartbeat

	b, err := hb.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, []byte{MsgTypeHeartbeat}, b)
}
