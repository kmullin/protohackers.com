package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"testing"
	"time"

	"github.com/kmullin/protohackers.com/2/message"
	"github.com/kmullin/protohackers.com/test"
	"github.com/stretchr/testify/assert"
)

func TestMean(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	inserts := insertCache{
		now.Add(-(60 * time.Minute)): 12093123,
		now.Add(-(30 * time.Minute)): 912378301,
		now.Add(-(30 * time.Minute)): 1902,
		now.Add(-(30 * time.Hour)):   99999999,
	}
	t.Logf("%v", inserts)

	d := 2 * time.Hour
	min := time.Now().Add(-d)
	max := time.Now().Add(d)
	i := inserts.Mean(min, max)
	assert.Equal(int32(6047513), i)
}

func TestHandler(t *testing.T) {
	/*
		<-- 49 00 00 30 39 00 00 00 65   I 12345 101
		<-- 49 00 00 30 3a 00 00 00 66   I 12346 102
		<-- 49 00 00 30 3b 00 00 00 64   I 12347 100
		<-- 49 00 00 a0 00 00 00 00 05   I 40960 5
		<-- 51 00 00 30 00 00 00 40 00   Q 12288 16384
		--> 00 00 00 65                  101
	*/
	assert := assert.New(t)

	t.Run("example request", func(t *testing.T) {
		t.Parallel()
		client, server := test.Conn(t)
		defer client.Close()
		go handler(server)

		requests := [][]byte{
			{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65},
			{0x49, 0x00, 0x00, 0x30, 0x3a, 0x00, 0x00, 0x00, 0x66},
			{0x49, 0x00, 0x00, 0x30, 0x3b, 0x00, 0x00, 0x00, 0x64},
			{0x49, 0x00, 0x00, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x05},
			{0x51, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x40, 0x00},
		}
		expected := []byte{0x00, 0x00, 0x00, 0x65}

		for _, r := range requests {
			n, err := client.Write(r)
			assert.Nil(err)
			assert.Equal(len(r), n)
		}

		var buf bytes.Buffer
		_, err := io.CopyN(&buf, client, int64(len(expected)))
		assert.Nil(err)
		assert.Equal(expected, buf.Bytes())
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()
		client, server := test.Conn(t)
		defer client.Close()
		go handler(server)

		now := time.Now()

		requests := []message.Message{
			message.Insert{now, -1827312},
			message.Query{now, now},
		}
		for _, r := range requests {
			b, err := r.MarshalBinary()
			assert.Nil(err)
			n, err := client.Write(b)
			assert.Nil(err)
			assert.Equal(n, 9)
			t.Logf("% x", b)
		}
		expected := int32(-1827312)

		var i int32
		err := binary.Read(client, message.ByteOrder, &i)
		assert.Nil(err)
		assert.Equal(expected, i)
	})

	t.Run("max int32", func(t *testing.T) {
		t.Parallel()
		client, server := test.Conn(t)
		defer client.Close()
		go handler(server)

		now := time.Now()
		var requests []message.Message
		for i := 0; i < 100; i++ {
			requests = append(requests, message.Insert{now, math.MaxInt32})
		}
		requests = append(requests, message.Query{now, now})

		for _, r := range requests {
			b, err := r.MarshalBinary()
			assert.Nil(err)
			n, err := client.Write(b)
			assert.Nil(err)
			assert.Equal(n, 9)
			t.Logf("% x", b)
		}
		expected := int32(math.MaxInt32)

		var i int32
		err := binary.Read(client, message.ByteOrder, &i)
		assert.Nil(err)
		assert.Equal(expected, i)
	})
}

func TestCalculateMean(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	m := insertCache{
		now.Add(-1): 100,
		now.Add(-2): 101,
		now.Add(-3): 102,
	}

	cases := []struct {
		Name     string
		Expected int32
		Min, Max time.Time
	}{
		{"min after max", 0, now.Add(1), now},
		{"no samples", 0, now, now},
		{"real mean", 101, time.Unix(0, 0), now},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(c.Expected, m.Mean(c.Min, c.Max))
		})
	}
}
