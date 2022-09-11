package main

import (
	"bytes"
	"io"
	"testing"
	"time"

	"git.kpmullin.com/kmullin/protocolhackers.com/2/message"
	"git.kpmullin.com/kmullin/protocolhackers.com/test"
	"github.com/stretchr/testify/assert"
)

func TestMean(t *testing.T) {
	assert := assert.New(t)

	inserts := []message.Insert{
		{time.Now().Add(-(60 * time.Minute)), 12093123},
		{time.Now().Add(-(30 * time.Minute)), 912378301},
		{time.Now().Add(-(30 * time.Minute)), 1902},
		{time.Now().Add(-(30 * time.Hour)), 99999999},
	}

	d := 2 * time.Hour
	min := time.Now().Add(-d)
	max := time.Now().Add(d)
	i := findMean(inserts, min, max)
	assert.Equal(int32(308157775), i)
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
}
