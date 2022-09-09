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

const onlyValidMethod = "isPrime"

type request struct {
	Method string      `json:"method"`
	Number json.Number `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed %s", conn.RemoteAddr())
	}()

	var input request
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		dec := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
		if err := dec.Decode(&input); err != nil {
			if err != io.EOF {
				// malformed
				log.Printf("err decoding request: %v", err)
				_ = sendMalformedResponse(conn)
				break
			}
			continue
		}
		log.Printf("received: %+v", input)
		if input.Method != onlyValidMethod {
			// malformed
			log.Printf("method is not %q", onlyValidMethod)
			_ = sendMalformedResponse(conn)
			break
		}
		_ = sendResponse(conn, isPrime(input.Number))
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanning input: %v", err)
	}
}

func sendMalformedResponse(w io.Writer) error {
	// return sendResponse(w, true)
	log.Printf("sending malformed response")
	art := struct {
		Error string `json:"error"`
	}{
		Error: "malformed request",
	}
	return jsonResponse(w, art)
}

func jsonResponse(w io.Writer, i any) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(i); err != nil {
		log.Printf("encoding err: %v", err)
		return err
	}
	return nil
}

func sendResponse(w io.Writer, isPrime bool) error {
	output := response{
		Method: onlyValidMethod,
		Prime:  isPrime,
	}
	return jsonResponse(w, output)
}

func isPrime(jn json.Number) bool {
	n, err := jn.Int64()
	if err != nil {
		log.Println(err)
		return false
	}

	for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
		if int(n)%i == 0 {
			return false
		}
	}
	return n > 1
}
