package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"git.kpmullin.com/kmullin/protocolhackers.com/test"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
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

	invalidCases := [][]byte{
		[]byte(`{"method":"isPrime","number":"1043398"}`),
		[]byte(`{"method":"isPrime"}`),
	}

	for i, b := range invalidCases {
		t.Run(fmt.Sprintf("%v invalid json number", i), func(t *testing.T) {
			client, server := test.Conn(t)
			go handleConn(server)

			client.Write(append(b, []byte("\n")...))
			scanner := bufio.NewScanner(client)
			for scanner.Scan() {
				var r map[string]interface{}
				dec := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
				err := dec.Decode(&r)
				if err != nil {
					t.Fatalf("err decoding response: %v", err)
				}
				// is valid response ?
				break
			}
			client.Close()
		})
	}
}
