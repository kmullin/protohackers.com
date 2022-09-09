package main

import (
	"encoding/json"
	"net"
	"testing"
	"time"
)

// taken from https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=37

func testServer(t *testing.T) (client, server net.Conn) {
	var err error
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer ln.Close()
		server, err = ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
	}()

	client, err = net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	return client, server
}

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
	client, server := testServer(t)
	handleConn(server)

	client.Write([]byte(`{"method":"isPrime","number":123}\n`))
	client.Write([]byte(`{"method":"isPrime","number":123.2}\n`))
	time.Sleep(5 * time.Second)
	client.Close()
}
