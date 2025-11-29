package message

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
