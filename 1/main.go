// Prime Time
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net"

	"git.kpmullin.com/kmullin/protocolhackers.com/server"
)

func main() {
	server.TCP(handleConn)
}

type inputRequest struct {
	Method string      `json:"method"`
	Number json.Number `json:"number"`
}

func (ir *inputRequest) IsValid() bool {
	if ir.Method != onlyValidMethod {
		return false
	}

	return true
}

type outputResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

const onlyValidMethod = "isPrime"

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed %s", conn.RemoteAddr())
	}()

	var input inputRequest
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		dec := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
		if err := dec.Decode(&input); err != nil {
			if err != io.EOF {
				log.Printf("err decoding: %v", err)
			}
			continue
		}
		log.Printf("received: %+v", input)
		if !input.IsValid() {
			continue
		}

		// Note that non-integers can not be prime.
		n, err := input.Number.Int64()
		if err != nil {
			log.Printf("err: %v", err)
			continue
		}

		output := outputResponse{
			Method: onlyValidMethod,
			Prime:  isPrime(int(n)),
		}
		enc := json.NewEncoder(conn)
		if err := enc.Encode(&output); err != nil {
			log.Printf("encoding err: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanning input: %v", err)
	}
}

func sendMalformedResponse(conn net.Conn) error {
	_, err := conn.Write([]byte(`{}`))
	return err
}

func isPrime(n int) bool {
	for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}
