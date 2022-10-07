package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/kmullin/protohackers.com/test"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// mock function for testinging only
func (r *response) IsValid() bool {
	return r.Method == onlyValidMethod
}

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

func TestHandler(t *testing.T) {
	// these are considered invalid input data
	invalidCases := [][]byte{
		[]byte(`{"method":"isPrime","number":"1043398"}`),
		[]byte(`{"method":"isPrime"}`),
		[]byte(`{"mthod":"isPrime"}`),
		[]byte(`{"method":"isrime"}`),
		[]byte(`{"method":"isPrime"}`),
	}

	for i, b := range invalidCases {
		t.Run(fmt.Sprintf("%v invalid request", i), func(t *testing.T) {
			client, server := test.Conn(t)
			go handleConn(server)

			client.Write(append(b, []byte("\n")...))
			scanner := bufio.NewScanner(client)
			scanner.Scan() // scan once

			var r response
			dec := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
			err := dec.Decode(&r)
			if err != nil {
				t.Logf("err decoding response: %v", err)
			}

			if r.IsValid() {
				t.Fatal("valid response from malformed request")
			}
			client.Close()

		})
	}
}
