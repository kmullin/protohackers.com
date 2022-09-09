package main

import (
	"encoding/json"
	"testing"

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
	client, server := test.Conn(t)
	go handleConn(server)

	p := []byte(`{"method":"isPrime","number":123}
{"method":"isPrime","number":11}
	`)
	n, err := client.Write(p)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(p) {
		t.Fatalf("wrong num bytes written, got %v wanted %v", n, len(p))
	}
	client.Close()
}
