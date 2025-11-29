package message

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

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

func TestMsgLength(t *testing.T) {
	e := Error{RandStringRunes(256)}
	_, err := e.MarshalBinary()
	assert.Error(t, err)
}

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
				Speed:      10000,
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
				Speed:      6000,
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

func TestClientMsg(t *testing.T) {
	// Plate msg
	r := bytes.NewReader([]byte{0x20, 0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x00, 0x03, 0xe8})
	msg, err := New(r)
	assert.Nil(t, err)

	t.Logf("%+v", msg)

	r = bytes.NewReader([]byte{0x20, 0x07, 0x52, 0x45, 0x30, 0x35, 0x42, 0x4b, 0x47, 0x00, 0x01, 0xe2, 0x40})

	msg, err = New(r)
	assert.Nil(t, err)

	t.Logf("%+v", msg)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
