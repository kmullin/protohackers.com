package message

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgLength(t *testing.T) {
	e := Error{RandStringRunes(256)}
	_, err := e.MarshalBinary()
	assert.Error(t, err)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
