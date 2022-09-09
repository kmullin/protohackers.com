package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMalformed(t *testing.T) {
	// A response is malformed if it is not a well-formed JSON object,
	// if any required field is missing, if the method name is not "isPrime",
	// or if the prime value is not a boolean.
	var buf bytes.Buffer
	var r response
	var isMalformed bool

	sendMalformedResponse(&buf)
	dec := json.NewDecoder(&buf)
	err := dec.Decode(&r)
	t.Run("invalid json", func(t *testing.T) {
		if err != nil {
			// invalid json
			isMalformed = true
		}
	})
	t.Run("wrong method name", func(t *testing.T) {
		if r.Method != onlyValidMethod {
			// method name is wrong
			isMalformed = true
		}
	})

	if !isMalformed {
		t.Fatalf("object returned is not malformed")
	}
}

/*
func TestHandler(t *testing.T) {
	client, server := test.Conn(t)
	go handleConn(server)

	n, err := client.Write(p)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(p) {
		t.Fatalf("wrong num bytes written, got %v wanted %v", n, len(p))
	}
	t.Log("scanning...")
	scanner := bufio.NewScanner(client)
	for scanner.Scan() {

	}

	client.Close()
}
*/
