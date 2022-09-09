package main

import (
	"encoding/json"
	"testing"
	"time"

	"git.kpmullin.com/kmullin/protocolhackers.com/test"
)

func TestMalformed(t *testing.T) {
	var blank interface{}
	input := []byte(`{"}`)
	err := json.Unmarshal(input, &blank)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandler(t *testing.T) {
	t.Skip()
	client, server := test.Conn(t)
	handleConn(server)

	client.Write([]byte(`{"method":"isPrime","number":123}\n`))
	client.Write([]byte(`{"method":"isPrime","number":123.2}\n`))
	time.Sleep(5 * time.Second)
	client.Close()
}
