package main

import (
	"testing"
	"time"

	"git.kpmullin.com/kmullin/protocolhackers.com/2/message"
	"github.com/stretchr/testify/assert"
)

func TestMean(t *testing.T) {
	assert := assert.New(t)

	inserts := []message.Insert{
		{time.Now().Add(-(60 * time.Minute)), 12093123},
		{time.Now().Add(-(30 * time.Minute)), 912378301},
		{time.Now().Add(-(30 * time.Minute)), 1902},
	}

	d := 2 * time.Hour
	min := time.Now().Add(-d)
	max := time.Now().Add(d)
	i := findMean(inserts, min, max)
	assert.Equal(int32(308157775), i)
}
