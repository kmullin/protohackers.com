package main

import (
	"bytes"
)

const keySeparator = "="

type messageType int

const (
	messageInsert messageType = iota
	messageRetrieve
)

type message struct {
	Type  messageType // insert or retrieve
	Key   string
	Value string
}

func NewMessage(b []byte) *message {
	var msg message
	if bytes.Contains(b, []byte(keySeparator)) {
		msg.Type = messageInsert
		m := bytes.SplitN(b, []byte(keySeparator), 2)

		if len(m) != 2 {
			panic("bad message length")
		}

		msg.Key = string(m[0])
		msg.Value = string(m[1])
	} else {
		msg.Type = messageRetrieve
		msg.Key = string(b)
	}

	return &msg
}
