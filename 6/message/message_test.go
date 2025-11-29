package message

import (
	"bytes"
	"fmt"
	"math/rand"
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

func TestMsgLength(t *testing.T) {
	e := Error{RandStringRunes(256)}
	_, err := e.MarshalBinary()
	assert.Error(t, err)
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
