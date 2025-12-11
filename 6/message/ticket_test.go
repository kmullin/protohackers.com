package message

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTicketMsg(t *testing.T) {
	cases := []struct {
		Ticket   Ticket
		Expected []byte
	}{
		{
			Ticket{
				Plate:      "UN1X",
				Road:       66,
				Mile1:      100,
				Timestamp1: time.Unix(123456, 0).UTC(),
				Mile2:      110,
				Timestamp2: time.Unix(123816, 0).UTC(),
				Speed:      100,
			}, []byte{0x21, 0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x42, 0x00, 0x64, 0x00, 0x01, 0xe2, 0x40, 0x00, 0x6e, 0x00, 0x01, 0xe3, 0xa8, 0x27, 0x10},
		},

		{
			Ticket{
				Plate:      "RE05BKG",
				Road:       368,
				Mile1:      1234,
				Timestamp1: time.Unix(1000000, 0).UTC(),
				Mile2:      1235,
				Timestamp2: time.Unix(1000060, 0).UTC(),
				Speed:      60,
			},
			[]byte{0x21, 0x07, 0x52, 0x45, 0x30, 0x35, 0x42, 0x4b, 0x47, 0x01, 0x70, 0x04, 0xd2, 0x00, 0x0f, 0x42, 0x40, 0x04, 0xd3, 0x00, 0x0f, 0x42, 0x7c, 0x17, 0x70},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			b, err := tc.Ticket.MarshalBinary()
			assert.Nil(t, err)
			assert.Equal(t, tc.Expected, b)
		})
	}
}
