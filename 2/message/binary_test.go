package message

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

const msgSize = 9

func TestMsgSize(t *testing.T) {
	msg := clientMessage{
		Timestamp: 4,
		Price:     10,
	}

	if s := binary.Size(msg); s != msgSize {
		t.Fatalf("incorrect msg size got %d wanted %d", s, msgSize)
	}
}

func TestMarshaling(t *testing.T) {
	assert := assert.New(t)

	t.Run("unmarshal", func(t *testing.T) {
		var cm clientMessage
		data := []byte{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}
		r := bytes.NewReader(data)
		err := binary.Read(r, byteOrder, &cm)

		assert.Nil(err)
		assert.Equal(insertByte, cm.Type)
		assert.Equal(int32(12345), cm.Timestamp)
		assert.Equal(int32(101), cm.Price)
	})

	t.Run("marshal", func(t *testing.T) {
		var buf bytes.Buffer
		cm := clientMessage{
			Type:      queryByte,
			Timestamp: 98723984,
			Price:     1293172,
		}
		expected := []byte{0x51, 0x05, 0xe2, 0x68, 0x90, 0x00, 0x13, 0xbb, 0x74}

		err := binary.Write(&buf, byteOrder, cm)
		assert.Nil(err)
		assert.Equal(len(buf.Bytes()), msgSize)
		assert.Equal(expected, buf.Bytes())
	})
}
