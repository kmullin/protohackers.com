package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAckMsg(t *testing.T) {
	cases := []struct {
		Input       []byte
		Expected    *Ack
		ShouldError bool
	}{
		{Input: []byte("/ack/1045021881/6/"), Expected: &Ack{SessionID: 1045021881, Length: 6}, ShouldError: false},
		{Input: []byte("/ack/1045021881/6"), Expected: &Ack{}, ShouldError: true},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			msg, err := New(tc.Input)
			if tc.ShouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Expected, msg)
			}
		})
	}
}
