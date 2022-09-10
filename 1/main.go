// Prime Time
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
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
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

// UnmarshalJSON implements a custom json.Unmarshaler that will detect missing fields
func (r *request) UnmarshalJSON(b []byte) error {
	m := struct {
		Method *string  `json:"method"`
		Number *float64 `json:"number"`
	}{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	if m.Method == nil || m.Number == nil {
		return errors.New("missing required fields")
	}
	r.Method = *m.Method
	r.Number = *m.Number
	return nil
}

func isIntegral(f float64) bool {
	return f == float64(int(f))
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func (r *response) IsValid() bool {
	return true
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		log.Printf("closed %s", conn.RemoteAddr())
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var r request
		log.Printf("received %v", scanner.Text())
		dec := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
		if err := dec.Decode(&r); err != nil {
			if err == io.EOF {
				continue
			}
			// malformed
			log.Printf("err decoding request: %v", err)
			_ = sendMalformedResponse(conn)
			break
		}
		if r.Method != onlyValidMethod {
			// malformed
			log.Printf("method is not %q", onlyValidMethod)
			_ = sendMalformedResponse(conn)
			break
		}

		var p bool
		// should have a compliant request
		if isIntegral(r.Number) {
			p = isPrime(int(r.Number))
		} else {
			p = false
		}
		_ = sendResponse(conn, p)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanning input: %v", err)
	}
}

func sendMalformedResponse(w io.Writer) error {
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

// isPrime implements the Sieve of Eratosthenes
func isPrime(n int) bool {
	for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
		if int(n)%i == 0 {
			return false
		}
	}
	return n > 1
}
